package handlers

import (
	"Contest/internal/services"
	"encoding/json"
	"net/http"
)

type CompileRequest struct {
	Code string `json:"code"`
}

type CompileResponse struct {
	Output string `json:"output"`
}

func CompileCPP(compileService services.ICompileService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CompileRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := compileService.CompileCPP(request.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := json.Marshal(&CompileResponse{Output: result})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
