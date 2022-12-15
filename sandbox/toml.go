//
//
//

/*
 */
package main

import (
	"github.com/pelletier/go-toml/v2"
)

func Marshal(r interface{}) (string, error) {
	serialized, err := toml.Marshal(r)
	if nil != err {
		return "", err
	}

	return string(serialized), nil
}

func UnMarshal(src string) (v interface{}) {
	err := toml.Unmarshal([]byte(src), &v)
	if nil != err {
		panic(err)
	}
	return v
}
