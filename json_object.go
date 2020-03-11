package njson

type JsonObject struct {
	_dict *JsonDictElement
}

// for parser
func newJsonObjectFromDictElement(dict *JsonDictElement) *JsonObject {
	return &JsonObject{
		_dict: dict,
	}
}

func (self *JsonObject) String() string {
	return self._dict.String()
}

func (self *JsonObject) Get(path string) (JsonElement, error) {
	return self._dict.Get(path)
}

func (self *JsonObject) DGet(path string) JsonElement {
	return self._dict.DGet(path)
}

func (self *JsonObject) ToDictElement() *JsonDictElement {
	return self._dict
}

func (self *JsonObject) ForEach(forfunc func(string, JsonElement)) {
	self._dict.ForEach(forfunc)
}
