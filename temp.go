package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func main() {
	r := gin.Default()

	r.GET("/download", func(c *gin.Context) {
		f := excelize.NewFile()

		f.SetSheetName("Sheet1", "Report 1")

		style, err := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold:   true,
				Size:   12,
				Family: "Times New Roman",
			},
			Border: []excelize.Border{
				{Type: "left", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
				{Type: "top", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})
		if err != nil {
			fmt.Println(err)
		}

		f.SetRowHeight("Report 1", 1, 32)
		f.SetCellValue("Report 1", "A1", "Hello")
		f.SetCellValue("Report 1", "B1", "World")
		f.MergeCell("Report 1", "A1", "B1")
		f.SetCellStyle("Report 1", "A1", "B1", style)

		f.SetCellValue("Report 2", "A1", "Hello")
		f.SetCellValue("Report 2", "B1", "World")

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=report.xlsx")
		c.Header("Content-Transfer-Encoding", "binary")

		if err := f.Write(c.Writer); err != nil {
			c.String(http.StatusInternalServerError, "Error creating Excel file: %v", err)
			return
		}
	})

	r.Run(":8080")
}
