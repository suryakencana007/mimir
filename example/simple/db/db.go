/*  db.go
*
* @Author:             Nanang Suryadi
* @Date:               April 05, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 05/04/20 09:29
 */

package db

import (
	"fmt"
	"strings"

	"github.com/suryakencana007/mimir"
	"github.com/suryakencana007/mimir/example/simple/config"
	"github.com/suryakencana007/tyr"
)

var sqlOpen = tyr.New

func PostgresDBConn(logger mimir.Logging, c *config.Config) (*tyr.DB, func(), error) {
	logger.Info("connecting to postgres database")

	connStr := strings.Join(
		[]string{
			c.DB.DsnMain,
			fmt.Sprintf(
				"application_name=%s",
				c.App.Name),
		}, " ")

	master := tyr.SqlConn{
		Driver:     tyr.POSTGRES,
		ConnStr:    connStr,
		RetryCount: c.CB.Retry,
		Timeout:    c.CB.Timeout,
		Concurrent: c.CB.Concurrent,
	}

	db, err := sqlOpen(master, master)

	cleanup := func() {
		_ = db.Close()
	}
	if err != nil {
		return nil, nil, err
	}

	db.Master.SetMaxOpenConns(c.DB.MaxOpenConnection)
	db.Master.SetMaxIdleConns(c.DB.MaxIdleConnection)

	db.Slave.SetMaxOpenConns(c.DB.MaxOpenConnection)
	db.Slave.SetMaxIdleConns(c.DB.MaxIdleConnection)

	if err = db.Master.Ping(); err != nil {
		return nil, cleanup, err
	}

	if err = db.Slave.Ping(); err != nil {
		return nil, cleanup, err
	}

	return db, cleanup, nil
}
