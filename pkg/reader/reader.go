// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

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
