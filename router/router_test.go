// router_test.go
package router

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/maeshinshin/go-multiapi/internal/database"
	"github.com/maeshinshin/go-multiapi/internal/util"
)

func TestMain(m *testing.M) {

	database.DbInfo.DB_DATABASE = "database"
	database.DbInfo.DB_PASSWORD = "password"
	database.DbInfo.DB_USERNAME = "user"

	teardown, err := util.MustStartMySQLContainer(database.DbInfo)

	if err != nil {
		log.Fatalf("could not start mysql container: %v", err)
	}

	if database.DbInfo.Db_HOST == "" || database.DbInfo.Db_PORT == "" {
		log.Fatalf("could not get mysql container Data: %v", database.DbInfo)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown mysql container: %v", err)
	}
}

func TestRouter_RegisterRoutes(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		path               string
		expectedStatusCode int
	}{
		{
			name:               "GET /assets/*",
			method:             http.MethodGet,
			path:               "/assets/test.txt",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "GET /web",
			method:             http.MethodGet,
			path:               "/web",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "POST /hello",
			method:             http.MethodPost,
			path:               "/hello",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "GET /",
			method:             http.MethodGet,
			path:               "/",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "GET /health",
			method:             http.MethodGet,
			path:               "/health",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "OPTIONS *",
			method:             http.MethodOptions,
			path:               "/api/resource",
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "PUT *",
			method:             http.MethodPut,
			path:               "/api/resource",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "PATCH *",
			method:             http.MethodPatch,
			path:               "/api/resource",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "DELETE *",
			method:             http.MethodDelete,
			path:               "/api/resource",
			expectedStatusCode: http.StatusNotFound,
		},
	}

	r := NewRouter(database.New())
	handler := r.RegisterRoutes()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(tt.method, tt.path, nil)
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			if diff := cmp.Diff(tt.expectedStatusCode, res.Code); diff != "" {
				t.Errorf("StatusCode mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
