package communication

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"manage-service/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	LoginBusy            = errors.New("Данный логин уже занят")
	LoginNotFound        = errors.New("Неправильный логин или пароль")
	PasswordIsNotCorrect = errors.New("Неправильный логин или пароль")
)

func RegistrationRequest(c *gin.Context, user models.User) (string, bool, error) {
	admin := false
	UUID := ""
	data, err := json.Marshal(user)
	if err != nil {
		return UUID, admin, fmt.Errorf("communication/Authentication json.Marshal: %v", err)
	}

	client := &http.Client{Timeout: 6 * time.Second}

	req, err := http.NewRequest(
		"POST",
		"http://auth_service:8081/registration",
		bytes.NewBuffer(data),
	)

	if err != nil {
		return UUID, admin, fmt.Errorf("communication/Authentication http.NewRequest: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return UUID, admin, fmt.Errorf("communication/Authentication client.Do: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UUID, admin, fmt.Errorf("communication/Authentication io.ReadAll: %v", err)
	}

	respAuth := models.ResponseAuth{}

	if err := json.Unmarshal(body, &respAuth); err != nil {
		return UUID, admin, fmt.Errorf("communication/Authentication json.Unmarhsal: %v", err)
	}

	if respAuth.Code != http.StatusCreated {
		if respAuth.Code == http.StatusConflict {
			return UUID, admin, LoginBusy
		}
		return UUID, admin, fmt.Errorf("Ошибка на стороне auth: %v", err)
	}

	return respAuth.UUID, respAuth.Admin, nil
}

func LoginRequest(user models.User) (string, bool, error) {
	admin := false
	UUID := ""
	data, err := json.Marshal(user)
	if err != nil {
		return UUID, admin, fmt.Errorf("communication/Authentication json.Marshal: %v", err)
	}

	client := &http.Client{Timeout: 6 * time.Second}

	req, err := http.NewRequest(
		"POST",
		"http://auth_service:8081/login",
		bytes.NewBuffer(data),
	)

	if err != nil {
		return UUID, admin, fmt.Errorf("communication/Authentication http.NewRequest: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return UUID, admin, fmt.Errorf("communication/Authentication client.Do: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UUID, admin, fmt.Errorf("communication/Authentication io.ReadAll: %v", err)
	}

	respAuth := models.ResponseAuth{}

	if err := json.Unmarshal(body, &respAuth); err != nil {
		return UUID, admin, fmt.Errorf("communication/Authentication json.Unmarhsal: %v", err)
	}

	if respAuth.Code != http.StatusCreated {
		if respAuth.Code == http.StatusNotFound {
			return UUID, admin, LoginNotFound
		}
		if respAuth.Code == http.StatusUnauthorized {
			return UUID, admin, PasswordIsNotCorrect
		}
		return UUID, admin, fmt.Errorf("Ошибка на стороне auth: %v", err)
	}

	return respAuth.UUID, respAuth.Admin, nil
}

func AuthorizationRequest(GUID string, admin bool) (models.Tokens, error) {
	var GUIDS struct {
		GUID  string `json:"id"`
		Admin bool   `json:"admin"`
	}

	GUIDS.GUID = GUID
	GUIDS.Admin = admin
	tokens := models.Tokens{}
	data, err := json.Marshal(&GUIDS)
	if err != nil {
		return tokens, fmt.Errorf("communication/Authentication json.Marshal: %v", err)
	}

	client := &http.Client{Timeout: 6 * time.Second}

	req, err := http.NewRequest(
		"POST",
		"http://autoriz_service:8083/authorization",
		bytes.NewBuffer(data),
	)

	if err != nil {
		return tokens, fmt.Errorf("communication/Authentication http.NewRequest: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return tokens, fmt.Errorf("communication/Authentication client.Do: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tokens, fmt.Errorf("communication/Authentication io.ReadAll: %v", err)
	}

	respAuthoriz := models.ResponseTokens{}

	if err := json.Unmarshal(body, &respAuthoriz); err != nil {
		return tokens, fmt.Errorf("communication/Authentication json.Unmarhsal: %v", err)
	}

	if respAuthoriz.Code != http.StatusOK {
		return tokens, fmt.Errorf("Ошибка на стороне authoriz: %v", err)
	}

	return respAuthoriz.Tokens, nil
}

func GUIDRequest(access string) (string, bool, error) {
	var ResponseAdmin struct {
		GUID  string `json:"id"`
		Admin bool   `json:"admin"`
		models.Response
	}
	{
	}
	resp, err := http.Get(
		"http://autoriz_service:8083/uuid/" + access)
	if err != nil {
		return "", false, fmt.Errorf("communication/Authentication client.Do: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("communication/Authentication io.ReadAll: %v", err)
	}

	if err := json.Unmarshal(body, &ResponseAdmin); err != nil {
		return "", false, fmt.Errorf("communication/Authentication json.Unmarhsal: %v", err)
	}

	if ResponseAdmin.Code != http.StatusOK {
		return "", false, fmt.Errorf("Ошибка на стороне authoriz: %v", err)
	}

	return ResponseAdmin.GUID, ResponseAdmin.Admin, nil
}

func RefreshRequest(access, refresh string) (models.Tokens, error) {
	tokens := models.Tokens{Access: access, Refresh: refresh}
	data, err := json.Marshal(&tokens)
	if err != nil {
		return tokens, fmt.Errorf("communication/Authentication json.Marshal: %v", err)
	}

	client := &http.Client{Timeout: 6 * time.Second}

	req, err := http.NewRequest(
		"POST",
		"http://autoriz_service:8083/refresh",
		bytes.NewBuffer(data),
	)

	if err != nil {
		return tokens, fmt.Errorf("communication/Authentication http.NewRequest: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return tokens, fmt.Errorf("communication/Authentication client.Do: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return tokens, fmt.Errorf("communication/Authentication io.ReadAll: %v", err)
	}

	respAuthoriz := models.ResponseTokens{}

	if err := json.Unmarshal(body, &respAuthoriz); err != nil {
		return tokens, fmt.Errorf("communication/Authentication json.Unmarhsal: %v", err)
	}

	if respAuthoriz.Code != http.StatusOK {
		return tokens, fmt.Errorf("Ошибка на стороне authoriz: %v", err)
	}

	return respAuthoriz.Tokens, nil
}

func LogoutRequest(access string) error {
	tokens := models.Tokens{Access: access}
	data, err := json.Marshal(&tokens)
	if err != nil {
		return fmt.Errorf("communication/Logout json.Marshal: %v", err)
	}

	client := &http.Client{Timeout: 6 * time.Second}

	req, err := http.NewRequest(
		"POST",
		"http://autoriz_service:8083/logout",
		bytes.NewBuffer(data),
	)

	if err != nil {
		return fmt.Errorf("communication/Logout http.NewRequest: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("communication/Logout client.Do: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("communication/Logout io.ReadAll: %v", err)
	}

	respLogout := models.Response{}

	if err := json.Unmarshal(body, &respLogout); err != nil {
		return fmt.Errorf("communication/Logout json.Unmarhsal: %v", err)
	}

	if respLogout.Code != http.StatusOK {
		return fmt.Errorf("Ошибка на стороне authoriz: %v", err)
	}

	return nil
}
