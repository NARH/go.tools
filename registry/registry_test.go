//
//
//

package registry

import (
	"fmt"
	"testing"
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
