package db

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

type DB struct {
	Cnct *pgx.Conn
}

func NewDataBase(DataBaseUrl string) (*DB, error) {
	logrus.Info("Connection to the database takes place . . .")
	connection, err := pgx.Connect(context.Background(), DataBaseUrl)
	if err != nil {
		logrus.Errorf("Failed to connect to Database: %v", err)
		return nil, err
	}
	logrus.Info("Starts database migration . . .")
	if err = runMigration(DataBaseUrl); err != nil {
		connection.Close(context.Background())
		logrus.Errorf("Fialed to run migration: %v . . .", err)
		return nil, err
	}
	logrus.Info("Database initialized successfully . . .")
	return &DB{Cnct: connection}, nil
}

func runMigration(DatabaseUrl string) error {
	pathForMigration, err := filepath.Abs("internal/db/migrations/")
	if err != nil {
		logrus.Errorf("Failed to determine absolute path for migrations: %v", err)
		return err
	}
	logrus.Infof("Using migration path: %s", pathForMigration)
	driver, err := (&file.File{}).Open(pathForMigration)
	if err != nil {
		logrus.Errorf("Failed to initialize migration source: %v", err)
		return err
	}
	migr, err := migrate.NewWithSourceInstance("file", driver, DatabaseUrl)
	if err != nil {
		logrus.Errorf("Database migration initialization error: %v", err)
		return err
	}
	if err = migr.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Errorf("Database migration error: %v . . .", err)
		return err
	}
	defer migr.Close()
	logrus.Info("Migrations applied successfully . . .")
	return nil
}

func (db *DB) CloseDatabase() {
	if err := db.Cnct.Close(context.Background()); err != nil {
		logrus.Errorf("Error closing database connection: %v", err)
	}
	logrus.Info("Database connection closed")
}
