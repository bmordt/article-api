package services

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bmordt/article-api/src/database"
	"github.com/bmordt/article-api/src/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	testLogger = newTestLogger()
)

func TestCreateArticle(t *testing.T) {
	t.Run("Given a valid create request, a new article is stored in the DB", func(t *testing.T) {
		dbMock := newDbClientMock(false, false)

		a := NewArticleService(dbMock, testLogger)

		dateString := "2016-09-22"
		expectedDateTime, _ := time.Parse(expectedDateFormatString, dateString)
		testReq := models.CreateArticleReq{
			Title: "latest science shows that potato chips are better for you than sugar",
			Date:  dateString,
			Body:  "some text, potentially containing simple markup about how potato chips are great",
			Tags:  []string{"health", "fitness", "science"},
		}
		testIncomingReq := &http.Request{
			Body: getBody(testReq),
		}
		w := httptest.NewRecorder()

		a.CreateArticle(w, testIncomingReq)

		resp := w.Result()
		t.Run("Response code is 201", func(t *testing.T) {
			assert.Equal(t, 201, resp.StatusCode)
		})
		t.Run("CreateArticleRow was Called once with the correct info", func(t *testing.T) {
			assert.Equal(t, 1, len(dbMock.CreateArticleRowCalls()))
			assert.Equal(t, testReq.Title, dbMock.CreateArticleRowCalls()[0].Title)
			assert.Equal(t, expectedDateTime, dbMock.CreateArticleRowCalls()[0].Date)
			assert.Equal(t, testReq.Body, dbMock.CreateArticleRowCalls()[0].Body)
			assert.Equal(t, testReq.Tags, dbMock.CreateArticleRowCalls()[0].Tags)
		})
	})
	t.Run("Given an invalid create request, the correct resp is returned with 400", func(t *testing.T) {
		dbMock := newDbClientMock(false, false)

		a := NewArticleService(dbMock, testLogger)

		dateString := "2016-09-22-12021"
		testReq := models.CreateArticleReq{
			Title: "latest science shows that potato chips are better for you than sugar",
			Date:  dateString,
			Body:  "some text, potentially containing simple markup about how potato chips are great",
			Tags:  []string{"health", "fitness", "science"},
		}
		testIncomingReq := &http.Request{
			Body: getBody(testReq),
		}
		w := httptest.NewRecorder()

		a.CreateArticle(w, testIncomingReq)

		resp := w.Result()
		t.Run("Response code is 400", func(t *testing.T) {
			assert.Equal(t, 400, resp.StatusCode)
		})
		t.Run("Response contains a message", func(t *testing.T) {
			actualResp := make(map[string]string)
			err := json.Unmarshal(w.Body.Bytes(), &actualResp)
			assert.NoError(t, err)

			assert.Equal(t, "Request date is not expected format \"2006-01-02\"", actualResp["Message"])
		})
		t.Run("CreateArticleRow was not Called", func(t *testing.T) {
			assert.Equal(t, 0, len(dbMock.CreateArticleRowCalls()))
		})
	})
	t.Run("Given a valid create request, with an error during the insert we respond with 500", func(t *testing.T) {
		dbMock := newDbClientMock(true, false)

		a := NewArticleService(dbMock, testLogger)

		dateString := "2016-09-22"
		expectedDateTime, _ := time.Parse(expectedDateFormatString, dateString)
		testReq := models.CreateArticleReq{
			Title: "latest science shows that potato chips are better for you than sugar",
			Date:  dateString,
			Body:  "some text, potentially containing simple markup about how potato chips are great",
			Tags:  []string{"health", "fitness", "science"},
		}
		testIncomingReq := &http.Request{
			Body: getBody(testReq),
		}
		w := httptest.NewRecorder()

		a.CreateArticle(w, testIncomingReq)

		resp := w.Result()
		t.Run("Response code is 500", func(t *testing.T) {
			assert.Equal(t, 500, resp.StatusCode)
		})
		t.Run("Response contains a message", func(t *testing.T) {
			actualResp := make(map[string]string)
			err := json.Unmarshal(w.Body.Bytes(), &actualResp)
			assert.NoError(t, err)

			assert.Equal(t, "Internal server error storing article", actualResp["Message"])
		})
		t.Run("CreateArticleRow was Called once with the correct info", func(t *testing.T) {
			assert.Equal(t, 1, len(dbMock.CreateArticleRowCalls()))
			assert.Equal(t, testReq.Title, dbMock.CreateArticleRowCalls()[0].Title)
			assert.Equal(t, expectedDateTime, dbMock.CreateArticleRowCalls()[0].Date)
			assert.Equal(t, testReq.Body, dbMock.CreateArticleRowCalls()[0].Body)
			assert.Equal(t, testReq.Tags, dbMock.CreateArticleRowCalls()[0].Tags)
		})
	})
}

func TestGetArticle(t *testing.T) {
	testIDInt := 111111
	testIDString := strconv.Itoa(testIDInt)
	t.Run("Given a valid get request, an article is returned from the DB", func(t *testing.T) {
		dbMock := newDbClientMock(false, false)

		a := NewArticleService(dbMock, testLogger)

		testIncomingReq := &http.Request{}
		pathVars := make(map[string]string)
		pathVars["id"] = testIDString
		testIncomingReq = mux.SetURLVars(testIncomingReq, pathVars)
		w := httptest.NewRecorder()

		a.GetArticle(w, testIncomingReq)

		resp := w.Result()
		t.Run("Response code is 200", func(t *testing.T) {
			assert.Equal(t, 200, resp.StatusCode)
		})
		t.Run("GetArticleRowByID was Called once with the correct info", func(t *testing.T) {
			assert.Equal(t, 1, len(dbMock.GetArticleRowByIDCalls()))
			assert.Equal(t, testIDInt, dbMock.GetArticleRowByIDCalls()[0].FindID)
		})
	})
	t.Run("Given an invalid get request, 400 and a message is returned", func(t *testing.T) {
		dbMock := newDbClientMock(false, false)

		a := NewArticleService(dbMock, testLogger)

		testIncomingReq := &http.Request{
			URL: &url.URL{
				Path: "blah",
			},
		}
		pathVars := make(map[string]string)
		pathVars["abdcassc"] = testIDString
		testIncomingReq = mux.SetURLVars(testIncomingReq, pathVars)
		w := httptest.NewRecorder()

		a.GetArticle(w, testIncomingReq)

		resp := w.Result()
		t.Run("Response code is 400", func(t *testing.T) {
			assert.Equal(t, 400, resp.StatusCode)
		})
		t.Run("Response contains a message", func(t *testing.T) {
			actualResp := make(map[string]string)
			err := json.Unmarshal(w.Body.Bytes(), &actualResp)
			assert.NoError(t, err)

			assert.Equal(t, "id path parameter is not provided", actualResp["Message"])
		})
		t.Run("GetArticleRowByID was not called", func(t *testing.T) {
			assert.Equal(t, 0, len(dbMock.GetArticleRowByIDCalls()))
		})
	})
	t.Run("Given a valid get request, with an error during the get from DB we respond with 500", func(t *testing.T) {
		dbMock := newDbClientMock(false, true)

		a := NewArticleService(dbMock, testLogger)

		testIncomingReq := &http.Request{
			URL: &url.URL{
				Path: "blah",
			},
		}
		pathVars := make(map[string]string)
		pathVars["id"] = testIDString
		testIncomingReq = mux.SetURLVars(testIncomingReq, pathVars)
		w := httptest.NewRecorder()

		a.GetArticle(w, testIncomingReq)

		resp := w.Result()
		t.Run("Response code is 500", func(t *testing.T) {
			assert.Equal(t, 500, resp.StatusCode)
		})
		t.Run("Response contains a message", func(t *testing.T) {
			actualResp := make(map[string]string)
			err := json.Unmarshal(w.Body.Bytes(), &actualResp)
			assert.NoError(t, err)

			assert.Equal(t, "Internal server error getting article", actualResp["Message"])
		})
		t.Run("GetArticleRowByID was Called once with the correct info", func(t *testing.T) {
			assert.Equal(t, 1, len(dbMock.GetArticleRowByIDCalls()))
			assert.Equal(t, testIDInt, dbMock.GetArticleRowByIDCalls()[0].FindID)
		})
	})
}

func newDbClientMock(createErr, getErr bool) *database.DBClientMock {
	return &database.DBClientMock{
		CreateArticleRowFunc: func(title, body string, date time.Time, tags []string) (int, error) {
			if createErr {
				return 0, errors.New("Create Error")
			}
			return 1, nil
		},
		GetArticleRowByIDFunc: func(findID int) (*models.Article, error) {
			if getErr {
				return &models.Article{}, errors.New("Get Error")
			}
			return &models.Article{
				ID: "1",
			}, nil
		},
	}
}

func newTestLogger() *logrus.Entry {
	testLogger := logrus.New()
	return testLogger.WithFields(logrus.Fields{})
}

func getBody(testReq interface{}) io.ReadCloser {
	testJsonBody, _ := json.Marshal(testReq)

	testBody := strings.NewReader(string(testJsonBody))
	return io.NopCloser(testBody)
}
