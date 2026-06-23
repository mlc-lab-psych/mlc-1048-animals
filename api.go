package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/WeatherGod3218/mlc-project-template/internal/airtable"
	"github.com/WeatherGod3218/mlc-project-template/internal/firebase"
	"github.com/WeatherGod3218/mlc-project-template/internal/redis"

	"github.com/WeatherGod3218/mlc-project-template/internal/logging"
	"github.com/sirupsen/logrus"
)

// func replaceNullWithString(obj []map[string]interface{}) {
// 	for key, value := range obj {
// 		if value == nil {
// 			obj[key] = "null"
// 		} else if nested, ok := value.(map[string]interface{}); ok {
// 			replaceNullWithString(nested)
// 		}
// 	}
// }

func getAirtableData() (*airtable.SavedData, error) {
	tableName, err := redis.GetNextTable()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "getAirtableData"}).Warn("error in firebase deciding next table!")
		tableName = airtable.GetNextAirtable()
	}

	airTable := airtable.GetAirtableData(tableName)

	if airTable == nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "getAirtableData"}).Warn("error deciding which airtable to use!")
		return nil, fmt.Errorf("Failed to find the tablename %s", tableName)
	}

	return airTable, nil
}

func GetHomepage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func GetData(c *gin.Context) {
	data, err := getAirtableData()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "GetData"}).Fatal("error fetching airtable!")
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func SubmitResults(c *gin.Context) {
	var results []map[string]any

	err := c.ShouldBindJSON(&results)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "SubmitResults"}).Warn("error casting")
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	// replaceNullWithString(results)

	err = firebase.PushToDatabase(c, results)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "SubmitResults"}).Warn("error updating database!")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving results.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully saved results!"})
}
