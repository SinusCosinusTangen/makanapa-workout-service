package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"workout-webservice/config"
	"workout-webservice/models"
)

func Authentication(userId string, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", config.AppConfig.AuthHost, nil)
	req.Header.Set("Authorization", token)
	req.Header.Add("Accept", "application/json")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("[ERROR] Failed to authenticate user: token")
	}

	if err != nil {
		defer resp.Body.Close()
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		defer resp.Body.Close()
		return nil, err
	}

	isMatch, err := checkUser(string(body), userId)

	if err != nil && !isMatch {
		defer resp.Body.Close()
		return nil, err
	}

	defer resp.Body.Close()
	return body, nil
}

func checkUser(resp string, userId string) (bool, error) {
	if strings.Contains(resp, "\"id\":"+userId) {
		return true, nil
	}
	return false, fmt.Errorf("[ERROR] Failed to authenticate user: user not match")
}

func GetUser(token string) (string, error) {
	req, err := http.NewRequest("GET", config.AppConfig.AuthHost, nil)
	req.Header.Set("Authorization", token)
	req.Header.Add("Accept", "application/json")

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return "", fmt.Errorf("[ERROR] Failed to authenticate user: token")
	}

	if err != nil {
		defer resp.Body.Close()
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		defer resp.Body.Close()
		return "", err
	}

	var user models.User
	err = json.Unmarshal(body, &user)

	if err != nil {
		defer resp.Body.Close()
		return "", err
	}

	return strconv.Itoa(user.Id), nil
}
