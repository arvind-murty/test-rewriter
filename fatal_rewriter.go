package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Need the file and only the file to be parsed as a command line argument")
	}
	fileName := os.Args[1]

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	parseFatalExpr(file)
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = printer.Fprint(f, fset, file)
	if err != nil {
		log.Fatal(err)
	}

	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
}

// checks that n represents an if statement with a t.Fatalf call within
func checkFatalSignature(n ast.Node) bool {
	if ifstmt, ok := n.(*ast.IfStmt); ok {
		if expr, ok1 := ifstmt.Body.List[0].(*ast.ExprStmt); ok1 && len(ifstmt.Body.List) == 1 {
			if call, ok2 := expr.X.(*ast.CallExpr); ok2 {
				if selector, ok3 := call.Fun.(*ast.SelectorExpr); ok3 {
					if id, ok4 := selector.X.(*ast.Ident); ok4 {
						if id.Name == "t" && selector.Sel.Name == "Fatalf" {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

// gets the signature for the new require function call
func getNewFatalSignature(fatalSignature *ast.CallExpr, funcName string, newArgs []ast.Expr) *ast.ExprStmt {
	// get rid of old position for old args (only for BasicLit, Ident, and SelectorExpr cases for now)
	for _, arg := range fatalSignature.Args {
		if basicLit, ok := arg.(*ast.BasicLit); ok {
			basicLit.ValuePos = token.NoPos
		} else if ident, ok := arg.(*ast.Ident); ok {
			ident.NamePos = token.NoPos
		} else if selectorExpr, ok := arg.(*ast.SelectorExpr); ok {
			if x, ok1 := selectorExpr.X.(*ast.Ident); ok1 {
				x.NamePos = token.NoPos
				selectorExpr.Sel.NamePos = token.NoPos
			}
		}
	}

	extraArgs := append(newArgs, fatalSignature.Args...)
	newFunc := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "require",
				},
				Sel: &ast.Ident{
					Name: funcName,
				},
			},
			Args:     extraArgs,
			Ellipsis: fatalSignature.Ellipsis,
		},
	}
	return newFunc
}

func parseFatalExpr(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if block, ok := n.(*ast.BlockStmt); ok {
			for i, stmt := range block.List {
				if checkFatalSignature(stmt) {
					fatalSignature := stmt.(*ast.IfStmt).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
					// ignore if statements with initializations
					if stmt.(*ast.IfStmt).Init == nil {
						if cond, ok1 := stmt.(*ast.IfStmt).Cond.(*ast.BinaryExpr); ok1 {
							x, ok2 := cond.X.(*ast.Ident)
							y, ok3 := cond.Y.(*ast.Ident)
							if ok2 && ok3 {
								// require.NoError
								if x.Name == "err" && cond.Op == token.NEQ && y.Name == "nil" {
									newArgs := []ast.Expr{
										&ast.Ident{
											Name: "t",
										},
										&ast.Ident{
											Name: "err",
										},
									}
									block.List[i] = getNewFatalSignature(fatalSignature, "NoError", newArgs)
									continue
								}
								// require.Error
								if x.Name == "err" && cond.Op == token.EQL && y.Name == "nil" {
									newArgs := []ast.Expr{
										&ast.Ident{
											Name: "t",
										},
										&ast.Ident{
											Name: "err",
										},
									}
									block.List[i] = getNewFatalSignature(fatalSignature, "Error", newArgs)
									continue
								}
							}
							// for some reason this adds a new line after these extraArgs and before the ones copied from Fatalf
							newArgs := []ast.Expr{
								&ast.Ident{
									Name: "t",
								},
								cond.X,
								cond.Y,
							}
							// require.Equal and require.NotEqual
							if cond.Op == token.NEQ {
								block.List[i] = getNewFatalSignature(fatalSignature, "Equal", newArgs)
							} else if cond.Op == token.EQL {
								block.List[i] = getNewFatalSignature(fatalSignature, "NotEqual", newArgs)
							}
						} else if cond, ok1 := stmt.(*ast.IfStmt).Cond.(*ast.UnaryExpr); ok1 {
							newArgs := []ast.Expr{
								&ast.Ident{
									Name: "t",
								},
								cond.X,
							}
							// require.True
							if cond.Op == token.NOT {
								block.List[i] = getNewFatalSignature(fatalSignature, "True", newArgs)
							}
						} else if cond, ok1 := stmt.(*ast.IfStmt).Cond.(*ast.Ident); ok1 {
							newArgs := []ast.Expr{
								&ast.Ident{
									Name: "t",
								},
								cond,
							}
							// require.False
							block.List[i] = getNewFatalSignature(fatalSignature, "False", newArgs)
						}
					}
				}
			}
		}

		return true
	})
}
