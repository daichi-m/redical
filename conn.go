package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	flag "github.com/spf13/pflag"
)

// DBConfig is the struct that encapsulates the user inputs
type DBConfig struct {
	host, user, password, client                     string
	port, db                                         int
	timeout, connectTO, readTO, writeTO, keepAliveTO time.Duration
	tls, skipVerifyTLS, debug                        bool
}

// ParseConfig parses the input flags and initializes the global Input struct
func ParseConfig() DBConfig {

	host := flag.StringP("host", "h", "localhost", "Hostname of the redis instance")
	port := flag.IntP("port", "p", 6379, "Port of the redis instance")
	db := flag.IntP("database", "d", 0, "Redis database to select")
	user := flag.StringP("username", "U", "", "Redis username in case Redis ACL's are enabled")
	password := flag.StringP("password", "P", "", "Redis password for auth command")

	defTO := 60 * time.Second
	timeout := flag.DurationP("timeout", "t", defTO, "Redis timeout")
	connectTO := flag.Duration("connectTO", defTO, "Connection timeout to redis")
	readTO := flag.Duration("readTO", defTO, "Read timeout from redis")
	writeTO := flag.Duration("writeTO", defTO, "Write timeout to redis")
	kaTO := flag.Duration("keepAlive", defTO, "Connection keep-alive timeout")
	tls := flag.BoolP("tls", "S", false, "Use TLS")
	skipTLS := flag.BoolP("skipVerifyTLS", "k", false, "Skip verifying TLS connection")
	clName := flag.StringP("client", "c", "", "Client name")

	debug := flag.Bool("debug", false, "Run in debug mode")

	flag.Parse()

	conf := DBConfig{
		host:          *host,
		port:          *port,
		db:            *db,
		user:          *user,
		password:      *password,
		timeout:       *timeout,
		connectTO:     *connectTO,
		readTO:        *readTO,
		writeTO:       *writeTO,
		keepAliveTO:   *kaTO,
		tls:           *tls,
		skipVerifyTLS: *skipTLS,
		client:        *clName,
		debug:         *debug,
	}
	return conf
}

/*
InitializeRedis initializes the redigo/redis client at startup based on the CLI
inputs in the DBConfig
*/
func (db *DBConfig) InitializeRedis() (redis.Conn, error) {
	dialOpts := db.createDialOpts()
	address := fmt.Sprintf("%s:%d", db.host, db.port)
	return redis.Dial("tcp", address, dialOpts...)
}

func (db *DBConfig) createDialOpts() []redis.DialOption {
	var dialOpts []redis.DialOption

	dialOpts = append(dialOpts, redis.DialDatabase(db.db))
	if len(db.password) > 0 {
		dialOpts = append(dialOpts, redis.DialPassword(db.password))
	}
	if len(db.user) > 0 {
		dialOpts = append(dialOpts, redis.DialUsername(db.user))
	}

	dialOpts = append(dialOpts, redis.DialConnectTimeout(db.connectTO))
	dialOpts = append(dialOpts, redis.DialReadTimeout(db.readTO))
	dialOpts = append(dialOpts, redis.DialWriteTimeout(db.writeTO))
	dialOpts = append(dialOpts, redis.DialKeepAlive(db.keepAliveTO))
	dialOpts = append(dialOpts, redis.DialUseTLS(db.tls))
	if db.skipVerifyTLS {
		dialOpts = append(dialOpts, redis.DialTLSSkipVerify(db.skipVerifyTLS))
	}
	if len(db.client) > 0 {
		dialOpts = append(dialOpts, redis.DialClientName(db.client))
	}
	return dialOpts
}

/*
Merge this instance of DBConfig with another DBConfig. This function
picks up all the non-zero values from other and assigns them to the
corresponding field in this DBConfig object
*/
func (db *DBConfig) Merge(other *DBConfig) {
	if other.host != "" {
		db.host = other.host
	}
	if other.port != 0 {
		db.port = other.port
	}
	if other.user != "" {
		db.user = other.user
	}
	if other.password != "" {
		db.password = other.password
	}
	if other.client != "" {
		db.client = other.client
	}
	if other.timeout != 0 {
		db.timeout = other.timeout
	}
	if other.readTO != 0 {
		db.readTO = other.readTO
	}
	if other.writeTO != 0 {
		db.writeTO = other.writeTO
	}
	if other.connectTO != 0 {
		db.connectTO = other.connectTO
	}
	if other.keepAliveTO != 0 {
		db.keepAliveTO = other.keepAliveTO
	}
	if other.tls == true {
		db.tls = true
	}
	if other.skipVerifyTLS == true {
		db.skipVerifyTLS = true
	}
}