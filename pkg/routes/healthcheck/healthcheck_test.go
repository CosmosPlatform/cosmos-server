package healthcheck

import (
	logMock "cosmos-server/pkg/log/mock"
	"cosmos-server/pkg/test"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
	"testing"
)

type mocks struct {
	loggerMock *logMock.MockLogger
}

func setUp(t *testing.T) (*gin.Engine, *mocks) {
	ctrl := gomock.NewController(t)

	mocks := &mocks{
		loggerMock: logMock.NewMockLogger(ctrl),
	}

	router := test.NewRouter(mocks.loggerMock)
	AddHealthcheckHandler(router.Group("/"))
	return router, mocks
}

func TestRouteHealthcheck(t *testing.T) {
	t.Run("healthcheck route - success", healthcheckSuccess)
}

func healthcheckSuccess(t *testing.T) {
	router, mocks := setUp(t)

	mocks.loggerMock.EXPECT().
		Infow(gomock.Any(), gomock.Any())

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
