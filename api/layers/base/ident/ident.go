package ident

import (
	"regexp"

	"github.com/google/uuid"
)

type Ident struct {
	r *regexp.Regexp
}

func NewIdent() *Ident {
	return &Ident{
		r: func() *regexp.Regexp {
			r, _ := regexp.Compile("/^[0-9A-F]{8}-[0-9A-F]{4}-[5][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i")
			return r
		}(),
	}
}

// Generation UUIDv5
func (i *Ident) GenerateUUID(salt string) string {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte(salt)).String()
}

// Check correct UUIDv5
func (i *Ident) CheckUUIDv5(uuid string) bool {
	return i.r.MatchString(uuid)
}
