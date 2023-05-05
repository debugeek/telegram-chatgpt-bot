package main

import "log"

type MemCache struct{}

func (c *MemCache) InitDatabase() {
	log.Println(`Memory cache initialized`)
}
func (c *MemCache) GetAccounts() ([]*Account, error) {
	return nil, nil
}
func (c *MemCache) GetAccount(id int64) (*Account, error) {
	return nil, nil
}
func (c *MemCache) SaveAccount(acc *Account) error {
	return nil
}
func (c *MemCache) GetMessage(acc *Account, msgId int) (*Message, error) {
	return nil, nil
}
func (c *MemCache) SaveMessage(acc *Account, msg *Message) error {
	return nil
}
