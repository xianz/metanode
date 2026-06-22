package utils

import (
	"fmt"
	"regexp"
	"runtime"
)

func LogWithLocation(msg string) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}
	fmt.Printf("[%s:%d] %s\n", file, line, msg)
}

func RegexpMatch(oriStr string, match string) bool {
	return regexp.MustCompile(match).MatchString(oriStr)
}
