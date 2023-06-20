package service

import (
	"encoding/json"
	"log"
	"net/http"

	"teodorsavin/ah-bonus/util"
)

func ReadResponseBodyLogin(response *http.Response) string {
	defer response.Body.Close()

	authData := &util.AuthData{}
	err := json.NewDecoder(response.Body).Decode(authData)
	if err != nil {
		log.Panic(err)
	}

	if response.StatusCode != http.StatusOK {
		log.Panic(response.Status)
	}

	return authData.AccessToken
}
