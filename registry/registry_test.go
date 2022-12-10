//
//
//

package registry

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/NARH/go.tools/logging"
)

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

// 新規レジストリデータを作成する
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
			name: "正常系(func追加)",
			key:  "test_6",
			data: f1,
			want: map[string]interface{}{"test_6": f1},
		}, {
			name: "正常系(nil追加)",
			key:  "test_7",
			data: nil,
			want: map[string]interface{}{"test_7": nil},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewRegistry()
			registry.Append(tt.key, tt.data)
			if !reflect.DeepEqual(registry.data, tt.want) {
				t.Errorf("registry.Append() = %v, but want %v", registry.data, tt.want)
			}
			logging.NewLogger().Debug("%v", registry.data)
		})
	}
}

// レジストリパッケージ関数 Add() の試験
func TestAdd(t *testing.T) {
	t.Fail()
}

// レジストリパッケージ関数 Lookup() の試験
func TestLookup(t *testing.T) {
	t.Fail()
}

// レジストリパッケージ関数 Delete() の試験
func TestDelete(t *testing.T) {
	t.Fail()
}
