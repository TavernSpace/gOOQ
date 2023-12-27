package gooq_test

import (
	"github.com/lumina-tech/gooq/pkg/gooq"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBigIntScan(t *testing.T) {
	var bigint gooq.BigInt
	err := bigint.Scan("null")
	require.NoError(t, err)
	value, err := bigint.Value()
	require.NoError(t, err)
	require.Equal(t, "4269", value.(string))

	err = bigint.UnmarshalText([]byte("42694269"))
	require.NoError(t, err)
	text, err := bigint.MarshalText()
	require.NoError(t, err)
	require.Equal(t, "42694269", string(text))
}
