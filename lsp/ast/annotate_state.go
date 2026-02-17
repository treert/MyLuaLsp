package ast

/*
参考对象：

lua 类型注释的语法：

LiteralValue ::= nil | true | false | <number> | <string>

TypeName ::= LiteralValue
TypeName ::= '{' {file_name:TypeName ,} '}'
TypeName ::= TypeName[]
TypeName ::= TypeName|TypeName
TypeName ::= Name | any | number | string

TypeName ::= fun( {param_name?:TypeName} ) [:TypeName{,TypeName} ]

TypeName ::= Name<TypeName{,TypeName}>
TypeName ::= table<TypeName,TypeName>  # 模板的一种特例

泛型

---@generic T[: TypeName] {, T[: TypeName]}


---@type TypeName{,TypeName}


---@param param_name? TypeName

---@return TypeName {, TypeName}

---@class TypeName : TypeName {, TypeName}
---@field [public|protected|private] field_name? TypeName @ 这个感觉没有意义呀，lua只能注释性提示下了。

---@alias new_name TypeName

---@enum new_name
new_name = {
	XX = Value, -- Value 最好是字面值，当然也可以是 exp.
}

*/

type TypeBase any

type LiteralValueType int8

const (
	LiteralValueNil LiteralValueType = iota
	LiteralValueTrue
	LiteralValueFalse
	LiteralValueNumber
	LiteralValueString
)

// LiteralValue ::= nil | true | false | <number> | <string>
type Type_LiteralValue struct {
	Num  float64
	Str  string
	Bool bool
	Type LiteralValueType
}

// name: TypeName 用于非常多地方
type Type_KeyValue struct {
	NameAndLoc NameAndLoc
	Type       TypeBase
}

// TypeName ::= '{' {file_name:TypeName ,} '}'
type Type_Map struct {
	FieldList []Type_KeyValue
}

// TypeName ::= TypeName[]
type Type_Array struct {
	ElementType TypeBase
}

// TypeName ::= TypeName|TypeName
type Type_Union struct {
	TypeList []TypeBase
}

// 最终要通过这个名字找到具体的类型
type Type_Identifier struct {
	NameAndLoc     NameAndLoc // 例如：any, number, string, 或者用户定义的类型
	IsGenericParam bool       // 是否是泛型类型名，例如：T, U
}

type Type_FunParam struct {
	NameAndLoc NameAndLoc
	Type       TypeBase
	IsOptional bool // 是否是可选参数
	Comment    string
}

type Type_FunReturn struct {
	Type    TypeBase
	Comment string
}

// 有两种定义方式：1. 注释里定义 2. lua语法定义的函数。
type Type_Fun struct {
	ParamList  []Type_FunParam
	ReturnList []Type_FunReturn
	Comment    string
	// todo 关联 lua 语法定义的函数信息
}

// TypeName ::= Name<TypeName{,TypeName}>
type Type_GenericInstance struct {
	NameAndLoc    NameAndLoc // 通过这个名字找到泛型定义，Type_Class
	ParamTypeList []TypeBase // 泛型实例化的参数类型列表
}

type Type_ClassField struct {
	NameAndLoc NameAndLoc
	Type       TypeBase
	Comment    string
	IsOptional bool // 是否是可选字段
}

// @class 构建的类型
type Type_Class struct {
	NameAndLoc       NameAndLoc
	ParentTypeList   []TypeBase
	GenericParamList []Type_KeyValue // 泛型类有这个
	FieldList        []Type_ClassField
	Comment          string
}

type Type_Alias struct {
	NameAndLoc NameAndLoc
	Type       TypeBase
	Comment    string
}

type Type_Enum struct {
	NameAndLoc NameAndLoc
	Comment    string
	// todo 关联 lua 定义的table，其中是实际定义的枚举成员
}

/////////////////// 以下是行注释语句片段 /////////////////////////

// ---@generic T[: TypeName] {, T[: TypeName]}
//
//	例如：
//	---@generic T, U: number
type AnnotateGenericState struct {
	ParamList []Type_KeyValue // 泛型参数列表，NameAndLoc.Name 是泛型参数的名字，Type 是泛型参数的约束类型，可以为 nil，表示没有约束
}

// ---@type TypeName{,TypeName}
//
//	例如：
//	---@type number, stirng
//	local a, b
type AnnotateTypeState struct {
	TypeList []TypeBase
	Comment  string
}

// ---@param param_name? TypeName
//
//	例如：
//	---@param a number
type AnnotateParamState struct {
	NameAndLoc NameAndLoc // 参数的名字和位置，如果是可变参数，NameAndLoc.Name == ...
	ParamType  TypeBase
	IsOptional bool // 是否是可选参数
	Comment    string
}

// ---@return TypeName {, TypeName}
//
//	例如：
//	---@return number, string
type AnnotateReturnState struct {
	ReturnTypeList []TypeBase
	Comment        string
}

// ---@class TypeName : TypeName {, TypeName}
//
//	例如：
//	---@class Person : Human, Animal
type AnnotateClassState struct {
	NameAndLoc     NameAndLoc
	ParentTypeList []TypeBase
	Comment        string
}

// ---@field field_name? TypeName
//
//	例如：
//	---@field a number
type AnnotateFieldState struct {
	NameAndLoc NameAndLoc
	FieldType  TypeBase
	Comment    string
	IsOptional bool
}

// ---@alias new_name TypeName
//
//	例如：
//	---@alias MyType number | string
type AnnotateAliasState struct {
	NameAndLoc NameAndLoc
	Type       TypeBase
	Comment    string
}

// ---@enum new_name
//
//	例如：
//	---@enum Color
//	Color = {
//		Red = 1,
//		Green = 2,
//		Blue = 3,
//	}
type AnnotateEnumState struct {
	NameAndLoc NameAndLoc
	Comment    string
}
