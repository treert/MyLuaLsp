package ast

import "mylua-lsp/lsp/common"

// LuaAstBase 语法树基础结构，提供父节点指针
type LuaAstBase struct {
	Loc    common.Location
	Parent *LuaAstBase
}

type ExpBase struct {
	LuaAstBase
	Type any // 关联的类型信息
}

type Stat LuaAstBase

/*
lua 语法相关的一些零碎放这儿。
*/

// LocalAttr  local var attribute
type LocalAttr uint8 // attribute type for local varible

const (
	// VDKREG no attr
	VDKREG LocalAttr = 0
	// RDKTOCLOSE close attr
	RDKTOCLOSE LocalAttr = 1
	// RDKCONST const attr
	RDKCONST LocalAttr = 2
)
