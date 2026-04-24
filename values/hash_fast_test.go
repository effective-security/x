package values

import (
	"encoding/hex"
	"testing"

	"github.com/effective-security/x/enum"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeebo/xxh3"
)

func TestHash64(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"nil", nil},
		{"empty", []byte{}},
		{"ascii", []byte("hello")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, xxh3.Hash(tc.data), XXH3Hash64(tc.data))
		})
	}
}

func TestHashString64(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		s    string
	}{
		{"empty", ""},
		{"hello", "hello"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, xxh3.HashString(tc.s), XXH3HashString64(tc.s))
			assert.Equal(t, XXH3Hash64([]byte(tc.s)), XXH3HashString64(tc.s))
		})
	}
}

func TestHash128(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"nil", nil},
		{"empty", []byte{}},
		{"payload", []byte("abc")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			want := xxh3.Hash128(tc.data).Bytes()
			got := XXH3Hash128(tc.data)
			require.Len(t, got, 16)
			assert.Equal(t, want[:], got)
		})
	}
}

func TestHash128Hex(t *testing.T) {
	t.Parallel()
	data := []byte("payload")
	h := XXH3Hash128Hex(data)
	assert.Len(t, h, 32)
	dec, err := hex.DecodeString(h)
	require.NoError(t, err)
	assert.Equal(t, XXH3Hash128(data), dec)
	assert.Equal(t, hex.EncodeToString(XXH3Hash128(data)), h)
}

func TestHashArgs128Hex(t *testing.T) {
	t.Parallel()
	t.Run("hex_format", func(t *testing.T) {
		t.Parallel()
		h := XXH3HashArgs128Hex("alpha")
		assert.Len(t, h, 32)
		_, err := hex.DecodeString(h)
		require.NoError(t, err)
	})
	t.Run("deterministic", func(t *testing.T) {
		t.Parallel()
		args := []any{"id", uint64(42), int(42), []byte{1, 2, 3}, []string{"a", "b"}}
		a := XXH3HashArgs128Hex(args...)
		b := XXH3HashArgs128Hex(args...)
		assert.Equal(t, a, b)
	})
	t.Run("order_sensitive", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, XXH3HashArgs128Hex("a", "b"), XXH3HashArgs128Hex("b", "a"))
	})
	t.Run("string_vs_bytes", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, XXH3HashArgs128Hex("x"), XXH3HashArgs128Hex([]byte("x")))
	})
	t.Run("numeric_and_bool_slices", func(t *testing.T) {
		t.Parallel()
		h := XXH3HashArgs128Hex(
			[]uint64{1, 2},
			[]int64{-1},
			[]uint32{3},
			[]int32{-4},
			[]int{5},
			[]bool{true, false},
		)
		assert.Len(t, h, 32)
	})
	t.Run("scalars", func(t *testing.T) {
		t.Parallel()
		h := XXH3HashArgs128Hex(
			uint64(9), int64(-9), uint32(8), int32(-8), int(7), true, false,
		)
		assert.Len(t, h, 32)
	})
	t.Run("strings_delimited", func(t *testing.T) {
		t.Parallel()
		assert.NotEqual(t, XXH3HashArgs128Hex("ab", "c"), XXH3HashArgs128Hex("a", "bc"))
	})
	t.Run("proto_enum", func(t *testing.T) {
		t.Parallel()
		h := XXH3HashArgs128Hex(enumLike_Low, []enum.ProtoEnum{enumLike_Medium})
		assert.Len(t, h, 32)
	})
	t.Run("default_string_fallback", func(t *testing.T) {
		t.Parallel()
		type other int
		h := XXH3HashArgs128Hex(other(42))
		assert.Len(t, h, 32)
	})
}

func TestUint64ToBytes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 0}, Uint64ToBytes(0))
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 1}, Uint64ToBytes(1))
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Uint64ToBytes(0xffffffffffffffff))
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8}, Uint64ToBytes(0x0102030405060708))
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 1}, IntToBytes(1))
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, IntToBytes(-1))
}

func TestUint32ToBytes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []byte{0, 0, 0, 0}, Uint32ToBytes(0))
	assert.Equal(t, []byte{0, 0, 0, 1}, Uint32ToBytes(1))
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0xff}, Uint32ToBytes(0xffffffff))
	assert.Equal(t, []byte{0x12, 0x34, 0x56, 0x78}, Uint32ToBytes(0x12345678))
}

func TestBoolToBytes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, []byte{0}, BoolToBytes(false))
	assert.Equal(t, []byte{1}, BoolToBytes(true))
}
