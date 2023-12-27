package nullable_test

import (
	"github.com/lumina-tech/gooq/pkg/nullable"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBigIntScan(t *testing.T) {
	var bigint nullable.BigInt
	err := bigint.Scan("4269")
	require.NoError(t, err)
	value, err := bigint.Value()
	require.NoError(t, err)
	require.Equal(t, "4269", value.(string))

	err = bigint.UnmarshalText([]byte("42694269"))
	require.NoError(t, err)
	text, err := bigint.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "42694269", string(text))

	err = bigint.Scan(nil)
	require.NoError(t, err)
	value, err = bigint.Value()
	require.NoError(t, err)
	require.Equal(t, nil, value)

	err = bigint.UnmarshalText([]byte(""))
	require.NoError(t, err)
	text, err = bigint.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "", string(text))
}
