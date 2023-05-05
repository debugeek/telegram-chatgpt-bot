package main

type Account struct {
	Id     int64  `firestore:"id"`
	Kind   int    `firestore:"kind"`
	Status int    `firestore:"status"`
	Key    string `firestore:"key"`
	Model  string `firestore:"model"`
}

type Message struct {
	Id        int    `firestore:"id"`
	ParentId  int    `firestore:"parentId"`
	Text      string `firestore:"text"`
	Role      string `firestore:"role"`
	Timestamp int64  `firestore:"timestamp"`
}

type DatabaseProtocol interface {
	InitDatabase()
	GetAccounts() ([]*Account, error)
	GetAccount(id int64) (*Account, error)
	SaveAccount(account *Account) error
	GetMessage(account *Account, messageId int) (*Message, error)
	SaveMessage(account *Account, message *Message) error
}

func InitDatabase() {
	dbOnce.Do(func() {
		if len(args.FirebaseConf) != 0 || len(args.FirebaseConfEnvKey) != 0 {
			db = &Firebase{}
		} else {
			db = &MemCache{}
		}
		db.InitDatabase()
	})
}
