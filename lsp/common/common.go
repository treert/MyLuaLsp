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

// GetRangeLoc 获取两个位置的范围，为[]
func GetRangeLoc(beginLoc, endLoc Location) Location {
	return Location{
		Start: beginLoc.Start,
		End:   endLoc.End,
	}
}
