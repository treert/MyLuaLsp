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

	defer func() {
		if err1 := recover(); err1 != nil {
			// 太多出错了，丢弃分析结果，只保留错误列表
			errList = parser.parseErrs
			return
		}
	}()

	blockBeginLoc := lexer.GetNowTokenLoc()
	block = parser.parseBlock() // block
	blockEndLoc := parser.l.GetNowTokenLoc()
	block.Loc = common.GetRangeLoc(&blockBeginLoc, &blockEndLoc)

	p.l.NextTokenKind(ast.TkEOF)
	p.l.SetEnd()

	return block, lexer.GetCommentMap(), parser.parseErrs
}

type Parser struct {
	// 词法分析器对象
	l *Lexer

	nowToken   Token
	aheadToken Token

	parseErrs []ast.ParseError
}

// NextToken 读取下一个单词
func (p *Parser) NextToken() {
	if p.aheadToken.valid {
		p.nowToken = p.aheadToken
		p.aheadToken.valid = false
		return
	}
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

// NextIdentifier 下一个标识
func (p *Parser) NextIdentifier() {
	p.NextTokenKind(ast.TkIdentifier)
}

// NextTokenKind 尝试读取下一个单词，并且检验类型，如果不满足条件
func (p *Parser) NextTokenKind(kind ast.TkKind) {
	p.NextToken()
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
		// Stats:   p.parseStats(),
		// RetExps: p.parseRetExps(),
	}
}
