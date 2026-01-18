package ast

/*
lua 类型注释的语法：

TypeName ::= Name | nil | boolean | number | string | any | void | '{' {file_name:TypeName ,} '}'
TypeName ::= TypeName|TypeName
TypeName ::= TypeName[]
TypeName ::= table<TypeName,TypeName>
TypeName ::= fun( {param_name?:TypeName} ) [:TypeName{,TypeName} ]


---@type TypeName{,TypeName}


---@param param_name? TypeName
---@vararg TypeName

---@return TypeName {, TypeName}

---@class TypeName : TypeName {, TypeName}
---@field [public|protected|private] field_name? TypeName @ 这个感觉没有意义呀，lua只能注释性提示下了。

---@alias new_name TypeName

---@alias new_name
---| XX
---| XX


泛型

---@generic T[: TypeName] {, T[: TypeName]}

TypeName ::= Name<TypeName{,TypeName}>

其他

---@language LangID
---@lang LangID



*/

// ATokenType 类型
type ATokenType int8

const (
	ATokenEOF          ATokenType = iota // end-of-file
	ATokenSepComma                       // ,
	ATokenSepColon                       // :
	ATokenVararg                         // ... 函数的可变参数
	ATokenVSepLparen                     // (
	ATokenVSepRparen                     // )
	ATokenVSepLbrack                     // [
	ATokenVSepRbrack                     // ]
	ATokenBor                            // |
	ATokenLt                             // <
	ATokenGt                             // >
	ATokenAt                             // @
	ATokenOption                         // ?
	ATokenString                         // 定义的其他字符串
	ATokenKwFun                          // fun
	ATokenKwTable                        // table
	ATokenKwType                         // type
	ATokenKwParam                        // param
	ATokenKwField                        // field
	ATokenKwClass                        // class
	ATokenKwReturn                       // return
	ATokenKwOverload                     // overload
	ATokenKwAlias                        // alias
	ATokenKwGeneric                      // generic
	ATokenKwPubic                        // public
	ATokenKwProtected                    // protected
	ATokenKwPrivate                      // private
	ATokenKwVararg                       // vararg
	ATokenKwIdentifier                   // identifier
	ATokenKwConst                        // const
	ATokenKwOther                        // other token， not valid
	ATokenKwEnum                         // enum 枚举段关键值
	ATokenKwEnumStart                    // start enum后面跟着的开始关键字，例如完整的为enum start
	ATokenKwEnumEnd                      // end enum后面跟着的结束关键字，例如完整的为enum end
)

var Annotate_Keywords = map[string]ATokenType{
	"fun":       ATokenKwFun,
	"table":     ATokenKwTable,
	"type":      ATokenKwType,
	"param":     ATokenKwParam,
	"field":     ATokenKwField,
	"class":     ATokenKwClass,
	"return":    ATokenKwReturn,
	"overload":  ATokenKwOverload,
	"alias":     ATokenKwAlias,
	"generic":   ATokenKwGeneric,
	"public":    ATokenKwPubic,
	"protected": ATokenKwProtected,
	"private":   ATokenKwPrivate,
	"vararg":    ATokenKwVararg,
	"const":     ATokenKwConst,
	"enum":      ATokenKwEnum,
}
