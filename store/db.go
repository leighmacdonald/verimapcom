package store

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
	"strconv"
	"strings"
)

func MustConnectDB(ctx context.Context) *pgxpool.Pool {
	dsn := viper.GetString("dsn")
	conn, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return conn
}

func dsnToConnConfig(dsn string) (pgx.ConnConfig, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return pgx.ConnConfig{}, err
	}
	h := u.Host
	prt := uint16(5432)
	if u.Port() != "" {
		port64, err := strconv.ParseUint(u.Port(), 10, 16)
		if err != nil {
			return pgx.ConnConfig{}, err
		}
		prt = uint16(port64)
	}
	d := strings.ReplaceAll(u.Path, "/", "")
	usr := u.User.Username()
	p, _ := u.User.Password()

	return pgx.ConnConfig{
		Host:     h,
		Port:     prt,
		Database: d,
		User:     usr,
		Password: p,
	}, nil
}

func Migrate(dsn string) error {
	log.Infof("Performing migration")
	cfg, err := dsnToConnConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to do migration: %v", err)
	}
	db := stdlib.OpenDB(cfg)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://./store/schema", "postgres", driver)
	if err != nil {
		return errors.Wrapf(err, "failed to create new migrate instance")
	}
	//if err := m.Steps(2); err != nil {
	//	return errors.Wrapf(err, "failed to .Steps migration")
	//}
	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			return nil
		}
		return errors.Wrapf(err, "Failed to .Up migration")
	}

	log.Infof("Migrations complete")
	return nil
}
