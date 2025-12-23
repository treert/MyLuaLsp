package compiler

import (
	"mylua-lsp/lsp/ast"
	"mylua-lsp/lsp/common"
)

func ParseLuaSource(source *common.LuaSource) (block *ast.Block, commentMap LuaCommentMap, errList []ast.ParseError) {
	parser := Parser{}
	lexer := NewLexer(source, parser.insertErr)
	parser.l = lexer
	parser.nowToken = lexer.GetNowToken()
	parser.aheadToken.valid = false
	parser.preToken = parser.nowToken // 一个 trick，保证 preToken 有意义

	defer func() {
		if err1 := recover(); err1 != nil {
			// 太多简单的语法错误，丢弃分析结果，只保留错误列表
			errList = parser.parseErrs
			return
		}
	}()

	blockBeginLoc := parser.nowToken.Loc
	block = parser.parseBlock() // block
	blockEndLoc := parser.preToken.Loc
	block.Loc = common.GetRangeLoc(&blockBeginLoc, &blockEndLoc)

	parser.CheckNowTokenKind(ast.TkEOF)

	return block, lexer.GetCommentMap(), parser.parseErrs
}

type Parser struct {
	// 词法分析器对象
	l *Lexer

	preToken   Token
	nowToken   Token
	aheadToken Token

	parseErrs []ast.ParseError
}

// NextToken 读取下一个单词
func (p *Parser) NextToken() {
	if p.aheadToken.valid {
		p.preToken = p.nowToken
		p.nowToken = p.aheadToken
		p.aheadToken.valid = false
		return
	}
	p.preToken = p.nowToken
	p.nowToken = p.l.NextToken()
}

// LookAheadToken 查看下一个单词，但不移动位置
func (p *Parser) LookAheadToken() Token {
	if p.aheadToken.valid {
		return p.aheadToken
	}
	p.aheadToken = p.l.NextToken()
	return p.aheadToken
}

// LookAheadKind 查看下一个单词的类型，但不移动位置
func (p *Parser) LookAheadKind() ast.TkKind {
	return p.LookAheadToken().tokenKind
}

// NextIdentifier 下一个标识
func (p *Parser) NextIdentifier() {
	p.NextTokenKind(ast.TkIdentifier)
}

// NextTokenKind 读取下一个单词，并且检验类型，如果不满足条件，记录错误
func (p *Parser) NextTokenKind(kind ast.TkKind) {
	p.NextToken()
	if p.nowToken.tokenKind != kind {
		p.l.errorPrint(p.nowToken.Loc, "expected %s, found '%s'", kind.String(), p.nowToken.tokenKind.String())
	}
}

// CheckNowTokenKind 检查当前单词类型是否符合要求，否则记录错误
func (p *Parser) CheckNowTokenKind(kind ast.TkKind) {
	if p.nowToken.tokenKind != kind {
		p.l.errorPrint(p.nowToken.Loc, "expected %s, found '%s'", kind.String(), p.nowToken.tokenKind.String())
	}
}

func (p *Parser) insertErr(oneErr ast.ParseError) {
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
	return &ast.Block{
		Stats: p.parseStats(),
	}
}

func isBlockEnd(tokenKind ast.TkKind) bool {
	switch tokenKind {
	case ast.TkEOF, ast.TkKwEnd,
		ast.TkKwElse, ast.TkKwElseif, ast.TkKwUntil:
		return true
	}
	return false
}

func (p *Parser) parseStats() []ast.Stat {
	stats := make([]ast.Stat, 0, 1)
	for !isBlockEnd(p.LookAheadKind()) {
		stat := p.parseStat()
		if stat != nil {
			stats = append(stats, stat)
		}
	}
	return stats
}

func (p *Parser) parseStat() ast.Stat {
	switch p.LookAheadKind() {
	case ast.TkKwBreak:
		p.NextToken()
		return &ast.BreakStat{}
	default:
		return nil
	}
}
