package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Firebase struct {
	app       *firebase.App
	firestore *firestore.Client
	ctx       context.Context
}

func (fb *Firebase) InitDatabase() {
	var conf []byte
	if len(args.FirebaseConf) != 0 {
		conf, _ = base64.StdEncoding.DecodeString(args.FirebaseConf)
	} else if len(args.FirebaseConfEnvKey) != 0 {
		conf, _ = base64.StdEncoding.DecodeString(os.Getenv(args.FirebaseConfEnvKey))
	} else {
		panic("firebase credential not found")
	}
	opt := option.WithCredentialsJSON(conf)

	fb.ctx = context.Background()

	if app, err := firebase.NewApp(fb.ctx, nil, opt); err != nil {
		panic(err)
	} else {
		fb.app = app
	}

	if firestore, err := fb.app.Firestore(fb.ctx); err != nil {
		panic(err)
	} else {
		fb.firestore = firestore
	}

	log.Println(`Firebase initialized`)
}

// Account

func (fb *Firebase) GetAccounts() ([]*Account, error) {
	accounts := make([]*Account, 0)

	iter := fb.firestore.Collection("accounts").Documents(fb.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var account Account
		doc.DataTo(&account)

		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func (fb *Firebase) GetAccount(id int64) (*Account, error) {
	iter := fb.firestore.Collection("accounts").Where("id", "==", id).Documents(fb.ctx)

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var account *Account
	err = doc.DataTo(&account)

	return account, err
}

func (fb *Firebase) SaveAccount(account *Account) error {
	id := strconv.FormatInt(account.Id, 10)

	_, err := fb.firestore.Collection("accounts").Doc(id).Set(fb.ctx, account)

	return err
}
