package aliaser

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

type Aliaser interface {
	GenerateByStr(str string) Alias
}

type md5Aliaser struct{}

func NewMd5Aliaser() Aliaser {
	return &md5Aliaser{}
}

var charset = generateCharset()

const aliasLen = 10

func (s *md5Aliaser) GenerateByStr(str string) Alias {
	hash := md5.Sum([]byte(str))
	encoded := encodeToCharset(binary.BigEndian.Uint64(hash[:8]))

	return Alias{val: string(encoded)}
}
