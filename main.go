package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	server.POST("/create", func(res *gin.Context) {
		requestData, _ := ioutil.ReadAll(res.Request.Body)
		var jsonBody map[string]interface{}
		json.Unmarshal(requestData, &jsonBody)
		username := jsonBody["username"].(string)
		password := jsonBody["password"].(string)

		if username == "" || password == "" {
			res.JSON(404, gin.H{
				"result": false,
				"text":   "Please enter your username or password",
			})
		}

		if _, err := os.Stat("./account/" + username + ".json"); err == nil {
			res.JSON(404, gin.H{
				"result": false,
				"text":   "Account already created",
			})
		} else {
			h := sha256.New()
			h.Write([]byte(username + password))
			result := fmt.Sprintf("%x", h.Sum(nil))
			data := map[string]interface{}{
				"password": result,
				"join":     time.Now(),
			}
			file, _ := json.MarshalIndent(data, "", " ")
			os.WriteFile("./account/"+username+".json", file, 0644)
			res.JSON(200, gin.H{
				"result": true,
				"text":   "Account creation complete",
			})
		}
	})

	server.POST("/login", func(res *gin.Context) {
		requestData, _ := ioutil.ReadAll(res.Request.Body)
		var jsonBody map[string]interface{}
		json.Unmarshal(requestData, &jsonBody)
		username := jsonBody["username"].(string)
		password := jsonBody["password"].(string)

		if username == "" || password == "" {
			res.JSON(200, gin.H{
				"result": false,
				"text":   "Please enter your username or password",
			})
		}

		if _, err := os.Stat("./account/" + username + ".json"); err != nil {
			res.JSON(200, gin.H{
				"result": false,
				"text":   "This account does not exist",
			})
		} else {
			h := sha256.New()
			h.Write([]byte(username + password))
			result := fmt.Sprintf("%x", h.Sum(nil))
			jsonFile, err := os.Open("./account/" + username + ".json")
			if err != nil {
				fmt.Println(err)
			}
			defer jsonFile.Close()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			var r map[string]interface{}
			json.Unmarshal([]byte(byteValue), &r)
			if result == r["password"] {
				res.JSON(200, gin.H{
					"result":   true,
					"username": username,
					"text":     "Login successful",
				})
			} else {
				res.JSON(200, gin.H{
					"result": false,
					"text":   "The username or password is incorrect.",
				})
			}
		}
	})

	server.POST("/delete", func(res *gin.Context) {
		requestData, _ := ioutil.ReadAll(res.Request.Body)
		var jsonBody map[string]interface{}
		json.Unmarshal(requestData, &jsonBody)
		username := jsonBody["username"].(string)
		password := jsonBody["password"].(string)

		if username == "" || password == "" {
			res.JSON(404, gin.H{
				"result": false,
				"text":   "Please enter your username or password",
			})
		}

		if _, err := os.Stat("./account/" + username + ".json"); err != nil {
			res.JSON(404, gin.H{
				"result": false,
				"text":   "This account does not exist",
			})
		} else {
			h := sha256.New()
			h.Write([]byte(username + password))
			result := fmt.Sprintf("%x", h.Sum(nil))
			jsonFile, err := os.Open("./account/" + username + ".json")
			if err != nil {
				fmt.Println(err)
			}
			defer jsonFile.Close()
			byteValue, _ := ioutil.ReadAll(jsonFile)
			var r map[string]interface{}
			json.Unmarshal([]byte(byteValue), &r)
			if result == r["password"] {
				res.JSON(200, gin.H{
					"result":   true,
					"username": username,
					"text":     "Deleted successful",
				})
				go os.Remove("./account/" + username + ".json")
			} else {
				res.JSON(404, gin.H{
					"result": false,
					"text":   "The username or password is incorrect.",
				})
			}
		}
	})

	server.Run()
}
