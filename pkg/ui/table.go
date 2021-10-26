package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InitTable initializes TUI table component with the given data and retuns a tview table primitive.
func (tui *TUI) InitTable(headers []string, data [][]string, isSelectable bool, isFirstColSelectable bool, title string) *tview.Table {

	table := tview.NewTable().SetFixed(1, 1)

	for k, v := range headers {
		color := tcell.ColorYellow

		table.SetCell(
			0,
			k,
			&tview.TableCell{
				Text:          v,
				Color:         color,
				Align:         tview.AlignLeft,
				NotSelectable: true,
				Transparent:   true,
				Expansion:     1,
			})

	}

	for i, row := range data {

		for j, col := range row {

			color := tcell.ColorWhite

			tableCell := tview.NewTableCell(col).
				SetTextColor(color).
				SetAlign(tview.AlignLeft).
				SetExpansion(1)

			if !isFirstColSelectable && j == 0 {
				tableCell.NotSelectable = true
			}

			table.SetCell(
				i+1,
				j,
				tableCell,
			)
		}

	}

	table.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetBorderColor(BorderColor).
		SetBorderAttributes(tcell.AttrDim)

	table.SetTitle(fmt.Sprintf(TitleFmt, title))

	if isSelectable {
		table.SetSelectable(true, false)
	}

	return table
}
