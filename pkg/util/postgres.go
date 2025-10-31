package util

import (
	"fmt"

	"github.com/quantsmithapp/datastation-backend/pkg/errors"
)

func ParsePostgresError(errPath string, err error) error {
	obj := map[string]interface{}{}
	if err := Recast(err, &obj); err != nil {
		return err
	}

	if code, ok := obj["Code"].(string); ok {
		if val, ok := postgresErrorMap[code]; ok {
			switch val {
			case "20000":
				return errors.NewNotFound(fmt.Sprintf("%v%v", errPath, 404), "record not found")
			default:
				return errors.NewBadRequest(fmt.Sprintf("%v%v", errPath, 400), val)
			}
		}
	}
	return err
}

var postgresErrorMap = map[string]string{
	"23505": "unique violation",
	"23502": "not null violation",
	"23503": "foreign key violation",
	"42703": "undefined column",
	"42P01": "undefined table",
	"23514": "check violation",
	"22001": "string data right truncation",
	"42701": "duplicate column",
	"22003": "numeric value out of range",
	"22007": "invalid datetime format",
	"42P09": "ambiguous alias",
	"42P18": "indeterminate datatype",
	"22023": "invalid parameter value",
	"42P02": "undefined parameter",
}
