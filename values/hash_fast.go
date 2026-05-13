package values

import (
	"encoding/hex"

	"github.com/effective-security/x/enum"
	"github.com/zeebo/xxh3"
)

// XXH3Hash64 returns the 64-bit XXH3 digest of data. It is fast and non-cryptographic.
func XXH3Hash64(data []byte) uint64 {
	return xxh3.Hash(data)
}

// XXH3HashString64 returns the 64-bit XXH3 digest of the string's UTF-8 bytes. It is fast and non-cryptographic.
func XXH3HashString64(data string) uint64 {
	return xxh3.HashString(data)
}

// XXH3Hash128 returns the 128-bit XXH3 digest of data as 16 bytes in canonical big-endian form (see xxh3.Uint128.Bytes). It is fast and non-cryptographic.
func XXH3Hash128(data []byte) []byte {
	sum := xxh3.Hash128(data)
	b := sum.Bytes()
	return b[:]
}

// XXH3Hash128Hex returns the lowercase hexadecimal encoding of Hash128(data) (32 characters).
func XXH3Hash128Hex(data []byte) string {
	return hex.EncodeToString(XXH3Hash128(data))
}

// XXH3HashArgs128Hex returns a 128-bit XXH3 digest over variadic arguments, encoded as lowercase hex (32 characters).
//
// Supported types: string, []string, numeric scalars and slices (uint64, int64, uint32, int32, int), bool and []bool,
// []byte, enum.ProtoEnum and []enum.ProtoEnum. Other values are serialized with String(v) (see package values).
// String and []string elements are followed by a zero byte delimiter so adjacent fields do not merge ambiguously.
// Slices of fixed-width types are written without per-element delimiters; argument order and Go types distinguish fields.
func XXH3HashArgs128Hex(data ...any) string {
	return hex.EncodeToString(XXH3HashArgs128(data...))
}

// XXH3HashArgs128 returns a 128-bit XXH3 digest over variadic arguments.
//
// Supported types: string, []string, numeric scalars and slices (uint64, int64, uint32, int32, int), bool and []bool,
// []byte, enum.ProtoEnum and []enum.ProtoEnum. Other values are serialized with String(v) (see package values).
// String and []string elements are followed by a zero byte delimiter so adjacent fields do not merge ambiguously.
// Slices of fixed-width types are written without per-element delimiters; argument order and Go types distinguish fields.
func XXH3HashArgs128(data ...any) []byte {
	hash := xxh3.New128()
	hashWriteArgs(hash, data...)
	return hash.Sum(nil)
}

// XXH3HashArgs64 returns a 64-bit XXH3 digest over variadic arguments.
//
// Supported types: string, []string, numeric scalars and slices (uint64, int64, uint32, int32, int), bool and []bool,
// []byte, enum.ProtoEnum and []enum.ProtoEnum. Other values are serialized with String(v) (see package values).
// String and []string elements are followed by a zero byte delimiter so adjacent fields do not merge ambiguously.
// Slices of fixed-width types are written without per-element delimiters; argument order and Go types distinguish fields.
func XXH3HashArgs64(data ...any) uint64 {
	hash := xxh3.New()
	hashWriteArgs(hash, data...)
	return hash.Sum64()
}

func hashWriteArgs(hash hasher, data ...any) {
	for _, v := range data {
		switch v := v.(type) {
		case string:
			hashWriteUint64(hash, uint64(len(v)))
			hashWriteString(hash, v)
			hashWrite(hash, hashDelimiter)
		case []string:
			hashWrite(hash, IntToBytes(len(v)))
			for _, v := range v {
				hashWriteString(hash, v)
				hashWrite(hash, hashDelimiter)
			}
		case []uint64:
			hashWriteUint64(hash, uint64(len(v)))
			for _, v := range v {
				hashWriteUint64(hash, v)
			}
		case []int64:
			hashWriteUint64(hash, uint64(len(v)))
			for _, v := range v {
				hashWriteUint64(hash, uint64(v))
			}
		case []uint32:
			hashWriteUint32(hash, uint32(len(v)))
			for _, v := range v {
				hashWriteUint32(hash, v)
			}
		case []int32:
			hashWriteUint32(hash, uint32(len(v)))
			for _, v := range v {
				hashWriteUint32(hash, uint32(v))
			}
		case []int:
			hashWriteUint64(hash, uint64(len(v)))
			for _, v := range v {
				hashWriteUint64(hash, uint64(v))
			}
		case []bool:
			hashWriteUint64(hash, uint64(len(v)))
			for _, v := range v {
				hashWrite(hash, BoolToBytes(v))
			}
		case []enum.ProtoEnum:
			hashWriteUint64(hash, uint64(len(v)))
			for _, v := range v {
				hashWrite(hash, Uint64ToBytes(uint64(v.Number())))
			}
		case []byte:
			hashWriteUint64(hash, uint64(len(v)))
			hashWrite(hash, v)
			hashWrite(hash, hashDelimiter)
		case uint64:
			hashWriteUint64(hash, v)
		case int64:
			hashWriteUint64(hash, uint64(v))
		case uint32:
			hashWriteUint32(hash, v)
		case int32:
			hashWriteUint32(hash, uint32(v))
		case int:
			hashWriteUint64(hash, uint64(v))
		case bool:
			hashWriteUint64(hash, uint64(Select(v, 1, 0)))
		case enum.ProtoEnum:
			hashWriteUint64(hash, uint64(v.Number()))
		default:
			hashWriteString(hash, String(v))
		}
	}
}

// Uint64ToBytes returns v in big-endian order (8 bytes). It is used for deterministic hashing of integers.
func Uint64ToBytes(v uint64) []byte {
	buf := []byte{
		byte(v >> 56),
		byte(v >> 48),
		byte(v >> 40),
		byte(v >> 32),
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
	return buf
}

func IntToBytes(v int) []byte {
	return Uint64ToBytes(uint64(v))
}

// Uint32ToBytes returns v in big-endian order (4 bytes). It is used for deterministic hashing of integers.
func Uint32ToBytes(v uint32) []byte {
	buf := []byte{
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
	return buf
}

// BoolToBytes returns []byte{1} or []byte{0}. It is used for deterministic hashing of booleans.
func BoolToBytes(v bool) []byte {
	return []byte{byte(Select(v, 1, 0))}
}

var hashDelimiter = []byte{0xff}

func hashWrite(h hasher, p []byte) {
	_, _ = h.Write(p)
}

func hashWriteString(h hasher, s string) {
	_, _ = h.WriteString(s)
}

func hashWriteUint64(h hasher, p uint64) {
	data := [8]byte{
		byte(p >> 56),
		byte(p >> 48),
		byte(p >> 40),
		byte(p >> 32),
		byte(p >> 24),
		byte(p >> 16),
		byte(p >> 8),
		byte(p),
	}
	_, _ = h.Write(data[:])
}

func hashWriteUint32(h hasher, p uint32) {
	data := [4]byte{
		byte(p >> 24),
		byte(p >> 16),
		byte(p >> 8),
		byte(p),
	}
	_, _ = h.Write(data[:])
}

type hasher interface {
	Write(p []byte) (n int, err error)
	WriteString(buf string) (int, error)
}
