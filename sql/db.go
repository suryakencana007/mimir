/*  db.go
*
* @Author:             Nanang Suryadi
* @Date:               February 12, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-02-12 11:20
 */

package sql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/suryakencana007/mimir/breaker"
	"github.com/suryakencana007/mimir/log"
)

type ArgsTx struct {
	Query string
	Args  []interface{}
}

// DBFactory is an abstract for sql database
type DBFactory interface {
	OpenConnection(connString string, retry, timeout, concurrent int)
	Close() error
	GetDB() (*DB, error)
	QueryRow(query string, args ...interface{}) (*sql.Row, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	TxExecMany(args []ArgsTx) error
	Prepare(query string) (*sql.Stmt, error)
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	SetCommandBreaker(commandName string, timeout, maxConcurrent int, args ...interface{}) *DB
	SetConnMaxLifetime(int)
	GetQueryTimeout() int
	SetMaxIdleConn(int)
	GetDefaultMaxConcurrent() int
	SetMaxOpenConn(int)
}

type fallbackFunc func(error) error

type DB struct {
	*sql.DB
	cb           *breaker.CircuitBreaker
	retryCount   int
	timeout      int
	concurrent   int
	fallbackFunc func(error) error
}

func Open(driverName, connString string, retry, timeout, concurrent int) (*DB, error) {
	db, err := sql.Open(driverName, connString)
	if err != nil {
		panic(err.Error())
	}
	return &DB{
		DB:         db,
		retryCount: retry,
		timeout:    timeout,
		concurrent: concurrent,
	}, nil
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (r *DB) SetConnMaxLifetime(connMaxLifetime int) {
	r.DB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
func (r *DB) SetMaxIdleConn(maxIdleConn int) {
	r.DB.SetMaxIdleConns(maxIdleConn)
}

// SetMaxOpenConns sets the maximum amount of time a connection may be reused.
func (r *DB) SetMaxOpenConn(maxOpenConn int) {
	r.DB.SetMaxOpenConns(maxOpenConn)
}

func (r *DB) Close() error {
	return r.DB.Close()
}

func (r *DB) Begin() (*sql.Tx, error) {
	return r.DB.Begin()
}

func (r *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return r.DB.Exec(query, args...)
}

// QueryRows the fetch data rows
func (r *DB) Query(query string, args ...interface{}) (rs *sql.Rows, err error) {
	if err = r.callBreaker(func() error {
		if r.DB == nil {
			err = errors.New("the database connection is nil")
			log.Error(err.Error(), log.Field("query", query), log.Field("args", args))

			return err
		}
		if rs, err = r.DB.Query(query, args...); err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Error(err.Error(), log.Field("query", query), log.Field("args", args))
	}
	return rs, err
}

// QueryRow the fetch data row
func (r *DB) QueryRow(query string, args ...interface{}) (rs *sql.Row, err error) {
	if err = r.callBreaker(func() (err error) {
		if r.DB == nil {
			err = errors.New("the database connection is nil")
			log.Error(err.Error(), log.Field("query", query), log.Field("args", args))
			return err
		}
		rs = r.DB.QueryRow(query, args...)
		return nil
	}); err != nil {
		log.Error(err.Error(), log.Field("query", query), log.Field("args", args))
	}
	return rs, err
}

// Sql Transaction Tx Exec many
func (r *DB) TxExecMany(args []ArgsTx) error {
	return r.callBreaker(func() (err error) {
		var tx *sql.Tx
		tx, err = r.DB.Begin()
		if err != nil {
			return err
		}
		for _, arg := range args {
			var stmt *sql.Stmt
			if stmt, err = tx.Prepare(arg.Query); err != nil {
				log.Error("TxExecMany:",
					log.Field("error", err.Error()),
					log.Field("query", arg.Query),
					log.Field("args", arg.Args),
				)
				return err
			}
			var result sql.Result
			result, err = stmt.Exec(arg.Args...)
			if err != nil {
				log.Error("TxExecMany:",
					log.Field("error", err.Error()),
					log.Field("query", arg.Query),
					log.Field("args", arg.Args),
				)
				return err
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				log.Error("TxExecMany:",
					log.Field("error", err.Error()),
					log.Field("query", arg.Query),
					log.Field("args", arg.Args),
				)
				return err
			}
			log.Info("Commit Transaction TxExecMany", log.Field("RowsAffected", rowsAffected))
			err = stmt.Close()
			if err != nil {
				return err
			}
		}
		log.Info("Commit Transaction TxExecMany")
		// commit db transaction
		if err := tx.Commit(); err != nil {
			if err = tx.Rollback(); err != nil {
				log.Error("TxExecMany:",
					log.Field("error", err.Error()),
				)
				return err
			} // rollback if fail query statement
			return err
		}
		return nil
	})
}

// SetCommandBreaker the circuit breaker
func (r *DB) SetCommandBreaker(commandName string, timeout, maxConcurrent int, args ...interface{}) *DB {
	r.cb = breaker.NewBreaker(
		commandName,
		timeout,
		maxConcurrent,
		args...)
	return r
}

// callBreaker command circuit breaker
func (r *DB) callBreaker(fn func() error) (err error) {
	return r.cb.Execute(fn)
}

// GetQueryTimeout for circuit breaker
func (r *DB) GetQueryTimeout() int {
	if timeout := r.timeout; timeout > 1 {
		return timeout
	}
	return 1000
}

// GetDefaultMaxConcurrent circuit breaker
func (r *DB) GetDefaultMaxConcurrent() int {
	if concurrent := r.concurrent; concurrent > 1 {
		return concurrent
	}
	return 100
}
