package reader

type HTMLReader interface {
	LoadHTML(path string) (string, error)
}

type CSVReader interface {
	LoadCSV(path string) ([]map[string]interface{}, error)
}

type Reader interface {
	HTMLReader
	CSVReader
}
