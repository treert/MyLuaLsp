package ast

import "mylua-lsp/lsp/common"

type Position = common.Position
type Location = common.Location

// ParseError check的错误信息
type ParseError struct {
	ErrStr string          // 简单的错误信息
	Loc    common.Location // 错误的位置区域
}

// TooManyErr 当Parse太多语法错误的时候，终止
type TooManyErr struct {
	ErrNum int //  错误的数量
}

// 注释信息。如果是长注释，要求结尾没有其他token，否则当其不存在。
type CommentLine struct {
	Str       string          // 注释内容。如果是长注释，会做预处理。
	StartPos  common.Position // 注释开始的位置
	ShortFlag bool            // 是否是短注释，true表示短注释
	HeadFlag  bool            // 是否为头部注释，一行开头就是注释
}

// 连续的注释块。是后续注释解析的高层单位
type CommentBlock struct {
	List []CommentLine
}

// NameAndLoc 名字和位置，会是一个非常常用的结构
type NameAndLoc struct {
	Name string
	Loc  Location
}
