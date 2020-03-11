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
	tok, err := newTokenizer(fpath)
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

func ForEach(jsonElement JsonElement, forfunc func(JsonElement)) {
	switch jsonElement.(type) {

	case *JsonDictElement:
		o := jsonElement.(*JsonDictElement)
		for _, v := range o.keys {
			forfunc(o.dict[v])
		}
	case *JsonArrayElement:
		o := jsonElement.(*JsonArrayElement)
		for _, v := range o.array {
			forfunc(v)
		}
	default:
		panic("Only dict or array support for range.")

	}

}

func StringElement(o JsonElement) *JsonStringElement {
	if v, ok := o.(*JsonStringElement); ok {
		return v
	}
	return nil
}

func IntegerElement(o JsonElement) *JsonIntegerElement {
	if v, ok := o.(*JsonIntegerElement); ok {
		return v
	}
	return nil
}

func FloatElement(o JsonElement) *JsonFloatElement {
	if v, ok := o.(*JsonFloatElement); ok {
		return v
	}
	return nil
}

func ArrayElement(o JsonElement) *JsonArrayElement {
	if v, ok := o.(*JsonArrayElement); ok {
		return v
	}
	return nil
}

func DictElement(o JsonElement) *JsonDictElement {
	if v, ok := o.(*JsonDictElement); ok {
		return v
	}
	return nil
}

func NullElement(o JsonElement) *JsonNullElement {
	if v, ok := o.(*JsonNullElement); ok {
		return v
	}
	return nil
}

func BoolElement(o JsonElement) *JsonBoolElement {
	if v, ok := o.(*JsonBoolElement); ok {
		return v
	}
	return nil
}
