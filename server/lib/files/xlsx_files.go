package files

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/shakinm/xlsReader/helpers"
	"github.com/xuri/excelize/v2"
)

func ReadFileXlsxAutoHeader(byteContent []byte, totalHeader int, opts ...excelize.Options) ([]map[string]interface{}, error) {
	f, openErr := excelize.OpenReader(bytes.NewReader(byteContent))
	if openErr != nil {
		return nil, openErr
	}
	defer f.Close()

	// Get all the rows in the Sheet.
	rows, rowsErr := f.GetRows(f.GetSheetName(0), opts...)
	if rowsErr != nil {
		return nil, rowsErr
	}

	// Check if rows empty
	if len(rows) < 1 {
		return nil, errors.New("file rows empty")
	}

	// Set first row as headers
	headers := rows[0]

	result := make([]map[string]interface{}, 0)

	// Skipping first row directly read content
	for rowIndex := 1; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]

		tempMap := make(map[string]interface{})
		// Iterate column
		for colIndex := 0; colIndex < totalHeader; colIndex++ {
			content := row[colIndex]
			// if 1st cell of content is empty assume its EOF stop read file content
			if colIndex == 0 && len(content) == 0 {
				tempMap = nil
				break
			}

			// Make key to lower case
			header := strings.ToLower(headers[colIndex])
			// Replace space with underscore
			header = strings.ReplaceAll(header, " ", "_")
			// Assign value based on headers
			tempMap[header] = content
		}

		if tempMap != nil {
			result = append(result, tempMap)
		}
	}

	return result, nil
}

func ReadFileXlsxManualHeader(byteContent []byte, headers []*ExcelFileHeader) ([]map[string]interface{}, error) {
	f, openErr := excelize.OpenReader(bytes.NewReader(byteContent))
	if openErr != nil {
		return nil, openErr
	}
	defer f.Close()

	// Get all the rows in the Sheet.
	rows, rowsErr := f.GetRows(f.GetSheetName(0))
	if rowsErr != nil {
		return nil, rowsErr
	}

	// Check if rows empty
	if len(rows) < 1 {
		return nil, errors.New("file rows empty")
	}

	result := make([]map[string]interface{}, 0)

	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {
		tempMap := make(map[string]interface{})

		for _, header := range headers {
			var content string
			var contentErr error
			opts := excelize.Options{
				RawCellValue: true,
			}
			cellPos := fmt.Sprintf("%s%d", header.CellPosition, rowIndex+1)
			content, contentErr = f.GetCellValue(f.GetSheetName(0), cellPos, opts)
			if contentErr != nil {
				return nil, contentErr
			}
			if rowIndex > 0 {
				if strings.EqualFold(header.CellPosition, "A") && len(content) < 1 {
					tempMap = nil
					break
				} else {
					if strings.EqualFold(header.DataType, "date") {
						timeParse, parseErr := strconv.ParseFloat(content, 64)
						if parseErr != nil {
							return nil, parseErr
						}
						tempMap[header.TableColumn] = helpers.TimeFromExcelTime(timeParse, false)
						continue
					}
				}

				tempMap[header.TableColumn] = content
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

func WriteXlsx(sheetName string, headers []string, contents []map[string]string) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	err := f.SetSheetName("Sheet1", sheetName)
	if err != nil {
		return nil, err
	}

	for col, header := range headers {
		cellName := columnToLetter(col) + fmt.Sprint(1)
		f.SetCellValue(sheetName, cellName, header)
	}

	for rowIndex := 0; rowIndex < len(contents); rowIndex++ {
		for col, header := range headers {
			cellName := columnToLetter(col) + fmt.Sprint(rowIndex+2)
			f.SetCellValue(sheetName, cellName, contents[rowIndex][header])
		}
	}

	f.SetActiveSheet(0)

	return f.WriteToBuffer()
}

func columnToLetter(colIndex int) string {
	var result string
	for colIndex >= 0 {
		result = fmt.Sprintf("%c", 'A'+colIndex%26) + result
		colIndex = colIndex/26 - 1
	}
	return result
}
