package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
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

var wg sync.WaitGroup

func writeFile(ctx context.Context, filename string, ch <-chan []int) {
	wg.Add(1)
	defer wg.Done()

	primes := <-ch

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Невозможно создать файл:", err)
		os.Exit(1)
	}
	defer file.Close()

	select {
	case <-ctx.Done():
		fmt.Println("Выход из программы:", ctx.Err())
		os.Exit(1)
	default:
		file.WriteString(fmt.Sprintf("Простые числа из диапазона %d:%d это: %s\n", primes[0], primes[1], strings.Trim(strings.Join(strings.Split(fmt.Sprint(primes[2:]), " "), ", "), "[]")))
	}

}

func main() {
	filename := flag.String("file", "", "имя файла для вывода найденных простых чисел")
	timeout := flag.Int("timeout", 10, "значение в секундах, по истечении которого программа должна прекратить свое исполнение")
	var ranges arrRangeStr
	flag.Var(&ranges, "range", "диапазоны чисел, в пределах которого программа должна найти простые числа")

	flag.Parse()

	if *timeout < 0 {
		fmt.Println("Введен отрицательный таймаут")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	if *filename == "" {
		fmt.Println("Неверное имя файла")
		os.Exit(1)
	}

	primeCh := make(chan []int)

	for _, v := range ranges {

		go func(ctx context.Context, rangeStr string, ch chan<- []int) {
			wg.Add(1)
			defer wg.Done()

			strSl := strings.Split(rangeStr, ":")

			intSl := make([]int, 2)
			var err error
			intSl[0], err = strconv.Atoi(strSl[0])
			if err != nil {
				fmt.Println("Неверный формат диапазона", err)
				os.Exit(1)
			}

			intSl[1], err = strconv.Atoi(strSl[1])
			if err != nil {
				fmt.Println("Неверный формат диапазона", err)
				os.Exit(1)
			}

			if intSl[0] > intSl[1] {
				fmt.Println("Левая граница диапазона меньше правой")
				os.Exit(1)
			}

			for i := intSl[0]; i <= intSl[1]; i++ {
				if big.NewInt(int64(i)).ProbablyPrime(0) {
					intSl = append(intSl, i)
				}
			}

			select {
			case <-ctx.Done():
				fmt.Println("Выход из программы:", ctx.Err())
				os.Exit(1)
			default:
				ch <- intSl
			}

		}(ctx, v, primeCh)

		go writeFile(ctx, *filename, primeCh)

	}

	wg.Wait()

	select {
	case <-ctx.Done():
		fmt.Println("Main: Выход из программы. Окончен таймаут:", ctx.Err())
		os.Exit(1)
	default:
		fmt.Println("Main: Прогамма завершила работу")
	}
}
