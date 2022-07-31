package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	stations := getStations()

	stationMap := make(map[rune][]Station)
	for _, station := range stations {
		r, _ := utf8.DecodeRuneInString(station.Furigana)
		if stationMap[r] == nil {
			stationMap[r] = make([]Station, 0)
		}
		stationMap[r] = append(stationMap[r], station)
	}

	siritoris := make([]*Stack[Station], len(stations))
	for i, station := range stations {
		siritoris[i] = NewStack(station)
	}
	for {
		newSiritoris := make([]*Stack[Station], 0)
		for _, siritori := range siritoris {
			lastStation := siritori.Cur()
			lastRune, _ := utf8.DecodeLastRuneInString(lastStation.Furigana)
			candidates := stationMap[lastRune]
			for _, candidate := range candidates {
				doubled := false
				siritori.Each(func(s Station) {
					if s == candidate {
						doubled = true
					}
				})
				if doubled {
					continue
				}

				newSiritori := siritori.Push(candidate)
				newSiritoris = append(newSiritoris, newSiritori)
			}
		}

		if len(newSiritoris) == 0 {
			break
		}
		siritoris = newSiritoris
	}

	for _, siritori := range siritoris {
		siritoriArray := siritori.ToArray()
		stationNames := make([]string, len(siritoriArray))
		for i, station := range siritoriArray {
			stationNames[i] = station.Name
		}
		fmt.Println(strings.Join(stationNames, " "))
	}
}

type Station struct {
	Name     string
	Furigana string
}

func getStations() []Station {
	csvFile, err := os.Open("stations.csv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reader := csv.NewReader(csvFile)
	_, err = reader.Read() // header
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stations := make([]Station, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		station := Station{record[0], record[1]}
		stations = append(stations, station)
	}
	return stations
}

type Stack[T any] struct {
	parent *Stack[T]
	value  T
	length int
}

func NewStack[T any](rootValue T) *Stack[T] {
	return &Stack[T]{nil, rootValue, 1}
}

func (s *Stack[T]) Push(value T) *Stack[T] {
	return &Stack[T]{s, value, s.length + 1}
}

func (s *Stack[T]) Cur() T {
	return s.value
}

func (s *Stack[T]) Each(f func(T)) {
	current := s
	for current != nil {
		f(current.value)
		current = current.parent
	}
}

func (s *Stack[T]) ToArray() []T {
	a := make([]T, s.length)
	i := s.length - 1
	s.Each(func(value T) {
		a[i] = value
		i--
	})
	return a
}
