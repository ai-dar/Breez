package handlers

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func JSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Ответ для GET-запроса
		response := JSONResponse{Status: "success", Message: "GET request received successfully"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == http.MethodPost {
		var request map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		message, ok := request["message"].(string)
		if !ok || message == "" {
			response := JSONResponse{Status: "fail", Message: "Invalid JSON message"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := JSONResponse{Status: "success", Message: "Data successfully received"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Ответ на неподдерживаемые методы
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
