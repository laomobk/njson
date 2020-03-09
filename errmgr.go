package njson

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type NJsonError struct {
	message  string
	lno      int
	offset   int
	filepath string
	source   []byte
}

func (self *NJsonError) getLine() string {
	lno := self.lno

	invalid := "(invalid line number)\n"

	if lno <= 0 {
		return invalid
	}

	f := strings.NewReader(string(self.source))

	buf := bufio.NewReader(f)

	for lc := 1; true; lc++ {
		line, err := buf.ReadString('\n')

		if err != nil && err != io.EOF {
			break
		}

		if lno == lc {
			return fmt.Sprintf("file %s : %d : %d :\n   %s\n",
				self.filepath, self.lno, self.offset, line)
		}
	}

	return invalid
}

func (self *NJsonError) Error() string {
	ln1 := self.getLine()
	ln2 := "JsonFormatError : " + self.message

	return ln1 + ln2
}

func (self *NJsonError) ThrowError() {
	panic(self)
}
