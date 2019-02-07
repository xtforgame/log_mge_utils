package lmu

import (
	"fmt"
	"strconv"
)

func IterationNameToNumber(iterationName string) uint64 {
	var iterationNumber uint64
	var err error
	if iterationNumber, err = strconv.ParseUint(iterationName, 16, 64); err != nil {
		iterationNumber = 0
	}
	return iterationNumber
}

func IterationNumberToName(iterationNumber uint64) string {
	return fmt.Sprintf("%016x", iterationNumber)
}
