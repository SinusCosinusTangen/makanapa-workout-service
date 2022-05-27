package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"workout-webservice/models"
	"workout-webservice/services"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

const WorkoutNotFoundMessage = "System Error (Workout Not Found)."
const HistoryNotFoundMessage = "System Error (History Not Found)."

func FetchWorkoutByTime(c *gin.Context) {
	db := services.Database
	esClient := services.ESClient

	var input models.Request
	if err := c.ShouldBindJSON(&input); err != nil {
		logger(esClient, []string{"error", "FetchWorkoutByTime", "Failed to bind JSON", requestToString(&input)})
		c.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
		return
	}

	logger(esClient, []string{"info", "FetchWorkoutByTime", "Request: ", requestToString(&input)})

	userId := input.UserID
	targetTime := input.Time * 60
	var workouts []models.Workout
	var selectedId []int

	_, err := services.Authentication(userId, c.GetHeader("Authorization"))

	if err != nil {
		logger(esClient, []string{"error", "FetchWorkoutByTime", "Failed to authenticate. Reason: ", err.Error()})
		c.JSON(http.StatusUnauthorized, models.Response{Error: err.Error()})
		return
	}

	for time := 0; time < targetTime; time += 0 {
		var randomId = rand.Intn(6)
		randomId += 1

		if contains(selectedId, randomId) {
			var iter = 0
			for ok := true; ok; ok = contains(selectedId, randomId) {
				if iter == 50 {
					break
				}
				randomId = rand.Intn(6)
				randomId += 1
				iter += 1
			}
		}

		selectedId = append(selectedId, randomId)

		var workout models.Workout
		if result := db.First(&workout, "id = ?", randomId); result.Error != nil {
			logger(esClient, []string{"error", "FetchWorkoutByTime", "Failed to find workout: ", strconv.Itoa(randomId)})
			c.JSON(http.StatusNotFound, models.Response{Error: WorkoutNotFoundMessage})
			return
		}
		time += workout.Time

		workouts = append(workouts, workout)
	}

	history := models.History{
		UserID:   userId,
		Workouts: workouts,
	}
	db.Create(&history)

	var response models.Response = models.Response{Data: workouts}
	logger(esClient, []string{"info", "FetchWorkoutByTime", "Response: ", responseToString(&response)})

	jsonResponse, err := json.Marshal(workouts)
	if err != nil {
		panic(err)
	}

	resp, err := services.SendToRabitMQ(jsonResponse)
	if err != nil {
		logger(esClient, []string{"error", "FetchWorkoutByTime", "Failed to send to rabitmq", requestToString(&input)})
		c.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}

	logger(esClient, []string{"info", "FetchWorkoutByTime", "Response: ", resp})

	responseData := models.History{
		UserID:   userId,
		Workouts: workouts,
	}

	jsonResponse, err = json.Marshal(responseData)
	c.JSON(http.StatusOK, responseData)
}

func FetchWorkoutHistory(c *gin.Context) {
	db := services.Database
	esClient := services.ESClient

	logger(esClient, []string{"info", "FetchWorkoutHistory", "Request: ", "Fetching workout history"})

	userId, err := services.GetUser(c.GetHeader("Authorization"))

	if err != nil {
		logger(esClient, []string{"error", "FetchWorkoutByTime", "Cannot get User: ", err.Error()})
		c.JSON(http.StatusUnauthorized, models.Response{Error: err.Error()})
		return
	}

	res, err := services.Authentication(userId, c.GetHeader("Authorization"))
	logger(esClient, []string{"info", "FetchWorkoutHistory", "Authentication: ", string(res)})

	if err != nil {
		logger(esClient, []string{"error", "FetchWorkoutByTime", "Failed to authenticate", err.Error()})
		c.JSON(http.StatusUnauthorized, models.Response{Error: err.Error()})
		return
	}

	var history []models.History
	if result := db.Preload("Workouts").Find(&history, "user_id = ?", userId); result.Error != nil {
		logger(esClient, []string{"error", "FetchWorkoutHistory", "Failed to find history of user: ", userId})
		c.JSON(http.StatusNotFound, models.Response{Error: HistoryNotFoundMessage})
		return
	}

	var response models.Response = models.Response{Data: history}
	logger(esClient, []string{"info", "FetchWorkoutHistory", "Response: ", responseToString(&response)})
	c.JSON(http.StatusOK, history)
}

func responseToString(response *models.Response) string {
	dataJSON, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return string(dataJSON)
}

func requestToString(request *models.Request) string {
	dataJSON, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	return string(dataJSON)
}

func logger(esClient *elastic.Client, args []string) {
	ctx := context.Background()

	var level string = args[0]
	var method string = args[1]
	var message string = args[2]
	if args[3] != "" {
		message = message + " " + args[3]
	}

	source := models.Source{
		Timestamp: time.Now().Local(),
		Level:     level,
		Message:   message,
		Version:   "1",
		Method:    method,
		Type:      "workout-webservice",
	}

	fields := models.Fields{
		Timestamp: []time.Time{time.Now().Local()},
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
	ind, err := esClient.Index().
		Index(logIndex).
		BodyJson(js).
		Do(ctx)

	if err != nil {
		fmt.Println(ind)
		panic(err)
	}
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
