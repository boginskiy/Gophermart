package pkg

import (
	"encoding/hex"
	"errors"
	"runtime"

	"github.com/aileron-projects/go/zcrypto/zargon2"
)

var (
	ErrNotMatch  = errors.New("zargon2: password hash value not match")
	ErrShortHash = errors.New("zargon2: password hash is too short")
)

type Hasher interface {
}

type HashByArgon2 struct {
	SaltLen int
	H       *zargon2.Argon2id
}

func NewHashByArgon2(saltLen int, keyLen uint32) (*HashByArgon2, error) {
	// Automatic params
	threads := uint8(runtime.NumCPU()) // Parallelism (number of threads)
	memory := uint32(64 << 10)         // Memory usage in KB (~64 MB)
	time := uint32(3)                  // Iterations/time cost

	tmpH, err := zargon2.NewArgon2id(saltLen, time, memory, threads, keyLen)

	if err != nil {
		return nil, err
	}
	return &HashByArgon2{SaltLen: saltLen, H: tmpH}, nil
}

func (h *HashByArgon2) GenerateHashedPW(password []byte) (hashedPW string, err error) {
	// Выполняем хеширование
	hashByte, err := h.H.Sum(password)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hashByte), nil
}

func (h *HashByArgon2) SplitHashedPW(hashedPW string) (salt, hash string) {
	hashedPWByte := []byte(hashedPW)
	return hex.EncodeToString(hashedPWByte[:h.SaltLen]), hex.EncodeToString(hashedPWByte[:h.SaltLen])
}
