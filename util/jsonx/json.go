package jsonx

import "github.com/bytedance/sonic"

var API sonic.API

func init() {
	cfg := sonic.Config{
		EscapeHTML:              false,
		SortMapKeys:             true,
		CompactMarshaler:        true,
		NoValidateJSONMarshaler: true,
	}

	API = cfg.Froze()
}

func Marshal(v any) ([]byte, error) {
	return API.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	return API.Unmarshal(data, v)
}

func MarshalString(v any) (string, error) {
	return API.MarshalToString(v)
}
