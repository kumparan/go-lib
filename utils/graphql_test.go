package utils

import (
	"fmt"
	"math"
	"testing"
)

func Test_Int642GraphQLID(t *testing.T) {
	g := Int642GraphQLID(int64(math.MaxInt64))
	res := fmt.Sprintf("%d", math.MaxInt64)

	if res != string(g) {
		t.Error("different value")
	}
}

func Test_GraphQLID2String(t *testing.T) {
	g := GraphQLID2Int64(Int642GraphQLID(math.MaxInt64))
	if int64(math.MaxInt64) != g {
		t.Error("different value")
	}
}
