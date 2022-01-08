// Package database provides support for access the database.
package database

import (
	"context"
	"errors"
	"net/url"
	"reflect"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Calls init function.
	"go.uber.org/zap"
)

var (
	ErrNotFound              = errors.New("not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrForbidden             = errors.New("attempted action is not allowed")

	ErrDBNotFound = errors.New("not found")
)

type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxIDLEConns int
	MaxOpenConns int
	DisableTLS   bool
	// SSLMode string
}

func Connect(cfg Config) (*sqlx.DB, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Connect("postgres", u.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIDLEConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

// NamedExecContext is a helper function to execute a CUD operation with
// logging and tracing.
func NamedExecContext(ctx context.Context, log *zap.SugaredLogger, sqlxDB *sqlx.DB, query string, data interface{}) error {
	// q := queryString(query, data)
	// log.Infow("database.NamedExecContext", "traceid", web.GetTraceID(ctx), "query", q)

	if _, err := sqlxDB.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	return nil
}

// NamedQuerySlice is a helper function for executing queries that return a
// collection of data to be unmarshaled into a slice.
func NamedQuerySlice(ctx context.Context, log *zap.SugaredLogger, sqlxDB *sqlx.DB, query string, data interface{}, dest interface{}) error {
	// q := queryString(query, data)
	// log.Infow("database.NamedQuerySlice", "traceid", web.GetTraceID(ctx), "query", q)

	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("must provide a pointer to a slice")
	}

	rows, err := sqlxDB.NamedQueryContext(ctx, query, data)
	if err != nil {
		return err
	}

	slice := val.Elem()
	for rows.Next() {
		v := reflect.New(slice.Type().Elem())
		if err := rows.StructScan(v.Interface()); err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, v.Elem()))
	}

	return nil
}

// NamedQueryStruct is a helper function for executing queries that return a
// single value to be unmarshalled into a struct type.
func NamedQueryStruct(ctx context.Context, log *zap.SugaredLogger, sqlxDB *sqlx.DB, query string, data interface{}, dest interface{}) error {
	// q := queryString(query, data)
	// log.Infow("database.NamedQueryStruct", "traceid", web.GetTraceID(ctx), "query", q)

	rows, err := sqlxDB.NamedQueryContext(ctx, query, data)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return ErrDBNotFound
	}

	if err := rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}

// queryString provides a pretty print version of the query and parameters.
// func queryString(query string, args ...interface{}) string {
// 	query, params, err := sqlx.Named(query, args)
// 	if err != nil {
// 		return err.Error()
// 	}

// 	for _, param := range params {
// 		var value string
// 		switch v := param.(type) {
// 		case string:
// 			value = fmt.Sprintf("%q", v)
// 		case []byte:
// 			value = fmt.Sprintf("%q", string(v))
// 		default:
// 			value = fmt.Sprintf("%v", v)
// 		}
// 		query = strings.Replace(query, "?", value, 1)
// 	}

// 	query = strings.Replace(query, "\t", "", -1)
// 	query = strings.Replace(query, "\n", " ", -1)

// 	return strings.Trim(query, " ")
// }

// QuerySlice is a helper function for executing queries that return a
// collection of data to be unmarshaled into a slice.
func QuerySlice(ctx context.Context, log *zap.SugaredLogger, sqlxDB *sqlx.DB, query string, dest interface{}) error {
	// q := queryString(query, data)
	// log.Infow("database.NamedQuerySlice", "traceid", web.GetTraceID(ctx), "query", q)

	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("must provide a pointer to a slice")
	}

	rows, err := sqlxDB.QueryxContext(ctx, query)
	if err != nil {
		return err
	}

	slice := val.Elem()
	for rows.Next() {
		v := reflect.New(slice.Type().Elem())
		if err := rows.StructScan(v.Interface()); err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, v.Elem()))
	}

	return nil
}