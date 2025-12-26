package main

import (
	"mylua-lsp/lsp/ast"
	"mylua-lsp/lsp/common"
	"strings"
)

func main() {
	println("hello world")
	println(common.StringToBytes("hi"))
	println(common.BytesToString([]byte("hi")))
	var builder strings.Builder
	println(builder.String())

	var stat ast.Stat = &ast.Block{}
	var stat2 ast.Stat = &ast.BreakStat{}
	stat2.SetParent(stat)
	var stat3 ast.Stat = &ast.TrueExp{}
	var exp ast.Exp = &ast.TrueExp{}
	var exp2 ast.Exp = &ast.NilExp{}
	exp2.SetParent(exp)
	println("exp info:", exp, exp.IsExp(), exp.GetParentExp())
	println("exp2 info:", exp2.IsExp(), exp2.GetParentExp())
	println("stat2 info:", stat2.IsExp(), stat2.GetParent())
	x1, _ := stat.(ast.Exp)
	println(x1, x1 == nil)
	x2, _ := stat3.(ast.Exp)
	println(x2, x2 == nil)
	x3, _ := exp.(ast.Stat)
	println(x3, x3 == nil, &x3, &stat, stat)
}
