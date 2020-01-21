package utils

import (
	"os"

	"github.com/jedib0t/go-pretty/table"
)

// NewTable returns a new table with neutral style.
func NewTable() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.Style().Color = table.ColorOptions{}

	return t
}
