package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func Test_AddDeleteGet(t *testing.T) {

	// Create a new task with a POST
	taskName := fmt.Sprintf("task %d", rand.Intn(999))
	bodyText := fmt.Sprintf(`{"title":"%s"}`, taskName)

	statusCode, responseBody, err := callService("POST", "/tasks/", bodyText, nil, createTask)
	if err != nil {
		t.Errorf("post failed. %s", err.Error())
	}

	if statusCode != http.StatusOK {
		t.Errorf("failed to post a task. Status: %v", statusCode)
	}

	var newTask = convStringToTask(responseBody)
	if newTask.Title != taskName {
		t.Errorf("response body should contain a task with name %s, got %s", taskName, newTask.Title)
	}

	// get the task
	statusCode, responseBody, err = callService("GET", "/tasks/", "", map[string]string{"id": strconv.FormatInt(newTask.ID, 10)}, getTask)
	if err != nil {
		t.Errorf("get failed. %s", err.Error())
	}
	// make sure the task matches
	var retTask = convStringToTask(responseBody)
	if !reflect.DeepEqual(retTask, newTask) {
		t.Errorf("returned task is not what is expected. got %v, expected %v", retTask, newTask)
	}

	// delete the task
	statusCode, responseBody, err = callService("DELETE", "/tasks/", "", map[string]string{"id": strconv.FormatInt(newTask.ID, 10)}, deleteTask)
	if err != nil {
		t.Errorf("delete failed. %s", err.Error())
	}
	if statusCode != http.StatusOK {
		t.Errorf("delete failed. Expected status code %v, got %v", http.StatusOK, statusCode)
	}
	if strings.Contains(responseBody, newTask.Title) {
		t.Errorf("delete failed. Response body included '%s', the name of the task which should have been deleted.", newTask.Title)
	}

	// get the same task again to make sure it's empty
	statusCode, responseBody, err = callService("GET", "/tasks/", "", map[string]string{"id": strconv.FormatInt(newTask.ID, 10)}, getTask)
	if err != nil {
		t.Errorf("get failed. %s", err.Error())
	}
	// the task should not exist, a blank task should be returned
	var retTask2 = convStringToTask(responseBody)
	if !reflect.DeepEqual(retTask2, Task{}) {
		t.Errorf("returned task is not what is expected. got %v, expected %v", retTask2, Task{})
	}
}

func callService(method string, url string, body string, urlVars map[string]string, f func(w http.ResponseWriter, r *http.Request)) (statusCode int, responseBody string, err error) {
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return
	}
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(f)
	request = mux.SetURLVars(request, urlVars)

	handler.ServeHTTP(recorder, request)
	statusCode = recorder.Code
	responseBody = recorder.Body.String()
	return
}

func convStringToTask(body string) (retTask Task) {
	_ = json.NewDecoder(strings.NewReader(body)).Decode(&retTask)
	return
}
