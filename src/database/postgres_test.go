package database

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	DBUSER     string
	DBNAME     string
	DBPASSWORD string
	DBHOST     string
	DBPORT     string

	expectedDateFormatString = "2006-01-02"
)

//These tests are assuming being run on a clean DB. i.e. first test returns id 1
func TestPostgresDBClient(t *testing.T) {
	initDBEnvVars()
	testLogger := newTestLogger()

	idsToDelete := []int{}
	t.Run("CreateArticleRow", func(t *testing.T) {
		var testID int
		testTitle := "testTitle"
		testBody := "testBody"
		testDate, _ := time.Parse(expectedDateFormatString, "1991-01-01")
		testTags := []string{"TestTag1", "TestTag2", "TestTag3"}

		dbClient := NewArticleDBClient(DBUSER, DBNAME, DBPASSWORD, DBHOST, DBPORT, testLogger)
		t.Run("Given valid input a row gets created and the row ID returned without errors", func(t *testing.T) {

			resultRow, err := dbClient.CreateArticleRow(testTitle, testBody, testDate, testTags)

			t.Run("No error occured", func(t *testing.T) {
				assert.NoError(t, err)
			})

			t.Run("The data returned is correct", func(t *testing.T) {
				assert.Equal(t, 1, resultRow)
			})

			testID = resultRow
			idsToDelete = append(idsToDelete, testID)
		})
		t.Run("Given valid id an article can be returned without errors", func(t *testing.T) {
			resultArticle, err := dbClient.GetArticleRowByID(testID)

			t.Run("No error occured", func(t *testing.T) {
				assert.NoError(t, err)
			})

			t.Run("The data returned is correct", func(t *testing.T) {
				assert.Equal(t, testTitle, resultArticle.Title)
				assert.Equal(t, testBody, resultArticle.Body)
				assert.Equal(t, "1991-01-01", resultArticle.Date.Format(expectedDateFormatString))
				assert.Equal(t, testTags, resultArticle.Tags)
			})
		})
		t.Run("Given valid tag and date the correct stats are returned without errors", func(t *testing.T) {
			resultArticles, err := dbClient.GetArticleRowByTagAndDate("TestTag1", testDate.Format(expectedDateFormatString))

			t.Run("No error occured", func(t *testing.T) {
				assert.NoError(t, err)
			})

			t.Run("The data returned is correct", func(t *testing.T) {
				assert.Equal(t, 1, len(*resultArticles))
			})
		})
		t.Run("Given valid id the artcile can be deleted without errors", func(t *testing.T) {
			for _, i := range idsToDelete {
				err := dbClient.DeleteArticleByID(i)

				t.Run(fmt.Sprintf("No error occured for id %d", i), func(t *testing.T) {
					assert.NoError(t, err)
				})
			}
		})
	})

}

func newTestLogger() *logrus.Entry {
	testLogger := logrus.New()
	return testLogger.WithFields(logrus.Fields{})
}

func initDBEnvVars() {
	DBUSER = os.Getenv("DBUSER")
	if DBUSER == "" {
		DBUSER = "postgres"
	}
	DBNAME = os.Getenv("DBNAME")
	if DBNAME == "" {
		DBNAME = "nine"
	}
	DBPASSWORD = os.Getenv("DBPASSWORD")
	if DBPASSWORD == "" {
		DBPASSWORD = "12345"
	}
	DBHOST = os.Getenv("DBHOST")
	if DBHOST == "" {
		DBHOST = "localhost"
	}
	DBPORT = os.Getenv("DBPORT")
	if DBPORT == "" {
		DBPORT = "5432"
	}
}
