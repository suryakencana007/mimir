/*  config.go
*
* @Author:             Nanang Suryadi
* @Date:               April 04, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 04/04/20 22:54
 */

package config

import (
	"time"
)

type Config struct {
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
	CB struct {
		Retry      int `mapstructure:"retry_count"`
		Timeout    int `mapstructure:"db_timeout"`
		Concurrent int `mapstructure:"max_concurrent"`
	}
	DB struct {
		DsnMain           string `mapstructure:"dsn_main" toml:"dsn_main,omitempty"`
		MaxLifeTime       int    `mapstructure:"max_life_time"`
		MaxIdleConnection int    `mapstructure:"max_idle_connection"`
		MaxOpenConnection int    `mapstructure:"max_open_connection"`
	}
}
