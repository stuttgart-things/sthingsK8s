/*
Copyright © 2025 PATRICK HERMANN patrick.hermann@sva.de
*/

package k8s

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestFetchYAML tests the FetchYAML function
func TestFetchYAML(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		statusCode     int
		responseBody   string
		expectedResult string
		expectedError  string
	}{
		{
			name:           "successful fetch",
			url:            "/test/success",
			statusCode:     http.StatusOK,
			responseBody:   "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: redis-deployment\n  labels:\n    app: redis\n",
			expectedResult: "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: redis-deployment\n  labels:\n    app: redis\n",
			expectedError:  "",
		},
		{
			name:           "http error (404)",
			url:            "/test/404",
			statusCode:     http.StatusNotFound,
			responseBody:   "",
			expectedResult: "",
			expectedError:  "FAILED TO FETCH YAML, STATUS CODE: 404",
		},
		{
			name:           "http error (500)",
			url:            "/test/500",
			statusCode:     http.StatusInternalServerError,
			responseBody:   "",
			expectedResult: "",
			expectedError:  "FAILED TO FETCH YAML, STATUS CODE: 500",
		},
		{
			name:           "network error",
			url:            "/test/network-error",
			statusCode:     http.StatusOK,
			responseBody:   "",
			expectedResult: "",
			expectedError:  "FAILED TO FETCH YAML: net/http: request canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new test server with a custom handler for this test
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Respond based on the test case's status code and body
				w.WriteHeader(tt.statusCode)
				_, err := w.Write([]byte(tt.responseBody))
				if err != nil {
					t.Fatalf("Failed to write response: %v", err)
				}
			}))
			defer ts.Close()

			// Call FetchYAML with the test server URL
			result, err := FetchYAML(ts.URL + tt.url)

			// Check for the expected result and error
			if result != tt.expectedResult {
				t.Errorf("Expected result: %s, got: %s", tt.expectedResult, result)
			}

			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("Expected error: %s, got: %s", tt.expectedError, err.Error())
			}
		})
	}
}
