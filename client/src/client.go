package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type User struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Status bool   `json:"status"`
}

const URL string = "http://localhost:8080"
const PATH string = "user/"

func main() {
	variation := "1"
	var err error
	var customUrl *url.URL

	for {
		fmt.Scan(&variation)

		customUrl, err = url.Parse(URL)
		if err != nil {
			fmt.Println(err)
			return
		}
		customUrl.Path += PATH

		var resp *http.Response
		var reqt *http.Request
		switch variation {
		case "1": // GET /user?status=true&name=alex
			params := url.Values{}
			params.Add("status", "true")
			params.Add("name", "test")
			customUrl.RawQuery = params.Encode()

			resp, err = http.Get(customUrl.String())
		case "2": // GET /user/{id}
			customUrl.Path += "1"
			resp, err = http.Get(customUrl.String())
		case "3": // DELETE /user/{id}
			customUrl.Path += "1"
			reqt, err = http.NewRequest("DELETE", customUrl.String(), bytes.NewBuffer([]byte{}))
			reqt.Header.Set("Content-Type", "application/json; charset=UTF-8")
			client := &http.Client{}
			resp, err = client.Do(reqt)

			defer resp.Body.Close()

		case "4": // PUT /user/{id}
			var u = User{
				Name:   "Alex",
				Status: true,
			}

			uInBytes, err := json.Marshal(u)
			if err != nil {
				fmt.Println("Can't encode user")
			}

			customUrl.Path += "2"
			reqt, err = http.NewRequest("PUT", customUrl.String(), bytes.NewBuffer(uInBytes))
			reqt.Header.Set("Content-Type", "application/json; charset=UTF-8")
			client := &http.Client{}
			resp, err = client.Do(reqt)

			defer resp.Body.Close()

		case "5": // POST /user/
			var u = User{
				Name:   "alex",
				Status: true,
			}

			uInBytes, err := json.Marshal(u)
			if err != nil {
				fmt.Println("Can't encode user")
			}

			resp, err = http.Post(customUrl.String(), "application/json", bytes.NewBuffer(uInBytes))
			if err != nil {
				fmt.Println("Can't post user")
			}
		default:
			break
		}

		if err != nil {
			fmt.Println(err)
		}

		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp.Status)
		if string(data) != "null" {
			fmt.Println(string(data))
		}
	}
}
