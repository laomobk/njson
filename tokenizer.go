package njson

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"unicode"
)

type tokenizer struct {
	source   []byte
	filepath string
	_cp      int
	_ln      int
	_ofs     int
}

func (self *tokenizer) peek(step int) (byte, bool) {
	if len(self.source) > self._cp+step {
		return self.source[self._cp+step], true
	}
	return 0, false
}

func (self *tokenizer) moveCp(step int) {
	self._cp += step
}

func (self *tokenizer) readNext() (byte, bool) {
	self._cp++
	if len(self.source) > self._cp {
		return self.source[self._cp], true
	}
	return 0, false
}

func (self *tokenizer) peekString(size int) string {
	sb := make([]byte, size)

	for i := 0; i+self._cp < len(self.source) && i < size; i++ {
		sb[i] = self.source[self._cp+i]
	}

	return string(sb)
}

func (self *tokenizer) makeToken(value string, type_ int) *token {
	return &token{
		lineno:    self._ln,
		offset:    self._ofs,
		tokenType: type_,
		value:     string(value),
	}
}

func (self *tokenizer) getEscape() ([]byte, int, error) {
	cur := &self._cp
	nxtcp := *cur + 1

	if nxtcp >= len(self.source) {
		return nil, 0, fmt.Errorf("invalid escape character")
	}

	nxtch := self.source[nxtcp]

	var hexString string

	switch nxtch {

	case 'n':
		return []byte{'\n'}, 1, nil
	case 't':
		return []byte{'\t'}, 1, nil
	case '"':
		return []byte{'"'}, 1, nil
	case '/':
		return []byte{'/'}, 1, nil
	case 'u': // unicode
		hexString = self.peekString(6)[2:]
		u32, err := strconv.ParseUint(hexString, 16, 32)
		if err != nil {
			return nil, 0, err
		}

		return []byte(string([]rune{rune(u32)})), 5, nil
	}

	return nil, 0, fmt.Errorf("invalid escape character : '" + string(nxtch) + "'")
}

func (self *tokenizer) parseString() (string, error) {
	cur := &self._cp
	buf := []byte{}

	self.moveCp(1) // eat "
	var ch byte

	for ; *cur < len(self.source); *cur++ {
		ch = self.source[*cur]

		switch ch {

		case '\n':
			return "", fmt.Errorf("do not support muitline string")
		case '\\':
			ech, jump, err := self.getEscape()

			if err != nil {
				self.handleError(err)
			}

			buf = append(buf, ech...)

			*cur += jump
			continue

		case '"':
			// *cur++
			goto outside
		}

		buf = append(buf, self.source[*cur])
	}

outside:

	return string(buf), nil
}

func (self *tokenizer) parseNumber() (string, int) {
	cur := &self._cp
	buf := []byte{}

	hasDot := false
	tokType := _T_INTEGER

	var ch byte

	if self.source[*cur] == '-' {
		buf = append(buf, '-')
		*cur++
	}

	for ; *cur < len(self.source); *cur++ {
		ch = self.source[*cur]

		if ch == '.' && !hasDot {
			hasDot = true
			tokType = _T_FLOAT
		} else if !unicode.IsNumber(rune(ch)) {
			break
		}

		buf = append(buf, ch)
	}
	*cur-- // for next token
	return string(buf), tokType
}

func (self *tokenizer) handleError(err error) {
	msg := err.Error()

	njerr := &NJsonError{
		filepath: self.filepath,
		lno:      self._ln,
		offset:   self._ofs,
		message:  msg,
		source:   self.source,
	}

	njerr.ThrowError()
}

func (self *tokenizer) run() *tokenStream {
	self._cp = -1
	self._ofs = 1
	self._ln = 1

	stream := newTokenStream(self.filepath)
	stream.source = self.source

	for ch, ok := self.readNext(); ok; ch, ok = self.readNext() {
		switch ch {
		case '\n':
			self._ofs = 1
			self._ln++
			continue
		case '{':
			stream.addToken(self.makeToken(string(ch), _T_LLBASKET))
		case '}':
			stream.addToken(self.makeToken(string(ch), _T_LRBASKET))
		case '[':
			stream.addToken(self.makeToken(string(ch), _T_MLBASKET))
		case ']':
			stream.addToken(self.makeToken(string(ch), _T_MRBASKET))
		case '"':
			str, err := self.parseString()
			if err != nil {
				self.handleError(err)
			}
			stream.addToken(self.makeToken(str, _T_STRING))
			self._ofs += len(str) + 1
		case ':':
			stream.addToken(self.makeToken(string(ch), _T_COLON))
		case ',':
			stream.addToken(self.makeToken(string(ch), _T_COMMA))
		case 32, 9, 13:
			self._ofs++
		default:
			if self.peekString(4) == "true" {
				stream.addToken(self.makeToken("true", _T_TRUE))
				self.moveCp(3)
				self._ofs += 3

			} else if self.peekString(5) == "false" {
				stream.addToken(self.makeToken("false", _T_FALSE))
				self.moveCp(4)
				self._ofs += 4

			} else if self.peekString(4) == "null" {
				stream.addToken(self.makeToken("null", _T_NULL))
				self.moveCp(3)
				self._ofs += 3

			} else if unicode.IsNumber(rune(ch)) || ch == '-' {
				numstr, tokType := self.parseNumber()
				stream.addToken(self.makeToken(numstr, tokType))
				self._ofs += len(numstr) - 1

			} else {
				self.handleError(fmt.Errorf("invalid character : " + string(ch)))
			}
		}
		self._ofs++
	}

	stream.addToken(self.makeToken("", _T_EOF))

	return stream
}

func (self *tokenizer) RunForTest() *tokenStream {
	return self.run()
}

func newTokenizer(fpath string) (*tokenizer, error) {
	b, e := ioutil.ReadFile(fpath)

	if e != nil {
		return nil, e
	}

	return &tokenizer{
		source:   b,
		filepath: fpath,
	}, nil
}

func NewTokenizerForTest(fpath string) (*tokenizer, error) {
	return newTokenizer(fpath)
}
