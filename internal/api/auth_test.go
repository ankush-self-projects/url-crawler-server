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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLogin(t *testing.T) {
	e := echo.New()
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	testDB.AutoMigrate(&model.User{})

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	testUser := model.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}
	testDB.Create(&testUser)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid credentials",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "testpassword",
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Invalid password",
			requestBody: map[string]interface{}{
				"username": "testuser",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name: "Non-existent user",
			requestBody: map[string]interface{}{
				"username": "nonexistent",
				"password": "testpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name: "Empty request",
			requestBody: map[string]interface{}{
				"username": "",
				"password": "",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			originalDB := db.DB
			db.DB = testDB
			defer func() { db.DB = originalDB }()

			err := Login(c)

			if tt.expectedError {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedStatus, he.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)

				var response map[string]interface{}
				err = json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
				assert.NotEmpty(t, response["token"])
			}
		})
	}
}
