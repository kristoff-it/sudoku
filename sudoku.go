package sudoku

import (
	"bytes"
	"fmt"
)

// Size is the number of rows and columns of a sudoku puzzle.
const Size = 9
// SubSize is the number of rows and columns of a sudoku subgrid.
const SubSize = 3
// Unknown is the glyph that indicates a masked or unknown value.
const Unknown byte = ' '

// Puzzle represents a 9×9 sudoku grid.
type Puzzle [Size][Size]byte

// Read in a puzzle definition from a slice of bytes.
//
// The input format should contain one line for each puzzle row, with
// lines containing puzzle glyphs separated by a single byte.  Unknown
// (masked) values should be indicated by the underscore character 0x5f.
//
// E.g., the following is valid puzzle input:
// 1 _ 3 _ _ 6 _ 8 _
// _ 5 _ _ 8 _ 1 2 _
// 7 _ 9 1 _ 3 _ 5 6
// _ 3 _ _ 6 7 _ 9 _
// 5 _ 7 8 _ _ _ 3 _
// 8 _ 1 _ 3 _ 5 _ 7
// _ 4 _ _ 7 8 _ 1 _
// 6 _ 8 _ _ 2 _ 4 _
// _ 1 2 _ 4 5 _ 7 8
func (puz *Puzzle) Read(input []byte) error {
	ending := []byte("\n")
	lines := bytes.Split(bytes.TrimSpace(input), ending)
	if len(lines) != Size {
		return fmt.Errorf("malformed input: expected %v lines, got %v", Size, len(lines))
	}
	for i, line := range lines {
		// Expecting one byte separating each puzzle glyph.
		length := (Size * 2) - 1
		if len(line) != length {
			return fmt.Errorf("malformed input on line %v: expected %v bytes, got %v", i+1, len(line), length)
		}
		for j := 0; j < Size; j++ {
			glyph := line[j*2]
			if glyph == '_' {
				// Masked value, set cell to space for "unknown".
				puz[i][j] = Unknown
			} else if glyph >= '1' && glyph <= '9' {
				// Known value.
				puz[i][j] = glyph
			} else {
				return fmt.Errorf("malformed input on line %v: expected underscore or digit 1-9 in column %v, got %v", i+1, j, glyph)
			}
		}
	}
	return nil
}

// findDuplicate searches the argument for duplicate glyphs, and returns the
// first glyph which occurs more than once.  It returns the null byte 0x00 if
// no duplicates exist.  Duplicates of the Unknown byte are disregarded.
func findDuplicate(input []byte) byte {
	// Use an empty struct mapping as a poor man's "set" type.
	var seen map[byte]struct{} = make(map[byte]struct{})
	for _, ch := range input {
		if ch == Unknown {
			continue
		}
		if _, ok := seen[ch]; ok {
			return ch
		}
		seen[ch] = struct{}{}
	}
	return 0
}

// Validate a puzzle for correctness.
//
// A puzzle is incorrect if it contains the same glyph more than once on any
// line, any column, or in any of the nine 3×3 subgrids.
func (puz *Puzzle) Validate() error {
	// Rows
	for i, row := range puz {
		dup := findDuplicate(row[:])
		if dup != 0 {
			return fmt.Errorf("invalid puzzle: duplicate %v in row %v", dup, i+1)
		}
	}
	// Columns
	for i := 0; i < Size; i++ {
		var col [Size]byte
		for j := 0; j < Size; j++ {
			col[j] = puz[j][i]
		}
		dup := findDuplicate(col[:])
		if dup != 0 {
			return fmt.Errorf("invalid puzzle: duplicate %v in col %v", dup, i+1)
		}
	}
	// Subgrids
	for i := 0; i < Size; i++ {
		r := (i / SubSize) * SubSize
		c := (i % SubSize) * SubSize
		subgrid := []byte{}
		for j := 0; j < SubSize; j++ {
			for k := 0; k < SubSize; k++ {
				subgrid = append(subgrid, puz[r+j][c+k])
			}
		}
		dup := findDuplicate(subgrid)
		if dup != 0 {
			return fmt.Errorf("invalid puzzle: duplicate %v in subgrid of R%vC%v", dup, r+1, c+1)
		}
	}
	return nil
}
