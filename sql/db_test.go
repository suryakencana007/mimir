/*  db_test.go
*
* @Author:             Nanang Suryadi
* @Date:               July 25, 2019
* @Last Modified by:   @suryakencana007
* @Last Modified time: 2019-07-25 02:52
 */

package sql

import (
    "context"
    "fmt"
    "testing"
    "time"

    "github.com/lib/pq"
    "github.com/ory/dockertest"
    "github.com/pkg/errors"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/suryakencana007/mimir/constant"
)

const (
    pgSchema = `CREATE TABLE IF NOT EXISTS users (
id integer NOT NULL, 
name varchar(255) NOT NULL
);`
)

type ConnPGSuite struct {
    suite.Suite
    DB       *DB
    pool     *dockertest.Pool
    Resource *dockertest.Resource
}

func (s *ConnPGSuite) GetResource() *dockertest.Resource {
    return s.Resource
}

func (s *ConnPGSuite) SetResource(resource *dockertest.Resource) {
    s.Resource = resource
}

func (s *ConnPGSuite) GetPool() *dockertest.Pool {
    return s.pool
}

func (s *ConnPGSuite) SetPool(pool *dockertest.Pool) {
    s.pool = pool
}

func (s *ConnPGSuite) GetDB() *DB {
    return s.DB
}

func (s *ConnPGSuite) SetDB(db *DB) {
    s.DB = db
}

func (s *ConnPGSuite) SetupTest() {
    var err error
    s.pool, err = dockertest.NewPool("")
    if err != nil {
        panic(fmt.Sprintf("could not connect to docker: %s\n", err))
    }
    err = NewPoolPG(s)
    if err != nil {
        panic(fmt.Sprintf("prepare pg with docker: %v\n", err))
    }
}

func (s *ConnPGSuite) TearDownTest() {
    if err := s.DB.Close(); err != nil {
        panic(fmt.Sprintf("could not db close: %v\n", err))
    }
    if err := s.pool.RemoveContainerByName("pg_test"); err != nil {
        panic(fmt.Sprintf("could not remove postgres container: %v\n", err))
    }
}

func (s *ConnPGSuite) TestMainCommitInFailedTransaction() {
    t := s.T()
    txn, err := s.DB.Begin()
    assert.NoError(t, err)
    rows, err := txn.Query("SELECT error")
    assert.Error(t, err)
    if err == nil {
        rows.Close()
        t.Fatal("expected failure")
    }
    err = txn.Commit()
    assert.Error(t, err)
    if err != pq.ErrInFailedTransaction {
        t.Fatalf("expected ErrInFailedTransaction; got %#v", err)
    }
}

func (s *ConnPGSuite) TestGetSchema() {
    t := s.T()
    tx, err := s.DB.Begin()
    assert.NoError(t, err)

    rows, err := tx.Query("SELECT id FROM users")
    assert.NoError(t, err)

    defer rows.Close()
    names := make([]string, 0)
    for rows.Next() {
        var name string
        if err = rows.Scan(&name); err != nil {
            t.Fatal(err)
        }
        names = append(names, name)
    }
    assert.NoError(t, rows.Err())
    assert.IsType(t, []string{}, names)
    err = tx.Commit()
    assert.NoError(t, err)
}

func (s *ConnPGSuite) TestTxFunc() {
    t := s.T()
    txArgs := make([]ArgsTx, 0)
    args := []interface{}{
        1001,
        "TEST Tx Func",
    }
    query := `INSERT INTO users (id, name) VALUES ($1, $2)`
    txArgs = append(txArgs, ArgsTx{
        Query: query,
        Args:  args,
    })

    err := s.DB.TxExecContextMany(txArgs)
    assert.NoError(t, err)
}

func (s *ConnPGSuite) TestTxFuncFail() {
    t := s.T()
    txArgs := make([]ArgsTx, 0)
    args := []interface{}{
        1001,
        "TEST Tx Func",
    }
    query := `INSERT INTO users (id, name) VALUES (?, $2)`
    txArgs = append(txArgs, ArgsTx{
        Query: query,
        Args:  args,
    })

    err := s.DB.TxExecContextMany(txArgs)
    assert.Error(t, err)
}

func (s *ConnPGSuite) TestMainDB() {
    assert.IsType(s.T(), &DB{}, s.GetDB())
    assert.Equal(s.T(), 100008, getServerVersion(s.T(), s.GetDB()))
}

func TestMainPGSuite(t *testing.T) {
    suite.Run(t, new(ConnPGSuite))
}

type ConnectionSuite interface {
    T() *testing.T
    GetResource() *dockertest.Resource
    SetResource(resource *dockertest.Resource)
    GetPool() *dockertest.Pool
    SetPool(pool *dockertest.Pool)
    GetDB() *DB
    SetDB(factory *DB)
}

func NewPoolPG(c ConnectionSuite) (err error) {
    t := c.T()
    resource, err := c.GetPool().RunWithOptions(
        &dockertest.RunOptions{
            Name:       "pg_test",
            Repository: "postgres",
            Tag:        "10-alpine",
            Env: []string{
                "POSTGRES_PASSWORD=root",
                "POSTGRES_USER=root",
                "POSTGRES_DB=dev",
            },
        })
    c.SetResource(resource)
    if err != nil {
        return errors.Wrap(err, "start postgres")
    }
    err = c.GetResource().Expire(5)
    assert.NoError(t, err)
    purge := func() error {
        return c.GetPool().Purge(c.GetResource())
    }

    if err := c.GetPool().Retry(func() error {
        connInfo := fmt.Sprintf(`postgresql://%s:%s@%s:%s/%s?sslmode=disable`,
            "root",
            "root",
            c.GetResource().GetBoundIP("5432/tcp"),
            c.GetResource().GetPort("5432/tcp"),
            "dev",
        )
        db, err := Open(
            constant.POSTGRES,
            connInfo,
            3,
            5,
            500,
        )
        if err != nil {
            panic(err.Error())
        }
        c.SetDB(db)
        return c.GetDB().Ping()
    }); err != nil {
        _ = purge()
        return errors.Wrap(err, "check connection")
    }
    if _, err := c.GetDB().Exec(pgSchema); err != nil {
        _ = purge()
        return errors.Wrap(err, "failed to create schema")
    }

    return nil
}

func getServerVersion(t *testing.T, db *DB) int {
    var (
        version int
    )
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    row := db.QueryRowContext(ctx, `SHOW server_version_num;`)
    err := row.Scan(&version)
    if err != nil {
        t.Log(err)
    }
    return version
}
