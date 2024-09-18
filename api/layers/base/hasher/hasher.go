package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) HashByte(b []byte) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	return string(bytes), err
}

func (h *Hasher) HashStr(str string) (string, error) {
	return h.HashByte([]byte(str))
}

func (h *Hasher) CompareHashAndStr(hash string, str string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str)) == nil
}
