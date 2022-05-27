package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"workout-webservice/config"
)

func Authentication(userId string, token string) ([]byte, error) {
	var bearer = "Bearer " + token
	req, err := http.NewRequest("GET", config.AppConfig.AuthHost, nil)
	req.Header.Set("Authorization", bearer)
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
