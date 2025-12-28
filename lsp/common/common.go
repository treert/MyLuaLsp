package common

// 实际是 Range [Start,End)
type Location struct {
	Start Position
	End   Position
}

type Position struct {
	Line   int32 // from 0
	Column int32 // from 0
}

func (pos Position) GetLine() int {
	return int(pos.Line)
}

func (pos Position) GetColumn() int {
	return int(pos.Column)
}

// GetRangeLoc 获取两个位置的范围，为[]
func GetRangeLoc(beginLoc, endLoc Location) Location {
	return Location{
		Start: beginLoc.Start,
		End:   endLoc.End,
	}
}
