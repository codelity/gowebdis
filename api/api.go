package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/codelity/gowebdis/internal/gowebdis"
)

type ApiResponse struct {
}

func StartServer() {

	router := gin.Default()
	router.GET("/healthz", pingCommand)
	router.POST("/:command", apiCommand)
	router.Run()
}

func pingCommand(context *gin.Context) {
	var jsonPayload gowebdis.JsonPayload
	var commandResponse gowebdis.CommandResponse
	commandResponse = gowebdis.RunRedisCommand("ping", jsonPayload)
	if commandResponse.Success {
		context.JSON(200, gin.H{"boolVal": true})
	} else {
		context.JSON(400, gin.H{"errorMessage": commandResponse.ErrorMessage})
	}
	return
}

func apiCommand(context *gin.Context) {

	var jsonPayload gowebdis.JsonPayload

	command, _ := context.Params.Get("command")
	var commandResponse gowebdis.CommandResponse

	err := context.BindJSON(&jsonPayload)
	if err != nil {
		log.Error("[ERROR] " + err.Error())
		context.JSON(400, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	err = validateJsonPayload(command, jsonPayload)
	if err != nil {
		log.Error("[ERROR] " + err.Error())
		context.JSON(400, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	commandResponse = gowebdis.RunRedisCommand(command, jsonPayload)
	if commandResponse.Success {
		var responsePayload map[string]interface{}
		switch commandResponse.Name {
		case "hset":
			responsePayload = gin.H{
				"boolValue": commandResponse.BoolVal,
			}
		case "hgetall":
			responsePayload = gin.H{
				"stringArrayValue": getStringArrayValue(commandResponse.MapVal),
			}
		case "hdel":
			responsePayload = gin.H{
				"intValue": commandResponse.IntVal,
			}
		}
		context.JSON(200, responsePayload)
	} else {
		context.JSON(400, gin.H{
			"errorMessage": commandResponse.ErrorMessage,
		})
	}
	return
}

func getStringArrayValue(m map[string]string) []string {
	v := make([]string, len(m), len(m))
	idx := 0
	for _, value := range m {
		v[idx] = value
		idx++
	}
	return v
}

func validateJsonPayload(command string, jsonPayload gowebdis.JsonPayload) error {
	var err error
	switch command {
	case "hset":
		if len(jsonPayload.Key) == 0 {
			err = errors.New("'key' attribute cannot be found in payload")
		} else if len(jsonPayload.Field) == 0 {
			err = errors.New("'field' attribute cannot be found in payload")
		} else if len(jsonPayload.Value) == 0 {
			err = errors.New("'value' attribute cannot be found in payload")
		}
	case "hgetall":
		if len(jsonPayload.Key) == 0 {
			err = errors.New("'key' attribute cannot be found in payload")
		}
	case "hdel":
		if len(jsonPayload.Key) == 0 {
			err = errors.New("'key' attribute cannot be found in payload")
		} else if len(jsonPayload.Fields) == 0 {
			err = errors.New("'fields' attribute is empty in payload")
		}
	}
	return err
}
