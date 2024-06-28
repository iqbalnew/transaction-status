package files

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/utils"
)

func ReadFileCsv(byteContent []byte, csvFileOptions *CsvFileOptions) ([]map[string]interface{}, error) {
	csvReader := csv.NewReader(bytes.NewReader(byteContent))

	if !strings.EqualFold(csvFileOptions.Delimiter, ",") {
		csvReader.Comma = []rune(csvFileOptions.Delimiter)[0]
	}
	csvReader.FieldsPerRecord = len(csvFileOptions.CsvHeader)
	csvReader.ReuseRecord = true

	result := make([]map[string]interface{}, 0)
	index := 0
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		tempMap := make(map[string]interface{})
		for i, header := range csvFileOptions.CsvHeader {
			if index > 0 {
				tempMap[header.TableColumn] = rec[i]
				if strings.EqualFold(header.DataType, "date") {
					timeParse, parseErr := time.Parse("02/01/2006", utils.PaddedDateString(rec[i], "/"))
					if parseErr != nil {
						return nil, parseErr
					}
					tempMap[header.TableColumn] = timeParse
				}
			} else {
				tempMap = nil
				if !strings.EqualFold(strings.Trim(rec[i], "\uFEFF"), header.HeaderName) {
					return nil, fmt.Errorf("file header %s difference from registered template %s", rec[i], header.HeaderName)
				}
			}
		}

		if tempMap != nil {
			result = append(result, tempMap)
		}
		index += 1
	}

	return result, nil
}

func WriteCsv(delimiter string, headers []string, contents []map[string]string) (*bytes.Buffer, error) {
	var csvBuffer bytes.Buffer

	csvWriter := csv.NewWriter(&csvBuffer)
	if !strings.EqualFold(delimiter, ",") {
		csvWriter.Comma = []rune(delimiter)[0]
	}
	defer csvWriter.Flush()

	writeHeaderErr := csvWriter.Write(headers)
	if writeHeaderErr != nil {
		return nil, writeHeaderErr
	}

	for _, dataContent := range contents {
		record := make([]string, 0)
		for _, header := range headers {
			record = append(record, dataContent[header])
		}
		writeContentErr := csvWriter.Write(record)
		if writeContentErr != nil {
			return nil, writeContentErr
		}
	}

	csvWriter.Flush()
	if flushErr := csvWriter.Error(); flushErr != nil {
		return nil, flushErr
	}

	return &csvBuffer, nil
}
