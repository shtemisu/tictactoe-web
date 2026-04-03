package domain

type Board struct {
	Cells [3][3]uint8
}

func InitBoard() Board {
	return Board{
		Cells: [3][3]uint8{
			{0, 0, 0},
			{0, 0, 0},
			{0, 0, 0},
		},
	}
}

func (b Board) GetEmptyCells() map[[2]uint8]bool {
	emptyCells := make(map[[2]uint8]bool)

	for i := range 3 {
		for j := range 3 {
			if b.Cells[i][j] == 0 {
				key := [2]uint8{uint8(i), uint8(j)}
				emptyCells[key] = true
			}
		}
	}

	return emptyCells
}

func (b Board) ToOneDimensionArray() []int {
	result := make([]int, 9)
	for i := range 3 {
		for j := range 3 {
			result[i*3+j] = int(b.Cells[i][j])
		}
	}
	return result
}
