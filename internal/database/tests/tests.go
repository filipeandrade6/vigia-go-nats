package tests

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/filipeandrade6/vigia-go/internal/database/migration"
	"github.com/filipeandrade6/vigia-go/internal/sys/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"go.uber.org/zap"
)

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// TODO no services ele joga todo o log dos containers no final dos tests

func New(t *testing.T) (*zap.SugaredLogger, *sqlx.DB, func()) {
	var err error
	log, err := logger.New("TEST")
	if err != nil {
		t.Fatalf("logger error: %s", err)
	}

	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("postgres", "secret"),
		Path:   "vigia",
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("could not connect to docker: %v", err)
	}

	pw, _ := pgURL.User.Password()
	runOpts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=" + pgURL.User.Username(),
			"POSTGRES_PASSWORD=" + pw,
			"POSTGRES_DB=" + pgURL.Path,
		},
	}

	resource, err := pool.RunWithOptions(&runOpts)
	if err != nil {
		t.Fatalf("could not start postgres container: %v", err)
	}

	pgURL.Host = resource.Container.NetworkSettings.IPAddress

	resource.Expire(60)
	pool.MaxWait = 15 * time.Second
	var db *sqlx.DB
	err = pool.Retry(func() error {
		db, err = sqlx.Open("postgres", pgURL.String())
		if err != nil {
			return err
		}
		return db.Ping()
	})
	if err != nil {
		t.Fatalf("could not connect to postgres server: %v", err)
	}

	migration.Migrate(context.Background(), 0, pgURL.String())

	teardown := func() {
		if err := pool.Purge(resource); err != nil {
			t.Errorf("could not purge the resources: %v", err)
		}
	}

	return log, db, teardown
}

// TODO - está mostrando logs... nos tests - é certo?

func StringPointer(s string) *string {
	return &s
}

func IntPointer(i int) *int {
	return &i
}

func BoolPointer(b bool) *bool {
	return &b
}
