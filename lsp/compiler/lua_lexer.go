package compiler

import (
	"fmt"
	"mylua-lsp/lsp/ast"
	"mylua-lsp/lsp/common"
	"strings"
)

/*
解析的逻辑是基于行的。触发换行的地方
1. 空行
2. raw_string
3. block_comment
*/

type Token = ast.Token
type TkKind = ast.TkKind
type Position = common.Position
type Location = common.Location
type ParseError = ast.ParseError

// ErrorHandler 词法分析上报错误
type ErrorHandler func(oneErr ParseError)

type LuaCommentMap map[int]*ast.CommentBlock

// Lexer 词法分析的结构
type Lexer struct {
	source   *common.LuaSource
	cur_line string
	nextPos  Position

	tokenStartPos Position

	preToken   Token
	nowToken   Token
	aheadToken Token

	commentMap LuaCommentMap // 收集注释块信息，key时每块的最后一行。

	errHandler ErrorHandler // error reporting; or nil
}

// NewLexer 创建一个词法分析器
func NewLexer(source *common.LuaSource, errHandler ErrorHandler) *Lexer {
	var lex = &Lexer{
		source:   source,
		cur_line: source.GetOneLine(0),
		nextPos:  Position{Line: 0, Column: 0},
		preToken: Token{
			Valid: false,
		},
		nowToken: Token{
			Valid: false,
		},
		aheadToken: Token{
			Valid: false,
		},
		commentMap: LuaCommentMap{},
		errHandler: errHandler,
	}
	lex.NextToken()
	lex.NextToken()
	return lex
}

func (l *Lexer) isEndOfFile() bool {
	return l.nextPos.GetLine() >= l.source.GetLineNum()
}

func (l *Lexer) isEndOfLine() bool {
	return l.nextPos.GetColumn() >= len(l.cur_line)
}

func (l *Lexer) lookChar() byte {
	if l.nextPos.GetColumn() < len(l.cur_line) {
		return l.cur_line[l.nextPos.GetColumn()]
	}
	return '\n'
}

// 只检测，不 next
func (l *Lexer) test_str(s string) bool {
	if l.nextPos.GetColumn() < len(l.cur_line) {
		var next_str = l.cur_line[l.nextPos.GetColumn():]
		return strings.HasPrefix(next_str, s)
	}
	return false
}

// 读取当前字符并往前移动一个字符
func (l *Lexer) next() {
	l.nextPos.Column++
}

func (l *Lexer) next_until(f func(byte) bool, check bool) bool {
	for ; l.nextPos.GetColumn() < len(l.cur_line); l.nextPos.Column++ {
		c := l.cur_line[l.nextPos.GetColumn()]
		if f(c) == check {
			return true
		}
	}
	return false
}

func (l *Lexer) next_line() {
	common.Assert(!l.isEndOfFile())
	l.nextPos.Line++
	l.nextPos.Column = 0
	l.cur_line = l.source.GetOneLine(l.nextPos.GetLine())
}

func (l *Lexer) backOneChar() {
	common.Assert(l.nextPos.Column > 0)
	l.nextPos.Column--
}

// GetCommentMap 获取所有的注释map
func (l *Lexer) GetCommentMap() map[int]*ast.CommentBlock {
	return l.commentMap
}

// GetNowToken get now token
func (l *Lexer) GetNowToken() Token {
	return l.nowToken
}

func (l *Lexer) getLineEndLoc(line int) Location {
	var str = l.source.GetOneLine(line)
	var pos = Position{Line: int32(line), Column: int32(len(str))}
	return Location{
		Start: pos,
		End:   pos,
	}
}

func (l *Lexer) getFileEndLoc() Location {
	return l.getLineEndLoc(l.source.GetLineNum() - 1)
}

// NextToken 下一个单词
func (l *Lexer) NextToken() Token {
	l.preToken = l.nowToken
	l.nowToken = l.aheadToken
	l.nextTokenStruct()
	l.nowToken, l.aheadToken = l.aheadToken, l.nowToken
	// mylua 词法特殊，可以根据前后关键字来决定某些 token 的类型。这么些只是为了预留。
	return l.nowToken
}

// setNowToken 设置当前的单词
func (l *Lexer) setNowToken(kind TkKind, tokenStr string) {
	l.nowToken.Loc.Start = l.tokenStartPos
	l.nowToken.Loc.End = l.nextPos
	l.nowToken.TokenKind = kind
	l.nowToken.TokenStr = tokenStr
}

// nextTokenStruct 获取下一个单词结构
func (l *Lexer) nextTokenStruct() {
	l.skipWhiteSpaces()
	if l.isEndOfFile() {
		l.tokenStartPos = l.nextPos
		l.setNowToken(ast.TkEOF, "EOF")
		return
	}
	l.tokenStartPos = l.nextPos
	c := l.lookChar()
	l.next()
	switch c {
	case ';':
		l.setNowToken(ast.TkSepSemi, ";")
		return
	case ',':
		l.setNowToken(ast.TkSepComma, ",")
		return
	case '(':
		l.setNowToken(ast.TkSepLparen, "(")
		return
	case ')':
		l.setNowToken(ast.TkSepRparen, ")")
		return
	case ']':
		l.setNowToken(ast.TkSepRbrack, "]")
		return
	case '{':
		l.setNowToken(ast.TkSepLcurly, "{")
		return
	case '}':
		l.setNowToken(ast.TkSepRcurly, "}")
		return
	case '+':
		l.setNowToken(ast.TkOpAdd, "+")
		return
	case '-':
		l.setNowToken(ast.TkOpMinus, "-")
		return
	case '*':
		l.setNowToken(ast.TkOpMul, "*")
		return
	case '^':
		l.setNowToken(ast.TkOpPow, "^")
		return
	case '%':
		l.setNowToken(ast.TkOpMod, "%")
		return
	case '&':
		l.setNowToken(ast.TkOpBand, "&")
		return
	case '|':
		l.setNowToken(ast.TkOpBor, "|")
		return
	case '#':
		l.setNowToken(ast.TkOpNen, "#")
		return
	case ':':
		if l.test_then_next(':') {
			l.setNowToken(ast.TkSepLabel, "::")
		} else {
			l.setNowToken(ast.TkSepColon, ":")
		}
		return
	case '/':
		if l.test_then_next('/') {
			l.setNowToken(ast.TkOpIdiv, "//")
		} else {
			l.setNowToken(ast.TkOpDiv, "/")
		}
		return
	case '~':
		if l.test_then_next('=') {
			l.setNowToken(ast.TkOpNe, "~=")
		} else {
			l.setNowToken(ast.TkOpWave, "~")
		}
		return
	case '=':
		if l.test_then_next('=') {
			l.setNowToken(ast.TkOpEq, "==")
		} else {
			l.setNowToken(ast.TkOpAssign, "=")
		}
		return
	case '<':
		if l.test_then_next('<') {
			l.setNowToken(ast.TkOpShl, "<<")
		} else if l.test_then_next('=') {
			l.setNowToken(ast.TkOpLe, "<=")
		} else {
			l.setNowToken(ast.TkOpLt, "<")
		}
		return
	case '>':
		if l.test_then_next('>') {
			l.setNowToken(ast.TkOpShr, ">>")
		} else if l.test_then_next('=') {
			l.setNowToken(ast.TkOpGe, ">=")
		} else {
			l.setNowToken(ast.TkOpGt, ">")
		}
		return
	case '.':
		if l.test_then_next('.') {
			if l.test_then_next('.') {
				l.setNowToken(ast.TkVararg, "...")
			} else {
				l.setNowToken(ast.TkOpConcat, "..")
			}
		} else {
			l.setNowToken(ast.TkSepDot, ".")
		}
	case '[':
		if l.test2('[', '=') {
			l.setNowToken(ast.TkString, l.scanLongString(0))
		} else {
			l.setNowToken(ast.TkSepLbrack, "[")
		}
		return
	case '\'', '"':
		l.setNowToken(ast.TkString, l.scanShortString(c))
		return
	}

	if c == '.' || common.IsDigit(c) {
		token := l.scanNumber()
		l.setNowToken(ast.TkNumber, token)
		return
	}

	if c == '_' || common.IsLetterChar(c) {
		token := l.scanIdentifier()
		if kind, ok := ast.Keywords[token]; ok {
			l.setNowToken(kind, token)
		} else {
			l.setNowToken(ast.TkIdentifier, token)
		}
		return
	}

	illegalStr := l.scanIllegalToken()
	l.setNowToken(ast.IKIllegal, illegalStr)
	l.errorPrint(l.nowToken.Loc, "unexpected token:%s", illegalStr)
}

func (l *Lexer) scanIllegalToken() string {
	for !l.isEndOfLine() {
		ch := l.lookChar()
		if isWhiteSpace(ch) {
			break
		}
		l.next()
	}
	str := l.cur_line[l.tokenStartPos.Column:l.nextPos.Column]
	return str
}

func (l *Lexer) test1(c1 byte) bool {
	var c = l.lookChar()
	return c == c1
}

func (l *Lexer) test2(c1 byte, c2 byte) bool {
	var c = l.lookChar()
	return c == c1 || c == c2
}

func (l *Lexer) test_then_next(ch byte) bool {
	if l.lookChar() == ch {
		l.next()
		return true
	}
	return false
}

// errorPrint 错误打印，词法分析报异常
func (l *Lexer) errorPrint(loc Location, f string, a ...any) {
	err := fmt.Sprintf(f, a...)
	paseError := ParseError{
		ErrStr: err,
		Loc:    loc,
	}

	if l.errHandler != nil {
		l.errHandler(paseError)
	}
}

func (l *Lexer) _skipWhiteSpacesUntilNextLine() bool {
	for !l.isEndOfLine() {
		if isWhiteSpace(l.lookChar()) {
			l.next()
		} else {
			return false
		}
	}
	return true
}

// skipWhiteSpaces 跳过空格，并搜集注释
func (l *Lexer) skipWhiteSpaces() {
	var commentInfo *ast.CommentBlock
	lastLine := -1
	add_last_comment_block := func() {
		if commentInfo != nil {
			l.commentMap[lastLine] = commentInfo
			commentInfo = nil
		}
	}
	for !l.isEndOfFile() {
		if l.test_then_next('\n') {
			continue
		} else if isWhiteSpace(l.lookChar()) {
			l.next()
			continue
		} else if !l.test_str("--") {
			break
		}

		// 下面的逻辑是--开头的注释

		// headFlag 表示是否为首行注释还是尾部注释
		var headFlag bool
		if l.nowToken.Valid {
			headFlag = l.nowToken.Loc.End.Line != l.nextPos.Line
		} else {
			headFlag = true // 其他的都当成是首行注释，特别的是 块注释
		}

		startPos := l.nextPos
		shortFlag, skipComment := l.skipComment()

		if !shortFlag {
			// 长注释，需要满足要类似 单行注释的要求。
			var ok = l._skipWhiteSpacesUntilNextLine()
			if !ok {
				// 丢弃掉 结尾还有其他内容的长注释。当其不存在, 并且中断注释块。
				add_last_comment_block()
				continue
			}

			skipComment = strings.TrimSpace(skipComment)
			// 块注释的一种使用技巧，方便快捷取消注释
			skipComment = strings.TrimSuffix(skipComment, "\n--")
		} else {
			common.Assert(startPos.Line+1 == l.nextPos.Line)
		}

		// 空行会隔开注释块
		if startPos.GetLine() != lastLine+1 {
			add_last_comment_block()
		}

		lastLine = l.nextPos.GetLine() - 1 // 合适的注释已经吃掉最后一个 \n 了
		if commentInfo == nil {
			commentInfo = &ast.CommentBlock{}
		}

		commentInfo.List = append(commentInfo.List, ast.CommentLine{
			Str:       skipComment,
			StartPos:  startPos,
			ShortFlag: shortFlag,
			HeadFlag:  headFlag,
		})

		// 尾注释独立成块
		if !headFlag {
			add_last_comment_block()
		}
	}

	add_last_comment_block()
}

// skipComment 跳过注释, 获取注释内容
// shortFlag 是否是短注释
// strComment 为注释的内容
func (l *Lexer) skipComment() (shortFlag bool, strComment string) {
	var start_pos = l.nextPos
	l.next() // skip --
	l.next()

	// long comment ?
	if l.test_then_next('[') {
		count := l.matchLongStringBacket()
		if count > 0 {
			shortFlag = false
			strComment = l.scanLongString(count)
			return
		}
	}
	l.next_line()
	shortFlag = true
	strComment = l.cur_line[start_pos.Column:]
	return
}

// scanIdentifier 获取下个标识符
func (l *Lexer) scanIdentifier() string {
	var start_pos = l.nextPos.Column - 1
	for {
		if common.IsNameChar(l.lookChar()) {
			l.next()
			continue
		}
		break
	}
	str := l.cur_line[start_pos:l.nextPos.Column]
	return str
}

// scanNumber 扫描数字
func (l *Lexer) scanNumber() string {
	l.backOneChar()
	start_idx := l.nextPos.Column
	// l.test_then_next('-') // 这个可以优化来着，目前是把 - 当成了单目运算符了

	if l.test_then_next('0') { // 特别的数字模式，现在只支持 0x
		if l.test2('x', 'X') { // 0x
			l.next()
			l.next_until(common.IsHexChar, false)
			if l.test_then_next('.') {
				l.next_until(common.IsHexChar, false)
			}
			if l.test2('p', 'P') {
				l.next()
				if l.test2('+', '-') {
					l.next()
					l.next_until(common.IsDigit, false)
				}
			}
			goto finish_all
		}
	}
	// 普通的数字。lua 没有对 011 这种做8进制支持。
	l.next_until(common.IsDigit, false)
	if l.test_then_next('.') {
		l.next_until(common.IsDigit, false)
	}
	if l.test2('e', 'E') {
		l.next()
		if l.test2('+', '-') {
			l.next()
			l.next_until(common.IsDigit, false)
		}
	}

finish_all:

	// 切分出来的字符串
	str := l.cur_line[start_idx:l.nextPos.Column]
	return str
}

// scanLongString 扫描长字符串。之前已经被吃掉一个 [ 。返回的字符串不包含 --[[ ]]
func (l *Lexer) scanLongString(count int) string {
	if count <= 0 {
		count = l.matchLongStringBacket()
	}
	if count < 0 {
		l.errorPrint(Location{
			Start: l.tokenStartPos,
			End:   l.nextPos,
		}, "invalid long string delimiter")
		return ""
	}
	var start_idx = l.nextPos.Column
	var builder strings.Builder
	// 往后找 ]=*]
	for {
		if l.test1('\n') {
			builder.WriteString(l.cur_line[start_idx:])
			builder.WriteByte('\n')
			l.next_line()
			if l.isEndOfFile() {
				goto failed
			}
			start_idx = 0
			continue
		}
		if l.test_then_next(']') {
			// 尝试匹配
			var end_idx = l.nextPos.Column - 1
			var cnt = count - 2
			for cnt > 0 {
				if !l.test_then_next('=') {
					break
				}
				cnt--
			}
			if cnt == 0 && l.test_then_next(']') {
				// 匹配成功
				builder.WriteString(l.cur_line[start_idx:end_idx])
				return builder.String()
			}
			continue
		}
		l.next()
	}
failed:
	// 查找匹配失败了，但是到当前位置的字符当成注释。
	l.errorPrint(l.getFileEndLoc(), "reach EOF, missing `]={%d}]`", count-2)
	return builder.String()
}

// 读取并计算长字符串的前缀 [=*[ 的总长度。格式不对返回 负 <= 0, 前面已经吃掉一个 [ 了
func (l *Lexer) matchLongStringBacket() int {
	var count int = 1
	for l.test_then_next('=') {
		count++
	}
	if l.test_then_next('[') {
		return count + 1
	}
	return -count
}

// 短字符串，词法分析时处理少量转义字符。主要是 行尾的 \ \z
func (l *Lexer) scanShortString(delimiter byte) string {
	var builder strings.Builder

	start_idx := l.nextPos.GetColumn()
	for {
		var end_idx = start_idx
		var line_len = len(l.cur_line)
		for ; end_idx < line_len; end_idx++ {
			c := l.cur_line[end_idx]
			if c == delimiter {
				builder.WriteString(l.cur_line[start_idx:end_idx])
				l.nextPos.Column = int32(end_idx + 1)
				goto finish // 正确结束
			}
			if c == '\\' {
				if end_idx == line_len-1 {
					// 行尾的 \
					builder.WriteString(l.cur_line[start_idx:end_idx])
					l.next_line()
					if l.isEndOfFile() {
						l.errorPrint(l.getFileEndLoc(), "unfinished string, expect %c but got \\", delimiter)
						goto finish
					}
					start_idx = 0
					goto next_line
				} else {
					if (end_idx+1 == line_len-1) && l.cur_line[end_idx+1] == 'z' {
						// 行尾的 \z。吃掉空白符直到遇到非空白符
						builder.WriteString(l.cur_line[start_idx:end_idx])
						var find_idx int
						for {
							l.next_line()
							if l.isEndOfFile() {
								l.errorPrint(l.getFileEndLoc(), "unfinished string, missing %c", delimiter)
								goto finish
							}
							for find_idx = 0; find_idx < len(l.cur_line); find_idx++ {
								c1 := l.cur_line[find_idx]
								if isWhiteSpace(c1) == false {
									goto has_find
								}
							}
						}
					has_find:
						start_idx = find_idx
						goto next_line
					} else {
						// 遇到其他转义了。之前往前多吃掉一个字符
						end_idx++
					}
				}
			}
		}
		builder.WriteString(l.cur_line[start_idx:])
		l.errorPrint(l.getLineEndLoc(l.nextPos.GetLine()), "unfinished string, missing %c", delimiter)
		l.next_line() // 不管如何进入下一行
		break
	next_line:
		continue
	finish:
		break
	}
	return builder.String()
}

func isWhiteSpace(c byte) bool {
	switch c {
	case '\t', '\v', '\f', ' ':
		return true
	}
	return false
}
