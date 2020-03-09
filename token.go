package njson

import (
	"fmt"
	"strings"
)

type token struct {
	value     string
	tokenType int
	lineno    int
	offset    int
}

func (self *token) String() string {
	return fmt.Sprintf("< Token '%s' type = %d line = %d offset = %d >",
		self.value, self.tokenType, self.lineno, self.offset)
}

func (self *token) Equals(value string) bool {
	if self.value == value && self.tokenType != _T_STRING {
		return true
	}
	return false
}

type tokenStream struct {
	stream   []*token
	filepath string
	source   []byte
}

func newTokenStream(fpath string) *tokenStream {
	return &tokenStream{
		filepath: fpath,
	}
}

func (self *tokenStream) addToken(t *token) {
	self.stream = append(self.stream, t)
}

func (self *tokenStream) Stream() []*token {
	return self.stream
}

func (self *tokenStream) String() string {
	sbuf := []string{}

	for _, tok := range self.stream {
		sbuf = append(sbuf, tok.String())
	}

	return strings.Join(sbuf, "\n")
}

func (self *tokenStream) Size() int {
	return len(self.stream)
}

func (self *tokenStream) At(index int) *token {
	if index < len(self.stream) {
		return self.stream[index]
	}
	return nil
}
