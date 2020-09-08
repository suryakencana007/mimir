/*  db_test.go
*
* @Author:             Nanang Suryadi
* @Date:               April 05, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 05/04/20 21:39
 */

package db

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/suryakencana007/mimir"
	"github.com/suryakencana007/mimir/example/simple/config"
	"github.com/suryakencana007/tyr"
)

func TestPostgresDBConn(t *testing.T) {
	logger := mimir.With(mimir.Field("Headless", "listen and serve"))
	masterDB, mockMaster, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		assert.Failf(t, "failed to open stub db", "%v", err)
	}
	defer masterDB.Close()

	slaveDB, mockSlave, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		assert.Failf(t, "failed to open stub db", "%v", err)
	}
	defer slaveDB.Close()

	sqlOpen = func(master, slave tyr.SqlConn) (*tyr.DB, error) {
		logger.Info("SQL Open")
		return &tyr.DB{
			Master: masterDB,
			Slave:  slaveDB,
		}, nil
	}

	mockMaster.ExpectPing()
	mockSlave.ExpectPing()

	cfg := &config.Config{}
	cfg.DB.MaxOpenConnection = 1
	cfg.DB.MaxIdleConnection = 1
	db, cleanup, err := PostgresDBConn(logger, cfg)

	// Asserts
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.NotNil(t, cleanup)

	mockMaster.ExpectClose()
	mockSlave.ExpectClose()
	cleanup()

	if err := mockMaster.ExpectationsWereMet(); err != nil {
		assert.Failf(t, "there were unfulfilled expectations", "%v", err)
	}

	if err := mockSlave.ExpectationsWereMet(); err != nil {
		assert.Failf(t, "there were unfulfilled expectations", "%v", err)
	}
}

func TestPostgresDBConnConnectFail(t *testing.T) {
	logger := mimir.With(mimir.Field("Headless", "listen and serve"))
	sqlOpen = func(master, slave tyr.SqlConn) (*tyr.DB, error) {
		return nil, fmt.Errorf("expected connection failure")
	}

	cfg := &config.Config{}
	cfg.DB.MaxOpenConnection = 1
	cfg.DB.MaxIdleConnection = 1
	db, cleanup, err := PostgresDBConn(logger, cfg)

	// Asserts
	assert.Nil(t, db)
	assert.Nil(t, cleanup)
	assert.Error(t, err)

}
