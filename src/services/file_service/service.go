package file_service

import (
	"backend/src/domains/entities"
	"backend/src/services"
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/xuri/excelize/v2"
)

const (
	defaultSheetName = "Sheet1"
)

type service struct {
}

func NewService() services.IFileService {
	return &service{}
}

func (s *service) ReadFile(file *multipart.FileHeader) ([]string, [][]*string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, nil, ErrorCannotReadFile{err: err}
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".xlsx", ".xls":
		f, err := excelize.OpenReader(src)
		if err != nil {
			return nil, nil, ErrorCannotReadFile{err: err}
		}
		defer f.Close()

		return s.ReadExcel(f)

	case ".csv":
		reader := csv.NewReader(src)
		reader.LazyQuotes = true

		return s.ReadCSV(reader)

	default:
		return nil, nil, ErrorWrongFileFormat{}
	}
}

func (s *service) ReadExcel(f *excelize.File) ([]string, [][]*string, error) {
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, nil, err
	}

	if len(rows) < 2 {
		return nil, nil, ErrorEmptyFile{}
	}

	columns := rows[0]
	data := make([][]*string, 0, len(rows)-1)
	var errs []error
	for i := 1; i < len(rows); i++ {
		if len(rows[i]) > len(columns) {
			errs = append(errs, &ErrorInvalidRow{
				Row:             i,
				ExpectedColumns: len(columns),
				ActualColumns:   len(rows[i]),
			})
			continue
		}

		row := make([]*string, len(columns))
		for j, v := range rows[i] {
			if v == "" {
				continue
			}

			row[j] = pointer.To(v)
		}

		data = append(data, row)
	}

	if len(errs) > 0 {
		return nil, nil, errors.Join(errs...)
	}

	return columns, data, nil
}

func (s *service) ReadCSV(r *csv.Reader) ([]string, [][]*string, error) {
	r.FieldsPerRecord = -1

	records, err := r.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	if len(records) < 2 {
		return nil, nil, ErrorEmptyFile{}
	}

	columns := records[0]
	if len(columns) > 0 {
		columns[0] = strings.TrimPrefix(columns[0], "\uFEFF")
	}

	data := make([][]*string, 0, len(records)-1)
	var errs []error
	for i := 1; i < len(records); i++ {
		rec := records[i]

		if len(rec) != len(columns) {
			errs = append(errs, &ErrorInvalidRow{
				Row:             i,
				ExpectedColumns: len(columns),
				ActualColumns:   len(rec),
			})
			continue
		}

		row := make([]*string, 0, len(columns))
		for _, v := range rec {
			if v == "" {
				row = append(row, nil)
				continue
			}

			row = append(row, pointer.To(v))
		}

		data = append(data, row)
	}

	if len(errs) > 0 {
		return nil, nil, errors.Join(errs...)
	}

	return columns, data, nil
}

func (s *service) CreateExcel(table *entities.Table, data []entities.TableRow) (f *excelize.File, err error) {
	f = excelize.NewFile()
	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	_, err = f.NewSheet(defaultSheetName)
	if err != nil {
		return nil, err
	}

	header := make([]string, 0, len(table.Columns))
	headerIDs := make([]string, 0, len(table.Columns))
	for _, col := range table.Columns {
		if col.DeletedAt != nil {
			continue
		}
		header = append(header, col.Name)
		headerIDs = append(headerIDs, col.ID)
	}

	err = f.SetSheetRow(defaultSheetName, "A1", &header)
	if err != nil {
		return nil, err
	}

	for i, row := range data {
		rowValues := make([]interface{}, 0, len(headerIDs))
		for _, id := range headerIDs {
			rowValues = append(rowValues, row[id])
		}

		err = f.SetSheetRow(defaultSheetName, fmt.Sprintf("A%d", 2+i), &rowValues)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}
