//go:build debug

package common

// go 里面没有 assert 费解
func Assert(condition bool, msg string) {
	if !condition {
		panic("Something Wrong, check stack.")
	}
}
