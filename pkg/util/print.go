package util

import (
	"fmt"

	"github.com/TylerBrock/colorjson"
)

func PrintJSON(obj interface{}) {
	var mapData map[string]interface{}
	if err := Recast(obj, &mapData); err != nil {
		return
	}

	f := colorjson.NewFormatter()
	f.Indent = 4
	s, _ := f.Marshal(mapData)
	fmt.Println(string(s))
}
