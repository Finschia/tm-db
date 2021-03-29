package db

import (
	"bytes"
	"os"
)

func Cp(bz []byte) (ret []byte) {
	ret = make([]byte, len(bz))
	copy(ret, bz)
	return ret
}

func Concat(bz1 []byte, bz2 []byte) (ret []byte) {
	bz1len := len(bz1)
	if bz1len == 0 {
		return Cp(bz2)
	}
	bz2len := len(bz2)
	if bz2len == 0 {
		return Cp(bz1)
	}

	ret = make([]byte, bz1len+bz2len)
	copy(ret, bz1)
	copy(ret[bz1len:], bz2)
	return ret
}

// Returns a slice of the same length (big endian)
// except incremented by one.
// Returns nil on overflow (e.g. if bz bytes are all 0xFF)
// CONTRACT: len(bz) > 0
func CpIncr(bz []byte) (ret []byte) {
	if len(bz) == 0 {
		panic("cpIncr expects non-zero bz length")
	}
	ret = Cp(bz)
	for i := len(bz) - 1; i >= 0; i-- {
		if ret[i] < byte(0xFF) {
			ret[i]++
			return
		}
		ret[i] = byte(0x00)
		if i == 0 {
			// Overflow
			return nil
		}
	}
	return nil
}

// See DB interface documentation for more information.
func IsKeyInDomain(key, start, end []byte) bool {
	if bytes.Compare(key, start) < 0 {
		return false
	}
	if end != nil && bytes.Compare(end, key) <= 0 {
		return false
	}
	return true
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// Path
func MakePath(path string) error {
	if len(path) == 0 {
		return nil
	}
	return os.MkdirAll(path, 0755)
}
