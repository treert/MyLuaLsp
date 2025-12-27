package compiler

import (
	"mylua-lsp/lsp/ast"
)

func (p *Parser) parseStat() ast.Stat {
	switch p.LookAheadKind() {
	case ast.TkSepSemi:
		return nil
	case ast.TkKwBreak:
		return p.parseBreakStat()
	case ast.TkSepLabel:
		return p.parseLabelStat()
	case ast.TkKwGoto:
		return p.parseGotoStat()
	case ast.TkKwDo:
		return p.parseDoStat()
	case ast.TkKwWhile:
		return p.parseWhileStat()
	case ast.TkKwRepeat:
		return p.parseRepeatStat()
	case ast.TkKwIf:
		return p.parseIfStat()
	case ast.TkKwFor:
		return p.parseForStat()
	case ast.TkKwFunction:
		return p.parseFuncDefStat()
	case ast.TkKwLocal:
		return p.parseLocalAssignOrFuncDefStat()
	case ast.TkKwReturn:
		return p.parseRetStat()
	case ast.IKIllegal:
		return p.parseIKIllegalStat()
	default:
		return p.parseAssignOrFuncCallStat()
	}
}

func (p *Parser) parseBreakStat() ast.Stat {
	p.NextTokenKind(ast.TkKwBreak)
	var stat = &ast.BreakStat{}
	stat.Loc = p.nowToken.Loc
	return stat
}

// ‘::’ Name ‘::’
func (p *Parser) parseLabelStat() ast.Stat {
	p.NextTokenKind(ast.TkSepLabel)
	p.NextIdentifier()
	var token = p.nowToken
	p.NextTokenKind(ast.TkSepLabel)
	var stat = &ast.LabelStat{
		Name: token,
	}
	return stat
}

// goto Name
func (p *Parser) parseGotoStat() ast.Stat {
	p.NextTokenKind(ast.TkKwGoto)
	p.NextIdentifier()
	var stat = &ast.GotoStat{
		Name: p.nowToken,
	}
	return stat
}

// do block end
func (p *Parser) parseDoStat() ast.Stat {
	p.NextTokenKind(ast.TkKwDo)
	block := p.parseBlock()
	p.NextTokenKind(ast.TkKwEnd)
	var stat = &ast.DoStat{
		Block: block,
	}
	return stat
}

// while exp do block end
func (p *Parser) parseWhileStat() ast.Stat {
	p.NextTokenKind(ast.TkKwWhile)
	exp := p.parseExp()
	p.NextTokenKind(ast.TkKwDo)
	block := p.parseBlock()
	p.NextTokenKind(ast.TkKwEnd)
	var stat = &ast.WhileStat{
		Exp:   exp,
		Block: block,
	}
	return stat
}

// repeat block until exp
// 这个语法不好，最好不要使用
func (p *Parser) parseRepeatStat() ast.Stat {
	p.NextTokenKind(ast.TkKwRepeat)
	block := p.parseBlock()
	p.NextTokenKind(ast.TkKwUntil)
	exp := p.parseExp()
	var stat = &ast.RepeatStat{
		Block: block,
		Exp:   exp,
	}
	return stat
}

// if exp then block {elseif exp then block} [else block] end
func (p *Parser) parseIfStat() ast.Stat {
	exps := make([]ast.Exp, 0, 1)
	blocks := make([]*ast.Block, 0, 1)
	p.NextTokenKind(ast.TkKwIf)
	exps = append(exps, p.parseExp())
	p.NextTokenKind(ast.TkKwThen)
	blocks = append(blocks, p.parseBlock())
	for p.LookAheadKind() == ast.TkKwElseIf {
		p.NextToken()
		exps = append(exps, p.parseExp())
		p.NextTokenKind(ast.TkKwThen)
		blocks = append(blocks, p.parseBlock())
	}
	if p.LookAheadKind() == ast.TkKwElse {
		p.NextToken()
		blocks = append(blocks, p.parseBlock())
	}
	p.NextTokenKind(ast.TkKwEnd)
	var stat = &ast.IfStat{
		Exps:   exps,
		Blocks: blocks,
	}
	return stat
}

func (p *Parser) parseForStat() ast.Stat {
	p.NextTokenKind(ast.TkKwFor)
	p.NextIdentifier()
	var varName = p.nowToken
	if p.LookAheadKind() == ast.TkOpAssign {
		return p.finishForNumStat(&varName)
	} else {
		return p.finishForInStat(&varName)
	}
}

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
func (p *Parser) finishForNumStat(varName *Token) *ast.ForNumStat {
	p.NextTokenKind(ast.TkOpAssign)
	initExp := p.parseExp()
	p.NextTokenKind(ast.TkSepComma)
	limitExp := p.parseExp()
	var stepExp ast.Exp = nil
	if p.LookAheadKind() == ast.TkSepComma {
		p.NextToken()
		stepExp = p.parseExp()
	}
	p.NextTokenKind(ast.TkKwDo)
	block := p.parseBlock()
	p.NextTokenKind(ast.TkKwEnd)
	var stat = &ast.ForNumStat{
		VarName:  *varName,
		InitExp:  initExp,
		LimitExp: limitExp,
		StepExp:  stepExp,
		Block:    block,
	}
	return stat
}

// for namelist in explist do block end
func (p *Parser) finishForInStat(varName *Token) *ast.ForInStat {
	nameList := p._parseNameList(varName)
	p.NextTokenKind(ast.TkKwIn)
	expList := p.parseExpList()
	p.NextTokenKind(ast.TkKwDo)
	block := p.parseBlock()

	p.NextTokenKind(ast.TkKwEnd)
	var stat = &ast.ForInStat{
		NameList: nameList,
		ExpList:  expList,
		Block:    block,
	}
	return stat
}

// functiondef ::= function funcname funcbody
func (p *Parser) parseFuncDefStat() ast.Stat {
	p.NextTokenKind(ast.TkKwFunction)
	funcNameExp, isColon := p._parseFuncName()
	funcDef := p.parseFuncBodyExp()
	funcDef.IsColon = isColon
	var stat = &ast.AssignStat{
		VarList: []ast.Exp{funcNameExp},
		ExpList: []ast.Exp{funcDef},
	}
	return stat
}

// funcname ::= Name {‘.’ Name} [‘:’ Name]
func (p *Parser) _parseFuncName() (exp ast.Exp, isColon bool) {
	exp = p.parseNameExp()
	for p.LookAheadKind() == ast.TkSepDot {
		p.NextToken()
		keyExp := p.parseNameExpAsStringExp()
		exp = p._createTableAccessExp(exp, keyExp)
	}
	if p.LookAheadKind() == ast.TkSepColon {
		p.NextToken()
		keyExp := p.parseNameExpAsStringExp()
		exp = p._createTableAccessExp(exp, keyExp)
		return exp, true
	}
	return exp, false
}

func (p *Parser) parseLocalAssignOrFuncDefStat() ast.Stat {
	p.NextTokenKind(ast.TkKwLocal)
	if p.LookAheadKind() == ast.TkKwFunction {
		return p._finishLocalFuncDefStat()
	} else {
		return p._finishLocalVarDeclStat()
	}
}

// local function Name funcbody
func (p *Parser) _finishLocalFuncDefStat() *ast.LocalFuncDefStat {
	p.NextTokenKind(ast.TkKwFunction)
	p.NextIdentifier()
	var funcName = p.nowToken
	funcDef := p.parseFuncBodyExp()
	var stat = &ast.LocalFuncDefStat{
		Name:    funcName,
		FuncDef: funcDef,
	}
	return stat
}

// local name [<attribute>] {, name [<attribute>]} [‘=’ explist]
func (p *Parser) _finishLocalVarDeclStat() *ast.LocalVarDeclStat {
	nameList := make([]ast.Token, 0, 1)
	for {
		p.NextIdentifier()
		var token = p.nowToken
		token.LocalAttr = p._getLocalAttribute()
		nameList = append(nameList, token)
		if p.LookAheadKind() == ast.TkSepComma {
			p.NextToken()
		} else {
			break
		}
	}
	var expList []ast.Exp = nil
	if p.LookAheadKind() == ast.TkOpAssign {
		p.NextToken()
		expList = p.parseExpList()
	}
	var stat = &ast.LocalVarDeclStat{
		NameList: nameList,
		ExpList:  expList,
	}
	return stat
}

func (p *Parser) _getLocalAttribute() ast.LocalAttr {
	if p.LookAheadKind() == ast.TkOpLt {
		p.NextToken()
		p.NextIdentifier()
		var attrName = p.nowToken.TokenStr
		if attrName == "const" {
			p.NextTokenKind(ast.TkOpGt)
			return ast.RDKCONST
		} else if attrName == "close" {
			p.NextTokenKind(ast.TkOpGt)
			return ast.RDKTOCLOSE
		} else {
			p.l.errorPrint(p.nowToken.Loc, "unknown local attribute '%s'", attrName)
			p.NextTokenKind(ast.TkOpGt)
		}
	}
	return ast.VDKREG
}

func (p *Parser) parseIKIllegalStat() ast.Stat {
	p.NextToken()
	return nil // 忽略掉
}

func (p *Parser) parseAssignOrFuncCallStat() ast.Stat {
	prefixExp := p.parsePrefixExp()
	if _, ok := prefixExp.(*ast.BadExpr); ok {
		return nil // 直接忽略掉
	}

	if fc, ok := prefixExp.(*ast.FuncCallExp); ok {
		return fc
	}

	assignStat := p.parseAssignStat(prefixExp)
	return assignStat
}

// varlist ‘=’ explist
func (p *Parser) parseAssignStat(var0 ast.Exp) *ast.AssignStat {
	var varList = []ast.Exp{var0}
	for p.LookAheadKind() == ast.TkSepComma {
		p.NextToken()
		varList = append(varList, p.parseExp())
	}
	p.NextTokenKind(ast.TkOpAssign)
	var expList = p.parseExpList()
	return &ast.AssignStat{
		VarList: varList,
		ExpList: expList,
	}
}

// retstat ::= return [explist]
func (p *Parser) parseRetStat() ast.Stat {
	p.NextTokenKind(ast.TkKwReturn)
	var stat = &ast.RetStat{}

	switch p.LookAheadKind() {
	case ast.TkEOF, ast.TkKwEnd,
		ast.TkKwElse, ast.TkKwElseIf, ast.TkKwUntil:
		break
	case ast.TkSepSemi:
		p.NextToken()
		break
	default:
		exps := p.parseExpList()
		if p.LookAheadKind() == ast.TkSepSemi {
			p.NextToken()
		}
		stat.ExpList = exps
	}
	return stat
}

// namelist ::= Name {‘,’ Name}
func (p *Parser) _parseNameList(firstName *Token) []Token {
	nameList := make([]Token, 0, 1)
	if firstName != nil {
		nameList = append(nameList, *firstName)
	} else {
		p.NextIdentifier()
		nameList = append(nameList, p.nowToken)
	}
	for p.LookAheadKind() == ast.TkSepComma {
		p.NextToken()
		p.NextIdentifier()
		nameList = append(nameList, p.nowToken)
	}
	return nameList
}
