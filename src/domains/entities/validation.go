package entities

import (
	"strconv"
	"strings"
	"time"
)

var timestampLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.000",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05.000-0700",
	"2006-01-02T15:04:05.000-07:00",
}

func TryParseTimestamp(s string) (time.Time, string, error) {
	s = strings.TrimSpace(s)
	var zero time.Time

	var lastErr error
	for _, layout := range timestampLayouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, layout, nil
		}
		lastErr = err
		if !strings.Contains(layout, "Z07") && !strings.Contains(layout, "-07") {
			if t, err2 := time.ParseInLocation(layout, s, time.UTC); err2 == nil {
				return t, layout, nil
			}
		}
	}
	return zero, "", lastErr
}

func (c *TableColumn) ValidateColumnValue(value *string) bool {
	if value == nil || *value == "" {
		return true
	}

	switch c.Type {
	case ColumnTypeText:
		return true
	case ColumnTypeNumeric:
		_, err := strconv.ParseFloat(*value, 64)
		return err == nil
	case ColumnTypeEnum:
		for _, v := range c.Enum {
			if v == *value {
				return true
			}
		}
		return false
	case ColumnTypeTimestamp:
		_, _, err := TryParseTimestamp(*value)
		return err == nil
	default:
		return false
	}
}
