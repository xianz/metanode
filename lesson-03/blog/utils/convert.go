package utils

import "strconv"

func StrToUint(s string) (uint, error) {
	u, err := strconv.ParseUint(s, 10, 64)
	return uint(u), err
}
