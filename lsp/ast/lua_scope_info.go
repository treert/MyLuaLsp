package ast

// ScopeInfo 作用域信息
type ScopeInfo struct {
	Parent    *ScopeInfo   // 当前作用域的父作用域
	SubScopes []*ScopeInfo // 所有子的ScopeInfos

	// 当前作用域内的变量信息
	VarInfoList []*VarInfo

	// 当前作用域内的label标签信息
	LabelInfoList []*LabelInfo

	// 当前作用域内的函数信息
	FuncInfoList []*FuncInfo
}
