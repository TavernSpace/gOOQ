package gooq

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

type BigInt big.Int

// Scan implements the Scanner interface.
func (b *BigInt) Scan(value interface{}) error {
	bint := (*big.Int)(b)
	switch str := value.(type) {
	case string:
		if _, ok := bint.SetString(str, 10); ok {
			return nil
		}
	}
	return fmt.Errorf("could not scan type %T into BigInt", value)
}

// Value implements the driver Valuer interface.
func (b BigInt) Value() (driver.Value, error) {
	bint := (*big.Int)(&b)
	return bint.String(), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank BigInt when this BigInt is null.
func (b BigInt) MarshalText() ([]byte, error) {
	return (*big.Int)(&b).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null BigInt if the input is a null BigInt.
func (b *BigInt) UnmarshalText(text []byte) error {
	return (*big.Int)(b).UnmarshalText(text)
}
