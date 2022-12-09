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
}

// Registry Store レジストリの格納場所
var store = map[hive]*registry{}

// レジストリの基本機能
type Registry interface {
	// レジストリの追加を行う
	Add(h hive, r ...registry) error
	// レジストリの検索を行う
	Lookup(h hive, key ...string) (*[]registry, error)
	// レジストリの削除を行う
	Delete(h hive, r ...registry) error
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
func hiveCreate(uri string) *hive {
	if 0 >= len(uri) {
		panic(fmt.Errorf("hive name is empty."))
	}

	r := regexp.MustCompile(HIVE_PTN)

	//if !strings.HasSuffix(uri, HIVE_NOT_END_WITH) || r.MatchString(uri) {
	if r.MatchString(uri) {
		return &hive{uri: uri}
	}

	panic(fmt.Errorf("Illegal hive name [%s]", uri))
}

// レジストリパッケージ関数 Add()
// Registry.Add() のラッパー関数
//
// ex:
// registory.Add(h, r) のように使う
func Add(h string, r ...registry) error {
	err := fmt.Errorf("Add error")
	return err
}

// レジストリパッケージ関数 Lookup()
// Registry.Lookup() のラッパー関数
//
// ex:
// registry.Lookup(r, s) のように使う
func Lookup(h string, key ...string) (*[]registry, error) {
	err := fmt.Errorf("Lookup error")
	return nil, err
}

// レジストリパッケージ関数 Delete()
// Registry.Delete() のラッパー関数
//
// ex:
// registry.Delete() のように使う
func Delete(h string, r ...registry) error {
	err := fmt.Errorf("Delete error")
	return err
}
