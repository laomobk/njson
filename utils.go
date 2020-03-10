package njson

func inRuneArray(r []rune, rune_ rune) bool {
	for _, v := range r {
		if v == rune_ {
			return true
		}
	}
	return false
}
