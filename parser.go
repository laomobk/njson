package njson

import (
	"fmt"
	"strconv"
)

type parser struct {
	stream   *tokenStream
	filename string
	source   []byte

	_tp int
}

func newParser(tokenStream *tokenStream) *parser {
	return &parser{
		stream:   tokenStream,
		filename: tokenStream.filepath,
		source:   tokenStream.source,
	}
}

func (self *parser) nextTok() *token {
	self._tp++
	return self.stream.At(self._tp)
}

func (self *parser) nowTok() *token {
	return self.stream.At(self._tp)
}

func (self *parser) setTp(index int) {
	self._tp = index
}

func (self *parser) movTp(step int) {
	self._tp += step
}

func (self *parser) peek(step int) *token {
	return self.stream.At(self._tp + step)
}

func (self *parser) parseInteger() *JsonIntegerElement {
	nt := self.nowTok()
	if nt.tokenType == _T_INTEGER {
		self.nextTok()

		v, err := strconv.ParseInt(nt.value, 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		return &JsonIntegerElement{
			slot: &_JsonNumberSlot{
				vint: int64(v),
			},
		}
	}
	return nil
}

func (self *parser) parseFloat() *JsonFloatElement {
	nt := self.nowTok()
	if nt.tokenType == _T_FLOAT {
		self.nextTok()

		v, err := strconv.ParseFloat(nt.value, 64)
		if err != nil {
			fmt.Println(err)
		}
		return &JsonFloatElement{
			slot: &_JsonNumberSlot{
				vfloat: float64(v),
			},
		}
	}
	return nil
}

func (self *parser) parseString() *JsonStringElement {
	nt := self.nowTok()
	if nt.tokenType == _T_STRING {
		self.nextTok()

		v := nt.value
		return &JsonStringElement{
			value: v,
		}
	}
	return nil
}

func (self *parser) parseBool() *JsonBoolElement {
	nt := self.nowTok()
	v := true
	if nt.tokenType == _T_FALSE {
		v = false
	}
	self.nextTok()
	return &JsonBoolElement{
		value: v,
	}
}

func (self *parser) handleError(err error) {
	msg := err.Error()
	nt := self.nowTok()

	njerr := &NJsonError{
		filepath: self.filename,
		lno:      nt.lineno,
		offset:   nt.offset,
		message:  msg,
		source:   self.source,
	}
	njerr.ThrowError()
}

func (self *parser) syntaxError(msg string) {
	self.handleError(fmt.Errorf(msg))
}

func (self *parser) parseArray() *JsonArrayElement {
	itemList := []JsonElement{}

	self.nextTok() // eat '['

	if self.nowTok().Equals("]") {
		self.nextTok() // eat ']'
		return &JsonArrayElement{
			array: []JsonElement{},
		}
	}

	fitem := self.parseElement()
	if fitem == nil {
		self.syntaxError("except Element or ']'")
	}

	itemList = append(itemList, fitem)

	var item JsonElement

	for self.nowTok().Equals(",") {
		self.nextTok() // eat  ','
		item = self.parseElement()
		if item == nil {
			self.syntaxError("except Element or ']'")
		}
		itemList = append(itemList, item)
	}

	if self.nowTok().Equals("]") {
		self.nextTok() // eat ']'
	} else {
		self.syntaxError("except ']'")
	}

	return &JsonArrayElement{
		array: itemList,
	}
}

func (self *parser) parseNull() *JsonNullElement {
	self.nextTok() // eat 'null'
	return &JsonNullElement{}
}

func (self *parser) parseKVPair() (string, JsonElement) {
	key := self.parseString()

	if key == nil {
		self.syntaxError("except string")
	}

	if !self.nowTok().Equals(":") {
		self.syntaxError("except ':'")
	}

	self.nextTok() // eat ':'

	ele := self.parseElement()

	if ele == nil {
		self.syntaxError("except JsonElement")
	}

	return key.value, ele
}

func (self *parser) parseDict() *JsonDictElement {
	self.nextTok() // eat '{'

	m := map[string]JsonElement{}
	keys := []string{}

	if self.nowTok().Equals("}") {
		self.nextTok()
		return &JsonDictElement{
			dict: m,
		}
	}

	fk, fv := self.parseKVPair()
	m[fk] = fv

	keys = append(keys, fk)

	for self.nowTok().Equals(",") {
		self.nextTok() // eat ','
		k, v := self.parseKVPair()
		m[k] = v
		keys = append(keys, k)
	}

	if !self.nowTok().Equals("}") {
		self.syntaxError("except '}'")
	}

	self.nextTok() // eat '}'

	return &JsonDictElement{
		dict: m,
		keys: keys,
	}
}

func (self *parser) parseElement() JsonElement {
	nt := self.nowTok()

	switch nt.tokenType {

	case _T_STRING:
		return self.parseString()
	case _T_INTEGER:
		return self.parseInteger()
	case _T_FLOAT:
		return self.parseFloat()
	case _T_FALSE, _T_TRUE:
		return self.parseBool()
	case _T_MLBASKET:
		return self.parseArray()
	case _T_LLBASKET:
		return self.parseDict()
	case _T_NULL:
		return self.parseNull()
	}

	return nil
}

func (self *parser) parseJson() *JsonObject {
	return newJsonObjectFromDictElement(self.parseDict())
}
