//go:build postgresqldb
// +build postgresqldb

// Initializes a PostgreSQL DB for testing purposes

package storage

import (
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/theupdateframework/notary"
)

func init() {
	// Get the PostgreSQL connection string from an environment variable
	dburl := os.Getenv("DBURL")
	if dburl == "" {
		logrus.Fatal("PostgreSQL environment variable not set")
	}

	for i := 0; i <= 30; i++ {
		gormDB, err := gorm.Open(notary.PostgresBackend, dburl)
		if err == nil {
			err := gormDB.DB().Ping()
			if err == nil {
				break
			}
		}
		if i == 30 {
			logrus.Fatalf("Unable to connect to %s after 60 seconds", dburl)
		}
		time.Sleep(2 * time.Second)
	}

	sqldbSetup = func(t *testing.T) (*SQLStorage, func()) {
		var dropTables = func(gormDB *gorm.DB) {
			// drop all tables, if they exist
			gormDB.DropTable(&TUFFile{})
			gormDB.DropTable(&SQLChange{})
		}
		gormDB, err := gorm.Open(notary.PostgresBackend, dburl)
		require.NoError(t, err)
		dropTables(gormDB)
		gormDB.Close()
		dbStore := SetupSQLDB(t, notary.PostgresBackend, dburl)
		return dbStore, func() {
			dropTables(&dbStore.DB)
			dbStore.Close()
		}
	}
}
