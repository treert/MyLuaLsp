package ast

import "mylua-lsp/lsp/common"

// LuaAstBase 语法树基础结构，提供父节点指针
type LuaAstBase struct {
	Loc    common.Location
	Parent Stat
}

type ExpBase struct {
	LuaAstBase
	Type any // 关联的类型信息
}

type Stat interface {
	GetLoc() common.Location
	GetParent() Stat
	IsExp() bool

	SetLoc(loc common.Location)
	SetParent(parent Stat)
}

type Exp interface {
	Stat
	GetParentExp() Exp
}

func (b *LuaAstBase) GetLoc() common.Location {
	return b.Loc
}

func (b *LuaAstBase) SetLoc(loc common.Location) {
	b.Loc = loc
}

func (b *LuaAstBase) SetParent(parent Stat) {
	b.Parent = parent
}

func (b *LuaAstBase) GetParent() Stat {
	return b.Parent
}

func (b *LuaAstBase) IsExp() bool {
	return false
}

func (b *ExpBase) GetParentExp() Exp {
	var p, _ = b.Parent.(Exp)
	return p
}

func (b *ExpBase) IsExp() bool {
	return true
}

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
