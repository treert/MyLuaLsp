package ast

import "mylua-lsp/lsp/common"

// LabelInfo 定义label标签，lua goto语法用到
type LabelInfo struct {
	Name    string          // 名称
	ScopeLv int             // 当前的作用域层级，初始值为0
	Loc     common.Location // 位置信息
}

type ReturnInfo struct {
}

type FuncInfo struct {
	Parent        *FuncInfo     // 父函数
	LabelInfoList []*LabelInfo  // 函数内的label标签信息
	ReturnList    []*ReturnInfo // 函数的返回值列表，支持多返回值

}
