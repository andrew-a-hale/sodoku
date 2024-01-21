package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	SIZE  = 9
	SSIZE = 3
)

type (
	Cell struct {
		row int
		col int
	}
	Moves []Move
	Move  struct {
		cell Cell
		val  int
	}
)

type Grid struct {
	rows     [][]int
	cols     [][]int
	subgrids [][]int
	size     int
	ssize    int
}

func (grid Grid) nextCell() (Cell, bool) {
	for r, row := range grid.rows {
		for c, cell := range row {
			if cell == 0 {
				return Cell{row: r, col: c}, true
			}
		}
	}

	return Cell{}, false
}

func removeDigit(digits []int, x int) []int {
	tmp := digits[:0]
	for _, d := range digits {
		if d != x {
			tmp = append(tmp, d)
		}
	}

	return tmp
}

func (grid Grid) getMoves(cell Cell) Moves {
	remaining := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	// check row
	for _, val := range grid.rows[cell.row] {
		if val != 0 && slices.Contains(remaining, val) {
			remaining = removeDigit(remaining, val)
		}
	}

	// check col
	for _, val := range grid.cols[cell.col] {
		if val != 0 && slices.Contains(remaining, val) {
			remaining = removeDigit(remaining, val)
		}
	}

	// check subgrid
	subgrid, _ := grid.cellToSubgridIndexes(cell)
	for _, val := range grid.subgrids[subgrid] {
		if val != 0 && slices.Contains(remaining, val) {
			remaining = removeDigit(remaining, val)
		}
	}

	var moves Moves
	for _, move := range remaining {
		moves = append(moves, Move{cell: cell, val: move})
	}

	return moves
}

func (grid Grid) makeMove(move Move) {
	grid.rows[move.cell.row][move.cell.col] = move.val
	grid.cols[move.cell.col][move.cell.row] = move.val

	subgrid, element := grid.cellToSubgridIndexes(move.cell)
	grid.subgrids[subgrid][element] = move.val
}

func (grid Grid) clearCell(cell Cell) {
	grid.rows[cell.row][cell.col] = 0
	grid.cols[cell.col][cell.row] = 0

	subgrid, element := grid.cellToSubgridIndexes(cell)
	grid.subgrids[subgrid][element] = 0
}

func (grid Grid) solve() (bool, error) {
	cell, ok := grid.nextCell()
	if !ok {
		return true, nil
	}

	moves := grid.getMoves(cell)
	for _, move := range moves {
		if len(moves) == 1 || grid.isValidMove(move) {
			grid.makeMove(move)
			if solved, _ := grid.solve(); solved {
				return true, nil
			}
			grid.clearCell(cell)
		}
	}

	return false, fmt.Errorf("failed to solve to sodoku")
}

// converts (r, c) => nth subgrid, ith cell
// eg. (0, 0) => 0, 0 (first)
// eg. (3, 3) => 4, 0
// eg. (8, 8) => 8, 8 (last)
func (grid Grid) cellToSubgridIndexes(cell Cell) (int, int) {
	subgrid := cell.row/grid.ssize*grid.ssize + cell.col/grid.ssize
	element := cell.row%grid.ssize*grid.ssize + cell.col%grid.ssize
	return subgrid, element
}

func (grid Grid) isValidMove(move Move) bool {
	// row check
	row := grid.rows[move.cell.row]
	for i, val := range row {
		if val == move.val && i != move.cell.col {
			return false
		}
	}

	// col check
	col := grid.cols[move.cell.col]
	for i, val := range col {
		if val == move.val && i != move.cell.row {
			return false
		}
	}

	// subgrid check
	idx, el := grid.cellToSubgridIndexes(move.cell)
	subgrid := grid.subgrids[idx]
	for i, val := range subgrid {
		if val == move.val && i != el {
			return false
		}
	}

	return true
}

func (grid Grid) print() {
	for i := 0; i < grid.size; i++ {
		for j := 0; j < grid.size; j++ {
			fmt.Printf("%d ", grid.rows[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func hasDuplication(xs []int) bool {
	tmp := slices.Clone(xs)
	slices.Sort(tmp)
	for i := 1; i < len(tmp); i++ {
		if tmp[i] != 0 && tmp[i] == tmp[i-1] {
			return true
		}
	}

	return false
}

func (grid Grid) hasDuplication() error {
	for i := 0; i < grid.size; i++ {
		// check rows
		if hasDuplication(grid.rows[i]) {
			return fmt.Errorf("row %d contains a duplicate", i+1)
		}

		// check cols
		if hasDuplication(grid.cols[i]) {
			return fmt.Errorf("col %d contains a duplicate", i+1)
		}

		// check subgrids
		if hasDuplication(grid.subgrids[i]) {
			return fmt.Errorf("subgrid %d contains a duplicate", i+1)
		}
	}

	return nil
}

func newGrid(s string) (Grid, error) {
	// make empty grid
	rows := make([][]int, SIZE)
	cols := make([][]int, SIZE)
	subgrids := make([][]int, SIZE)
	for i := 0; i < SIZE; i++ {
		rows[i] = []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
		cols[i] = []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
		subgrids[i] = []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	}
	grid := Grid{rows: rows, cols: cols, subgrids: subgrids, size: SIZE, ssize: SSIZE}

	// parse sodoku string
	flatGrid := strings.Split(strings.ReplaceAll(s, ".", "0"), "")
	if len(flatGrid) != grid.size*grid.size {
		return grid, fmt.Errorf("unable to parse input, had size %d, expected %d", len(flatGrid), grid.size*grid.size)
	}

	// populate grid
	for i := 0; i < grid.size; i++ {
		for j := 0; j < grid.size; j++ {
			parsed, err := strconv.ParseInt(flatGrid[i*grid.size+j], 10, 0)
			if err != nil {
				return grid, fmt.Errorf("unable to parse digit in input: %s", flatGrid[i*grid.size+j])
			}
			grid.rows[i][j] = int(parsed)
			grid.cols[j][i] = int(parsed)

			idx, el := grid.cellToSubgridIndexes(Cell{row: i, col: j})
			grid.subgrids[idx][el] = int(parsed)
		}
	}

	// check for valid input
	if err := grid.hasDuplication(); err != nil {
		return grid, err
	}

	return grid, nil
}

func main() {
	var gridString string
	if len(os.Args) > 0 {
		gridString = os.Args[1]
	}

	grid, err := newGrid(gridString)
	if err != nil {
		panic(err)
	}

	grid.print()
	grid.solve()
	grid.print()
}
