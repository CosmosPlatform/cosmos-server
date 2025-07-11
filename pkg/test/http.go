package test

import (
	"bytes"
	"cosmos-server/pkg/log"
	"cosmos-server/pkg/server/middleware"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func NewHTTPRequest(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder, error) {
	var request *http.Request
	var err error

	if body != nil {
		requestBody, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		request, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, nil, err
		}
	} else {
		request, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	return request, recorder, nil
}

func NewRouter(logger log.Logger) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.ErrorMiddleware(middleware.NewTranslator()))

	return router
}
