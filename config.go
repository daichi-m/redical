package main

import "github.com/gomodule/redigo/redis"

// RedicalConf is the global configuration struct to encapsulate all global parameters
type RedicalConf struct {
	config    DBConfig
	supported CommandList
	redis     *redis.Conn
}

// ModifyConfig modifies the DBConfig for redis and refreshes the global redis client.
func (rc *RedicalConf) ModifyConfig(mod *DBConfig) error {
	tmp := rc.config
	rc.config.Merge(mod)
	r, err := rc.config.InitializeRedis()
	if err != nil {
		rc.config = tmp
		return err
	}
	rc.redis = &r
	return nil
}
