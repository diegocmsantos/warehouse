package whhttp

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHandleIndex(t *testing.T) {

	cases := []struct {
		description string
		path        string
		statusCode  int
		body        string
	}{
		{
			description: "successful",
			path:        "../../testdata/test_file.csv",
			statusCode:  http.StatusOK,
			body: `
            {
				"children": {
					"category_1": {
						"children": {
							"category_2": {
								"children": {
									"item_1": {
										"item": true
									}
								}
							},
							"category_3": {
								"children": {
									"item_2": {
										"item": true
									}
								}
							}
						}
					}
				}
            }
        `,
		},
		{
			description: "empty category",
			path:        "../../testdata/specialcase.csv",
			statusCode:  http.StatusOK,
			body: `
            {
				"children": {
					"category_1": {
						"children": {
							"item_1": {
								"item": true
							}
						}
					},
					"category_2": {
						"children": {
							"category_3": {
								"children": {
									"item_2": {
										"item": true
									}
								}
							}
						}
					}
				}
            }
            `,
		},
		{
			description: "invalid list",
			path:        "../../testdata/wrongcase.csv",
			statusCode:  http.StatusBadRequest,
			body: `
            {
                "error": "could not get elements: line cannot have categories after a blank space"
            }
            `,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			file, err := os.Open(c.path)
			require.NoError(t, err)
			defer file.Close()

			req := httptest.NewRequest("POST", "/", file)
			req.Header.Add("Content-Type", "text/csv")
			rr := httptest.NewRecorder()
			s, err := New(ServerOptions{
				Addr:        "8080",
				Host:        "localhost",
				Port:        8080,
				ReadTimeout: 5 * time.Second,
				Logger:      zap.NewNop(),
			})
			require.NoError(t, err)
			handler := s.HandleIndex()
			handler.ServeHTTP(rr, req)
			require.Equal(t, c.statusCode, rr.Code)
			require.JSONEq(t, c.body, rr.Body.String())
		})
	}
}
