package common

import (
	"bytes"
	"strings"
)

type LuaSource struct {
	Name  string
	lines []string
}

func (s *LuaSource) GetOneLine(line int) string {
	if line >= 0 && line < len(s.lines) {
		return s.lines[line]
	}
	return ""
}

func (s *LuaSource) GetLineNum() int {
	return len(s.lines)
}

func NewLuaSource(chunk []byte, Name string) *LuaSource {
	var source = &LuaSource{
		Name: Name,
	}
	// try skip utf-8 bom
	chunk = bytes.TrimPrefix(chunk, []byte{0xEF, 0xBB, 0xBF})

	// 分解成行，只支持utf-8
	var last_idx = 0
	var builder strings.Builder
	for i := 0; i < len(chunk); i++ {
		var ch = chunk[i]
		if ch == '\r' { // 直接无视掉 \r
			continue
		}
		if ch == '\n' {
			source.lines = append(source.lines, builder.String())
			builder.Reset()
			last_idx = i
		} else {
			builder.WriteByte(ch)
		}
	}
	if last_idx < len(chunk) {
		source.lines = append(source.lines, builder.String())
	}
	if len(source.lines) > 0 {
		// 剔除 # 开头的第一行
		if strings.HasPrefix(source.lines[0], "#") {
			source.lines[0] = ""
		}
	}
	return source
}

func (s *LuaSource) GetOneLineLength(line int) int {
	if line < len(s.lines) && line >= 0 {
		return len(s.lines[line])
	}
	return 0
}

func (source *LuaSource) IsEOF(pos Position) bool {
	return int(pos.Line) >= len(source.lines)
}
