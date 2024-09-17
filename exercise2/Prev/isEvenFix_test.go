package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroFix(t *testing.T) {
	num := 0
	fmt.Printf("Проверка числа %d на четность\n", num)
	expected := true
	actual := isEvenFix(num)

	if !assert.Equal(t, expected, actual) {
		fmt.Printf("Тест с числом %d провален\nПолучили %t, должно быть %t", num, actual, expected)
	} else {
		fmt.Printf("Тест с числом %d пройден\nПолучили %t, должно быть %t", num, actual, expected)
	}
}

func TestTableFix(t *testing.T) {
	var tests = []struct {
		num      int
		expected bool
	}{
		{-243, false},
		{-100, true},
		{-3, false},
		{-2, true},
		{-1, false},
		{1, false},
		{2, true},
		{3, false},
		{8, true},
		{156, true},
		{1569, false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("Проверка числа %d на четность", tt.num)
		t.Run(testname, func(t *testing.T) {
			actual := isEvenFix(tt.num)
			if !assert.Equal(t, tt.expected, actual) {
				fmt.Printf("Тест с числом %d провален\nПолучили %t, должно быть %t", tt.num, actual, tt.expected)
			} else {
				fmt.Printf("Тест с числом %d пройден\nПолучили %t, должно быть %t", tt.num, actual, tt.expected)
			}
		})
	}
}
