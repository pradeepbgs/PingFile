package config

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type CookieData struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path,omitempty"`
	Domain   string `json:"domain,omitempty"`
	Expires  string `json:"expires,omitempty"`
	Secure   bool   `json:"secure"`
	HttpOnly bool   `json:"http_only"`
}

func ParseCookie (filename string) ([]*http.Cookie , error) {
	file,err := os.Open(filename)
	if err != nil {
		return nil,err
	}
	defer file.Close()

	var cookieList []CookieData
	if err := json.NewDecoder(file).Decode(&cookieList); err != nil {
		return nil,err
	}

	var cookies []*http.Cookie
	for _, cookieData := range cookieList {

		if cookieData.Name == "" || cookieData.Value == "" {
			continue
		}

		var expires time.Time
		if cookieData.Expires != ""{
			expires ,_ = time.Parse(time.RFC1123, cookieData.Expires)
		}

		cookies = append(cookies, &http.Cookie{
			Name: cookieData.Name,
			Value: cookieData.Value,
			Path: cookieData.Path,
			Domain: cookieData.Domain,
			Expires: expires,
			Secure: cookieData.Secure,
			HttpOnly: cookieData.HttpOnly,
		})
	}

	return cookies,nil
	
}
