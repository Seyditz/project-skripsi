package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

// IntArray custom type
type IntArray []int

// Scan implements the Scanner interface.
func (a *IntArray) Scan(value interface{}) error {
	if value == nil {
		*a = IntArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// PostgreSQL arrays are returned as strings like "{1,2,3}"
		str := string(v)
		str = strings.Trim(str, "{}")
		if str == "" {
			*a = IntArray{}
			return nil
		}

		strArr := strings.Split(str, ",")
		intArr := make(IntArray, len(strArr))
		for i, s := range strArr {
			var val int
			if err := json.Unmarshal([]byte(s), &val); err != nil {
				return err
			}
			intArr[i] = val
		}
		*a = intArr
		return nil
	case string:
		// Handle cases where the array is returned as a string
		str := strings.Trim(v, "{}")
		if str == "" {
			*a = IntArray{}
			return nil
		}

		strArr := strings.Split(str, ",")
		intArr := make(IntArray, len(strArr))
		for i, s := range strArr {
			var val int
			if err := json.Unmarshal([]byte(s), &val); err != nil {
				return err
			}
			intArr[i] = val
		}
		*a = intArr
		return nil
	default:
		return errors.New("unsupported type")
	}
}

// Value implements the Valuer interface.
func (a IntArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}
