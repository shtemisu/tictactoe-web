package model

type BoardModel struct {
	Cells []int `db:"cells"`
}

func (b BoardModel) ToTwoDimensionArray() [3][3]uint8 {
	var result [3][3]uint8

	for i := range 3 {
		for j := range 3 {
			result[i][j] = uint8(b.Cells[i*3+j])
		}
	}
	return result
}
