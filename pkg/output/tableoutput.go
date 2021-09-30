package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type table struct {
	writer *tablewriter.Table
	data   [][]string
}

// NewTable initializes a new table with defined table configuration.
func NewTable(mergeCol bool) *table {
	table := table{writer: tablewriter.NewWriter(os.Stdout)}

	table.writer.SetAutoWrapText(false)
	table.writer.SetHeaderLine(false)
	table.writer.SetAutoFormatHeaders(true)
	table.writer.SetAlignment(tablewriter.ALIGN_LEFT)
	table.writer.SetRowSeparator("")
	table.writer.SetCenterSeparator("")
	table.writer.SetColumnSeparator("")
	table.writer.SetTablePadding("\t")
	table.writer.SetBorder(false)

	// Choose to auto merge rows in first column
	if mergeCol {
		table.writer.SetAutoMergeCellsByColumnIndex([]int{0})
	}

	return &table
}

// AddRow adds a new data row to the table.
func (t *table) AddRow(row []string) {
	t.data = append(t.data, row)
}

// SetHeaders customizes and sets the table headers.
func (t *table) SetHeaders(headers []string) {
	t.writer.SetHeader(headers)

	var headerConfig []tablewriter.Colors

	for range headers {
		headerConfig = append(headerConfig, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor})
	}

	t.writer.SetHeaderColor(
		headerConfig...,
	)
}

// SetData sets all the data rows to table.
func (t *table) SetData() {
	t.writer.AppendBulk(t.data)
}

// PrintTable outputs the table to the console.
func (t *table) Print() {
	t.writer.Render()
}
