/*  copnfig_test.go
*
* @Author:             Nanang Suryadi
* @Date:               March 31, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 31/03/20 07:31
 */

package mimir

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var constants = struct {
	App struct {
		Name         string         `mapstructure:"name"`
		Port         int            `mapstructure:"port"`
		ReadTimeout  int            `mapstructure:"read_timeout"`
		WriteTimeout int            `mapstructure:"write_timeout"`
		Timezone     string         `mapstructure:"timezone"`
		Debug        bool           `mapstructure:"debug"`
		Env          string         `mapstructure:"env"`
		SecretKey    string         `mapstructure:"secret_key"`
		ExpireIn     *time.Duration `mapstructure:"expire_in"`
	}
	DB struct {
		DsnMain string `mapstructure:"dsn_main" toml:"dsn_main,omitempty"`
	}
}{}

func TestConfig(t *testing.T) {
	err := Config(ConfigOpts{
		Config:   ConfigConstants(&constants),
		Filename: "app.config.test",
		Paths:    []string{"."},
	}, func(viper *viper.Viper) error {
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 8778, constants.App.Port)
}

func TestConfigPathFail(t *testing.T) {
	err := Config(ConfigOpts{
		Config:   ConfigConstants(&constants),
		Filename: "app.config.test",
		Paths:    []string{"./configs"},
	}, func(viper *viper.Viper) error {
		return nil
	})

	assert.Error(t, err)
}

func TestConfigEnv(t *testing.T) {
	defer os.Clearenv()
	val := "host=localhost port=5999 user=root password=root123 dbname=dbroot sslmode=disable"

	err := Config(ConfigOpts{
		Config:   ConfigConstants(&constants),
		Filename: "app.config.test",
		Paths:    []string{"."},
	}, func(v *viper.Viper) error {
		return nil
	})

	t.Log(constants.App.Name)
	t.Log(constants.DB.DsnMain)
	assert.NoError(t, err)
	assert.Equal(t, val, constants.DB.DsnMain)
}

func TestConfigFunc(t *testing.T) {
	defer os.Clearenv()
	val := "host=localhost port=5999 user=root password=root123 dbname=dbroot sslmode=disable"
	_ = os.Setenv("ENV_DB_DSN_MAIN", val)

	err := Config(ConfigOpts{
		Config:   ConfigConstants(&constants),
		Filename: "app.config.test",
		Paths:    []string{"."},
	}, func(v *viper.Viper) error {
		return v.BindEnv("db.dsn_main")
	})

	assert.NoError(t, err)
	assert.Equal(t, val, constants.DB.DsnMain)
}

func TestFailConfigFunc(t *testing.T) {
	defer os.Clearenv()
	val := "host=localhost port=5999 user=root password=root123 dbname=dbroot sslmode=disable"
	_ = os.Setenv("ENV_DB_DSN_SECOND", val)

	err := Config(ConfigOpts{
		Config:   ConfigConstants(&constants),
		Filename: "app.config.test",
		Paths:    []string{"."},
	}, func(v *viper.Viper) error {
		return errors.New("bind env not found")
	})

	assert.Error(t, err)
}
