package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"theHeirophant/log"
)

func (server Server) doRequest(writer http.ResponseWriter,
	request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		server.doRequestPost(writer, request)

	default:
		writeJson(writer, http.StatusMethodNotAllowed, map[string]interface{}{
			"error": "Method not allowed.",
		})
	}
}

func extractResponseJson(writer http.ResponseWriter,
	request *http.Request) (*map[string]interface{}, bool) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		errorResponse(writer, request, err)
		return nil, false
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		errorResponse(writer, request, err)
		return nil, false
	}

	return &data, true
}

func generateRequest(writer http.ResponseWriter,
	request *http.Request, data *map[string]interface{}) (*http.Request, bool) {
	newRequestMethod, ok := (*data)["method"].(string)
	if !ok {
		errorResponseString(writer, request, "Error parsing request method.")
		return nil, false
	}
	newRequestURL, ok := (*data)["url"].(string)
	if !ok {
		errorResponseString(writer, request, "Error parsing request url.")
		return nil, false
	}
	newRequestBody, ok := (*data)["body"].(string)
	if !ok {
		errorResponseString(writer, request, "Error parsing request body.")
		return nil, false
	}

	newRequest, err := http.NewRequest(
		strings.ToUpper(newRequestMethod),
		newRequestURL,
		bytes.NewReader([]byte(newRequestBody)),
	)
	if err != nil {
		errorResponse(writer, request, err)
		return nil, false
	}
	return newRequest, true
}

func (server Server) doRequestPost(writer http.ResponseWriter,
	request *http.Request) {

	data, ok := extractResponseJson(writer, request)
	if !ok {
		return
	}

	newRequest, ok := generateRequest(writer, request, data)
	if !ok {
		return
	}

	newResponse, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		errorResponse(writer, request, err)
		return
	}

	body, err := io.ReadAll(newResponse.Body)
	if err != nil {
		errorResponse(writer, request, err)
		return
	}
	replacedBody := strings.ReplaceAll(string(body), "\\\"", "\"")
	log.Info(
		"Request made sucessfully",
		"url", newRequest.URL,
		"method", newRequest.Method,
		"headers", newRequest.Header,
		"status_code", newResponse.StatusCode,
		"body", replacedBody,
	)

	writeJson(writer, http.StatusOK, map[string]interface{}{
		"url":         newRequest.URL,
		"method":      newRequest.Method,
		"headers":     newRequest.Header,
		"status_code": newResponse.StatusCode,
		"body":        replacedBody,
	})
}

func genericError(writer http.ResponseWriter,
	_ *http.Request) {
	writeJson(writer, http.StatusInternalServerError, map[string]interface{}{
		"error": "Something went wrong. Try again.",
	})
}

func errorResponse(writer http.ResponseWriter,
	_ *http.Request, err error) {
	writeJson(writer, http.StatusInternalServerError, map[string]interface{}{
		"error": err.Error(),
	})
}

func errorResponseString(writer http.ResponseWriter,
	_ *http.Request, err string) {
	writeJson(writer, http.StatusInternalServerError, map[string]interface{}{
		"error": err,
	})
}
