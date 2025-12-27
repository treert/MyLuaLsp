package ast

// chunk ::= block
// type Chunk *Block

// Block code block
// block ::= {stat} [retstat]
// retstat ::= return [explist] [';']
// explist ::= exp {',' exp}}
type Block struct {
	LuaAstBase
	Stats []Stat
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

// BreakStat break语句
// break
type BreakStat struct {
	LuaAstBase
}

// ‘::’ Name ‘::’
type LabelStat struct {
	LuaAstBase
	Name Token
}

// goto Name
type GotoStat struct {
	LuaAstBase
	Name Token
}

// do block end
type DoStat struct {
	LuaAstBase
	Block *Block
}

// if exp then block {elseif exp then block} [else block] end
type IfStat struct {
	LuaAstBase
	Exps   []Exp
	Blocks []*Block // 如果最后时 else 的话，Blocks长度会比Exps多1
}

// while exp do block end
type WhileStat struct {
	LuaAstBase
	Exp   Exp
	Block *Block
}

// 不推荐使用 repeat 不是好的设计
// repeat block until exp
type RepeatStat struct {
	LuaAstBase
	Block *Block
	Exp   Exp
}

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
type ForNumStat struct {
	LuaAstBase
	VarName  Token
	InitExp  Exp
	LimitExp Exp
	StepExp  Exp
	Block    *Block
}

// for namelist in explist do block end
//
// namelist ::= Name {‘,’ Name}
//
// explist ::= exp {‘,’ exp}
type ForInStat struct {
	LuaAstBase
	NameList []Token
	ExpList  []Exp
	Block    *Block
}

// varlist ‘=’ explist
//
// varlist ::= var {‘,’ var}
//
// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
type AssignStat struct {
	LuaAstBase
	VarList []Exp
	ExpList []Exp
}

// local namelist [‘=’ explist]
//
// namelist ::= Name {‘,’ Name}
//
// explist ::= exp {‘,’ exp}
type LocalVarDeclStat struct {
	LuaAstBase
	NameList []Token
	ExpList  []Exp
}

// LocalFuncDefStat local function Name funcbody
type LocalFuncDefStat struct {
	LuaAstBase
	Name    Token
	FuncDef *FuncDefExp
}

type RetStat struct {
	LuaAstBase
	ExpList []Exp
}
