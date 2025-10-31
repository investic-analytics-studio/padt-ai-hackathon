package util

import "github.com/bytedance/sonic"

func Recast(from, to interface{}) error {
	switch v := from.(type) {
	case []byte:
		return sonic.Unmarshal(v, to)
	default:
		buf, err := sonic.Marshal(from)
		if err != nil {
			return err
		}

		return sonic.Unmarshal(buf, to)
	}
}
