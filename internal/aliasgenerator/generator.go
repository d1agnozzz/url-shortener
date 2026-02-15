package aliasgenerator

import (
	"crypto/md5"
	"encoding/binary"
)

type Alias struct {
	val string
}

func (s Alias) String() string {
	return s.val
}

type AliasGenerator struct{}

func NewAliasGenerator() *AliasGenerator {
	return &AliasGenerator{}
}

var charset = generateCharset()

const aliasLen = 10

func (s *AliasGenerator) GenerateByStr(str string) Alias {
	hash := md5.Sum([]byte(str))
	encoded := encodeToCharset(binary.BigEndian.Uint64(hash[:8]))

	return Alias{val: string(encoded)}
}
