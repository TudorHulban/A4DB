package a4db

import (
	"encoding/json"
	"io"
)

type Objects []*Object

func (o Objects) WriteTo(w io.Writer) error {
	var buf []byte

	for _, object := range o {
		bytes, errMa := json.Marshal(object)
		if errMa != nil {
			return errMa
		}

		bytes = append(bytes, []byte("\n")...)

		buf = append(buf, bytes...)
	}

	_, errWr := w.Write(buf)
	return errWr
}
