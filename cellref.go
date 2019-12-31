package sudoku

import "fmt"

// CellRef represents a reference to a puzzle cell
type CellRef struct {
	row, col int
}

func NewCellRef(row int, col int) CellRef {
	return CellRef{row, col}
}

func (ref *CellRef) String() string {
	return fmt.Sprintf("R%vC%v", ref.row+1, ref.col+1)
}
