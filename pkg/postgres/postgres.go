package pgsql

import (
	"log"
	"time"

	"github.com/go-pg/pg"
	// DB adapter
	_ "github.com/lib/pq"
)

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
		db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
			query, _ := event.FormattedQuery()
			log.Printf("%s | %s", time.Since(event.StartTime), query)
		})
	}

	return db, nil
}
