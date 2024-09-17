package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

type arrRange []string

func (i *arrRange) String() string {
	return ""
}

func (i *arrRange) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

func findPrimes(s string, ch chan<- []int) {

	ss := strings.Split(s, ":")

	var arrRange = [2]int{}
	var err error
	arrRange[0], err = strconv.Atoi(ss[0])
	if err != nil {
		fmt.Println("Неверный формат диапазона", err)
		os.Exit(1)
	}

	arrRange[1], err = strconv.Atoi(ss[1])
	if err != nil {
		fmt.Println("Неверный формат диапазона", err)
		os.Exit(1)
	}

	primes := []int{}
	for i := arrRange[0]; i < arrRange[1]; i++ {
		if big.NewInt(int64(i)).ProbablyPrime(0) {
			primes = append(primes, i)
		}
	}
	ch <- primes
}

func writeFile(file os.File, ch <-chan [][]int) {
	primes := <-ch

	file.
	defer file.Close()

	for i := 0; i < len(primes); i++ {
		file.WriteString(fmt.Sprintf("Простые числа из диапазона %d:\n", i+1))
		file.WriteString(strings.Trim(strings.Join(strings.Split(fmt.Sprint(primes[i]), " "), ","), "[]"))

		file.WriteString("\n")
	}
}

func main() {
	filename := flag.String("file", "", "имя файла для вывода найденных простых чисел")
	timeout := *flag.Int("timeout", 0, "значение в секундах, по истечении которого программа должна прекратить свое исполнение")
	var ranges arrRange
	flag.Var(&ranges, "range", "диапазоны чисел, в пределах которого программа должна найти простые числа")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	file, err := os.Create(*filename)
	if err != nil {
		fmt.Println("Невозможно создать файл:", err)
		os.Exit(1)
	}
	defer file.Close()

	primeCh := make(chan []int)
	allPrimesCh := make(chan [][]int)

	go writeFile(*filename, allPrimesCh)

	var primes [][]int
	for _, v := range ranges {
		go findPrimes(v, primeCh)
		primes = append(primes, <-primeCh)
	}
	fmt.Println(primes)

	allPrimesCh <- primes

	// time.Sleep(time.Second * 3)
}
