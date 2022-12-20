//
//
//

/*
レジストリ格納先を操作する
*/
package registry

import (
	"fmt"
	"sort"
	"time"

	"github.com/NARH/go.tools/logging"
	"github.com/c2fo/vfs/v6/vfssimple"
	"github.com/pelletier/go-toml/v2"
)

// レジストリの基本機能
type Registry interface {
	// レジストリの追加を行う
	Add(h hive, r registry)
	// レジストリの検索を行う
	Lookup(h hive, keys ...string) (*registry, error)
	// レジストリの削除を行う
	Delete(h hive, keys ...string) error
	// レジストリを保存する
	Store(f string) error
	// レジストリを復元する
	Restore(f string) error
}

// Registry Store の構造体
type Store struct {
	store map[hive]registry
}

// 新規レジストリストアを取得する
func NewStore() *Store {
	return &Store{
		store: map[hive]registry{},
	}
}

// レジストリストアにレジストリを追加する
func (s *Store) Add(h hive, r *registry) {
	s.store[h] = *r
}

// レジストリストアからレジストリを検索する
func (s *Store) Lookup(h hive, keys ...string) (*registry, error) {
	r, ok := s.store[h]

	if !ok {
		return nil, fmt.Errorf("No such hive [%v]", h.string())
	}

	// キー指定がある場合は対応するキーのみを返す
	if 0 < len(keys) {
		flg := false
		rr := NewRegistry()

		for _, key := range keys {
			if val, ok := r.data[key]; ok {
				rr.Append(key, val)

				flg = true
			}
		}

		if flg {
			return rr, nil
		} else {
			// 指定した全てのキーが存在しない場合はエラー
			return nil, fmt.Errorf("Not all keys exist. %v", keys)
		}
	} else {
		// キー指定がない場合は全件を返す
		return &r, nil
	}
}

// レジストリストアからレジストリを削除する
func (s *Store) Delete(h hive, keys ...string) error {
	r, ok := s.store[h]
	if !ok {
		return fmt.Errorf("No such hive [%v]", h.string())
	}

	if 0 < len(keys) {
		flg := false

		for _, key := range keys {
			if _, ok := r.data[key]; ok {
				delete(r.data, key)

				flg = true
			}
		}

		if flg {
			return nil
		} else {
			// 指定した全てのキーが存在しない場合はエラー
			return fmt.Errorf("Not all keys exist. %v", keys)
		}
	} else {
		// キー指定がない場合はhive を消す
		delete(s.store, h)
		return nil
	}
}

// レジストリを保存する
func (s *Store) Store(f string) error {
	file, err := vfssimple.NewFile(f)
	if nil != err {
		return err
	}

	back, err := vfssimple.NewFile(s.Concat(f, time.Now().Format(".20060102150405")))
	if nil != err {
		return err
	}

	toml, err := s.ToToml()
	if nil != err {
		return err
	}

	_, err = file.Write([]byte(toml))
	if nil != err {
		return err
	}

	file.Close()

	if err := file.CopyToFile(back); nil != err {
		return err
	}

	return nil
}

// レジストリを復元する
func (s *Store) Restore(f string) error {
	file, err := vfssimple.NewFile(f)
	if nil != err {
		return err
	}

	ok, err := file.Exists()
	if nil != err {
		return err
	}

	if !ok {
		return fmt.Errorf("toml file [%s] not exists.", f)
	} else {

		buf := make([]byte, 1000)
		n, err := file.Read(buf)
		if nil != err {
			return fmt.Errorf("cannot read toml file [%s].", f)
		}

		str := string(buf[:n])
		logging.NewLogger().Info(">>> \n%s", str)
		return s.FromToml(str)
	}
}

func (s *Store) Concat(s1, s2 string) string {
	c := len([]byte(s1))
	c += len([]byte(s2))
	ss := make([]byte, 0, c)
	ss = append(ss, s1...)
	ss = append(ss, s2...)

	return string(ss)
}

func (s *Store) ToToml() (string, error) {
	root := map[string]map[string]interface{}{}

	// mapのキーサイズは検討が付かないので適当な capacityを当てる
	hives := make([]hive, 0, 1000)
	for k := range s.store {
		hives = append(hives, k)
	}

	sort.Slice(hives, func(i, j int) bool {
		return hives[i].string() < hives[j].string()
	})

	for _, k := range hives {
		if _, ok := root[k.string()]; !ok {
			root[k.string()] = make(map[string]interface{})
		}

		var keys []string
		for kk := range s.store[k].data {
			keys = append(keys, kk)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, kk := range keys {
			root[k.string()][kk] = s.store[k].data[kk]
		}
	}

	toml, err := toml.Marshal(&root)
	if nil != err {
		return "", err
	}

	return string(toml), nil
}

func (s *Store) FromToml(t string) error {
	var root map[string]map[string]interface{}

	err := toml.Unmarshal([]byte(t), &root)
	if nil != err {
		return err
	}

	logging.NewLogger().Info(">>> \n%v", root)

	return nil
}
