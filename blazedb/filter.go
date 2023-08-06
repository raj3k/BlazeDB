package blazedb

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type Filter struct {
	blazedb *BlazeDB
}

func NewFilter(db *BlazeDB) *Filter {
	return &Filter{
		blazedb: db,
	}
}

func uint64Bytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

func uint64FromBytes(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func (f *Filter) Set(values Map) (uint64, error) {
	tx, err := f.blazedb.db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	collBucket, err := tx.CreateBucketIfNotExists([]byte("bucket"))
	if err != nil {
		return 0, err
	}
	id, err := collBucket.NextSequence()
	if err != nil {
		return 0, err
	}
	b, err := json.Marshal(values)
	if err != nil {
		return 0, err
	}

	fmt.Println(values)

	fmt.Println(string(b))

	if err := collBucket.Put(uint64Bytes(id), b); err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (f *Filter) Get(key uint64) (Map, error) {
	tx, err := f.blazedb.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte("bucket"))
	if bucket == nil {
		return nil, fmt.Errorf("bucket (%s) not found", "bucket")
	}

	v := bucket.Get(uint64Bytes(key))

	var map1 Map
	if err := json.Unmarshal(v, &map1); err != nil {
		return nil, err
	}

	// result := Map{
	// 	key: string(v),
	// }

	fmt.Println(map1)

	return nil, tx.Commit()
}
