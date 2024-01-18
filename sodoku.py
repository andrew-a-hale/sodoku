import sys

SIZE = 9
SSIZE = 3
GRID_T = list[list[int]]
CELL_T = tuple[int, int]
DIGITS = set(range(1, 10))


class InvalidInput(Exception):
    def __init__(self, msg: str):
        super().__init__(msg)


def validate_grid_string(grid_string: str) -> list[int]:
    digits = [int(x) for x in grid_string.replace(".", "0")]
    if len(digits) == SIZE * SIZE:
        return digits
    else:
        raise InvalidInput("malformed input string")


def check_duplication(grid: list[int]):
    for i in range(SIZE):
        tmp = [x for x in get_row(grid, i) if x > 0]
        if len(tmp) != len(set(tmp)):
            raise InvalidInput(f"duplication found in row {i+1}")

    for i in range(SIZE):
        tmp = [x for x in get_col(grid, i) if x > 0]
        if len(tmp) != len(set(tmp)):
            raise InvalidInput(f"duplication found in column {i+1}")

    for i in range(SIZE):
        sx, sy = i // SSIZE, i % SSIZE
        tmp = [x for x in get_subgrid(grid, (sx, sy)) if x > 0]
        if len(tmp) != len(set(tmp)):
            raise InvalidInput(f"duplication found in subgrid {i+1}")


def read_grid(grid_string: str) -> GRID_T:
    digits = validate_grid_string(grid_string)
    grid = []
    for i in range(SIZE):
        grid.append(digits[i * SIZE : (i + 1) * SIZE])

    check_duplication(grid)

    return grid


def solve(grid: GRID_T):
    """Backtrace algorithm
    look forward 1 step
    apply change
    if change is bad reset

    mutates grid globally"""
    cell = next_cell(grid)
    if cell is None:
        return True

    moves = get_moves_for_cell(grid, cell)
    for move in moves:
        if is_valid_move(grid, cell, move):
            grid[cell[0]][cell[1]] = move
            if solve(grid):
                return True
            grid[cell[0]][cell[1]] = 0

    return False


def next_cell(grid: GRID_T) -> CELL_T | None:
    # can add heuristic here
    for row in range(SIZE):
        for col in range(SIZE):
            if grid[row][col] == 0:
                return row, col


def is_valid_move(grid: GRID_T, cell: CELL_T, move: int) -> bool:
    # check row
    row = get_row(grid, cell[0])
    for i, val in enumerate(row):
        if val == move and i != cell[1]:
            return False

    # check column
    col = get_col(grid, cell[1])
    for i, val in enumerate(col):
        if i != cell[0] and val == move:
            return False

    # check subgrid
    subgrid = get_subgrid(grid, cell)
    cell_to_idx = (cell[0] % 3) * 3 + (cell[1] % 3)
    for i, val in enumerate(subgrid):
        if i == cell_to_idx and val == move:
            return False

    return True


def get_moves_for_cell(grid: GRID_T, cell: CELL_T) -> set[int]:
    x, y = cell[0], cell[1]
    candidates = DIGITS.copy()
    found = set()
    found = found.union(get_row(grid, x), get_col(grid, y), get_subgrid(grid, cell))
    candidates.difference_update(found)
    return candidates


def get_row(grid: GRID_T, row_idx: int) -> list[int]:
    return grid[row_idx]


def get_col(grid: GRID_T, col_idx: int) -> list[int]:
    return [grid[row][col_idx] for row in range(SIZE)]


def get_subgrid(grid: GRID_T, cell: CELL_T) -> list[int]:
    sg_x = cell[0] // 3 * 3
    sg_y = cell[1] // 3 * 3
    subgrid_size = SSIZE

    return [
        grid[row][col]
        for row in range(sg_x, sg_x + subgrid_size)
        for col in range(sg_y, sg_y + subgrid_size)
    ]


def print_grid(grid: GRID_T) -> None:
    for row in grid:
        print(" ".join(str(cell) for cell in row))


if __name__ == "__main__":
    grid_string = sys.argv[1]

    grid = read_grid(grid_string)
    solve(grid)
    print_grid(grid)
