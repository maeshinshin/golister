package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/labstack/echo/v4"

	"github.com/maeshinshin/go-multiapi/internal/database"
	"github.com/maeshinshin/go-multiapi/internal/testutil"
)

type dbInfo struct {
	dB_DATABASE string
	dB_PASSWORD string
	dB_USERNAME string
	dB_HOST     string
	dB_PORT     string
}

func TestMain(m *testing.M) {

	testDBData := &dbInfo{
		dB_DATABASE: "database",
		dB_PASSWORD: "password",
		dB_USERNAME: "user",
	}

	teardown, containerData, err := testutil.MustStartMySQLContainer(testDBData.dB_DATABASE, testDBData.dB_USERNAME, testDBData.dB_PASSWORD)

	if err != nil {
		log.Fatalf("could not start mysql container: %v", err)
	}

	if containerData == nil {
		log.Fatalf("could not get mysql container Data: %v", err)
	}

	tmpDBData := &dbInfo{
		dB_DATABASE: os.Getenv("DB_DATABASE"),
		dB_PASSWORD: os.Getenv("DB_PASSWORD"),
		dB_USERNAME: os.Getenv("DB_USERNAME"),
		dB_HOST:     os.Getenv("DB_PORT"),
		dB_PORT:     os.Getenv("DB_HOST"),
	}

	os.Setenv("DB_DATABASE", testDBData.dB_DATABASE)
	os.Setenv("DB_PASSWORD", testDBData.dB_PASSWORD)
	os.Setenv("DB_USERNAME", testDBData.dB_USERNAME)
	os.Setenv("DB_PORT", containerData.Port())
	os.Setenv("DB_HOST", containerData.Host())

	defer func() {
		os.Setenv("DB_DATABASE", tmpDBData.dB_DATABASE)
		os.Setenv("DB_PASSWORD", tmpDBData.dB_PASSWORD)
		os.Setenv("DB_USERNAME", tmpDBData.dB_USERNAME)
		os.Setenv("DB_PORT", tmpDBData.dB_PORT)
		os.Setenv("DB_HOST", tmpDBData.dB_HOST)
	}()

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown mysql container: %v", err)
	}
}

func TestHandlers(t *testing.T) {

	// Following types define the expected response structure for HelloWorldHandler tests.
	// Fields with concrete types are considered mandatory in the response body.
	// Fields with type `any` are used for optional or unexpected data, allowing for looser type checking during tests.
	type HelloWorldResponseCriteria struct {
		Message      string `json:"message"`
		Message2     string `json:"message2"`
		Message3     string `json:"message3"`
		ErrorMessage any    `json:"errormessage"`
		Error        any    `json:"errors"`
	}

	type HealthResponseCriteria struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Error   any    `json:"error"`
	}

	tests := []struct {
		name                      string
		req                       *http.Request
		rec                       *httptest.ResponseRecorder
		handler                   echo.HandlerFunc
		expectedStatus            int
		requiredAndUnexpectedBody any
	}{
		{
			name: "check HelloWorldHandler",
			req:  httptest.NewRequest(http.MethodGet, "/", nil),
			rec:  httptest.NewRecorder(),
			handler: func(c echo.Context) error {
				h := new(Handlers)
				return h.HelloWorldHandler(c)
			},
			expectedStatus: http.StatusOK,
			requiredAndUnexpectedBody: &HelloWorldResponseCriteria{
				Message:  "Hello World",
				Message2: "Hello World2",
				Message3: "Hello World3",
			},
		},
		{
			name: "check HealthHandler",
			req:  httptest.NewRequest(http.MethodGet, "/", nil),
			rec:  httptest.NewRecorder(),
			handler: func(c echo.Context) error {
				h := &(Handlers{db: database.New()})
				return h.HealthHandler(c)
			},
			expectedStatus: http.StatusOK,
			requiredAndUnexpectedBody: &HealthResponseCriteria{
				Status:  "up",
				Message: "It's healthy",
			},
		},
	}

	e := echo.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c := e.NewContext(tt.req, tt.rec)

			err := tt.handler(c)
			if err != nil {
				t.Fatalf("handler returned error: %v", err)
			}

			if diff := cmp.Diff(tt.expectedStatus, tt.rec.Code); diff != "" {
				t.Errorf("handler returned wrong status code (-want +got):\n%s", diff)
			}

			var actualBody any
			switch body := tt.requiredAndUnexpectedBody.(type) {
			case *HelloWorldResponseCriteria:
				actualBody = &HelloWorldResponseCriteria{}
			case *HealthResponseCriteria:
				actualBody = &HealthResponseCriteria{}
			default:
				t.Fatalf("unknown necessaryBody type: %T", body)
			}
			if err := json.NewDecoder(tt.rec.Body).Decode(actualBody); err != nil {
				t.Fatalf("failed to bind response body: %v", err)
				return
			}

			opts := []cmp.Option{
				cmpopts.EquateEmpty(),
			}

			if diff := cmp.Diff(tt.requiredAndUnexpectedBody, actualBody, opts...); diff != "" {
				t.Errorf("Handler response body is missing necessary elements or unexpected elements (-want +got):\n%s", diff)
			}

		})

	}
}
