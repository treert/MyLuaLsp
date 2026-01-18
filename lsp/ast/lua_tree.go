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

type VarInfo struct {
}
