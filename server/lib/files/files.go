package files

type ExcelFileHeader struct {
	HeaderName   string `json:"header_name"`
	CellPosition string `json:"cell_position"`
	DataType     string `json:"data_type"`
	TableColumn  string `json:"table_column"`
}

type CsvFileOptions struct {
	CsvHeader []*CsvHeader `json:"csv_header"`
	Delimiter string       `json:"delimiter"`
}

type CsvHeader struct {
	HeaderName  string `json:"header_name"`
	TableColumn string `json:"table_column"`
	DataType    string `json:"data_type"`
}
