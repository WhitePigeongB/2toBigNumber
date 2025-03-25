package main

import (
	"fmt"
	"time"
)

func main() {
	var start = time.Now()
	const times = 2025
	var result = "1"
	for i := 1; i <= times/3; i++ {
		result = mult8(result)
	}
	for i := 1; i <= times%3; i++ {
		result = twice(result)
	}
	fmt.Println(result)
	fmt.Printf("time in seconds %v", time.Since(start))
}

func twice(input string) string {
	var a = ""
	var remain = 0
	var num = 0
	var mult = 0
	var inputDigit int = 0
	var chars []byte = []byte(input)
	var length = len(chars)
	for i := 1; i <= length; i++ {
		inputDigit = int(chars[length-i]) - 48
		mult = inputDigit*2 + remain
		num = mult % 10
		remain = mult / 10
		a = string(num+48) + a
	}
	if remain > 0 {
		a = string(remain+48) + a
	}
	return a
}

func mult8(input string) string {
	var a = ""
	var remain = 0
	var num = 0
	var mult = 0
	var inputDigit int = 0
	var chars []byte = []byte(input)
	var length = len(chars)
	for i := 1; i <= length; i++ {
		inputDigit = int(chars[length-i]) - 48
		mult = inputDigit*8 + remain
		num = mult % 10
		remain = mult / 10
		a = string(num+48) + a
	}
	if remain > 0 {
		a = string(remain+48) + a
	}
	return a
}
