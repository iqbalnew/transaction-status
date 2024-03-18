package files

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/shakinm/xlsReader/helpers"
	"github.com/shakinm/xlsReader/xls"
)

func ReadFileXls(byteContent []byte, headers []*ExcelFileHeader) ([]map[string]interface{}, error) {
	workbook, workbookErr := xls.OpenReader(bytes.NewReader(byteContent))

	if workbookErr != nil {
		return nil, workbookErr
	}

	sheet, sheetErr := workbook.GetSheet(0)
	if sheetErr != nil {
		return nil, sheetErr
	}

	// Check if rows empty
	if sheet.GetNumberRows() < 1 {
		return nil, errors.New("file rows empty")
	}

	result := make([]map[string]interface{}, 0)

	for rowIndex := 0; rowIndex < sheet.GetNumberRows(); rowIndex++ {
		tempMap := make(map[string]interface{})

		row, rowErr := sheet.GetRow(rowIndex)
		if rowErr != nil {
			return nil, rowErr
		}

		for colIndex, header := range headers {
			cell, cellErr := row.GetCol(colIndex)
			if cellErr != nil {
				return nil, cellErr
			}

			content := cell.GetString()
			if len(content) <= 0 {
				tempMap = nil
				continue
			}
			if rowIndex > 0 {
				tempMap[header.TableColumn] = content

				if strings.EqualFold(header.DataType, "date") {
					tempDate := cell.GetFloat64()
					tempMap[header.TableColumn] = helpers.TimeFromExcelTime(tempDate, false)
				}
			} else {
				tempMap = nil
				if !strings.EqualFold(content, header.HeaderName) {
					return nil, fmt.Errorf("file header %s difference from registered template %s", content, header.HeaderName)
				}
			}
		}

		if tempMap != nil {
			result = append(result, tempMap)
		}
	}

	return result, nil
}
