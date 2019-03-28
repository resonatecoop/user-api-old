package pgsql

import (
	"fmt"
	"log"
	"time"

	"github.com/go-pg/pg"
)

type dbLogger struct { }

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {
	fmt.Println(q.FormattedQuery())
}

// New creates new database connection to a postgres database
// Function panics if it can't connect to database
func New(psn string, logQueries bool, timeout int) (*pg.DB, error) {
	u, err := pg.ParseURL(psn)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(u)

	_, err = db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	if timeout > 0 {
		db.WithTimeout(time.Second * time.Duration(timeout))
	}

	if logQueries {
		db.AddQueryHook(dbLogger{})
	}

	return db, nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
