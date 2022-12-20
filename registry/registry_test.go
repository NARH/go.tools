//
//
//

package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/NARH/go.tools/logging"
)

const SAMPLE_HIVE_NAME = "example.com/foo/bar"

var log = logging.NewLogger()

// hive を作成試験
func TestHiveCreate(t *testing.T) {
	/* 試験ケース */
	testCases := []struct {
		/* 試験ケースの構造体 */
		name string      // 試験名称
		hive string      // 試験対象のhive
		want interface{} // 期待値
	}{
		{
			name: "正常系(foo)",
			hive: "foo",
			want: hive{uri: "foo"},
		}, {
			name: "正常系(foo/bar)",
			hive: "foo/bar",
			want: hive{uri: "foo/bar"},
		}, {
			name: "正常系(example.com)",
			hive: "example.com",
			want: hive{uri: "example.com"},
		}, {
			name: "正常系(example.com/feature)",
			hive: "example.com/feature",
			want: hive{uri: "example.com/feature"},
		}, {
			name: "正常系(example.com/feature/example/foo)",
			hive: "example.com/feature/example/foo",
			want: hive{uri: "example.com/feature/example/foo"},
		}, {
			name: "正常系(example.com:8080/feature/example/foo)",
			hive: "example.com:8080/feature/example/foo",
			want: hive{uri: "example.com:8080/feature/example/foo"},
		}, {
			name: "異常系(空文字)",
			hive: "",
			want: fmt.Errorf("hive name is empty."),
		}, {
			name: "異常系(接尾が/で終わる)",
			hive: "example.com/feature/example/foo/",
			want: fmt.Errorf("Illegal hive name [%s]", "example.com/feature/example/foo/"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				e := recover()

				if want, ok := tt.want.(error); ok && nil != e {
					// want が error かつ panic 時のみ検証
					ee, _ := e.(error)
					if ee.Error() != want.Error() {
						t.Errorf("createHive() = %v, but want %v", e, tt.want)
						t.Fail()
					} else if (ok && nil == e) || (!ok && nil != e) {
						// want が error または panic いずれか一方だけの場合
						t.Errorf("createHive() returned not error")
					}
				}
			}()

			h := hiveCreate(tt.hive)
			if want, ok := tt.want.(hive); ok {
				if want.string() != h.string() {
					t.Errorf("createHive() = %v, but want %v", h, tt.want)
				}
			}
		})
	}
}

// エラーとなるキーサイズ
const WORD_SIZE_256 = "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456"

// 新規レジストリデータを作成する試験 / レジストリデータを登録する試験
func TestNewRegistry(t *testing.T) {
	t.Run("インスタンス取得正常系試験", func(t *testing.T) {
		registry := NewRegistry()
		if nil == registry {
			t.Errorf("NewRegistry() returned nil.")
			t.FailNow()
		}
	})

	// 試験用関数
	f1 := func() int { return 100 }

	// 新規レジストリデータ追加試験
	testCases := []struct {
		name string
		key  string
		data interface{}
		want map[string]interface{}
		err  error
	}{
		{
			name: "正常系(文字列追加)",
			key:  "test_1",
			data: "foo",
			want: map[string]interface{}{"test_1": "foo"},
		}, {
			name: "正常系(数値追加)",
			key:  "test_2",
			data: 100,
			want: map[string]interface{}{"test_2": 100},
		}, {
			name: "正常系(bool追加)",
			key:  "test_3",
			data: true,
			want: map[string]interface{}{"test_3": true},
		}, {
			name: "正常系(array追加)",
			key:  "test_4",
			data: [2]string{"foo", "bar"},
			want: map[string]interface{}{"test_4": [2]string{"foo", "bar"}},
		}, {
			name: "正常系(slice追加)",
			key:  "test_5",
			data: []int{100, 200},
			want: map[string]interface{}{"test_5": []int{100, 200}},
		}, {
			name: "正常系(map追加)",
			key:  "test_6",
			data: map[string]int{"foo": 100, "bar": 200},
			want: map[string]interface{}{"test_6": map[string]int{"foo": 100, "bar": 200.}},
		}, {
			name: "正常系(func追加)",
			key:  "test_7",
			data: &f1,
			want: map[string]interface{}{"test_7": &f1},
		}, {
			name: "正常系(nil追加)",
			key:  "test_8",
			data: nil,
			want: map[string]interface{}{"test_8": nil},
		}, {
			name: "正常系(nested array追加)",
			key:  "test_9",
			data: [2][2]string{{"foo", "bar"}, {"hoge", "fuga"}},
			want: map[string]interface{}{"test_9": [2][2]string{{"foo", "bar"}, {"hoge", "fuga"}}},
		}, {
			name: "正常系(nested slice追加)",
			key:  "test_10",
			data: [][]string{{"foo", "bar"}, {"hoge", "fuga"}},
			want: map[string]interface{}{"test_10": [][]string{{"foo", "bar"}, {"hoge", "fuga"}}},
		}, {
			name: "正常系(nested map追加)",
			key:  "test_10",
			data: map[string]map[int]int{"foo": {1: 10}, "bar": {2: 100}},
			want: map[string]interface{}{"test_10": map[string]map[int]int{"foo": {1: 10}, "bar": {2: 100}}},
		}, {
			name: "異常系(キーサイズエラー)",
			key:  WORD_SIZE_256,
			data: "foo",
			err:  fmt.Errorf("Illegal key length. [%s:%d]", WORD_SIZE_256, len(WORD_SIZE_256)),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			defer func() {
				e := recover()

				if nil != e && e.(error).Error() != tt.err.Error() {
					t.Errorf("registry.Add() = %v, but want %v", e, tt.err)
					t.Fail()
				}
			}()

			registry := NewRegistry()
			registry.Append(tt.key, tt.data)
			if !reflect.DeepEqual(registry.data, tt.want) {
				t.Errorf("registry.Append() = %v, but want %v", registry.data, tt.want)
			}
			log.Debug("%v", registry.data)
		})
	}

	t.Run("異常系(キー重複)", func(t *testing.T) {
		defer func() {
			e := recover()

			msg := fmt.Errorf("key [foo] is already registered.")
			if nil != e && e.(error).Error() != msg.Error() {
				t.Errorf("registry.Add() = %v, but want %v", e, msg)
				t.Fail()
			}
		}()

		registry := NewRegistry()
		registry.Append("foo", "bar")
		registry.Append("foo", "hoge")
	})
}

// レジストリデータからKeyに対応したValueを取得する試験
func TestGet(t *testing.T) {
	t.Run("正常系試験", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("foo", "bar")

		if result, err := registry.Get("foo"); nil == err {
			if "bar" != result.(string) {
				t.Errorf("regitry.Get() = %v, but want %v", result, "bar")
			}
		} else {
			t.Errorf(err.Error())
		}
	})

	t.Run("異常系試験(キーなし)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("foo", "bar")

		result, err := registry.Get("bar")
		want := fmt.Errorf("key [%s] is not registered.", "bar")

		if nil == err {
			t.Errorf("registry.Get() not registered key[%s], but got value[%v].", "bar", result)
		} else if want.Error() != err.Error() {
			t.Errorf("registry.Get() = %v, but want %v", err, want)
		}
	})

	t.Run("異常系試験(キーサイズエラー)", func(t *testing.T) {
		registry := NewRegistry()
		_, err := registry.Get(WORD_SIZE_256)

		want := fmt.Errorf("Illegal key length. [%s:%d]", WORD_SIZE_256, len(WORD_SIZE_256))

		if nil == err {
			t.Errorf("registry.Get() key length over, not errored.")
		} else if want.Error() != err.Error() {
			t.Errorf("registry.Get() = %v, but want %v", err, want)
		}
	})
}

// レジストリデータにKeyに対応したValueを設定する試験
func TestSet(t *testing.T) {
	t.Run("正常系試験(未登録)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Set("foo", "bar")

		val, err := registry.Get("foo")
		if nil != err {
			t.Fail()
		}

		result, ok := val.(string)
		if !ok {
			t.Fail()
		}

		if "bar" != result {
			t.Errorf("registry.Set() = %v, but want %v", result, "bar")
		}
	})

	t.Run("正常系試験(キー登録済み)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("foo", "bar")
		registry.Set("foo", "hoge")

		val, err := registry.Get("foo")
		if nil != err {
			t.Fail()
		}

		result, ok := val.(string)
		if !ok {
			t.Fail()
		}

		if "hoge" != result {
			t.Errorf("registry.Set() = %v, but want %v", result, "hoge")
		}
	})
}

// レジストリデータのKey=Valueを削除する試験
func TestRemove(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("foo", "bar")
		err := registry.Remove("foo")
		if nil != err {
			t.Error(err)
		}
	})

	t.Run("異常系(未登録のキー)", func(t *testing.T) {
		registry := NewRegistry()
		err := registry.Remove("foo")
		want := fmt.Errorf("key [%s] is not registered.", "foo")

		if nil == err {
			t.Fail()
		} else if want.Error() != err.Error() {
			t.Errorf("registry.Remove() = %v, but want %v", err, want)
		}
	})

	t.Run("異常系(キーサイズエラー)", func(t *testing.T) {
		registry := NewRegistry()
		err := registry.Remove(WORD_SIZE_256)
		want := fmt.Errorf("Illegal key length. [%s:%d]", WORD_SIZE_256, len(WORD_SIZE_256))

		if nil == err {
			t.Fail()
		} else if want.Error() != err.Error() {
			t.Errorf("registry.Remove() = %v, but want %v", err, want)
		}
	})
}

// レジストリパッケージ関数 Add() の試験
func TestAdd(t *testing.T) {
	t.Run("正常系試験", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("hoge", "fuga")
		Add(SAMPLE_HIVE_NAME, registry)

		h := hiveCreate(SAMPLE_HIVE_NAME)
		result, ok := store.store[h]
		want := map[string]interface{}{
			"hoge": "fuga",
		}

		if !ok {
			t.Errorf("hive [%v] not registered", SAMPLE_HIVE_NAME)
		}

		if !reflect.DeepEqual(result.data, want) {
			t.Errorf("Add() = %v, but want %v", result.data, want)
		}
	})
}

// レジストリパッケージ関数 Lookup() の試験
func TestLookup(t *testing.T) {
	t.Run("正常系試験(hive で検索)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("test_1", "value_1")
		Add(SAMPLE_HIVE_NAME, registry)

		want := map[string]interface{}{"test_1": "value_1"}
		result, err := Lookup(SAMPLE_HIVE_NAME)
		if nil != err {
			t.Errorf(err.Error())
		} else if !reflect.DeepEqual(result.data, want) {
			t.Errorf("Lookup() = %v, but want %v", result.data, want)
		}
	})

	t.Run("正常系試験(hiveと一つのキーで検索)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("test_1", "value_1")
		registry.Append("test_2", "value_2")
		Add(SAMPLE_HIVE_NAME, registry)

		want := map[string]interface{}{"test_1": "value_1"}
		result, err := Lookup(SAMPLE_HIVE_NAME, "test_1")
		log.Info("%v", result)
		if nil != err {
			t.Errorf(err.Error())
		} else if !reflect.DeepEqual(result.data, want) {
			t.Errorf("Lookup() = %v, but want %v", result.data, want)
		}
	})

	t.Run("異常系試験(hiveなし)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("test_1", "value_1")
		Add(SAMPLE_HIVE_NAME, registry)

		h := "example.com/hoge"
		want := fmt.Errorf("No such hive [%v]", h)

		_, err := Lookup(h)
		if nil == err {
			t.Errorf("Lookup() = %v, but want %v", nil, want)
		} else if want.Error() != err.Error() {
			t.Errorf("Lookup() = %v, but want %v", err, want)
		}
	})

	t.Run("異常系試験(該当キーなし)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("test_1", "value_1")
		Add(SAMPLE_HIVE_NAME, registry)

		keys := []string{"test_99"}
		want := fmt.Errorf("Not all keys exist. %v", keys)

		_, err := Lookup(SAMPLE_HIVE_NAME, keys...)
		if want.Error() != err.Error() {
			t.Errorf("Lookup() = %v, but want %v", err, want)
		}
	})
}

// レジストリパッケージ関数 Delete() の試験
func TestDelete(t *testing.T) {
	t.Run("正常系試験(hiveを消す)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("test_1", "value_1")
		Add(SAMPLE_HIVE_NAME, registry)

		if err := Delete(SAMPLE_HIVE_NAME); nil != err {
			t.Errorf("Delete() = %v, but want %v", err, nil)
		}

		want := fmt.Errorf("No such hive [%v]", SAMPLE_HIVE_NAME)
		if _, err := Lookup(SAMPLE_HIVE_NAME); want.Error() != err.Error() {
			t.Errorf("Lookup() = %v, but want %v", err, want)
		}
	})

	t.Run("正常系試験(keyを指定して消す)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("test_1", "value_1")
		registry.Append("test_2", "value_2")
		Add(SAMPLE_HIVE_NAME, registry)

		if err := Delete(SAMPLE_HIVE_NAME, "test_1"); nil != err {
			t.Errorf("Delete() = %v, but want %v", err, nil)
		}

		// 消えた事の確認
		keys := [1]string{"test_1"}
		want := fmt.Errorf("Not all keys exist. %v", keys)

		_, err := Lookup(SAMPLE_HIVE_NAME, "test_1")
		if want.Error() != err.Error() {
			t.Errorf("Lookup() = %v, but want %v", err, want)
		}

		// 消えていない事の確認
		r, err := Lookup(SAMPLE_HIVE_NAME, "test_2")
		if nil != err {
			t.Errorf("Lookup() = %v, but want %v", err, nil)
		}

		if _, ok := r.data["test_2"]; !ok {
			t.Errorf("Delete() deleted [%s]", "test_2")
		}

		log.Debug("%v", r.data)
	})

	t.Run("異常系試験(hiveなし)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("test_1", "value_1")
		Add(SAMPLE_HIVE_NAME, registry)

		want := fmt.Errorf("No such hive [%v]", "example.com/hoge")
		if err := Delete("example.com/hoge"); want.Error() != err.Error() {
			t.Errorf("Delete() = %v, but want %v", err, want)
		}

		r, err := Lookup(SAMPLE_HIVE_NAME)
		if nil != err {
			t.Errorf("Delete() deleted [%s]", SAMPLE_HIVE_NAME)
		}
		if 1 != len(r.data) {
			t.Errorf("Delete() = %v, but want %v", len(r.data), 1)
		}
		log.Debug("%v", r)
	})

	t.Run("異常系試験(該当キーなし)", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("test_1", "value_1")
		Add(SAMPLE_HIVE_NAME, registry)

		keys := []string{"test_99"}
		want := fmt.Errorf("Not all keys exist. %v", keys)

		if err := Delete(SAMPLE_HIVE_NAME, keys...); want.Error() != err.Error() {
			t.Errorf("Delete() = %v, but want %v", err, want)
		}

		r, err := Lookup(SAMPLE_HIVE_NAME, "test_1")
		if nil != err {
			t.Errorf("Delete() deleted [%v]", "test_1")
		}
		if 1 != len(r.data) {
			t.Errorf("Delete() = %v, but want %v", len(r.data), 1)
		}
	})
}

func TestSave(t *testing.T) {
	t.Run("正常系試験", func(t *testing.T) {
		registry := NewRegistry()
		registry.Append("str", "value_1")
		registry.Append("int", 99)
		registry.Append("bool", true)
		registry.Append("ary", []string{"foo", "bar"})
		registry.Append("nil", nil)
		Add(SAMPLE_HIVE_NAME, registry)

		defer func() {
			e := recover()
			if nil != e {
				t.Errorf("error %v", e)
			}
		}()

		f := filepath.Join("file:", os.TempDir(), "test.toml")
		logging.NewLogger().Info(">>> out of toml [%s]", f)
		Save(f)

		Restore(f)
	})
}
