package persistence

import (
	"os"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/seadiaz/katalog/src/utils"

	"github.com/emirpasic/gods/lists/arraylist"
)

// BoltDriverInterface ...
type BoltDriverInterface interface {
	Open(path string, mode os.FileMode, options interface{}) (boltDBInterface, error)
}

type boltDBInterface interface {
	Update(fn func(*bolt.Tx) error) error
}

// BoltPersistence ...
type BoltPersistence struct {
	driver interface{}
}

// Create ...
func (p *BoltPersistence) Create(kind string, id string, obj interface{}) {
	db := p.driver.(*bolt.DB)
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		objJSON := utils.Serialize(obj)
		b.Put([]byte(id), []byte(objJSON))
		return nil
	})
	if err != nil {
		glog.Error(err)
	}
}

// GetAll ...
func (p *BoltPersistence) GetAll(kind string) []interface{} {
	glog.Info("get all called")
	db := p.driver.(*bolt.DB)
	list := arraylist.New()
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(kind))
		glog.Info("get all called")
		b.ForEach(func(k, v []byte) error {
			glog.Info("get all called")
			// var obj map[string]interface{}
			// json.Unmarshal(v, &obj)
			list.Add(v)
			return nil
		})
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
func CreateBoltDriver(db *bolt.DB) Persistence {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("services"))
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
