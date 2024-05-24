package main


import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GreetingResponse struct {
	Message string `json:"message"`
}

type LengthResponse struct {
	NameLength map[string]int `json:"name_length"`
}


func main() {
	http.HandleFunc("/hello", HelloHandler)

	fmt.Println("Microservice listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {

	greeting := GreetingResponse{
		Message: "Hello, world!",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(greeting); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/greet/")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/hello?name=%s", name))
	if err != nil {
		http.Error(w, "Failed to fetch greeting", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var lengthResponse LengthResponse
	if err := json.NewDecoder(resp.Body).Decode(&lengthResponse); err != nil {
		http.Error(w, "Failed to decode length response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lengthResponse); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}