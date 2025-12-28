package compiler

import (
	"fmt"
	"mylua-lsp/lsp/ast"
	"mylua-lsp/lsp/common"
)

func ParseLuaSource(source *common.LuaSource) (block *ast.Block, commentMap LuaCommentMap, errList []ParseError) {
	parser := Parser{}
	lexer := NewLexer(source, parser.insertErr)
	parser.l = lexer
	parser.aheadToken = lexer.GetNowToken() // lexer 已经准备好了第一个token了

	defer func() {
		if err1 := recover(); err1 != nil {
			// 太多简单的语法错误，丢弃分析结果，只保留错误列表
			errList = parser.parseErrs
			return
		}
	}()

	block = parser.parseBlock() // block

	parser.NextTokenKind(ast.TkEOF)

	return block, lexer.GetCommentMap(), parser.parseErrs
}

type Parser struct {
	// 词法分析器对象
	l *Lexer

	nowToken   Token
	aheadToken Token

	parseErrs []ParseError
}

// NextToken 读取下一个单词
func (p *Parser) NextToken() {
	if p.aheadToken.Valid {
		p.nowToken = p.aheadToken
		p.aheadToken.Valid = false
		return
	}
	p.nowToken = p.l.NextToken()
}

// LookAheadToken 查看下一个单词，但不移动位置
func (p *Parser) LookAheadToken() Token {
	if p.aheadToken.Valid {
		return p.aheadToken
	}
	p.aheadToken = p.l.NextToken()
	return p.aheadToken
}

// LookAheadKind 查看下一个单词的类型，但不移动位置
func (p *Parser) LookAheadKind() TkKind {
	return p.LookAheadToken().TokenKind
}

// NextIdentifier 等同 NextTokenKind(ast.TkIdentifier)
func (p *Parser) NextIdentifier() {
	p.NextTokenKind(ast.TkIdentifier)
}

// NextTokenKind 读取下一个单词，并且检验类型，如果不满足条件，记录错误。【遇到特殊token，会卡住】
func (p *Parser) NextTokenKind(kind TkKind) {
	look_kind := p.LookAheadKind()
	if look_kind != kind {
		p.insertParserErr(p.aheadToken.Loc, "expected %s, found '%s'", kind.String(), look_kind.String())
		p.NextToken()
	} else {
		p.NextToken()
	}
}

func (p *Parser) NextTokenKindExpectByLoc(kind TkKind, beginTokenLoc Location) {
	look_kind := p.LookAheadKind()
	if look_kind != kind {
		p.insertParserErr(beginTokenLoc, "miss correspond %s", kind.String())
		p.insertParserErr(p.aheadToken.Loc, "expected %s, found '%s'", kind.String(), look_kind.String())
		p.NextToken()
	} else {
		p.NextToken()
	}
}

// CheckNowTokenKind 检查当前单词类型是否符合要求，否则记录错误
func (p *Parser) CheckNowTokenKind(kind TkKind) {
	if p.nowToken.TokenKind != kind {
		p.insertParserErr(p.nowToken.Loc, "expected %s, found '%s'", kind.String(), p.nowToken.TokenKind.String())
	}
}

// insert now token info
func (p *Parser) insertParserErr(loc Location, f string, a ...any) {
	err := fmt.Sprintf(f, a...)
	paseError := ast.ParseError{
		ErrStr: err,
		Loc:    loc,
	}

	p.insertErr(paseError)
}

func (p *Parser) insertErr(oneErr ParseError) {
	if len(p.parseErrs) < 30 {
		p.parseErrs = append(p.parseErrs, oneErr)
	} else {
		oneErr.ErrStr = oneErr.ErrStr + "(too many err...)"
		p.parseErrs = append(p.parseErrs, oneErr)

		panic("too many parse errors")
	}
}

// block ::= {stat} [retstat]
func (p *Parser) parseBlock() *ast.Block {
	var start_loc = p.aheadToken.Loc
	var block = &ast.Block{
		Stats: p.parseStats(),
	}
	var end_loc = p.nowToken.Loc
	block.Loc = common.GetRangeLoc(start_loc, end_loc)
	return block
}

func isBlockEnd(tokenKind TkKind) bool {
	switch tokenKind {
	case ast.TkEOF, ast.TkKwEnd,
		ast.TkKwElse, ast.TkKwElseIf, ast.TkKwUntil:
		return true
	}
	return false
}

func (p *Parser) parseStats() []ast.Stat {
	stats := make([]ast.Stat, 0, 1)
	for !isBlockEnd(p.LookAheadKind()) {
		var start_loc = p.aheadToken.Loc
		stat := p.parseStat()
		if stat != nil {
			var Loc = common.GetRangeLoc(start_loc, p.nowToken.Loc)
			stat.SetLoc(Loc)
			stats = append(stats, stat)
		}
	}
	return stats
}
