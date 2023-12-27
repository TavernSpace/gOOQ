package nullable

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math/big"
)

type BigInt struct {
	BigInt *big.Int
	Valid  bool //Valid is true if BigInt is not NULL
}

// Scan implements the Scanner interface.
func (b *BigInt) Scan(value interface{}) error {
	b.BigInt = new(big.Int).SetInt64(0)
	if value == nil {
		b.Valid = false
		return nil
	}
	var i sql.NullString
	if err := i.Scan(value); err != nil {
		return err
	}
	if _, ok := b.BigInt.SetString(i.String, 10); ok {
		b.Valid = true
		return nil
	}
	return fmt.Errorf("could not scan type %T into BigInt", value)
}

// Value implements the driver Valuer interface.
func (b BigInt) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.BigInt.String(), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank BigInt when this BigInt is null.
func (b BigInt) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	return b.BigInt.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null BigInt if the input is a null BigInt.
func (b *BigInt) UnmarshalText(text []byte) error {
	str := string(text)
	b.BigInt = new(big.Int).SetInt64(0)
	if str == "" {
		b.Valid = false
		return nil
	}
	return b.BigInt.UnmarshalText(text)
}
