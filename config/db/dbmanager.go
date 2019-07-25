package db

import (

	"github.com/yamamushi/durouter/config"
	"gopkg.in/mgo.v2"
	"time"
)

type DBManager struct {

	MongoDB     string
	MongoHost   string
	MongoUser   string
	MongoPass   string

}

func NewDBManager(config config.Config) (*DBManager) {
	newManager := &DBManager{}
	newManager.MongoDB =  config.DBConfig.MongoDB
	newManager.MongoHost = config.DBConfig.MongoHost
	newManager.MongoPass = config.DBConfig.MongoPass
	newManager.MongoUser = config.DBConfig.MongoUser
	return newManager
}

func (h *DBManager) GetCollection(collection string) (c *mgo.Collection, err error) {

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{h.MongoHost},
		Timeout:  30 * time.Second,
		Database: h.MongoDB,
		Username: h.MongoUser,
		Password: h.MongoPass,
	}

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c = session.DB(h.MongoDB).C(collection)

	return c, nil
}
