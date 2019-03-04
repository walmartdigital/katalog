package persistence

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/seadiaz/katalog/src/domain"

	"github.com/emirpasic/gods/lists/arraylist"
)

// BoltPersistence ...
type BoltPersistence struct {
	driver interface{}
}

// Create ...
func (p *BoltPersistence) Create(kind string, id string, obj interface{}) {
	db := p.driver.(*bolt.DB)
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		objJSON, err := json.Marshal(obj)
		if err != nil {
			glog.Error(err)
			return err
		}
		b.Put([]byte(id), objJSON)
		glog.Infof("%s id %s created", kind, id)
		return nil
	})
	if err != nil {
		glog.Error(err)
	}
}

// GetAll ...
func (p *BoltPersistence) GetAll(kind string) []interface{} {
	db := p.driver.(*bolt.DB)
	list := arraylist.New()
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		b.ForEach(func(k, v []byte) error {
			glog.Infof("%s: %s", k, v)
			var obj domain.Service
			json.Unmarshal(v, obj)
			glog.Infof("%s", obj.Name)
			list.Add(obj)
			return nil
		})
		glog.Infof("%s id %s created", kind, list)
		return nil
	})
	if err != nil {
		glog.Error(err)
	}
	return list.Values()
}

// Close ...
func (p *BoltPersistence) Close() {
	db := p.driver.(*bolt.DB)
	db.Close()
}

// CreateBoltDriver ...
func CreateBoltDriver() Persistence {
	db, err := bolt.Open("bolt.db", 0600, nil)
	if err != nil {
		glog.Error(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("services"))
		if err != nil {
			glog.Error(err)
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("operations"))
		if err != nil {
			glog.Error(err)
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte("endpoints"))
		if err != nil {
			glog.Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		glog.Error(err)
	}
	return &BoltPersistence{driver: db}
}
