package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

const SIZE int = 9
const SSIZE int = 3

type Cell struct {
	row int
	col int
}

type Moves []Move
type Move struct {
	cell Cell
	val  int
}

type Grid [][]int
type Sodoku interface {
	nextCell() (Cell, bool)
	getMoves(cell Cell) Moves
	makeMove(move Move)
	clearCell(cell Cell)
	solve() bool
	getSubGrid(cell Cell) []int
	isValidMove(move Move) bool
	isSolved() bool
	print()
}

func (grid Grid) nextCell() (Cell, bool) {
	for r, row := range grid {
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
	for _, val := range grid[cell.row] {
		if val != 0 && slices.Contains(remaining, val) {
			remaining = removeDigit(remaining, val)
		}
	}

	// check col
	for i := 0; i < SIZE; i++ {
		val := grid[i][cell.col]
		if val != 0 && slices.Contains(remaining, val) {
			remaining = removeDigit(remaining, val)
		}
	}

	// check subgrid
	subgrid := grid.getSubGrid(cell)
	for _, val := range subgrid {
		if slices.Contains(remaining, val) {
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
	grid[move.cell.row][move.cell.col] = move.val
}

func (grid Grid) clearCell(cell Cell) {
	grid[cell.row][cell.col] = 0
}

func (grid Grid) solve() bool {
	cell, ok := grid.nextCell()
	if !ok {
		return true
	}

	moves := grid.getMoves(cell)
	for _, move := range moves {
		if grid.isValidMove(move) {
			grid.makeMove(move)
			if grid.solve() {
				return true
			}
			grid.clearCell(cell)
		}
	}

	return false
}

func (grid Grid) getSubGrid(cell Cell) []int {
	sx := cell.row / SSIZE * SSIZE
	sy := cell.col / SSIZE * SSIZE
	ssize := int(SSIZE)

	var value int
	var subgrid []int
	for i := sx; i < sx+ssize; i++ {
		for j := sy; j < sy+ssize; j++ {
			value = grid[i][j]
			if value != 0 {
				subgrid = append(subgrid, value)
			}
		}
	}

	return subgrid
}

func (grid Grid) isValidMove(move Move) bool {
	// row check
	row := grid[move.cell.row]
	for i, val := range row {
		if val == move.val && i != int(move.cell.col) {
			return false
		}
	}

	// col check
	for i := 0; i < SIZE; i++ {
		val := grid[move.cell.row][i]
		if i == int(move.cell.col) && i != int(move.cell.row) && val == move.val {
			return false
		}
	}

	// subgrid check
	subgrid := grid.getSubGrid(move.cell)
	idx := move.cell.row%SSIZE*SSIZE + move.cell.col%SSIZE
	for i, val := range subgrid {
		if i == idx && val == move.val {
			return false
		}
	}

	return true
}

func (grid Grid) isSolved() bool {
	for cell := range grid {
		if cell == 0 {
			return false
		}
	}

	return true
}

func (grid Grid) print() {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			fmt.Printf("%d ", grid[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func newGrid() Grid {
	tmp := [][]int{}
	for i := 0; i < SIZE; i++ {
		tmp = append(tmp, []int{0, 0, 0, 0, 0, 0, 0, 0, 0})
	}

	return tmp
}

func parseGridString(s string) (grid Grid, err error) {
	tmp := newGrid()
	flatGrid := strings.Split(strings.Replace(s, ".", "0", -1), "")
	if len(flatGrid) != SIZE*SIZE {
		return grid, fmt.Errorf("Error: unable to parse input, had size %d, expected %d", len(flatGrid), SIZE*SIZE)
	}

	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			parsed, err := strconv.ParseInt(flatGrid[i*SIZE+j], 10, 0)
			tmp[i][j] = int(parsed)
			if err != nil {
				return grid, fmt.Errorf("Error: unable to parse digit in input: %s", flatGrid[i*SIZE+j])
			}
		}
	}

	return tmp, nil
}

func main() {
	var gridString string
	if len(os.Args) > 0 {
		gridString = os.Args[1]
	}

	grid, err := parseGridString(gridString)
	if err != nil {
		panic(err)
	}

	grid.solve()
	grid.print()
}
