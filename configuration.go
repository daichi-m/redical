package main

import (
	"fmt"

	"github.com/kpango/glg"
)

// RedicalConf is the global configuration struct to encapsulate all global parameters
type RedicalConf struct {
	redisDB   *RedisDB
	supported *CommandList
}

// SwitchDB switches the underlying RedisDB of RedicalConf to use the new DB number
func (rc *RedicalConf) SwitchDB(db int) error {
	return rc.modifyConfig(&DBConfig{database: db})
}

// Authenticate authenticates the underlying RedisDB of RedicalConf with the password
func (rc *RedicalConf) Authenticate(pass string) error {
	return rc.modifyConfig(&DBConfig{password: pass})
}

// SwitchRedis changes the underlying RedisDB of RedicalConf to a new Redis instance [Experimental].
func (rc *RedicalConf) SwitchRedis(host string, port int, db int, user string, pass string) error {
	return rc.modifyConfig(&DBConfig{
		host:     host,
		port:     port,
		database: db,
		username: user,
		password: pass,
	})
}

// modifyConfig modifies the DBConfig for redis and refreshes the global redis client.
func (rc *RedicalConf) modifyConfig(mod *DBConfig) error {
	tmp := rc.redisDB
	rc.redisDB.Merge(mod)
	if err := rc.redisDB.InitializeRedis(); err != nil {
		rc.redisDB = tmp
		return err
	}
	glg.Info("Redis client re-initialized with modified config %v", mod)
	return nil
}

// PromptPrefix returns the prefix to be displayed in the prompt for this RedicalConf
func (rc *RedicalConf) PromptPrefix() string {
	var serv string
	if rc.redisDB.redisConn == nil {
		serv = "NA"
	} else {
		serv = fmt.Sprintf("%s:%d/%d", rc.redisDB.host, rc.redisDB.port, rc.redisDB.database)
	}
	return fmt.Sprintf("[%s] >>> ", serv)
}

// Close closes the RedisConnection associated with this RedicalConf
func (rc *RedicalConf) Close() {
	rc.redisDB.TearDownRedis()
}

// NewRedicalConf creates an instance of RedicalConf that is being used across the system
func NewRedicalConf() (*RedicalConf, error) {
	db := ParseConfig()
	supp, err := InitCmds()
	if err != nil {
		return nil, err
	}
	return &RedicalConf{
		redisDB:   &RedisDB{db, nil},
		supported: supp,
	}, nil
}
