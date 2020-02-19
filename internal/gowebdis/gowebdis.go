package gowebdis

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Name         string            `json:"name"`
	Success      bool              `json:"success"`
	ErrorMessage string            `json:"errorMessage"`
	BoolVal      bool              `json:"boolValue"`
	MapVal       map[string]string `json:"mapValue"`
	IntVal       int64             `json:"intVal"`
	StringVal    string            `json:"stringVal"`
}

func InitConnectionSetting(cmd *cobra.Command) error {

	sentinelAddressString := viper.GetString("sentinel-address")
	if len(sentinelAddressString) > 0 {
		connFailoverOptions = initFailoverConnectionSetting(sentinelAddressString, cmd)
	} else {
		hostString := viper.GetString("host")
		if len(hostString) > 0 {
			hostArray := strings.Split(hostString, ",")
			if len(hostArray) > 1 {
				connClusterOptions = initClusterConnectionSetting(hostArray, cmd)
			} else {
				connHostOptions = initHostConnectionSetting(hostArray[0], cmd)
			}
		}
	}
	if connType == "cluster" {

	} else {
		var commandResponse = ping()
		if !commandResponse.Success {
			return errors.New(commandResponse.ErrorMessage)
		}
	}
	return nil
}

func initFailoverConnectionSetting(sentinelAddressString string, cmd *cobra.Command) redis.FailoverOptions {

	connType = "sentinel"
	var options redis.FailoverOptions

	sentinelAddress := strings.Split(sentinelAddressString, ",")
	masterName := viper.GetString("master-name")

	password := viper.GetString("password")
	db := viper.GetInt("db")

	maxRetries := viper.GetInt("max-retries")
	minRetryBackoff := viper.GetInt("min-retry-backoff")
	maxRetryBackoff := viper.GetInt("max-retry-backoff")
	dialTimeout := viper.GetInt("dial-timeout")
	readTimeout := viper.GetInt("read-timeout")
	writeTimeout := viper.GetInt("write-timeout")

	poolSize := viper.GetInt("pool-size")
	minIdleConns := viper.GetInt("min-idle-conns")
	maxConnAge := viper.GetInt("min-conn-age")
	poolTimeout := viper.GetInt("pool-timeout")
	idleTimeout := viper.GetInt("idle-timeout")
	idleCheckFrequencey := viper.GetInt("idle-check-frequency")

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

	password := viper.GetString("password")

	maxRetries := viper.GetInt("max-retries")
	minRetryBackoff := viper.GetInt("min-retry-backoff")
	maxRetryBackoff := viper.GetInt("max-retry-backoff")
	dialTimeout := viper.GetInt("dial-timeout")
	readTimeout := viper.GetInt("read-timeout")
	writeTimeout := viper.GetInt("write-timeout")

	poolSize := viper.GetInt("pool-size")
	minIdleConns := viper.GetInt("min-idle-conns")
	maxConnAge := viper.GetInt("min-conn-age")
	poolTimeout := viper.GetInt("pool-timeout")
	idleTimeout := viper.GetInt("idle-timeout")
	idleCheckFrequencey := viper.GetInt("idle-check-frequency")

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

	password := viper.GetString("password")
	db := viper.GetInt("db")

	maxRetries := viper.GetInt("max-retries")
	minRetryBackoff := viper.GetInt("min-retry-backoff")
	maxRetryBackoff := viper.GetInt("max-retry-backoff")
	dialTimeout := viper.GetInt("dial-timeout")
	readTimeout := viper.GetInt("read-timeout")
	writeTimeout := viper.GetInt("write-timeout")

	poolSize := viper.GetInt("pool-size")
	minIdleConns := viper.GetInt("min-idle-conns")
	maxConnAge := viper.GetInt("min-conn-age")
	poolTimeout := viper.GetInt("pool-timeout")
	idleTimeout := viper.GetInt("idle-timeout")
	idleCheckFrequencey := viper.GetInt("idle-check-frequency")

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
	case "ping":
		commandResponse = ping()
	case "hset":
		commandResponse = hSet(jsonPayload.Key, jsonPayload.Field, jsonPayload.Value)
	case "hgetall":
		commandResponse = hGetAll(jsonPayload.Key)
	case "hdel":
		commandResponse = hDel(jsonPayload.Key, jsonPayload.Fields)
	default:
		commandResponse.Success = false
		commandResponse.ErrorMessage = fmt.Sprintf(`Does not support %v command`, redisCommand)
		log.Error("[ERROR] " + commandResponse.ErrorMessage)
	}
	return commandResponse
}

func ping() CommandResponse {
	var statusCmd *redis.StatusCmd
	var commandResponse = CommandResponse{Name: "ping"}

	if connType == "cluster" {
		var conn = startConnection()
		if conn != nil {
			statusCmd = conn.Ping()
			defer conn.Close()
		} else {
			commandResponse.Success = false
			commandResponse.ErrorMessage = "Cannot make redis connection"
			log.Error("[ERROR] " + commandResponse.ErrorMessage)
			return commandResponse
		}
	} else {
		var conn = startConnection()
		if conn != nil {
			statusCmd = conn.Ping()
			defer conn.Close()
		} else {
			commandResponse.Success = false
			commandResponse.ErrorMessage = "Cannot make redis connection"
			log.Error("[ERROR] " + commandResponse.ErrorMessage)
			return commandResponse
		}
	}

	var err = statusCmd.Err()
	if err != nil {
		commandResponse.Success = false
		commandResponse.ErrorMessage = err.Error()
		log.Error("[ERROR] " + commandResponse.ErrorMessage)
	} else {
		commandResponse.Success = true
		commandResponse.StringVal = statusCmd.Val()
		log.Info("[INFO] " + statusCmd.String())
	}
	return commandResponse

}

func hSet(key string, field string, value string) CommandResponse {
	var boolCmd *redis.BoolCmd
	var commandResponse = CommandResponse{Name: "hset"}

	if connType == "cluster" {
		var conn = startConnection()
		if conn != nil {
			boolCmd = conn.HSet(key, field, value)
			defer conn.Close()
		} else {
			commandResponse.Success = false
			commandResponse.ErrorMessage = "Cannot make redis connection"
			log.Error("[ERROR] " + commandResponse.ErrorMessage)
			return commandResponse
		}
	} else {
		var conn = startConnection()
		if conn != nil {
			boolCmd = conn.HSet(key, field, value)
			defer conn.Close()
		} else {
			commandResponse.Success = false
			commandResponse.ErrorMessage = "Cannot make redis connection"
			log.Error("[ERROR] " + commandResponse.ErrorMessage)
			return commandResponse
		}
	}

	var err = boolCmd.Err()
	if err != nil {
		commandResponse.Success = false
		commandResponse.ErrorMessage = err.Error()

		log.Error("[ERROR] " + commandResponse.ErrorMessage)
	} else {
		commandResponse.Success = true
		commandResponse.BoolVal = true
		log.Info("[INFO] " + boolCmd.String())
	}
	return commandResponse

}

func hGetAll(key string) CommandResponse {
	var stringStringMapCmd *redis.StringStringMapCmd
	var commandResponse = CommandResponse{Name: "hgetall"}

	if connType == "cluster" {
		var conn = startConnection()
		if conn != nil {
			stringStringMapCmd = conn.HGetAll(key)
			defer conn.Close()
		} else {
			commandResponse.Success = false
			commandResponse.ErrorMessage = "Cannot make redis connection"
			log.Error("[ERROR] " + commandResponse.ErrorMessage)
			return commandResponse
		}
	} else {
		var conn = startConnection()
		if conn != nil {
			stringStringMapCmd = conn.HGetAll(key)
			defer conn.Close()
		} else {
			commandResponse.Success = false
			commandResponse.ErrorMessage = "Cannot make redis connection"
			log.Error("[ERROR] " + commandResponse.ErrorMessage)
			return commandResponse
		}
	}

	var err = stringStringMapCmd.Err()
	if err != nil {
		commandResponse.Success = false
		commandResponse.ErrorMessage = err.Error()
		log.Error("[ERROR] " + commandResponse.ErrorMessage)
	} else {
		commandResponse.Success = true
		commandResponse.MapVal = stringStringMapCmd.Val()
		log.Info("[INFO] " + stringStringMapCmd.String())
	}
	return commandResponse
}

func hDel(key string, fields []string) CommandResponse {
	var intCmd *redis.IntCmd
	var commandResponse = CommandResponse{Name: "hdel"}

	if connType == "cluster" {
		var conn = startConnection()
		if conn != nil {
			intCmd = conn.HDel(key, fields...)
			defer conn.Close()
		} else {
			commandResponse.Success = false
			commandResponse.ErrorMessage = "Cannot make redis connection"
			log.Error("[ERROR] " + commandResponse.ErrorMessage)
			return commandResponse
		}
	} else {
		var conn = startConnection()
		if conn != nil {
			intCmd = conn.HDel(key, fields...)
			defer conn.Close()
		} else {
			commandResponse.Success = false
			commandResponse.ErrorMessage = "Cannot make redis connection"
			log.Error("[ERROR] " + commandResponse.ErrorMessage)
			return commandResponse
		}
	}

	var err = intCmd.Err()
	if err != nil {
		commandResponse.Success = false
		commandResponse.ErrorMessage = err.Error()
		log.Error("[ERROR] " + commandResponse.ErrorMessage)
	} else {
		commandResponse.Success = true
		log.Info("[INFO] " + intCmd.String())
	}
	return commandResponse

}
