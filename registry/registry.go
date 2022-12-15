//
//
//

/*
任意の値をHiveに格納するレジストリとして利用
基本的なデータ構造は以下の通り

URI表記(schemeは省略): hive(巣箱)として利用
data: Key(string)=Value(interface) pairs

ex:

	registry := map[hive]*registry{}

	  "example.com/feature/example/foo":  map[string]{},
	  "example.com/feature/example/bar":  map[string]{},
	  "example.com/constant/name":        map[string]{},
	  "hoge.example.com/feature/example": map[string]{},
	  "fuga.example.com/constant/name":    map[string]{},
*/
package registry

import (
	"fmt"
	"regexp"
)

// Hive 構造体
type hive struct {
	uri string
}

// レジストリデータ構造体
type registry struct {
	data map[string]interface{}
}

// hive の Stringer method
func (h *hive) string() string {
	return h.uri
}

// hive チェックパターン
/*
 * RFC3986 準拠ではない。
 * とりあえず、以下の条件にマッチするものとする
 *
 *   host part: "[a-zA-z]([a-zA-Z0-9+\-.])*"
 *   port part: ":?"
 *   path part: "(\/[a-zA-z0-9_+\-.\/]*)*"
 */
const HIVE_PTN = `^[a-zA-z]([a-zA-Z0-9+\-.])*(:[0-9]+)?(\/[a-zA-z0-9_+\-.\/]*)*([^/])$`

// hive はスラッシュ / で終わってはいけない
const HIVE_NOT_END_WITH = `/`

// hive を作成する
func hiveCreate(uri string) hive {
	// 長さ 0 は禁止は禁止
	if 0 >= len(uri) {
		panic(fmt.Errorf("hive name is empty."))
	}

	r := regexp.MustCompile(HIVE_PTN)

	if r.MatchString(uri) {
		return hive{uri: uri}
	}

	// NGの場合はpanicでerrorを投げる
	panic(fmt.Errorf("Illegal hive name [%s]", uri))
}

// 新規レジストリデータを作成する
func NewRegistry() *registry {
	return &registry{data: map[string]interface{}{}}
}

// レジストリデータ最大サイズ長(以下)
const MAX_KEY_LENGTH = 255

// keyのチェック
func (r *registry) keyCheck(key string) {
	// Keyサイズチェック
	if 1 > len(key) || MAX_KEY_LENGTH < len(key) {
		panic(fmt.Errorf("Illegal key length. [%s:%d]", key, len(key)))
	}
}

// レジストリデータにKey=Valueを追加する
func (r *registry) Append(key string, val interface{}) {
	// keyのチェック
	r.keyCheck(key)

	// Appendなので登録済みのKeyはエラー
	if _, ok := r.data[key]; ok {
		panic(fmt.Errorf("key [%s] is already registered.", key))
	}

	r.data[key] = val
}

// レジストリデータからKeyに対応したValueを取得する
func (r *registry) Get(key string) (result interface{}, err error) {
	defer func() {
		// keyのチェック時にエラーとして返す
		e := recover()
		if nil != e {
			if ee, ok := e.(error); ok {
				err = ee
			}
		}
	}()

	r.keyCheck(key)

	if val, ok := r.data[key]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("key [%s] is not registered.", key)
	}
}

// レジストリデータにKeyに対応したValueを設定する
// 登録済みキーであってもエラーにならずに上書きする
func (r *registry) Set(key string, val interface{}) {
	r.keyCheck(key)
	r.data[key] = val
}

// レジストリデータのKey=Valueを削除する
func (r *registry) Remove(key string) (err error) {
	defer func() {
		// keyのチェック時にエラーとして返す
		e := recover()
		if nil != e {
			if ee, ok := e.(error); ok {
				err = ee
			}
		}
	}()

	r.keyCheck(key)

	if _, ok := r.data[key]; ok {
		delete(r.data, key)
		return nil
	} else {
		return fmt.Errorf("key [%s] is not registered.", key)
	}
}

// Registry Store レジストリの格納場所
var store *Store = NewStore()

// レジストリパッケージ関数 Add()
// Registry.Add() のラッパー関数
//
// ex:
// registory.Add(h, r) のように使う
func Add(h string, r *registry) {
	hive := hiveCreate(h)
	store.Add(hive, r)
}

// レジストリパッケージ関数 Lookup()
// Registry.Lookup() のラッパー関数
//
// ex:
// registry.Lookup(r, s) のように使う
func Lookup(h string, keys ...string) (*registry, error) {
	hive := hiveCreate(h)
	return store.Lookup(hive, keys...)
}

// レジストリパッケージ関数 Delete()
// Registry.Delete() のラッパー関数
//
// ex:
// registry.Delete() のように使う
func Delete(h string, keys ...string) error {
	hive := hiveCreate(h)
	return store.Delete(hive, keys...)
}
