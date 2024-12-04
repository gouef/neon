package neon

import "fmt"

type Position struct {
	Line   int
	Column int
	Offset int
}

func NewPosition(line, column, offset int) Position {
	if line == 0 {
		line = 1
	}
	if column == 0 {
		column = 1
	}
	return Position{
		Line:   line,
		Column: column,
		Offset: offset,
	}
}

func (p Position) String() string {
	if p.Column > 0 {
		return fmt.Sprintf("on line %d at column %d", p.Line, p.Column)
	}
	return fmt.Sprintf("on line %d", p.Line)
}
