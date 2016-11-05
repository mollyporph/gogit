package main

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

//PrintTable prints a table with the given header and data
func PrintTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
