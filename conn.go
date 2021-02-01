package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
)

// DBConfig is the struct that encapsulates the user inputs
type DBConfig struct {
	host     string
	username string
	password string
	client   string
	port     int
	database int

	timeout   time.Duration
	connectTO time.Duration
	readTO    time.Duration
	writeTO   time.Duration

	keepAliveTO   time.Duration
	tls           bool
	skipVerifyTLS bool
	debug         bool
	prod          bool
}

// RedisDB is an instance of redis database
type RedisDB struct {
	DBConfig
	redisConn redis.Conn
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
	prod := flag.Bool("prod", false, "Run in production mode")

	flag.Parse()

	conf := DBConfig{
		host:          *host,
		port:          *port,
		database:      *db,
		username:      *user,
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
		prod:          *prod,
	}
	return conf
}

/*
InitializeRedis initializes the redigo/redis client at startup based on the CLI
inputs in the DBConfig
*/
func (db *RedisDB) InitializeRedis() error {
	dialOpts := db.createDialOpts()
	address := fmt.Sprintf("%s:%d", db.host, db.port)
	r, err := redis.Dial("tcp", address, dialOpts...)
	if err != nil {
		return err
	}
	db.redisConn = r
	return nil
}

// TearDownRedis tears down the redis connection
func (db *RedisDB) TearDownRedis() {
	if db.redisConn != nil {
		db.redisConn.Close()
	}
}

func (db *RedisDB) createDialOpts() []redis.DialOption {
	var dialOpts []redis.DialOption

	dialOpts = append(dialOpts, redis.DialDatabase(db.database))
	if len(db.password) > 0 {
		dialOpts = append(dialOpts, redis.DialPassword(db.password))
	}
	if len(db.username) > 0 {
		dialOpts = append(dialOpts, redis.DialUsername(db.username))
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
	zap.S().Infow("Create redis with options", "options", db)
	return dialOpts
}

/*
Merge this instance of DBConfig with another DBConfig. This function
picks up all the non-zero values from other and assigns them to the
corresponding field in this DBConfig object
*/
func (db *DBConfig) Merge(other *DBConfig) {
	zap.S().Debugf("Request to merge config: %v", other)

	if other.host != "" {
		db.host = other.host
	}
	if other.port != 0 {
		db.port = other.port
	}
	if other.database != -1 {
		db.database = other.database
	}
	if other.username != "" {
		db.username = other.username
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
	if other.tls {
		db.tls = true
	}
	if other.skipVerifyTLS {
		db.skipVerifyTLS = true
	}
	zap.S().Infof("Merged DBConfig: %v", db)
}

func (db *DBConfig) String() string {
	var x map[string]string = make(map[string]string)

	x["redis"] = fmt.Sprintf("tcp://%s:%d/%d", db.host, db.port, db.database)
	x["credentials"] = fmt.Sprintf("%s:%s", db.username, "****(Redacted)")
	x["timeouts"] = fmt.Sprintf("%v", db.timeout)
	x["connect-timeout"] = fmt.Sprintf("%v", db.connectTO)
	x["read-timeout"] = fmt.Sprintf("%v", db.readTO)
	x["write-timeout"] = fmt.Sprintf("%v", db.writeTO)
	x["keep-alive-timeout"] = fmt.Sprintf("%v", db.keepAliveTO)
	x["client"] = db.client
	x["tls"] = fmt.Sprintf("%v", db.tls)
	x["tls-verify"] = fmt.Sprintf("%v", !db.skipVerifyTLS)

	j, _ := json.Marshal(x)
	return string(j)
}
