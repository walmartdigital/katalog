package persistence

import (
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/mitchellh/mapstructure"
	"github.com/seadiaz/katalog/src/utils"

	"github.com/emirpasic/gods/lists/arraylist"
)

// BoltWrapper ...
type BoltWrapper struct {
	DB *bolt.DB
}

// Update ...
func (b BoltWrapper) Update(fn func(BoltTxInterface) error) error {
	b.DB.Update(func(tx *bolt.Tx) error {
		return fn(tx)
	})
	return nil
}

// View ...
func (b BoltWrapper) View(fn func(BoltTxInterface) error) error {
	b.DB.View(func(tx *bolt.Tx) error {
		return fn(tx)
	})
	return nil
}

// Close ...
func (b BoltWrapper) Close() error {
	b.DB.Close()
	return nil
}

// BoltDBInterface ...
type BoltDBInterface interface {
	Update(fn func(BoltTxInterface) error) error
	View(fn func(BoltTxInterface) error) error
	Close() error
}

// BoltTxInterface ..
type BoltTxInterface interface {
	CreateBucketIfNotExists([]byte) (*bolt.Bucket, error)
	Bucket(name []byte) *bolt.Bucket
}

// BoltPersistence ...
type BoltPersistence struct {
	driver BoltDBInterface
}

// CreateBoltDriver ...
func CreateBoltDriver(db BoltDBInterface) Persistence {
	err := db.Update(initCollections)
	if err != nil {
		glog.Error(err)
	}
	return &BoltPersistence{driver: db}
}

func initCollections(tx BoltTxInterface) error {
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
}

// Create ...
func (p *BoltPersistence) Create(kind string, id string, obj interface{}) {
	db := p.driver.(*BoltWrapper)
	err := db.Update(func(tx BoltTxInterface) error {
		b := tx.Bucket([]byte(kind))
		var generic map[string]interface{}
		mapstructure.Decode(obj, &generic)
		objJSON := utils.Serialize(generic)
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
	list := arraylist.New()
	err := p.driver.View(func(tx BoltTxInterface) error {
		b := tx.Bucket([]byte(kind))
		b.ForEach(func(k, v []byte) error {
			obj := utils.Deserialize(string(v))
			list.Add(obj)
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
	p.driver.Close()
}
