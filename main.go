package main

import (
	"context"
	"encoding/json"
	"io/fs"
	"maps"
	"net/http"
	"os"
	"slices"

	"embed"
	"html/template"

	"github.com/WeatherGod3218/mlc-project-template/internal/airtable"
	"github.com/WeatherGod3218/mlc-project-template/internal/firebase"
	"github.com/WeatherGod3218/mlc-project-template/internal/logging"
	"github.com/WeatherGod3218/mlc-project-template/internal/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//go:embed templates/*
var embeddedFS embed.FS

//go:embed public/*
var staticFS embed.FS

func main() {
	redis.InitRedis()
	firebase.InitFirebase(context.Background())

	var tableMap map[string]string
	var tableList []string

	err := json.Unmarshal([]byte(os.Getenv("AIRTABLE_TABLES")), &tableMap)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("error loading airtables!")
	}

	tableList = slices.Collect(maps.Keys(tableMap))

	err = firebase.CreateCountTables(tableList)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "CreateCountTables"}).Fatal("error creating count tables!")
	}

	err = redis.InitQueue(tableList)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "LoadAllAirtables"}).Fatal("error queuing redis airtables!")
	}

	err = airtable.LoadAllAirtables(tableMap)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "LoadAllAirtables"}).Fatal("error loading airtables!")
	}

	router := gin.Default()
	router.Use(cors.Default())

	tmpl := template.Must(template.ParseFS(embeddedFS, "templates/*"))
	router.SetHTMLTemplate(tmpl)

	staticSub, err := fs.Sub(staticFS, "public")
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "main", "method": "main"}).Fatal("error embedding static files!")
	}
	router.StaticFS("/static", http.FS(staticSub))

	router.GET("/", redis.RedisRateLimiter(1, 50), GetHomepage)
	router.GET("/get-data", redis.RedisRateLimiter(1, 50), GetData)
	router.POST("/submit-results", redis.RedisRateLimiter(1, 50), SubmitResults)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router.Run(":" + port)
}
