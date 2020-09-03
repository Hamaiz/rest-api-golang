package session

import (
	"log"
	"os"

	"github.com/globalsign/mgo"
	"github.com/kidstuff/mongostore"
)

type AccountStore struct {
	store *mongostore.MongoStore
}

// DBConn - connects to mongodb
func DBConn() (*mgo.Session, error) {
	log.Println("connecting to mongodb...")
	dbsess, err := mgo.Dial(os.Getenv("SESSION_DB"))

	if err != nil {
		return nil, err
	}

	log.Println("connected to mongodb")
	return dbsess, nil
}

// StoreConn - connects store
func StoreConn(dbsess *mgo.Session) *AccountStore {
	log.Println("setting up session store...")

	conn := dbsess.DB("").C("sessions")
	url := []byte(os.Getenv("SESSION_KEY"))

	store := mongostore.NewMongoStore(conn, 3600, true, url)

	return &AccountStore{store}
}
