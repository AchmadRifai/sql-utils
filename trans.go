package sqlutils

import (
	"database/sql"
	"log"
	"runtime/debug"
	"sync"
)

type Transaction struct {
	sync.Mutex
	db        *sql.DB
	Singleton bool
}

func (t *Transaction) Execute(flow func(db *sql.DB)) {
	if t.Singleton {
		t.Lock()
	}
	tx, err := t.db.Begin()
	defer t.txDone(tx)
	if err != nil {
		panic(err)
	}
	flow(t.db)
}

func (t *Transaction) txDone(tx *sql.Tx) {
	defer func() {
		if t.Singleton {
			t.Unlock()
		}
	}()
	if r := recover(); r != nil {
		log.Println("Error catched", r)
		log.Println("stack trace", string(debug.Stack()))
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
	} else if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func NewTransaction(db *sql.DB, singleton bool) *Transaction {
	return &Transaction{db: db, Singleton: singleton}
}

func NormalError() {
	if r := recover(); r != nil {
		log.Println("Error catched", r)
		log.Println("stack trace", string(debug.Stack()))
	}
}
