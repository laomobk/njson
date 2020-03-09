package njson

import (
	"fmt"
	"strings"
)

type JsonElement interface {
	ToString() string
	ToInteger64() int64
	ToFloat64() float64
	ToElementArray() []JsonElement
	ToDict() map[string]JsonElement
	ToBool() bool
	Type() int
	String() string
	Raw() interface{}
}

type JsonBaseElement struct {
	// nothing
}

// Base method, sub class can overwrite these method.
func (self *JsonBaseElement) ToString() string               { return "" }
func (self *JsonBaseElement) ToInteger64() int64             { return 0 }
func (self *JsonBaseElement) ToFloat64() float64             { return 0.0 }
func (self *JsonBaseElement) ToElementArray() []JsonElement  { return nil }
func (self *JsonBaseElement) ToDict() map[string]JsonElement { return nil }
func (self *JsonBaseElement) ToBool() bool                   { return false }
func (self *JsonBaseElement) Type() int                      { return ELE_BASE }
func (self *JsonBaseElement) String() string                 { return "< JsonBaseElement >" }
func (self *JsonBaseElement) Raw() interface{}               { return nil }

type JsonStringElement struct {
	JsonBaseElement

	value string
}

func (self *JsonStringElement) Raw() interface{} { return self.ToString() }
func (self *JsonStringElement) ToString() string { return self.value }
func (self *JsonStringElement) Type() int        { return ELE_STRING }
func (self *JsonStringElement) String() string   { return fmt.Sprintf("'%s'", self.value) }

type _JsonNumberSlot struct {
	vfloat float64
	vint   int64
}

type JsonIntegerElement struct {
	JsonBaseElement

	slot *_JsonNumberSlot
}

func (self *JsonIntegerElement) Raw() interface{}   { return self.ToInteger64() }
func (self *JsonIntegerElement) ToInteger64() int64 { return self.slot.vint }
func (self *JsonIntegerElement) Type() int          { return ELE_INTEGER }
func (self *JsonIntegerElement) String() string     { return fmt.Sprint(self.slot.vint) }

type JsonFloatElement struct {
	JsonBaseElement

	slot *_JsonNumberSlot
}

func (self *JsonFloatElement) Raw() interface{}   { return self.ToFloat64() }
func (self *JsonFloatElement) ToFloat64() float64 { return self.slot.vfloat }
func (self *JsonFloatElement) Type() int          { return ELE_FLOAT }
func (self *JsonFloatElement) String() string     { return fmt.Sprint(self.slot.vfloat) }

type JsonArrayElement struct {
	JsonBaseElement

	array []JsonElement
}

func (self *JsonArrayElement) Raw() interface{}              { return self.ToElementArray() }
func (self *JsonArrayElement) ToElementArray() []JsonElement { return self.array }
func (self *JsonArrayElement) Type() int                     { return ELE_ARRAY }

func (self *JsonArrayElement) String() string {
	item := make([]string, len(self.array))

	for i, _ := range self.array {
		item[i] = self.array[i].String()
	}

	return "[" + strings.Join(item, ", ") + "]"
}

type JsonBoolElement struct {
	JsonBaseElement

	value bool
}

func (self *JsonBoolElement) Raw() interface{} { return self.ToBool() }
func (self *JsonBoolElement) ToBool() bool     { return self.value }
func (self *JsonBoolElement) Type() int        { return ELE_BOOL }
func (self *JsonBoolElement) String() string   { return fmt.Sprintf("%v", self.value) }

type JsonDictElement struct {
	JsonBaseElement

	dict map[string]JsonElement
}

func (self *JsonDictElement) Raw() interface{}               { return self.ToDict() }
func (self *JsonDictElement) ToDict() map[string]JsonElement { return self.dict }
func (self *JsonDictElement) Type() int                      { return ELE_DICT }
func (self *JsonDictElement) String() string {
	item := []string{}

	for k, v := range self.dict {
		item = append(item, fmt.Sprintf("%s : %s", k, v.String()))
	}

	return "{" + strings.Join(item, ", ") + "}"
}

func (self *JsonDictElement) Get(path string) (JsonElement, error) {
	parts := strings.Split(path, ".")
	fmt.Printf("%v\n", parts)

	var left []string
	var attr string

	if len(parts) > 1 {
		left = parts[:len(parts)-1]
		attr = parts[len(parts)-1]
	} else {
		left = []string{}
		attr = path
	}
	var last *JsonDictElement

	start := self
	last = start

	for _, v := range left {
		if _, ok := last.dict[v].(*JsonDictElement); !ok {
			return nil, fmt.Errorf(path + " : element '" + v + "' is not a dict.")
		}

		temp, ok := last.dict[v]
		if !ok {
			return nil, fmt.Errorf(path + " : key '" + v + "' is not exists")
		}

		last = temp.(*JsonDictElement)
	}

	v, ok := last.dict[attr]
	if !ok {
		return nil, fmt.Errorf(path + " : key '" + attr + "' is not exists")
	}
	return v, nil
}

func (self *JsonDictElement) DGet(path string) JsonElement {
	ele, err := self.Get(path)

	if err != nil {
		panic(err)
	}

	return ele
}

type JsonNullElement struct {
	JsonBaseElement
}

func (self *JsonNullElement) Type() int      { return ELE_NULL }
func (self *JsonNullElement) String() string { return "null" }

// ElementFactory

/*
 * this function can create object from the type of value
 */
func NewJsonElementByValue(value interface{}) JsonElement {
	switch value.(type) {

	case int:
		return &JsonIntegerElement{
			slot: &_JsonNumberSlot{
				vint: int64(value.(int)),
			},
		}

	case int64:
		return &JsonIntegerElement{
			slot: &_JsonNumberSlot{
				vint: value.(int64),
			},
		}

	case float32:
		return &JsonIntegerElement{
			slot: &_JsonNumberSlot{
				vfloat: float64(value.(float32)),
			},
		}

	case float64:
		return &JsonFloatElement{
			slot: &_JsonNumberSlot{
				vfloat: value.(float64),
			},
		}

	case string:
		return &JsonStringElement{
			value: value.(string),
		}

	}

	return nil
}
