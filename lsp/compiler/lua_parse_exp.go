package compiler

import (
	"mylua-lsp/lsp/ast"
	"mylua-lsp/lsp/common"
)

// explist ::= exp {‘,’ exp}
func (p *Parser) parseExpList() []ast.Exp {
	exps := make([]ast.Exp, 0, 1)
	exps = append(exps, p.parseExp())
	for p.LookAheadKind() == ast.TkSepComma {
		p.NextToken()
		exps = append(exps, p.parseExp())
	}
	return exps
}

/*
exp grammar:

	exp   ::= exp12
	exp12 ::= exp11 {or exp11}
	exp11 ::= exp10 {and exp10}
	exp10 ::= exp9 {(‘<’ | ‘>’ | ‘<=’ | ‘>=’ | ‘~=’ | ‘==’) exp9}
	exp9  ::= exp8 {‘|’ exp8}
	exp8  ::= exp7 {‘~’ exp7}
	exp7  ::= exp6 {‘&’ exp6}
	exp6  ::= exp5 {(‘<<’ | ‘>>’) exp5}
	exp5  ::= exp4 {‘..’ exp4}
	exp4  ::= exp3 {(‘+’ | ‘-’) exp3}
	exp3  ::= exp2 {(‘*’ | ‘/’ | ‘//’ | ‘%’) exp2}
	exp2  ::= {(‘not’ | ‘#’ | ‘-’ | ‘~’)} exp1
	exp1  ::= exp0 {‘^’ exp2}
	exp0  ::= nil | false | true | Numeral | LiteralString
			| ‘...’ | functiondef | prefixexp | tableconstructor
*/
func (p *Parser) parseExp() ast.Exp {
	return p.parseSubExp(0)
}

// 通过tokenKind 获取优先级
func getPriority(tokenKind ast.TkKind) int {
	switch tokenKind {
	case ast.TkOpPow: // ^
		return 12
	case ast.TkOpMul, ast.TkOpMod, ast.TkOpDiv, ast.TkOpIdiv: // *, %, /, //
		return 10
	case ast.TkOpAdd, ast.TkOpSub: // +, -
		return 9
	case ast.TkOpConcat: // ..
		return 8
	case ast.TkOpShl, ast.TkOpShr: // shift:  <<  >>
		return 7
	case ast.TkOpBand: // &
		return 6
	case ast.TkOpBxor: // x ~ y
		return 5
	case ast.TkOpBor: // x | y
		return 4
	case ast.TkOpLt, ast.TkOpGt, ast.TkOpNe,
		ast.TkOpLe, ast.TkOpGe, ast.TkOpEq: // (‘<’ | ‘>’ | ‘<=’ | ‘>=’ | ‘~=’ | ‘==’)
		return 3
	case ast.TkOpAnd: // x and y
		return 2
	case ast.TkOpOr: // x or y
		return 1
	}

	return 0
}

func (p *Parser) parseSubExp(limit int) ast.Exp {
	var exp ast.Exp
	tokenKind := p.LookAheadKind()
	start_loc := p.aheadToken.Loc
	// 单目运算符： not | # | - | ~
	if tokenKind == ast.TkOpNen || tokenKind == ast.TkOpNot || tokenKind == ast.TkOpSub || tokenKind == ast.TkOpBnot {
		p.NextToken()
		var unopExp = &ast.UnopExp{
			Op:  tokenKind,
			Exp: p.parseSubExp(10),
		}
		exp = unopExp
	} else {
		exp = p.parseExp0()
	}
	exp.SetLoc(common.GetRangeLoc(start_loc, p.nowToken.Loc))

	tokenKind = p.LookAheadKind()
	nowPriority := getPriority(tokenKind)
	for nowPriority > 0 {
		if nowPriority <= limit {
			break
		}

		if tokenKind == ast.TkOpPow || tokenKind == ast.TkOpConcat {
			nowPriority-- // 右结合
		}

		p.NextToken()
		subExp := p.parseSubExp(nowPriority)
		exp = &ast.BinopExp{
			Op:   tokenKind,
			Exp1: exp,
			Exp2: subExp,
		}
		exp.SetLoc(common.GetRangeLoc(start_loc, p.nowToken.Loc))

		tokenKind = p.LookAheadKind()
		nowPriority = getPriority(tokenKind)
	}

	return exp
}

func (p *Parser) parseExp0() ast.Exp {
	switch p.LookAheadKind() {
	case ast.TkVararg: // ...
		p.NextToken()
		return &ast.VarargExp{}
	case ast.TkKwNil: // nil
		p.NextToken()
		return &ast.NilExp{}
	case ast.TkKwTrue: // true
		p.NextToken()
		return &ast.TrueExp{}
	case ast.TkKwFalse: // false
		p.NextToken()
		return &ast.FalseExp{}
	case ast.TkString: // LiteralString
		p.NextToken()
		return &ast.StringExp{
			Str: p.nowToken.TokenStr,
		}
	case ast.TkNumber: // Numeral
		return p.parseNumberExp()
	case ast.TkSepLcurly: // tableconstructor
		return p.parseTableConstructorExp()
	case ast.TkKwFunction: // functiondef
		p.NextToken()
		beginLoc := p.nowToken.Loc
		return p.parseFuncBodyExp(beginLoc)
	default: // prefixexp
		return p.parsePrefixExp()
	}
}

func (p *Parser) parseNumberExp() ast.Exp {
	p.NextToken()
	token := p.nowToken.TokenStr
	if i, ok := parseInteger(token); ok {
		return &ast.IntegerExp{
			Val: i,
		}
	} else if f, ok := parseFloat(token); ok {
		return &ast.FloatExp{
			Val: f,
		}
	} else { // todo
		p.insertParserErr(p.nowToken.Loc, "not a number: %v", token)
		return &ast.FloatExp{
			Val: 0,
		}
	}
}

func (p *Parser) parseNameExp() *ast.NameExp {
	p.NextIdentifier()
	var exp = &ast.NameExp{
		Name: p.nowToken.TokenStr,
	}
	exp.SetLoc(p.nowToken.Loc)
	return exp
}

func (p *Parser) parseNameExpAsStringExp() *ast.StringExp {
	p.NextIdentifier()
	var exp = &ast.StringExp{
		Str: p.nowToken.TokenStr,
	}
	exp.SetLoc(p.nowToken.Loc)
	return exp
}

func (p *Parser) parseStringExp() *ast.StringExp {
	p.NextTokenKind(ast.TkString)
	var exp = &ast.StringExp{
		Str: p.nowToken.TokenStr,
	}
	exp.SetLoc(p.nowToken.Loc)
	return exp
}

func (p *Parser) _createTableAccessExp(prefixExp ast.Exp, keyExp ast.Exp) *ast.TableAccessExp {
	var exp = &ast.TableAccessExp{
		PrefixExp: prefixExp,
		KeyExp:    keyExp,
	}
	exp.SetLoc(common.GetRangeLoc(prefixExp.GetLoc(), keyExp.GetLoc()))
	return exp
}

// funcbody ::= ‘(’ [parlist] ‘)’ block end
func (p *Parser) parseFuncBodyExp(func_keyword_loc Location) *ast.FuncDefExp {
	start_loc := p.nowToken.Loc
	p.NextTokenKind(ast.TkSepLparen)
	parList, isVararg := p._parseParList()
	p.NextTokenKind(ast.TkSepRparen)
	var block = p.parseBlock()
	p.NextTokenKindExpectByLoc(ast.TkKwEnd, func_keyword_loc)
	var exp = &ast.FuncDefExp{
		ParList:  parList,
		Block:    block,
		IsVararg: isVararg,
	}
	exp.SetLoc(common.GetRangeLoc(start_loc, block.GetLoc()))
	return exp
}

// parlist ::= namelist [‘,’ ‘...’] | ‘...’
func (p *Parser) _parseParList() (parList []ast.Token, isVararg bool) {
	if p.LookAheadKind() == ast.TkSepRparen {
		return nil, false
	} else if p.LookAheadKind() == ast.TkVararg {
		p.NextToken()
		return nil, true
	}

	parList = make([]ast.Token, 0, 1)
	p.NextIdentifier()
	parList = append(parList, p.nowToken)

	for p.LookAheadKind() == ast.TkSepComma {
		p.NextToken()
		if p.LookAheadKind() == ast.TkIdentifier {
			p.NextToken()
			parList = append(parList, p.nowToken)
		} else {
			p.NextTokenKind(ast.TkVararg)
			isVararg = true
			break
		}
	}
	return parList, isVararg
}

// tableconstructor ::= ‘{’ [fieldlist] ‘}’
func (p *Parser) parseTableConstructorExp() *ast.TableConstructorExp {
	var start_loc = p.nowToken.Loc
	p.NextTokenKind(ast.TkSepLcurly)       // {
	keyExps, valExps := p.parseFieldList() // [fieldlist]
	p.NextTokenKind(ast.TkSepRcurly)       // }

	var exp = &ast.TableConstructorExp{
		KeyExps: keyExps,
		ValExps: valExps,
	}
	exp.SetLoc(common.GetRangeLoc(start_loc, p.nowToken.Loc))
	return exp
}

// fieldlist ::= field {fieldsep field} [fieldsep]
func (p *Parser) parseFieldList() (ks, vs []ast.Exp) {
	if p.LookAheadKind() != ast.TkSepRcurly {
		k, v := p.parseField()
		ks = append(ks, k)
		vs = append(vs, v)

		for _isFieldSep(p.LookAheadKind()) {
			p.NextToken()
			if p.LookAheadKind() != ast.TkSepRcurly {
				k, v := p.parseField()
				ks = append(ks, k)
				vs = append(vs, v)
			} else {
				break
			}
		}
	}
	return
}

// fieldsep ::= ‘,’ | ‘;’
func _isFieldSep(tokenKind ast.TkKind) bool {
	return tokenKind == ast.TkSepComma || tokenKind == ast.TkSepSemi
}

// field ::= ‘[’ exp ‘]’ ‘=’ exp | Name ‘=’ exp | exp
func (p *Parser) parseField() (k, v ast.Exp) {
	if p.LookAheadKind() == ast.TkSepLbrack {
		p.NextToken()                    // [
		k = p.parseExp()                 // exp
		p.NextTokenKind(ast.TkSepRbrack) // ]
		p.NextTokenKind(ast.TkOpAssign)  // =
		v = p.parseExp()                 // exp
		return
	}

	exp := p.parseExp()
	if nameExp, ok := exp.(*ast.NameExp); ok {
		//loc := l.GetHeardTokenLoc()
		if p.LookAheadKind() == ast.TkOpAssign {
			// Name ‘=’ exp => ‘[’ LiteralString ‘]’ = exp
			p.NextToken()

			k = &ast.StringExp{
				Str: nameExp.Name,
			}
			k.SetLoc(nameExp.GetLoc())
			v = p.parseExp()
			return
		}
	}

	return nil, exp
}

// prefixexp grammar:
//
//	prefixexp ::= var | functioncall | ‘(’ exp ‘)’
//	var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
//	functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
//
//	prefixexp ::= Name
//	| ‘(’ exp ‘)’
//	| prefixexp ‘[’ exp ‘]’
//	| prefixexp ‘.’ Name
//	| prefixexp [‘:’ Name] args
func (p *Parser) parsePrefixExp() ast.Exp {
	var exp ast.Exp
	beginLoc := p.aheadToken.Loc
	aheadKind := p.LookAheadKind()
	switch aheadKind {
	case ast.TkIdentifier:
		exp = p.parseNameExp()
	case ast.TkSepLparen: // ‘(’ exp ‘)’
		exp = p.parseParensExp()
	default:
		p.NextToken()
		exp = &ast.BadExpr{}
		exp.SetLoc(p.nowToken.Loc)
		p.insertParserErr(p.nowToken.Loc, "`%s` can not start prefixexp", aheadKind.String())
	}
	return p.finishPrefixExp(exp, beginLoc)
}

func (p *Parser) parseParensExp() ast.Exp {
	start_loc := p.aheadToken.Loc
	p.NextTokenKind(ast.TkSepLparen) // (
	exp := p.parseExp()              // exp
	p.NextTokenKind(ast.TkSepRparen) // )

	switch exp.(type) {
	case *ast.VarargExp, *ast.FuncCallExp, *ast.NameExp, *ast.TableAccessExp:
		loc := common.GetRangeLoc(start_loc, p.nowToken.Loc)
		exp = &ast.ParensExp{
			Exp: exp,
		}
		exp.SetLoc(loc)
	default:
		// do nothing
		// no need to keep parens
	}

	return exp
}

func (p *Parser) finishPrefixExp(exp ast.Exp, beginLoc Location) ast.Exp {
	for {
		switch p.LookAheadKind() {
		case ast.TkSepLbrack: // prefixexp ‘[’ exp ‘]’
			p.NextToken()                    // ‘[’
			keyExp := p.parseExp()           // exp
			p.NextTokenKind(ast.TkSepRbrack) // ‘]’
			exp = &ast.TableAccessExp{
				PrefixExp: exp,
				KeyExp:    keyExp,
			}
		case ast.TkSepDot: // prefixexp ‘.’ Name
			p.NextToken() // ‘.’
			var keyExp ast.Exp
			if p.LookAheadKind() == ast.TkIdentifier {
				keyExp = p.parseNameExpAsStringExp()
			} else {
				p.insertParserErr(p.nowToken.Loc, "missing field or attribute names")
				keyExp = &ast.BadExpr{}
				keyExp.SetLoc(p.nowToken.Loc)
			}
			exp = &ast.TableAccessExp{
				PrefixExp: exp,
				KeyExp:    keyExp,
			}
		case ast.TkSepColon, // prefixexp ‘:’ Name args
			ast.TkSepLparen, ast.TkSepLcurly, ast.TkString: // prefixexp args
			exp = p.finishFuncCallExp(exp)
		default:
			return exp
		}
		exp.SetLoc(common.GetRangeLoc(beginLoc, p.nowToken.Loc))
	}
	//return exp
}

// functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
func (p *Parser) finishFuncCallExp(prefixExp ast.Exp) *ast.FuncCallExp {
	var nameExp *ast.StringExp
	if p.LookAheadKind() == ast.TkSepColon { // prefixexp ‘:’ Name args
		p.NextToken() // :
		if p.LookAheadKind() != ast.TkIdentifier {
			p.insertParserErr(p.nowToken.Loc, "missing method name after ':'")
			nameExp = &ast.StringExp{
				Str: "",
			}
			nameExp.SetLoc(p.nowToken.Loc)
		} else {
			nameExp = p.parseNameExpAsStringExp()
		}
	}
	args := p.parseArgs()
	return &ast.FuncCallExp{
		PrefixExp: prefixExp,
		NameExp:   nameExp,
		Args:      args,
	}
}

// args ::=  ‘(’ [explist] ‘)’ | tableconstructor | LiteralString
func (p *Parser) parseArgs() (args []ast.Exp) {
	switch p.LookAheadKind() {
	case ast.TkSepLparen: // ‘(’ [explist] ‘)’
		p.NextToken()
		if p.LookAheadKind() != ast.TkSepRparen {
			args = p.parseExpList()
		}
		p.NextTokenKind(ast.TkSepRparen)
	case ast.TkSepLcurly: // ‘{’ [fieldlist] ‘}’
		args = []ast.Exp{p.parseTableConstructorExp()}
	case ast.TkString: // LiteralString
		args = []ast.Exp{p.parseStringExp()}
	default: // LiteralString
		p.insertParserErr(p.nowToken.Loc, "missing function call args")
	}
	return
}
