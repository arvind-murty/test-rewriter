package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func add(a int, b int) int {
	return a + b
}

func parseTest(t *testing.T) {
	fileName := "test.go"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, nil, 0)
	require.NoError(t, err, "Could not parse file: %v", fileName)

	a := make([]int, 5)
	require.NoError(t, err, "the values: %v", a...)

}

func secondParseTest(t *testing.T) {
	a := 5
	b := 6
	c := add(a, b)
	require.Equal(t, a, yuh.b, "a: %v and b: %v are not equal", a, yuh.b)

	if d := add(b, c); d == a {
		t.Fatalf("d: %v and a: %v are not equal", d, a)
	}

	if a != b || a != c {
		t.Fatalf("a: %v is not equal to b: %v or c: %v", a, b, c)
	}

	ye := true
	yo := false
	yeh := true
	require.True(t, ye, "ye: %v is not true", ye)
	require.False(t, ye, "ye: %v is not true", ye)

	if ye && yo && yeh {
		t.Fatalf("idk")
	}
	require.Equal(t, a, 5, "a: %v is not equal to 5", a)

	require.Equal(t, a, 5, "a: %v is not equal to 5", a)

	require.Equal(t, a, 5,
		"a: %v is not equal to 5", a)
}
