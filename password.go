package password

import (
	"bytes"
	"hash"
	"io"
)

// Password Lossy encryption used to store and verify passwords
type Password struct {
	NewHash        func() hash.Hash
	RandReader     io.Reader
	HashSize       int
	RandSize       int
	EncodeToString func([]byte) string
	DecodeString   func(string) ([]byte, error)
}

// Encrypt encryption adds random variables to make each encryption different
func (p *Password) Encrypt(word string) (code string) {
	sum := make([]byte, p.RandSize+p.HashSize)
	io.ReadFull(p.RandReader, sum[:p.RandSize])
	hashSum := p.NewHash()
	hashSum.Write([]byte(word))
	hashSum.Write(sum[:p.RandSize])
	sum = hashSum.Sum(sum[:p.RandSize])
	return p.EncodeToString(sum)
}

// Verify password
func (p *Password) Verify(word, code string) bool {
	sum, err := p.DecodeString(code)
	if err != nil {
		return false
	}
	if len(sum) != p.HashSize+p.RandSize {
		return false
	}
	hashSum := p.NewHash()
	hashSum.Write([]byte(word))
	hashSum.Write(sum[:p.RandSize])
	newSum := hashSum.Sum(nil)
	return bytes.Equal(sum[p.RandSize:], newSum)
}
