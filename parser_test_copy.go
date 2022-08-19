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
	if err != nil {
		t.Fatalf("Could not parse file: %v", fileName)
	}

	a := make([]int, 5)

	if err != nil {
		t.Fatalf("the values: %v", a...)
	}
}

func secondParseTest(t *testing.T) {
	a := 5
	b := 6
	c := add(a, b)
	if a != yuh.b {
		t.Fatalf("a: %v and b: %v are not equal", a, yuh.b)
	}

	if d := add(b, c); d == a {
		t.Fatalf("d: %v and a: %v are not equal", d, a)
	}

	if a != b || a != c {
		t.Fatalf("a: %v is not equal to b: %v or c: %v", a, b, c)
	}

	ye := true
	yo := false
	yeh := true
	if !ye {
		t.Fatalf("ye: %v is not true", ye)
	}

	if ye {
		t.Fatalf("ye: %v is not true", ye)
	}

	if ye && yo && yeh {
		t.Fatalf("idk")
	}

	if a != 5 {
		t.Fatalf("a: %v is not equal to 5", a)
	}

	require.Equal(t, a, 5, "a: %v is not equal to 5", a)

	require.Equal(t, a, 5,
		"a: %v is not equal to 5", a)
}
