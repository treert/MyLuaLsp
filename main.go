package main

import (
	"mylua-lsp/lsp/common"
	"strings"
)

func main() {
	println("hello world")
	println(common.StringToBytes("hi"))
	println(common.BytesToString([]byte("hi")))
	var builder strings.Builder
	println(builder.String())

}
