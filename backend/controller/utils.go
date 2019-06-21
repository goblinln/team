package controller

import (
	"errors"
	"strconv"
)

func atoi(param string) int64 {
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0
	}

	return id
}

func intArray(param []string) []int {
	ret := []int{}

	if param == nil || len(param) == 0 {
		return ret
	}

	for _, v := range param {
		id, err := strconv.Atoi(v)
		if err == nil {
			ret = append(ret, id)
		}
	}

	return ret
}

func assert(test bool, msg string) {
	if !test {
		panic(errors.New(msg))
	}
}
