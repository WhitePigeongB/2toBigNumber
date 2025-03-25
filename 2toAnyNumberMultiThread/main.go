package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

func main() {
	var start = time.Now()
	var power int = 1000000
	var result = powerOfTwo(power)
	fmt.Println(toString(result))
	fmt.Printf("time in seconds %v", time.Since(start))
}

func toString(a []byte) string {
	for i := 0; i < len(a); i++ {
		a[i] += 48
	}
	return string(a)
}

func powerOfTwo(power int) []byte {
	var len = int(math.Log((float64(power)))/math.Log(2)) + 1
	temp := power
	var digits []byte = make([]byte, len)
	var remain byte = 0
	realLen := 0
	for i := 0; i < len && temp > 0; i++ {
		remain = byte(temp % 2)
		digits[i] = remain
		temp /= 2
		if temp == 0 {
			realLen = i + 1
		}
	}

	var begingPower [1]byte = [1]byte{2}
	var beginNumber [1]byte = [1]byte{1}
	powerOf2toN := begingPower[:]
	number := beginNumber[:]
	c := make(chan []byte)
	d := make(chan []byte)
	for i := 0; i < realLen; i++ {
		if digits[i] == 1 {
			go Mult(number, powerOf2toN, c)
		}
		if i < realLen-1 {
			go Mult(powerOf2toN, powerOf2toN, d)
		}
		if digits[i] == 1 {
			number = <-c
		}
		if i < realLen-1 {
			powerOf2toN = <-d
		}
	}
	close(c)
	close(d)
	return number
}

func Sum(a, b []byte, d chan []byte) {
	var lengthA = len(a)
	var lengthB = len(b)
	var maxLen int
	var minLen int
	var moreDigitsNumber *[]byte
	if lengthA > lengthB {
		maxLen = lengthA
		minLen = lengthB
		moreDigitsNumber = &a
	} else {
		maxLen = lengthB
		minLen = lengthA
		moreDigitsNumber = &b
	}
	var result []byte = make([]byte, maxLen+1)
	var remain byte = 0
	var digit byte
	var sumDigits byte
	for i := 0; i < minLen; i++ {
		sumDigits = a[lengthA-i-1] + b[lengthB-i-1] + remain
		digit = sumDigits % 10
		remain = sumDigits / 10
		result[maxLen-i] = digit
	}
	if lengthA == lengthB {
		result[0] = remain
	} else {
		var biggerNumber = *moreDigitsNumber
		for i := minLen; i < maxLen; i++ {
			sumDigits = biggerNumber[maxLen-i-1] + remain
			digit = sumDigits % 10
			remain = sumDigits / 10
			result[maxLen-i] = digit
		}
		result[0] = remain
	}
	if result[0] == 0 {
		d <- result[1:]
	} else {
		d <- result
	}
}

func SumTable(multTable [][]byte, beginIndex, endIndex int, e chan []byte) {
	d := make(chan []byte)
	switch {
	case beginIndex >= endIndex || len(multTable)-1 < endIndex:
		log.Fatal(errors.New("begin and end indexes are not correct"))
	case beginIndex+1 == endIndex:
		go Sum(multTable[beginIndex], multTable[endIndex], d)
		e <- <-d
	case beginIndex+1 < endIndex:
		f := make(chan []byte)
		middle := (beginIndex + endIndex) / 2
		if beginIndex < middle && middle+1 < endIndex {
			go SumTable(multTable, beginIndex, middle, f)
			go SumTable(multTable, middle+1, endIndex, f)
			a, b := <-f, <-f
			go Sum(a, b, d)
			e <- <-d
		} else if beginIndex == middle {
			go SumTable(multTable, middle+1, endIndex, f)
			a := <-f
			go Sum(multTable[beginIndex], a, d)
			e <- <-d
		} else if middle+1 == endIndex {
			go SumTable(multTable, beginIndex, middle, f)
			a := <-f
			go Sum(a, multTable[endIndex], d)
			e <- <-d
		}
		close(f)
	}
	close(d)
}

func Mult(a, b []byte, r chan []byte) {
	const groupCount int = 1024
	var lengthB = len(b)
	var sum []byte
	var beginNumber = [1]byte{0}
	sum = beginNumber[:]
	var sumResult []byte
	var multTable [groupCount][]byte
	var digit byte
	c := make(chan []byte)
	d := make(chan []byte)
	e := make(chan []byte)
	if lengthB >= groupCount {
		count := lengthB / groupCount
		for j := 0; j < count; j++ {
			for i := 0; i < groupCount; i++ {
				digit = b[lengthB-i-j*groupCount-1]
				go MultipleByOneDigit(a, digit, i+j*groupCount, c)
			}
			for i := 0; i < groupCount; i++ {
				multTable[i] = <-c
			}
			go SumTable(multTable[:], 0, groupCount-1, d)
			sumResult = <-d
			go Sum(sum, sumResult, e)
			sum = <-e
		}
	}
	if lengthB%groupCount == 1 {
		digit = b[0]
		go MultipleByOneDigit(a, digit, lengthB-1, c)
		resultMultLastDigit := <-c
		go Sum(sum, resultMultLastDigit, e)
		sum = <-e
	} else if lengthB%groupCount > 1 {
		count := lengthB % groupCount
		for i := 0; i < count; i++ {
			digit = b[count-i-1]
			go MultipleByOneDigit(a, digit, lengthB-count+i, c)
		}
		for i := 0; i < count; i++ {
			multTable[i] = <-c
		}
		go SumTable(multTable[:], 0, count-1, d)
		sumResult = <-d
		go Sum(sum, sumResult, e)
		sum = <-e
	}
	close(c)
	close(d)
	close(e)
	r <- sum
}

func MultipleByOneDigit(a []byte, b byte, location int, c chan []byte) {
	var lengthA = len(a)
	var len = lengthA + location + 1
	var result []byte = make([]byte, len)
	var remain byte = 0
	var digit byte = 0
	var sumDigits byte
	for i := 0; i < location; i++ {
		result[len-i-1] = 0
	}
	for i := 0; i < len-location-1; i++ {
		sumDigits = a[lengthA-i-1]*b + remain
		digit = sumDigits % 10
		remain = sumDigits / 10
		result[len-i-1-location] = digit
	}
	if remain > 0 {
		result[0] = remain
		c <- result
	} else {
		c <- result[1:]
	}
}
