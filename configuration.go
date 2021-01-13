package main

import (
	"fmt"

	"github.com/kpango/glg"
)

// Redical is the global configuration struct to encapsulate all global parameters
type Redical struct {
	redisDB   *RedisDB
	supported *CommandList
}

// SwitchDB switches the underlying RedisDB of RedicalConf to use the new DB number
func (r *Redical) SwitchDB(db int) error {
	return r.modifyConfig(&DBConfig{database: db})
}

// Authenticate authenticates the underlying RedisDB of RedicalConf with the password
func (r *Redical) Authenticate(pass string) error {
	return r.modifyConfig(&DBConfig{password: pass})
}

// SwitchRedis changes the underlying RedisDB of RedicalConf to a new Redis instance [Experimental].
func (r *Redical) SwitchRedis(host string, port int, db int, user string, pass string) error {
	return r.modifyConfig(&DBConfig{
		host:     host,
		port:     port,
		database: db,
		username: user,
		password: pass,
	})
}

// modifyConfig modifies the DBConfig for redis and refreshes the global redis client.
func (r *Redical) modifyConfig(mod *DBConfig) error {
	tmp := r.redisDB
	r.redisDB.Merge(mod)
	if err := r.redisDB.InitializeRedis(); err != nil {
		r.redisDB = tmp
		return err
	}
	glg.Info("Redis client re-initialized with modified config %v", mod)
	return nil
}

// PromptPrefix returns the prefix to be displayed in the prompt for this RedicalConf
func (r *Redical) PromptPrefix() string {
	var serv string
	if r.redisDB.redisConn == nil {
		serv = "NA"
	} else {
		serv = fmt.Sprintf("%s:%d/%d", r.redisDB.host, r.redisDB.port, r.redisDB.database)
	}
	return fmt.Sprintf("[%s] >>> ", serv)
}

// Close closes the RedisConnection associated with this RedicalConf
func (r *Redical) Close() {
	r.redisDB.TearDownRedis()
}

// NewRedicalConf creates an instance of RedicalConf that is being used across the system
func NewRedicalConf() (*Redical, error) {
	db := ParseConfig()
	supp, err := InitCmds()
	if err != nil {
		return nil, err
	}
	return &Redical{
		redisDB:   &RedisDB{db, nil},
		supported: supp,
	}, nil
}
