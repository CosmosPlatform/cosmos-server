package healthcheck

import (
	"cosmos-server/pkg/test"
	"github.com/gin-gonic/gin"
	"testing"
)

func setUp() *gin.Engine {
	router := test.NewRouter()
	AddHealthcheckHandler(router.Group("/"))
	return router
}

func TestRouteHealthcheck(t *testing.T) {
	t.Run("healthcheck route - success", healthcheckSuccess)
}

func healthcheckSuccess(t *testing.T) {
	router := setUp()

	request, recorder, err := test.NewHTTPRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	router.ServeHTTP(recorder, request)

	if recorder.Code != 200 {
		t.Errorf("Expected status code 200, got %d", recorder.Code)
	}

	expectedBody := `{"status":"ok"}`
	if recorder.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, recorder.Body.String())
	}
}
