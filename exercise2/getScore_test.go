package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZero(t *testing.T) {
	fmt.Println("Проверка на offset = 0")
	offset := 0
	stamps := fillScores()
	expected := Score{0, 0}
	actual := getScore(*stamps, offset)

	if !assert.Equal(t, expected, actual) {
		fmt.Printf("Тест с offset %d провален\nПолучили %v, должно быть %v\n", offset, actual, expected)
	} else {
		fmt.Printf("Тест с offset %d пройден\nПолучили %v, должно быть %v\n", offset, actual, expected)
	}
}

func TestBeforeMatch(t *testing.T) {
	fmt.Println("Проверка на offset = -1")
	offset := -1
	stamps := fillScores()
	expected := Score{0, 0}
	actual := getScore(*stamps, offset)

	if !assert.Equal(t, expected, actual) {
		fmt.Printf("Тест с offset %d провален\nПолучили %v, должно быть %v\n", offset, actual, expected)
	} else {
		fmt.Printf("Тест с offset %d пройден\nПолучили %v, должно быть %v\n", offset, actual, expected)
	}
}

func TestAfterMatch(t *testing.T) {
	fmt.Println("Проверка на offset = 150001")
	stampsPointer := fillScores()
	stamps := *stampsPointer
	offset := 150001
	expected := stamps[TIMESTAMPS_COUNT-1].Score
	actual := getScore(stamps, offset)

	if !assert.Equal(t, expected, actual) {
		fmt.Printf("Тест с offset %d провален\nПолучили %v, должно быть %v\n", offset, actual, expected)
	} else {
		fmt.Printf("Тест с offset %d пройден\nПолучили %v, должно быть %v\n", offset, actual, expected)
	}
}

func TestDuringMatch(t *testing.T) {
	stampsPointer := fillScores()
	stamps := *stampsPointer
	startMatch := stamps[TIMESTAMPS_COUNT/4].Offset
	midMatch := stamps[TIMESTAMPS_COUNT/2].Offset
	endMatch := stamps[TIMESTAMPS_COUNT*3/4].Offset

	var tests = []struct {
		offset   int
		expected Score
	}{
		{startMatch, stamps[TIMESTAMPS_COUNT/4].Score},
		{midMatch, stamps[TIMESTAMPS_COUNT/2].Score},
		{endMatch, stamps[TIMESTAMPS_COUNT*3/4].Score},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("Проверка на offset = %d", tt.offset)
		t.Run(testname, func(t *testing.T) {
			actual := getScore(stamps, tt.offset)
			if !assert.Equal(t, tt.expected, actual) {
				fmt.Printf("Тест с offset %d провален\nПолучили %v, должно быть %v\n", tt.offset, actual, tt.expected)
			} else {
				fmt.Printf("Тест с offset %d пройден\n", tt.offset)
			}
		})
	}
}
