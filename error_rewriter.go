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

	parseErrorExpr(file)
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

// checks that n represents an if statement with a t.Errorf call within
func checkErrorSignature(n ast.Node) bool {
	if ifstmt, ok := n.(*ast.IfStmt); ok {
		if expr, ok1 := ifstmt.Body.List[0].(*ast.ExprStmt); ok1 && len(ifstmt.Body.List) == 1 {
			if call, ok2 := expr.X.(*ast.CallExpr); ok2 {
				if selector, ok3 := call.Fun.(*ast.SelectorExpr); ok3 {
					if id, ok4 := selector.X.(*ast.Ident); ok4 {
						if id.Name == "t" && selector.Sel.Name == "Errorf" {
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
func getNewErrorSignature(errorSignature *ast.CallExpr, funcName string, newArgs []ast.Expr) *ast.ExprStmt {
	// get rid of old position for old args (only for BasicLit, Ident, and SelectorExpr cases for now)
	for _, arg := range errorSignature.Args {
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

	extraArgs := append(newArgs, errorSignature.Args...)
	newFunc := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "assert",
				},
				Sel: &ast.Ident{
					Name: funcName,
				},
			},
			Args:     extraArgs,
			Ellipsis: errorSignature.Ellipsis,
		},
	}
	return newFunc
}

func parseErrorExpr(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if block, ok := n.(*ast.BlockStmt); ok {
			for i, stmt := range block.List {
				if checkErrorSignature(stmt) {
					errorSignature := stmt.(*ast.IfStmt).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
					// ignore if statements with initializations
					if stmt.(*ast.IfStmt).Init == nil {
						if cond, ok1 := stmt.(*ast.IfStmt).Cond.(*ast.BinaryExpr); ok1 {
							x, ok2 := cond.X.(*ast.Ident)
							y, ok3 := cond.Y.(*ast.Ident)
							if ok3 {
								if ok2 {
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
										block.List[i] = getNewErrorSignature(errorSignature, "NoError", newArgs)
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
										block.List[i] = getNewErrorSignature(errorSignature, "Error", newArgs)
										continue
									}
								}

								// require.Nil
								if cond.Op == token.NEQ && y.Name == "nil" {
									newArgs := []ast.Expr{
										&ast.Ident{
											Name: "t",
										},
										cond.X,
									}
									block.List[i] = getNewErrorSignature(errorSignature, "Nil", newArgs)
									continue
								}
								// require.NotNil
								if cond.Op == token.EQL && y.Name == "nil" {
									newArgs := []ast.Expr{
										&ast.Ident{
											Name: "t",
										},
										cond.X,
									}
									block.List[i] = getNewErrorSignature(errorSignature, "NotNil", newArgs)
									continue
								}
							}

							newArgs := []ast.Expr{
								&ast.Ident{
									Name: "t",
								},
								cond.Y,
								cond.X,
							}
							// require.Equal and require.NotEqual
							if cond.Op == token.NEQ {
								block.List[i] = getNewErrorSignature(errorSignature, "Equal", newArgs)
							} else if cond.Op == token.EQL {
								block.List[i] = getNewErrorSignature(errorSignature, "NotEqual", newArgs)
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
								block.List[i] = getNewErrorSignature(errorSignature, "True", newArgs)
							}
						} else if cond, ok1 := stmt.(*ast.IfStmt).Cond.(*ast.Ident); ok1 {
							newArgs := []ast.Expr{
								&ast.Ident{
									Name: "t",
								},
								cond,
							}
							// require.False
							block.List[i] = getNewErrorSignature(errorSignature, "False", newArgs)
						}
					}
				}
			}
		}

		return true
	})
}
