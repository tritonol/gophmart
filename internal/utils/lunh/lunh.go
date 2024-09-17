package lunh

import (
	"errors"
	"strconv"
)

var ErrNotValid = errors.New("luhn: data not valid")

func Validate(str string) (int64, error) {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return -1, ErrNotValid
	}

	if !valid(num) {
		return -1, ErrNotValid
	}

	return num, nil
}

func valid(number int64) bool {
	return (number%10+check(number/10))%10 == 0
}

func check(number int64) int64 {
	var luhn int64

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
