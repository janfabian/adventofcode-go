package lib

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func EmptyGrid(X, Y int) [][]rune {

	grid := [][]rune{}

	for i := 0; i < Y; i++ {
		row := make([]rune, X)
		for j := 0; j < X; j++ {
			row[j] = ' '
		}
		grid = append(grid, row)
	}

	return grid

}

func PrintGridTable(grid [][]rune) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	for _, row := range grid {
		for _, cell := range row {
			fmt.Fprintf(w, "%c\t", cell)
		}
		fmt.Fprintln(w)
	}

	w.Flush()
}
