package main

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
)

const DATE_TIME_FORMAT = "20060102150405"

func main() {
	/*
	   f := func() int { return 100 }
	   fmt.Printf("%v\n", f())
	   fmt.Printf("%v\n", &f)

	   map1 := map[string]interface{}{}
	   map1["foo"] = "bar"
	   map1["hoge"] = "fuga"
	   map1["int"] = 10
	   fmt.Printf("%v\n", map1)

	   map2 := map[string]interface{}{}
	   map2["foo"] = "bar"
	   map2["hoge"] = "fuga"
	   map2["int"] = 10
	   fmt.Printf("%v\n", map2)

	   var comp bool

	   comp = reflect.DeepEqual(map1, map2)
	   fmt.Printf("%v\n", comp)

	   map1["func"] = f
	   map2["func"] = f
	   comp = reflect.DeepEqual(map1, map2)
	   fmt.Printf("%v\n", comp)

	   f1 := map1["func"]
	   f2 := map2["func"]
	   fmt.Printf("f1 = %t, f2 = %t\n", f1, f2)
	   fmt.Printf("f1() = %v, f2() = %v\n", f1.(func() int)(), f2.(func() int)())

	   map1["func"] = &f
	   map2["func"] = &f
	   comp = reflect.DeepEqual(map1, map2)
	   fmt.Printf("%v\n", comp)

	   f3 := map1["func"]
	   f4 := map2["func"]
	   fmt.Printf("f3 = %t, f4 = %t\n", *f3.(*func() int), *f4.(*func() int))

	   f5 := *f3.(*func() int)
	   f6 := *f4.(*func() int)
	   fmt.Printf("f5() = %v, f6() = %v\n", f5(), f6())

	   fmt.Println("=== toml Mershal() ===")

	   	var valueSection = &ValueSection{
	   		IntVal:  999,
	   		StrVal:  "foo",
	   		BoolVal: false,
	   		ListVal: []string{"foo", "bar", "hoge", "fuga"},
	   		MapVal:  map[string]interface{}{},
	   	}

	   valueSection.MapVal["data"] = map[string]interface{}{"foo": "bar", "hoge": 99}

	   root := &Root{Value: *valueSection}

	   serialized, err := toml.Marshal(root)

	   	if nil != err {
	   		fmt.Printf("%v", err)
	   	} else {

	   		result := string(serialized)
	   		fmt.Printf("%v", result)

	   		defer func() {
	   			e := recover()
	   			if nil != e {
	   				fmt.Printf("%v\n", e)
	   			}
	   		}()

	   		fmt.Println("=== toml UnMershal() ===")

	   		var v interface{}
	   		err := toml.Unmarshal([]byte(result), &v)
	   		if nil != err {
	   			panic(err)
	   		}
	   		fmt.Printf("%v\n", v)
	   	}

	   fmt.Println("=== string Concat() ===")
	   s1 := "あいうえお"
	   s2 := "かきくけこ"
	   r1 := Concat(s1, s2)
	   fmt.Printf(">>> %s\n", r1)

	   fmt.Println("=== Date Format ===")
	   d1 := time.Now()
	   fmt.Printf(d1.Format(DATE_TIME_FORMAT))
	*/

	// toml Unmarshal
	data := []byte(`['example.com/foo/bar']
ary = ['foo', 'bar']
bool = true
int = 99
str = 'value_1'
`)

	var v map[string]map[string]interface{}
	if err := toml.Unmarshal(data, &v); nil != err {
		panic(err)
	}
	fmt.Printf("%v\n", v)
}

type Root struct {
	Value ValueSection `toml:"STRUCT_VALUES"`
}

type ValueSection struct {
	IntVal  int                    `toml:"val_1"`                        // 数値
	StrVal  string                 `toml:"val_2"`                        // 文字列
	BoolVal bool                   `toml:"val_3"`                        // 論理値
	ListVal []string               `toml:"val_4, multiline, omitempty"`  // リスト
	MapVal  map[string]interface{} `toml:"mapVal", multiline, omitempty` // map
}
