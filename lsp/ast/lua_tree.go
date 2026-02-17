package ast

import (
	"mylua-lsp/lsp/common"
)

// FileInfo 文件信息
type FileInfo struct {
	Source      *common.LuaSource   // 源码信息
	Block       *Block              // 整个ast结构
	ParseErrors []ParseError        // 解析错误列表
	MainFunc    *FuncInfo           // ast生成的主function
	GlobalMaps  map[string]*VarInfo // 所有的全局信息, 包含没有_G的与含有_G前缀的变量

}

// VarInfo 变量信息。lua 里定义的变量
type VarInfo struct {
}

// Lua 表达式的值，语义分析的结果。主要是抽取常量级别的信息，语义分析就能推导出来的值。
type LuaValue struct {
}
