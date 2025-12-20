package ast

import (
	"mylua-lsp/lsp/common"
)

// chunk ::= block
// type Chunk *Block

// Block code block
// block ::= {stat} [retstat]
// retstat ::= return [explist] [';']
// explist ::= exp {',' exp}}
type Block struct {
	LuaAstBase
	Stats   []*Stat
	RetExps []*ExpBase
}

/*
stat ::=  ‘;’ |
	 varlist ‘=’ explist |
	 functioncall |
	 label |
	 break |
	 goto Name |
	 do block end |
	 while exp do block end |
	 repeat block until exp |
	 if exp then block {elseif exp then block} [else block] end |
	 for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end |
	 for namelist in explist do block end |
	 function funcname funcbody |
	 local function Name funcbody |
	 local namelist [‘=’ explist]
*/

// EmptyStat 空表达式或者错误。不会出现在 Ast 中，会被忽略掉
type EmptyStat struct{}

// BreakStat break语句
// break
type BreakStat struct {
	LuaAstBase
}

// ‘::’ Name ‘::’
type LabelStat struct {
	LuaAstBase
	Name string
}

// goto Name
type GotoStat struct {
	LuaAstBase
	Name string
}

// do block end
type DoStat struct {
	LuaAstBase
	Block *Block
}

// if exp then block {elseif exp then block} [else block] end
type IfStat struct {
	LuaAstBase
	Exps   []*ExpBase
	Blocks []*Block
}

// while exp do block end
type WhileStat struct {
	LuaAstBase
	Exp   *ExpBase
	Block *Block
}

// 不推荐使用 repeat 不是好的设计
// repeat block until exp
type RepeatStat struct {
	LuaAstBase
	Block *Block
	Exp   *ExpBase
}

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
type ForNumStat struct {
	LuaAstBase
	VarName  string
	VarLoc   common.Location
	InitExp  *ExpBase
	LimitExp *ExpBase
	StepExp  *ExpBase
	Block    *Block
}

// for namelist in explist do block end
//
// namelist ::= Name {‘,’ Name}
//
// explist ::= exp {‘,’ exp}
type ForInStat struct {
	LuaAstBase
	NameList    []string
	NameLocList []common.Location // 所有变量的位置信息
	ExpList     []*ExpBase
	Block       *Block
}

// varlist ‘=’ explist
//
// varlist ::= var {‘,’ var}
//
// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
type AssignStat struct {
	LuaAstBase
	VarList []*ExpBase
	ExpList []*ExpBase
}

// local namelist [‘=’ explist]
//
// namelist ::= Name {‘,’ Name}
//
// explist ::= exp {‘,’ exp}
type LocalVarDeclStat struct {
	NameList   []string
	VarLocList []common.Location // 所有变量的位置信息
	AttrList   []LocalAttr       // 变量的属性
	ExpList    []*ExpBase
}

// LocalFuncDefStat local function Name funcbody
type LocalFuncDefStat struct {
	LuaAstBase
	Name    string
	NameLoc common.Location // 函数名的位置信息
	Exp     *FuncDefExp
}
