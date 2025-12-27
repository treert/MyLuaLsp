package ast

/*
exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef |
	 prefixexp | tableconstructor | exp binop exp | unop exp

prefixexp ::= var | functioncall | ‘(’ exp ‘)’

var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name

functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
*/

// NilExp nil
type NilExp struct {
	ExpBase
}

// A BadExpr node is a placeholder for an expression containing
// syntax errors for which a correct expression node cannot be
// created.
type BadExpr struct {
	ExpBase
}

type TrueExp struct {
	ExpBase
}

// FalseExp false
type FalseExp struct {
	ExpBase
}

// VarargExp ...
type VarargExp struct {
	ExpBase
}

// IntegerExp 整数
type IntegerExp struct {
	ExpBase
	Val int64
}

// FloatExp 浮点数
type FloatExp struct {
	ExpBase
	Val float64
}

// StringExp 字符串
type StringExp struct {
	ExpBase
	Str string
}

// UnopExp := unop exp
type UnopExp struct {
	ExpBase
	Op  TkKind // operator
	Exp Exp
}

// BinopExp := exp1 op exp2
type BinopExp struct {
	ExpBase
	Op   TkKind // operator
	Exp1 Exp
	Exp2 Exp
}

// tableconstructor ::= ‘{’ [fieldlist] ‘}’
//
// fieldlist ::= field {fieldsep field} [fieldsep]
//
// field ::= ‘[’ exp ‘]’ ‘=’ exp | Name ‘=’ exp | exp
//
// fieldsep ::= ‘,’ | ‘;’
type TableConstructorExp struct {
	ExpBase
	KeyExps []Exp
	ValExps []Exp
}

// funcbody ::= ‘(’ [parlist] ‘)’ block end
//
// parlist ::= namelist [‘,’ ‘...’] | ‘...’
//
// namelist ::= Name {‘,’ Name}
type FuncDefExp struct {
	ExpBase
	ParList  []Token
	Block    *Block
	IsVararg bool // 是否是...可变参数
	IsColon  bool // 是否为: 这样的函数
}

/*
prefixexp ::= Name |
              ‘(’ exp ‘)’ |
              prefixexp ‘[’ exp ‘]’ |
              prefixexp ‘.’ Name |
              prefixexp ‘:’ Name args |
              prefixexp args
*/

// NameExp 引用其他变量
type NameExp struct {
	ExpBase
	Name string
}

// ParensExp 括号包含表达式或值
type ParensExp struct {
	ExpBase
	Exp Exp
}

// TableAccessExp 成员变量获取
type TableAccessExp struct {
	ExpBase
	PrefixExp  Exp
	KeyExp     Exp
	IsWriteExp bool
}

// FuncCallExp 函数调用
// 当调用这样的函数时 aaa:bb("1", "2")
// 其中aaa 为 PrefixExp， bb 为NameExp，括号内的为参数
type FuncCallExp struct {
	ExpBase
	PrefixExp Exp
	NameExp   *StringExp
	Args      []Exp
}
