package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func Creater(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		error := map[string]interface{}{
			"result": false,
			"text":   "Please enter your username or password",
		}
		jsonData, err := json.Marshal(error)
		if err != nil {
			fmt.Printf("Could not marshal json: %s\n", err)
			return
		}
		fmt.Fprintf(w, "%s\n", jsonData)
		return
	}

	if _, err := os.Stat("./account/" + username + ".json"); err == nil {
		error := map[string]interface{}{
			"result": false,
			"text":   "Account already created",
		}
		jsonData, err := json.Marshal(error)
		if err != nil {
			fmt.Printf("Could not marshal json: %s\n", err)
			return
		}
		fmt.Fprintf(w, "%s\n", jsonData)
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
		res := map[string]interface{}{
			"result": true,
			"text":   "Account creation complete",
		}
		jsonData, err := json.Marshal(res)
		if err != nil {
			fmt.Printf("Could not marshal json: %s\n", err)
			return
		}
		fmt.Fprintf(w, "%s\n", jsonData)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		error := map[string]interface{}{
			"result": false,
			"text":   "Please enter your username or password",
		}
		jsonData, err := json.Marshal(error)
		if err != nil {
			fmt.Printf("Could not marshal json: %s\n", err)
			return
		}
		fmt.Fprintf(w, "%s\n", jsonData)
		return
	}

	if _, err := os.Stat("./account/" + username + ".json"); err != nil {
		error := map[string]interface{}{
			"result": false,
			"text":   "This account does not exist",
		}
		jsonData, err := json.Marshal(error)
		if err != nil {
			fmt.Printf("Could not marshal json: %s\n", err)
			return
		}
		fmt.Fprintf(w, "%s\n", jsonData)
	} else {
		h := sha256.New()
		h.Write([]byte(username + password))
		res := fmt.Sprintf("%x", h.Sum(nil))

		jsonFile, err := os.Open("./account/" + username + ".json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)

		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)
		if res == result["password"] {
			data := map[string]interface{}{
				"result":   true,
				"username": username,
				"text":     "Login successful",
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("Could not marshal json: %s\n", err)
				return
			}
			fmt.Fprintf(w, "%s\n", jsonData)
		} else {
			error := map[string]interface{}{
				"result": false,
				"text":   "The username or password is incorrect.",
			}
			jsonData, err := json.Marshal(error)
			if err != nil {
				fmt.Printf("Could not marshal json: %s\n", err)
				return
			}
			fmt.Fprintf(w, "%s\n", jsonData)
		}
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		error := map[string]interface{}{
			"result": false,
			"text":   "Please enter your username or password",
		}
		jsonData, err := json.Marshal(error)
		if err != nil {
			fmt.Printf("Could not marshal json: %s\n", err)
			return
		}
		fmt.Fprintf(w, "%s\n", jsonData)
		return
	}

	if _, err := os.Stat("./account/" + username + ".json"); err != nil {
		error := map[string]interface{}{
			"result": false,
			"text":   "This account does not exist",
		}
		jsonData, err := json.Marshal(error)
		if err != nil {
			fmt.Printf("Could not marshal json: %s\n", err)
			return
		}
		fmt.Fprintf(w, "%s\n", jsonData)
	} else {
		h := sha256.New()
		h.Write([]byte(username + password))
		res := fmt.Sprintf("%x", h.Sum(nil))

		jsonFile, err := os.Open("./account/" + username + ".json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)

		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)
		if res == result["password"] {
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
			data := map[string]interface{}{
				"result":   true,
				"username": username,
				"text":     "Deleted successful",
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("Could not marshal json: %s\n", err)
				return
			}
			fmt.Fprintf(w, "%s\n", jsonData)
			go os.Remove("./account/" + username + ".json")
			return
		} else {
			error := map[string]interface{}{
				"result": false,
				"text":   "The username or password is incorrect.",
			}
			jsonData, err := json.Marshal(error)
			if err != nil {
				fmt.Printf("Could not marshal json: %s\n", err)
				return
			}
			fmt.Fprintf(w, "%s\n", jsonData)
			return
		}
	}
}

func main() {
	http.HandleFunc("/login", Login)
	http.HandleFunc("/create", Creater)
	http.HandleFunc("/delete", Delete)

	fmt.Println("localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
