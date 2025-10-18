package file_service

import "fmt"

type ErrorEmptyFile struct{}

func (e ErrorEmptyFile) Error() string {
	return "Empty file"
}

type ErrorInvalidRow struct {
	Row             int
	ExpectedColumns int
	ActualColumns   int
}

func (e *ErrorInvalidRow) Error() string {
	return fmt.Sprintf("Invalid row %d: expected %d columns, got %d", e.Row, e.ExpectedColumns, e.ActualColumns)
}

type ErrorCannotReadFile struct {
	err error
}

func (e ErrorCannotReadFile) Error() string {
	return "Cannot read file: " + e.err.Error()
}

type ErrorWrongFileFormat struct{}

func (e ErrorWrongFileFormat) Error() string {
	return "Wrong file format"
}
