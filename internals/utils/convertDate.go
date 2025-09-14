package utils

import (
	"database/sql/driver"
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

type DateOnly time.Time

// // Untuk bind dari form-data (bukan JSON)
// func (d *DateOnly) UnmarshalText(text []byte) error {
// 	s := string(text)
// 	if s == "" {
// 		return nil
// 	}
// 	t, err := time.Parse("2006-01-02", s)
// 	if err != nil {
// 		return err
// 	}
// 	*d = DateOnly(t)
// 	return nil
// }

// var _ encoding.TextUnmarshaler = (*DateOnly)(nil)

const dateFormat = "2006-01-02"

// UnmarshalText digunakan oleh Gin (form binding) dan JSON untuk parsing input string
func (d *DateOnly) UnmarshalText(text []byte) error {
	str := strings.TrimSpace(string(text))
	if str == "" {
		return nil // Biarkan kosong
	}

	t, err := time.Parse(dateFormat, str)
	if err != nil {
		return fmt.Errorf("invalid date format, must be YYYY-MM-DD: %w", err)
	}

	*d = DateOnly(t)
	return nil
}

// Value digunakan untuk menyimpan ke database (driver.Valuer)
func (d DateOnly) Value() (driver.Value, error) {
	return time.Time(d).Format(dateFormat), nil
}

// Scan digunakan untuk membaca dari database (sql.Scanner)
func (d *DateOnly) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*d = DateOnly(v)
		return nil
	case string:
		t, err := time.Parse(dateFormat, v)
		if err != nil {
			return err
		}
		*d = DateOnly(t)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into DateOnly", value)
	}
}

// MarshalText agar bisa kembali ke string saat dibutuhkan
func (d DateOnly) MarshalText() ([]byte, error) {
	return []byte(time.Time(d).Format(dateFormat)), nil
}

// String supaya gampang diprint
func (d DateOnly) String() string {
	return time.Time(d).Format(dateFormat)
}

// Untuk unmarshalling dari JSON
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	fmt.Println(">>> UnmarshalJSON called with:", string(b))
	debug.PrintStack() // Ini akan cetak stack trace ke log
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return []byte(`"` + t.Format("2006-01-02") + `"`), nil
}

// Implementasi sql.Scanner untuk scan dari DB
// func (d *DateOnly) Scan(value any) error {
// 	switch v := value.(type) {
// 	case time.Time:
// 		*d = DateOnly(v)
// 		return nil
// 	case []byte:
// 		t, err := time.Parse("2006-01-02", string(v))
// 		if err != nil {
// 			return err
// 		}
// 		*d = DateOnly(t)
// 		return nil
// 	case string:
// 		t, err := time.Parse("2006-01-02", v)
// 		if err != nil {
// 			return err
// 		}
// 		*d = DateOnly(t)
// 		return nil
// 	default:
// 		return fmt.Errorf("cannot scan type %T into DateOnly", value)
// 	}
// }

// // Implementasi driver.Valuer untuk insert ke DB
// func (d DateOnly) Value() (driver.Value, error) {
// 	t := time.Time(d)
// 	return t.Format("2006-01-02"), nil
// }
