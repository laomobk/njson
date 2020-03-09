package njson

import (
	"fmt"
)

func makeObject(tok *tokenizer) *JsonObject {
	stream := tok.run()
	p := newParser(stream)
	return p.parseJson()
}

func Load(fpath string) (*JsonObject, error) {
	tok, err := NewTokenizerForTest(fpath)
	if err != nil {
		return nil, err
	}

	return makeObject(tok), nil
}

func DLoad(fpath string) *JsonObject {
	obj, err := Load(fpath)

	if err != nil {
		panic(err)
	}

	return obj
}

func Loads(source string) (*JsonObject, error) {
	tok := &tokenizer{
		filepath: "<source>",
		source:   []byte(source),
	}

	return makeObject(tok), nil
}

func DLoads(source string) *JsonObject {
	obj, err := Loads(source)

	if err != nil {
		panic(err)
	}

	return obj
}

func DGet(element JsonElement, indexOrKey interface{}) JsonElement {
	ele, err := Get(element, indexOrKey)

	if err != nil {
		panic(err)
	}

	return ele
}

func Get(element JsonElement, indexOrKey interface{}) (JsonElement, error) {
	switch indexOrKey.(type) {

	case int: // may be array
		index := indexOrKey.(int)
		a, ok := element.(*JsonArrayElement)

		if !ok {
			return nil, fmt.Errorf("Element is not an array")
		}
		size := len(a.array)

		if index > size || index < 0 {
			return nil, fmt.Errorf("index out if range")
		}

		return a.array[index], nil

	case string: // may be map
		key := indexOrKey.(string)
		m, ok := element.(*JsonDictElement)

		if !ok {
			return nil, fmt.Errorf("Element is not a dict")
		}

		v, ok := m.dict[key]

		if !ok {
			return nil, fmt.Errorf("key '" + key + "' is not exists.")
		}

		return v, nil
	}

	return nil, fmt.Errorf("Integer or string excepted")
}
