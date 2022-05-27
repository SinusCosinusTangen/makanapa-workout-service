package services

import (
	"bytes"
	"net/http"
	"workout-webservice/config"
)

func SendToRabitMQ(json []byte) (string, error) {
	req, err := http.NewRequest("POST", config.AppConfig.MQHost, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return "failed", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "failed", err
	}
	defer resp.Body.Close()

	return "success", nil
}
