/*  config.go
*
* @Author:             Nanang Suryadi
* @Date:               March 31, 2020
* @Last Modified by:   @suryakencana007
* @Last Modified time: 31/03/20 07:11
 */

package mimir

import (
	"strings"

	"github.com/spf13/viper"
)

type (
	ConfigConstants interface{}
	ConfigFunc      func(*viper.Viper) error
	ConfigOpts      struct {
		Config   ConfigConstants
		Filename string
		Paths    []string
	}
)

func Config(opts ConfigOpts, configFunc ConfigFunc) error {
	v := viper.New()
	// Search the root directory for the configuration file
	for _, path := range opts.Paths {
		v.AddConfigPath(path)
	}

	v.SetConfigName(opts.Filename) // Configuration fileName without the .TOML or .YAML extension

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&opts.Config); err != nil {
		return err
	}

	v.SetEnvPrefix("env")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	if err := configFunc(v); err != nil {
		return err
	}

	if err := v.MergeInConfig(); err != nil {
		return err
	}

	if err := v.UnmarshalExact(&opts.Config); err != nil {
		return err
	}

	for _, path := range opts.Paths {
		v.AddConfigPath(path)
	}
	v.SetConfigName(".env")
	v.SetConfigType("env")

	if err := v.MergeInConfig(); err == nil {
		if err := v.UnmarshalExact(&opts.Config); err != nil {
			return err
		}
	}

	return nil
}
