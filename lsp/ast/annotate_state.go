package ast

import "mylua-lsp/lsp/common"

type AnnotateClassState struct {
	Name    string
	NameLoc common.Location
}
