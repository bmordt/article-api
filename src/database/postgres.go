package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bmordt/article-api/src/models"

	"github.com/lib/pq"

	"github.com/sirupsen/logrus"
)

//DBClient interface for the DB packages
//go:generate moq -out dBClient_mock.go . DBClient
type DBClient interface {
	CreateArticleRow(title, body string, date time.Time, tags []string) (int, error)
	GetArticleRowByID(findID int) (*models.Article, error)
	GetArticleRowByTagAndDate(tag, date string) (*[]models.Article, error)
}

type ArticleDBClient struct {
	DB     *sql.DB
	Logger *logrus.Entry
}

//NewArticleDBClient initiates the connection to the DB
func NewArticleDBClient(dbUser, dbName, password, dbHost, dbPort string, logger *logrus.Entry) *ArticleDBClient {
	connStr := fmt.Sprintf("host=%s port=%s password=%s user=%s dbname=%s sslmode=disable", dbHost, dbPort, password, dbUser, dbName)
	log.Printf("NewArticleDBClient %s\n", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("NewArticleDBClient :: Error opening up connStr %s : %v", connStr, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("NewArticleDBClient :: Error pinging up connStr %s : %v", connStr, err)
	}
	log.Println("NewArticleDBClient connected")

	return &ArticleDBClient{
		DB:     db,
		Logger: logger,
	}
}

//CreateArticleRow inserts new article row
func (d *ArticleDBClient) CreateArticleRow(title, body string, date time.Time, tags []string) (int, error) {
	query := fmt.Sprintf(`INSERT INTO ARTICLES(TITLE, ARTICLE_DATE, BODY, TAGS) VALUES ($1, $2, $3, $4) RETURNING ID`)
	d.Logger.Debugf("CreateArticleRow :: %s", query)

	var temp int
	err := d.DB.QueryRow(query, title, date, body, pq.Array(tags)).Scan(&temp)
	return temp, err
}

//GetArticleRowByID queries db for article by its ID
func (d *ArticleDBClient) GetArticleRowByID(findID int) (*models.Article, error) {
	query := fmt.Sprintf(`SELECT ID, TITLE, ARTICLE_DATE, BODY, TAGS FROM ARTICLES WHERE ID=$1`)
	d.Logger.Infof("GetArticleRowByID :: %s ID: %d", query, findID)

	article := &models.Article{}
	err := d.DB.QueryRow(query, findID).Scan(&article.ID, &article.Title, &article.Date, &article.Body, pq.Array(&article.Tags))
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, nil
		}
		fmt.Printf("GetArticleRowByID :: Error finding db article %v", err)
		return nil, err
	}
	return article, nil
}

//GetArticleRowByID queries db for article by its ID
func (d *ArticleDBClient) GetArticleRowByTagAndDate(tag, date string) (*[]models.Article, error) {
	query := fmt.Sprintf(`SELECT ID, TITLE, ARTICLE_DATE, BODY, TAGS FROM ARTICLES WHERE TAGS && ARRAY[$1] and article_date = $2 order by CREATEDDATE desc`)

	d.Logger.Infof("GetArticleRowByTagAndDate :: %s tag %s date %s", query, tag, date)

	rows, err := d.DB.Query(query, tag, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []models.Article{}
	for rows.Next() {
		article := models.Article{}
		err = rows.Scan(&article.ID, &article.Title, &article.Date, &article.Body, pq.Array(&article.Tags))
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &articles, nil
}
