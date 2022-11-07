# test-rewriter

This tool uses the `go/ast` package to rewrite the old `t.Fatalf` and `t.Errorf` testing formats to the new `require` and `assert` packages.

To run it on a particular file `file_name.go`, do `./fatal_rewriter file_name.go` and `./error_rewriter file_name.go`

## Limitations

I do not think `if` statements with initialization clauses can be rewritten in this way because of scope issues.

## Issues

Comments sometimes get moved inside the new `require` or `assert` statements.

## TODO

Deal with more complicated `if` statements involving `||` and `&&`.
