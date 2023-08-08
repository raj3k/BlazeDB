package blazedb

import (
	"fmt"
	"os"

	bolt "go.etcd.io/bbolt"
)

const (
	defaultDBName = "blaze"
	extension     = "db"
)

type Map map[string]any

type BlazeDB struct {
	currentDatabase string
	*Options
	db *bolt.DB
}

func New(options ...OptFunc) (*BlazeDB, error) {
	opts := &Options{
		DBName: defaultDBName,
	}

	dbname := fmt.Sprintf("%s.%s", opts.DBName, extension)
	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &BlazeDB{
		currentDatabase: dbname,
		Options:         opts,
		db:              db,
	}, nil
}

func (b *BlazeDB) DropDatabase(name string) error {
	dbname := fmt.Sprintf("%s.%s", name, extension)
	return os.Remove(dbname)
}

func (b *BlazeDB) Set(fn func(*bolt.Tx) error) error {
	t, err := b.db.Begin(true)
	if err != nil {
		return err
	}

	defer t.Rollback()

	err = fn(t)
	if err != nil {
		_ = t.Rollback()
		return err
	}

	return t.Commit()
}

func (b *BlazeDB) Get(fn func(*bolt.Tx) error) error {
	t, err := b.db.Begin(false)
	if err != nil {
		return err
	}
	defer t.Rollback()

	err = fn(t)
	if err != nil {
		_ = t.Rollback()
		return err
	}
	return t.Rollback()
}
