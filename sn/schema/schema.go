package schema

import "github.com/gorilla/schema"

var decoder = schema.NewDecoder()

func Decode(dst interface{}, src map[string][]string) error {
	return decoder.Decode(dst, src)
}

func RegisterConverter(value interface{}, converterFunc schema.Converter) {
	decoder.RegisterConverter(value, converterFunc)
}
