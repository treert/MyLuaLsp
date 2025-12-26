package compiler

import (
	"mylua-lsp/lsp/ast"
	"mylua-lsp/lsp/common"
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
	var start_loc = p.aheadToken.Loc
	p.NextTokenKind(ast.TkSepLabel)
	p.NextIdentifier()
	var token = p.nowToken
	p.NextTokenKind(ast.TkSepLabel)
	var stat = &ast.LabelStat{
		Name: token,
	}
	stat.Loc = common.GetRangeLoc(&start_loc, &p.nowToken.Loc)
	return stat
}

// goto Name
func (p *Parser) parseGotoStat() ast.Stat {
	var start_loc = p.aheadToken.Loc
	p.NextTokenKind(ast.TkKwGoto)
	p.NextIdentifier()
	var stat = &ast.GotoStat{
		Name: p.nowToken,
	}
	stat.Loc = common.GetRangeLoc(&start_loc, &p.nowToken.Loc)
	return stat
}

// do block end
func (p *Parser) parseDoStat() ast.Stat {
	var start_loc = p.aheadToken.Loc
	p.NextTokenKind(ast.TkKwDo)
	block := p.parseBlock()
	p.NextTokenKind(ast.TkKwEnd)
	var stat = &ast.DoStat{
		Block: block,
	}
	stat.Loc = common.GetRangeLoc(&start_loc, &p.nowToken.Loc)
	return stat
}

func (p *Parser) parseWhileStat() ast.Stat {
	panic("unimplemented")
}

func (p *Parser) parseRepeatStat() ast.Stat {
	panic("unimplemented")
}

func (p *Parser) parseIfStat() ast.Stat {
	panic("unimplemented")
}

func (p *Parser) parseForStat() ast.Stat {
	panic("unimplemented")
}

func (p *Parser) parseFuncDefStat() ast.Stat {
	panic("unimplemented")
}

func (p *Parser) parseLocalAssignOrFuncDefStat() ast.Stat {
	panic("unimplemented")
}

func (p *Parser) parseIKIllegalStat() ast.Stat {
	panic("unimplemented")
}

func (p *Parser) parseAssignOrFuncCallStat() ast.Stat {
	panic("unimplemented")
}
