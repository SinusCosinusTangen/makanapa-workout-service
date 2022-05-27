package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"workout-webservice/config"
	"workout-webservice/models"

	elastic "github.com/olivere/elastic/v7"
)

var ESClient *elastic.Client

func InitializeLogger() (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(config.AppConfig.ESHost),
		elastic.SetBasicAuth(config.AppConfig.ESUsername, config.AppConfig.ESPassword),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("ES initialized...")

	ctx := context.Background()

	source := models.Source{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "Workout Webservice started",
		Version:   "1",
		Method:    "main",
		Type:      "workout-webservice",
	}

	fields := models.Fields{
		Timestamp: []time.Time{time.Now()},
	}

	var date string = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")
	var logIndex = fmt.Sprintf("workout-service-log-%s", date)

	log := models.Log{
		Index:   logIndex,
		Type:    "_doc",
		Version: 1,
		Source:  source,
		Fields:  fields,
	}

	dataJSON, err := json.Marshal(log)
	if err != nil {
		panic(err)
	}

	js := string(dataJSON)
	ind, err := client.Index().
		Index(logIndex).
		BodyJson(js).
		Do(ctx)

	if err != nil {
		fmt.Println(ind)
		panic(err)
	}

	ESClient = client

	return client, err
}
