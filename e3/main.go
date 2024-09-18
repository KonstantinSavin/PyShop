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

type arrRangeStr []string

func (i *arrRangeStr) String() string {
	return ""
}

func (i *arrRangeStr) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}

func rangeToArr(rangeStr string) (arrRangeInt [2]int) {
	ss := strings.Split(rangeStr, ":")

	var err error
	arrRangeInt[0], err = strconv.Atoi(ss[0])
	if err != nil {
		fmt.Println("rangeToArr -> Неверный формат диапазона", err)
		os.Exit(1)
	}

	arrRangeInt[1], err = strconv.Atoi(ss[1])
	if err != nil {
		fmt.Println("rangeToArr -> Неверный формат диапазона", err)
		os.Exit(1)
	}

	fmt.Println("rangeToArr -> arrRange: ", arrRangeInt)

	return arrRangeInt
}

func findPrimes(ctx context.Context, arrRangeInt [2]int, ch chan<- []int) {
	primes := []int{}
	for i := arrRangeInt[0]; i <= arrRangeInt[1]; i++ {
		if big.NewInt(int64(i)).ProbablyPrime(0) {
			primes = append(primes, i)
		}
	}

	fmt.Println("findPrimes -> primes: ", primes)

	ch <- primes
}

func writeFile(ctx context.Context, filename string, arrRangeInt [2]int, ch <-chan []int) {
	primes := <-ch
	fmt.Println("writeFile -> primes:", primes)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("writeFile -> Невозможно создать файл:", err)
		os.Exit(1)
	}
	defer file.Close()

	file.WriteString(fmt.Sprintf("Простые числа из диапазона %d:%d это %s\n", arrRangeInt[0], arrRangeInt[1], strings.Trim(strings.Join(strings.Split(fmt.Sprint(primes), " "), ", "), "[]")))
}

func main() {
	filename := flag.String("file", "", "имя файла для вывода найденных простых чисел")
	timeout := flag.Int("timeout", 0, "значение в секундах, по истечении которого программа должна прекратить свое исполнение")
	var ranges arrRangeStr
	flag.Var(&ranges, "range", "диапазоны чисел, в пределах которого программа должна найти простые числа")

	flag.Parse()

	fmt.Println("filename: ", *filename)
	fmt.Println("timeout: ", *timeout)
	fmt.Println("ranges: ", ranges)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)

	defer func() {
		cancel()
		fmt.Println("Main Defer: canceling context")
	}()

	primeCh := make(chan []int)

	for _, v := range ranges {
		arrRangeInt := rangeToArr(v)
		go writeFile(ctx, *filename, arrRangeInt, primeCh)
		go findPrimes(ctx, arrRangeInt, primeCh)
	}

	time.Sleep(time.Second * 10)
}
