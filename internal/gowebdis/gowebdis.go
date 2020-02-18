package gowebdis

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var connType string

var connFailoverOptions redis.FailoverOptions
var connClusterOptions redis.ClusterOptions
var connHostOptions redis.Options

var client interface{}

type JsonPayload struct {
	Key    string   `json:"key" binding:"required"`
	Field  string   `json:"field"`
	Fields []string `json:"fields"`
	Value  string   `json:"value"`
}

type CommandResponse struct {
	Name         string `json:"name"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
	BoolVal      bool   `json:"boolValue"`
	MapVal       map[string]string
	IntVal       int64
}

func InitConnectionSetting(cmd *cobra.Command) {
	sentinelAddressString, _ := cmd.Flags().GetString("sentinel-address")
	if len(sentinelAddressString) > 0 {
		connFailoverOptions = initFailoverConnectionSetting(sentinelAddressString, cmd)
	} else {
		hostString, _ := cmd.Flags().GetString("host")
		if len(hostString) > 0 {
			hostArray := strings.Split(hostString, ",")
			if len(hostArray) > 1 {
				connClusterOptions = initClusterConnectionSetting(hostArray, cmd)
			} else {
				connHostOptions = initHostConnectionSetting(hostArray[0], cmd)
			}
		}

	}

}

func initFailoverConnectionSetting(sentinelAddressString string, cmd *cobra.Command) redis.FailoverOptions {

	connType = "sentinel"
	var options redis.FailoverOptions

	sentinelAddress := strings.Split(sentinelAddressString, ",")
	masterName, _ := cmd.Flags().GetString("master-name")

	password, _ := cmd.Flags().GetString("password")
	db, _ := cmd.Flags().GetInt("db")

	maxRetries, _ := cmd.Flags().GetInt("max-retries")
	minRetryBackoff, _ := cmd.Flags().GetInt("min-retry-backoff")
	maxRetryBackoff, _ := cmd.Flags().GetInt("max-retry-backoff")
	dialTimeout, _ := cmd.Flags().GetInt("dial-timeout")
	readTimeout, _ := cmd.Flags().GetInt("read-timeout")
	writeTimeout, _ := cmd.Flags().GetInt("write-timeout")

	poolSize, _ := cmd.Flags().GetInt("pool-size")
	minIdleConns, _ := cmd.Flags().GetInt("min-idle-conns")
	maxConnAge, _ := cmd.Flags().GetInt("min-conn-age")
	poolTimeout, _ := cmd.Flags().GetInt("pool-timeout")
	idleTimeout, _ := cmd.Flags().GetInt("idle-timeout")
	idleCheckFrequencey, _ := cmd.Flags().GetInt("idle-check-frequency")

	options.MasterName = masterName
	options.SentinelAddrs = sentinelAddress
	if len(password) > 0 {
		options.Password = password
	}
	options.DB = db
	if maxRetries > -1 {
		options.MaxRetries = maxRetries
	}
	if minRetryBackoff > -1 {
		options.MinRetryBackoff = time.Duration(minRetryBackoff) * time.Second
	}
	if maxRetryBackoff > -1 {
		options.MaxRetryBackoff = time.Duration(maxRetryBackoff) * time.Second
	}
	if dialTimeout > -1 {
		options.DialTimeout = time.Duration(dialTimeout) * time.Second
	}
	if readTimeout > -1 {
		options.ReadTimeout = time.Duration(readTimeout) * time.Second
	}
	if writeTimeout > -1 {
		options.WriteTimeout = time.Duration(writeTimeout) * time.Second
	}
	options.PoolSize = poolSize
	options.MinIdleConns = minIdleConns
	if maxConnAge > -1 {
		options.MaxConnAge = time.Duration(maxConnAge) * time.Second
	}
	if poolTimeout > -1 {
		options.PoolTimeout = time.Duration(poolTimeout) * time.Second
	}
	if idleTimeout > -1 {
		options.IdleTimeout = time.Duration(idleTimeout) * time.Second
	}
	if idleCheckFrequencey > -1 {
		options.IdleCheckFrequency = time.Duration(idleCheckFrequencey) * time.Second
	}

	// Dialer    func(ctx context.Context, network, addr string) (net.Conn, error)
	// OnConnect func(*Conn) error
	// TLSConfig *tls.Config

	return options

}

func initClusterConnectionSetting(hostAddresses []string, cmd *cobra.Command) redis.ClusterOptions {

	connType = "cluster"
	var options redis.ClusterOptions

	password, _ := cmd.Flags().GetString("password")

	maxRetries, _ := cmd.Flags().GetInt("max-retries")
	minRetryBackoff, _ := cmd.Flags().GetInt("min-retry-backoff")
	maxRetryBackoff, _ := cmd.Flags().GetInt("max-retry-backoff")
	dialTimeout, _ := cmd.Flags().GetInt("dial-timeout")
	readTimeout, _ := cmd.Flags().GetInt("read-timeout")
	writeTimeout, _ := cmd.Flags().GetInt("write-timeout")

	poolSize, _ := cmd.Flags().GetInt("pool-size")
	minIdleConns, _ := cmd.Flags().GetInt("min-idle-conns")
	maxConnAge, _ := cmd.Flags().GetInt("min-conn-age")
	poolTimeout, _ := cmd.Flags().GetInt("pool-timeout")
	idleTimeout, _ := cmd.Flags().GetInt("idle-timeout")
	idleCheckFrequencey, _ := cmd.Flags().GetInt("idle-check-frequency")

	options.Addrs = hostAddresses
	if len(password) > 0 {
		options.Password = password
	}
	if maxRetries > -1 {
		options.MaxRetries = maxRetries
	}
	if minRetryBackoff > -1 {
		options.MinRetryBackoff = time.Duration(minRetryBackoff) * time.Second
	}
	if maxRetryBackoff > -1 {
		options.MaxRetryBackoff = time.Duration(maxRetryBackoff) * time.Second
	}
	if dialTimeout > -1 {
		options.DialTimeout = time.Duration(dialTimeout) * time.Second
	}
	if readTimeout > -1 {
		options.ReadTimeout = time.Duration(readTimeout) * time.Second
	}
	if writeTimeout > -1 {
		options.WriteTimeout = time.Duration(writeTimeout) * time.Second
	}
	options.PoolSize = poolSize
	options.MinIdleConns = minIdleConns
	if maxConnAge > -1 {
		options.MaxConnAge = time.Duration(maxConnAge) * time.Second
	}
	if poolTimeout > -1 {
		options.PoolTimeout = time.Duration(poolTimeout) * time.Second
	}
	if idleTimeout > -1 {
		options.IdleTimeout = time.Duration(idleTimeout) * time.Second
	}
	if idleCheckFrequencey > -1 {
		options.IdleCheckFrequency = time.Duration(idleCheckFrequencey) * time.Second
	}
	return options
}

func initHostConnectionSetting(hostAddress string, cmd *cobra.Command) redis.Options {

	connType = "host"
	var options redis.Options

	password, _ := cmd.Flags().GetString("password")
	db, _ := cmd.Flags().GetInt("db")

	maxRetries, _ := cmd.Flags().GetInt("max-retries")
	minRetryBackoff, _ := cmd.Flags().GetInt("min-retry-backoff")
	maxRetryBackoff, _ := cmd.Flags().GetInt("max-retry-backoff")
	dialTimeout, _ := cmd.Flags().GetInt("dial-timeout")
	readTimeout, _ := cmd.Flags().GetInt("read-timeout")
	writeTimeout, _ := cmd.Flags().GetInt("write-timeout")

	poolSize, _ := cmd.Flags().GetInt("pool-size")
	minIdleConns, _ := cmd.Flags().GetInt("min-idle-conns")
	maxConnAge, _ := cmd.Flags().GetInt("min-conn-age")
	poolTimeout, _ := cmd.Flags().GetInt("pool-timeout")
	idleTimeout, _ := cmd.Flags().GetInt("idle-timeout")
	idleCheckFrequencey, _ := cmd.Flags().GetInt("idle-check-frequency")

	options.Addr = hostAddress
	if len(password) > 0 {
		options.Password = password
	}
	options.DB = db
	if maxRetries > -1 {
		options.MaxRetries = maxRetries
	}
	if minRetryBackoff > -1 {
		options.MinRetryBackoff = time.Duration(minRetryBackoff) * time.Second
	}
	if maxRetryBackoff > -1 {
		options.MaxRetryBackoff = time.Duration(maxRetryBackoff) * time.Second
	}
	if dialTimeout > -1 {
		options.DialTimeout = time.Duration(dialTimeout) * time.Second
	}
	if readTimeout > -1 {
		options.ReadTimeout = time.Duration(readTimeout) * time.Second
	}
	if writeTimeout > -1 {
		options.WriteTimeout = time.Duration(writeTimeout) * time.Second
	}
	options.PoolSize = poolSize
	options.MinIdleConns = minIdleConns
	if maxConnAge > -1 {
		options.MaxConnAge = time.Duration(maxConnAge) * time.Second
	}
	if poolTimeout > -1 {
		options.PoolTimeout = time.Duration(poolTimeout) * time.Second
	}
	if idleTimeout > -1 {
		options.IdleTimeout = time.Duration(idleTimeout) * time.Second
	}
	if idleCheckFrequencey > -1 {
		options.IdleCheckFrequency = time.Duration(idleCheckFrequencey) * time.Second
	}
	return options
}

func startConnection() *redis.Client {
	var conn *redis.Client
	if connType == "sentinel" {
		conn = redis.NewFailoverClient(&connFailoverOptions)
	} else if connType == "host" {
		conn = redis.NewClient(&connHostOptions)
	}
	return conn
}

func startClusterConnection() *redis.ClusterClient {
	return redis.NewClusterClient(&connClusterOptions)
}

func RunRedisCommand(redisCommand string, jsonPayload JsonPayload) CommandResponse {
	var commandResponse = CommandResponse{}
	switch redisCommand {
	case "hset":
		commandResponse = hSet(jsonPayload.Key, jsonPayload.Field, jsonPayload.Value)
	case "hgetall":
		commandResponse = hGetAll(jsonPayload.Key)
	case "hdel":
		commandResponse = hDel(jsonPayload.Key, jsonPayload.Fields)
	default:
		commandResponse.Success = false
		commandResponse.ErrorMessage = fmt.Sprintf(`Does not support %v command`, redisCommand)
		log.Error(commandResponse.ErrorMessage)
	}
	return commandResponse
}

func hSet(key string, field string, value string) CommandResponse {
	var boolCmd *redis.BoolCmd

	if connType == "cluster" {
		var conn = startConnection()
		boolCmd = conn.HSet(key, field, value)
		conn.Close()
	} else {
		var conn = startConnection()
		boolCmd = conn.HSet(key, field, value)
		conn.Close()
	}

	var commandResponse = CommandResponse{Name: boolCmd.Name()}
	var err = boolCmd.Err()
	if err != nil {
		commandResponse.Success = false
		commandResponse.ErrorMessage = err.Error()
		log.Error(err)
	} else {
		commandResponse.Success = true
		commandResponse.BoolVal = true
	}
	return commandResponse

}

func hGetAll(key string) CommandResponse {
	var stringStringMapCmd *redis.StringStringMapCmd

	if connType == "cluster" {
		var conn = startConnection()
		stringStringMapCmd = conn.HGetAll(key)
		conn.Close()
	} else {
		var conn = startConnection()
		stringStringMapCmd = conn.HGetAll(key)
		conn.Close()
	}

	var commandResponse = CommandResponse{Name: stringStringMapCmd.Name()}
	var err = stringStringMapCmd.Err()
	if err != nil {
		commandResponse.Success = false
		commandResponse.ErrorMessage = err.Error()
		log.Error(err)
	} else {
		commandResponse.Success = true
		commandResponse.MapVal = stringStringMapCmd.Val()
		log.Info(commandResponse.MapVal)
	}
	return commandResponse

}

func hDel(key string, fields []string) CommandResponse {
	var intCmd *redis.IntCmd

	if connType == "cluster" {
		var conn = startConnection()
		intCmd = conn.HDel(key, fields...)
		conn.Close()
	} else {
		var conn = startConnection()
		intCmd = conn.HDel(key, fields...)
		conn.Close()
	}

	var commandResponse = CommandResponse{Name: intCmd.Name()}
	var err = intCmd.Err()
	if err != nil {
		commandResponse.Success = false
		commandResponse.ErrorMessage = err.Error()
		log.Error(err)
	} else {
		commandResponse.Success = true
		commandResponse.IntVal = intCmd.Val()
	}
	return commandResponse

}
