package main

import "log"

type Account struct {
	Id     int64  `firestore:"id"`
	Kind   int    `firestore:"kind"`
	Status int    `firestore:"status"`
	Key    string `firestore:"key"`
	Model  string `firestore:"model"`
}

type DatabaseProtocol interface {
	InitDatabase()
	GetAccounts() ([]*Account, error)
	GetAccount(id int64) (*Account, error)
	SaveAccount(account *Account) error
}

func InitDatabase() {
	dbOnce.Do(func() {
		if len(args.FirebaseConf) != 0 || len(args.FirebaseConfEnvKey) != 0 {
			db = &Firebase{}
		} else {
			db = &MemoryDatabase{}
		}
		db.InitDatabase()
	})
}

type MemoryDatabase struct{}

func (db *MemoryDatabase) InitDatabase() {
	log.Println(`Memory cache initialized`)
}
func (db *MemoryDatabase) GetAccounts() ([]*Account, error) {
	return nil, nil
}
func (db *MemoryDatabase) GetAccount(id int64) (*Account, error) {
	return nil, nil
}
func (db *MemoryDatabase) SaveAccount(account *Account) error {
	return nil
}
