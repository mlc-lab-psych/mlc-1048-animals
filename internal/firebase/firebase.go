package firebase

import (
	"context"
	"fmt"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"

	"github.com/WeatherGod3218/mlc-project-template/internal/logging"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var database *db.Client

func InitFirebase(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	options := option.WithAuthCredentialsJSON(option.ServiceAccount, []byte(os.Getenv("FIREBASE_CREDENTIALS_JSON")))

	app, err := firebase.NewApp(ctx, &firebase.Config{
		DatabaseURL: os.Getenv("FIREBASE_DATABASE_URL"),
	}, options)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "firebase", "method": "initFirebase"}).Fatal("error initializing firebase")
	}

	database, err = app.Database(ctx)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "firebase", "method": "initFirebase"}).Fatal("error connecting to firebase database")
	}

	logging.Logger.WithFields(logrus.Fields{"module": "firebase", "method": "initFirebase"}).Info("connected to firebase!")
}

func CreateCountTables(tableNames []string) error {

	ref := database.NewRef("count")

	for _, name := range tableNames {

		var count *int
		childRef := ref.Child(name)

		err := childRef.Get(context.Background(), &count)
		if err != nil {
			return err
		}

		if count == nil {
			err = childRef.Set(context.Background(), 0)
		}
	}

	return nil
}

func PushToDatabase(ctx context.Context, cleanedData []map[string]any) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	id, err := uuid.NewV7()
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "firebase", "method": "PushToDatabase"}).Warn("error generating the UUID7")
	}
	err = database.NewRef(fmt.Sprintf("users/%s", id.String())).Set(ctx, cleanedData)

	return err
}

// func GetLowestTable(ctx context.Context) (string, error) {
// 	ref := database.NewRef("count")

// 	var countData map[string]int
// 	err := ref.Get(ctx, &countData)
// 	if err != nil {
// 		return "", err
// 	}

// 	if len(countData) == 0 {
// 		return "", fmt.Errorf("no count data available")
// 	}

// 	var lowestKey string
// 	first := true

// 	for key, value := range countData {
// 		if first || value < countData[lowestKey] {
// 			lowestKey = key
// 			first = false
// 		}
// 	}

// 	updates := map[string]interface{}{
// 		"count/" + lowestKey: countData[lowestKey] + 1,
// 	}

// 	err = database.NewRef("").Update(ctx, updates)
// 	if err != nil {
// 		return "", err
// 	}

// 	return lowestKey, nil
// }
