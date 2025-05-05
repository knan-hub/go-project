package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-project/logger"
	"go-project/model"
	"go-project/setting"
	"io"
	"net/http"
)

func makeRequest(method, url string, data interface{}, headers map[string]string) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("json.Marshal failed, err:%v\n", err))
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("http.NewRequest failed, err:%v\n", err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	return client.Do(req)
}

func ValidateUserToken(token string) (valid bool, user model.AuthUser) {
	url := fmt.Sprintf("%s/internal/authenticated", setting.Config.Servers.CoreserviceUrl)

	resp, err := makeRequest("POST", url, struct {
		Token string `json:"token"`
	}{Token: token}, map[string]string{
		"x-api-key": setting.Config.Self.INTERNAL_API_KEY,
	})

	if err != nil {
		return false, user
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("io.ReadAll failed, err:%v\n", err))
		return false, user
	}

	if resp.StatusCode == http.StatusOK {
		if err := json.Unmarshal(responseData, &user); err != nil {
			logger.Logger.Error(fmt.Sprintf("json.Unmarshal failed, err:%v\n", err))
			return false, user
		}
		return true, user
	}

	return false, user
}

func GetUsernameByUserId(userId int) (string, error) {
	url := fmt.Sprintf("%s/admins/users/search", setting.Config.Servers.CoreserviceUrl)

	resp, err := makeRequest("POST", url, map[string]int{
		"id": userId,
	}, map[string]string{
		"x-api-key": setting.Config.Self.INTERNAL_API_KEY,
	})

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("io.ReadAll failed, err:%v\n", err))
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(string(responseData))
	}

	var result struct {
		Rows []struct {
			Username string `json:"username"`
		} `json:"rows"`
	}

	if err := json.Unmarshal(responseData, &result); err != nil {
		logger.Logger.Error(fmt.Sprintf("json.Unmarshal failed, err:%v\n", err))
		return "", err
	}

	if len(result.Rows) == 0 {
		return "", errors.New("user not found")
	}

	return result.Rows[0].Username, nil
}
