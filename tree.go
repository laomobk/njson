package njson

type jsonBaseAST interface {
	Execute() JsonElement
}

type jsonCellSlotAST struct {
	vfloat  float64
	vint    int64
	vstring string
}

type jsonDictAST struct {
	vmap map[string]*JsonElement
}

type jsonArrayAST struct {
	varray []*JsonElement
}

type jsonBoolAST struct {
	vbool bool
}
