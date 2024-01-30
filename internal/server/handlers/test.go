package handlers

import (
	. "Contest/internal/domain"
	"Contest/internal/enums"
	"Contest/internal/services"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type TestRequest struct {
	TaskID   int            `json:"task_id,string"`
	Language enums.Language `json:"language"`
	Code     string         `json:"code"`
}

type TestResponse struct {
	ResultCode  string `json:"result_code"`
	Description string `json:"description"`
	Points      int    `json:"points,string"`
}

type GetTestsResponse struct {
	Tests []Test `json:"tests"`
}

func RunTest(testService services.ITestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request TestRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := testService.RunTest(request.TaskID, request.Language, request.Code)
		if err != nil && !errors.Is(err, services.DeleteFileError) {
			if errors.Is(err, services.UnknownLanguage) {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errors.Is(err, services.DeleteFileError) {
			fmt.Printf("File Not Deleted: %w", err)
		}

		response, err := json.Marshal(&TestResponse{
			ResultCode:  string(result.ResultCode),
			Description: result.Description,
			Points:      result.Points,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func AddTest(testService services.ITestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var test Test
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&test); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err.Error())
			return
		}

		if err := testService.AddTest(test); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteTest(testService services.ITestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err.Error())
			return
		}

		err = testService.DeleteTest(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func UpdateTest(testService services.ITestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err.Error())
			return
		}

		var test Test
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&test); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err.Error())
			return
		}

		err = testService.UpdateTest(id, test)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetTest(testService services.ITestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err.Error())
			return
		}

		test, err := testService.GetTest(id)
		if err != nil {
			fmt.Println(err.Error())

			if errors.Is(err, services.ErrNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(test)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func GetTests(testService services.ITestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tests, err := testService.GetTests()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}

		response, err := json.Marshal(tests)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func GetTestsByTaskID(testService services.ITestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID, err := strconv.Atoi(mux.Vars(r)["task_id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err.Error())
			return
		}

		tests, err := testService.GetTestsByTaskID(taskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err.Error())
			return
		}

		response, err := json.Marshal(tests)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
