package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"url-crawler-backend/internal/api"
	"url-crawler-backend/internal/db"
	"url-crawler-backend/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestEnvironment() (*echo.Echo, *gorm.DB) {
	e := echo.New()

	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	testDB.AutoMigrate(&model.URL{}, &model.User{})

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	testUser := model.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}
	testDB.Create(&testUser)

	return e, testDB
}

func TestLoginFlow(t *testing.T) {
	e, testDB := setupTestEnvironment()

	// Temporarily replace the global DB
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Test login with valid credentials
	loginData := map[string]interface{}{
		"username": "testuser",
		"password": "testpassword",
	}

	jsonBody, _ := json.Marshal(loginData)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := api.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
	assert.NotEmpty(t, response["token"])
}

func TestURLManagementFlow(t *testing.T) {
	e, testDB := setupTestEnvironment()

	// Temporarily replace the global DB
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Test adding a URL
	urlData := map[string]interface{}{
		"url": "https://example.com",
	}

	jsonBody, _ := json.Marshal(urlData)
	req := httptest.NewRequest(http.MethodPost, "/api/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := api.AddURL(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var urlResponse model.URL
	err = json.Unmarshal(rec.Body.Bytes(), &urlResponse)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", urlResponse.URL)
	assert.Equal(t, "queued", urlResponse.Status)

	// Test getting URLs
	req = httptest.NewRequest(http.MethodGet, "/api/urls", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = api.GetURLs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var urls []model.URL
	err = json.Unmarshal(rec.Body.Bytes(), &urls)
	assert.NoError(t, err)
	assert.Len(t, urls, 1)
	assert.Equal(t, "https://example.com", urls[0].URL)
}

func TestCrawlFlow(t *testing.T) {
	e, testDB := setupTestEnvironment()

	// Temporarily replace the global DB
	originalDB := db.DB
	db.DB = testDB
	defer func() { db.DB = originalDB }()

	// Create a URL first
	url := model.URL{
		URL:    "https://example.com",
		Status: "queued",
	}
	testDB.Create(&url)

	// Test starting crawl
	jsonBody, _ := json.Marshal(map[string]interface{}{"ids": []uint{url.ID}})
	req := httptest.NewRequest(http.MethodPost, "/api/urls/crawl", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := api.StartBulkCrawl(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Bulk crawl started", response["message"])
}
