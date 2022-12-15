//
//
//

/*
レジストリ格納先を操作する
*/
package registry

import (
	"fmt"
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
