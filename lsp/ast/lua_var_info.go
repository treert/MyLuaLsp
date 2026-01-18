package ast

// LuaType 变量定义的类型
type LuaType uint8

const (
	_ LuaType = iota
	// LuaTypeNil nil
	LuaTypeNil

	// LuaTypeBool bool
	LuaTypeBool

	// LuaTypeNumber int64, float64
	LuaTypeNumber

	// LuaTypeInter int64
	LuaTypeInter

	// LuaTypeFloat float64
	LuaTypeFloat

	// LuaTypeString string
	LuaTypeString

	// LuaTypeTable table
	LuaTypeTable

	// LuaTypeArray array
	LuaTypeArray

	// LuaTypeFunc 函数
	LuaTypeFunc

	// LuaTypeRefer 引用其他的
	LuaTypeRefer

	// LuaTypeAll 什么都有可能是
	LuaTypeAll
)

// ReferExpFlag 值lua变量关联的exp标记，// 默认值为0，指向的ReferExp是否为empty，例如定义的时候 a = nil，值为1, 当被赋值后，值为2
type ReferExpFlag uint8

const (
	// ReferExpFlagDefault 初始值为0
	ReferExpFlagDefault ReferExpFlag = 0

	// ReferExpFlagEmpty 当lua变量的ReferExp定义为nil时候，此时标记为1
	ReferExpFlagEmpty ReferExpFlag = 1

	// ReferExpFlagAsign 当lua变量的ReferExp之前定义为nil，这是又进行了赋值，值为2
	ReferExpFlagAsign ReferExpFlag = 2
)

// VarInfoList 列表存放
type VarInfoList struct {
	VarVec []*VarInfo // 整体引用的vector，按先后顺序插入到其中
}
