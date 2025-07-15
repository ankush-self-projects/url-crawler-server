package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-crawler-backend/internal/db"
	"url-crawler-backend/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAddURL(t *testing.T) {
	e := echo.New()
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	testDB.AutoMigrate(&model.URL{})

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid URL",
			requestBody: map[string]interface{}{
				"url": "https://example.com",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "Empty URL",
			requestBody: map[string]interface{}{
				"url": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "Missing URL",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/urls", bytes.NewBuffer(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			originalDB := db.DB
			db.DB = testDB
			defer func() { db.DB = originalDB }()

			err := AddURL(c)

			if tt.expectedError {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestGetURLs(t *testing.T) {
	e := echo.New()
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	testDB.AutoMigrate(&model.URL{})

	testURLs := []model.URL{
		{URL: "https://example1.com", Status: "done"},
		{URL: "https://example2.com", Status: "queued"},
	}

	for _, url := range testURLs {
		testDB.Create(&url)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/urls", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	err = GetURLs(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response []model.URL
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}

func TestStartCrawl(t *testing.T) {
	e := echo.New()
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	testDB.AutoMigrate(&model.URL{})

	testURL := model.URL{URL: "https://example.com", Status: "queued"}
	testDB.Create(&testURL)

	tests := []struct {
		name           string
		urlID          string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "Valid URL ID",
			urlID:          "1",
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Invalid URL ID",
			urlID:          "999",
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/urls/"+tt.urlID+"/start", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.urlID)

			originalDB := db.DB
			db.DB = testDB
			defer func() { db.DB = originalDB }()

			err := StartCrawl(c)

			if tt.expectedError {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}
		})
	}
}
