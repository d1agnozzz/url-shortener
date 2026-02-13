package aliasgenerator

import (
	"crypto/md5"
	"encoding/binary"
)

type AliasGenerator interface {
	GenerateByStr(str string) string
}

type md5AliasGenerator struct{}

func NewMd5AliasGenerator() AliasGenerator {
	return &md5AliasGenerator{}
}

var charset = generateCharset()

const aliasLen = 10

func (s *md5AliasGenerator) GenerateByStr(str string) string {
	hash := md5.Sum([]byte(str))

	return string(encodeToCharset(binary.BigEndian.Uint64(hash[:8])))
}
