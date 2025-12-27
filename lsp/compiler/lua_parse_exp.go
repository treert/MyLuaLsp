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
	//return parseExp12(l)
	// return p.parseSubExp(0)
	var exp = &ast.BadExpr{}
	exp.Loc = p.nowToken.Loc
	return exp
}

func (p *Parser) parsePrefixExp() ast.Exp {
	panic("")
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

func (p *Parser) _createTableAccessExp(prefixExp ast.Exp, keyExp ast.Exp) *ast.TableAccessExp {
	var exp = &ast.TableAccessExp{
		PrefixExp: prefixExp,
		KeyExp:    keyExp,
	}
	exp.SetLoc(common.GetRangeLoc(prefixExp.GetLoc(), keyExp.GetLoc()))
	return exp
}

func (p *Parser) parseFuncBodyExp() *ast.FuncDefExp {
	panic("")
}
