package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	// "strconv"
)


type GreetingResponse struct {
	Message string `json:"message"`
}

type LengthResponse struct {
	NameLength map[string]int `json:"name_length"`
}

func main() {

	http.HandleFunc("/greet", GreetHandler)

	fmt.Println("Microservice listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	// req, err := http.NewResquest("GET", "http://localhost:8080/hello", nil)
	// if err != nil {
	// 	http.Error(w, "Failed to fetch greeting", http.StatusInternalServerError)
	// 	return
	// }

	// req.Header.Add("sdf", "sdfd")

	// client := http.Client{
	// 	Timeout: 5 * time.Second,
	// }

	// resp, err := client.Do(req)

	resp, err := http.Get("http://localhost:8080/hello")
	if err != nil {
		http.Error(w, "Failed to fetch greeting", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var greeting GreetingResponse
	if err := json.NewDecoder(resp.Body).Decode(&greeting); err != nil {
		http.Error(w, "Failed to decode greeting response", http.StatusInternalServerError)
		return
	}

	time.Sleep(500 * time.Millisecond)

	fmt.Fprintf(w, "Received greeting from first microservice: %s\n", greeting.Message)


}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	nameLength := map[string]int{name: len(name)}

	lengthResponse := LengthResponse{
		NameLength: nameLength,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lengthResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}