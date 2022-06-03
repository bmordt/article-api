package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bmordt/article-api/src/database"
	"github.com/bmordt/article-api/src/middleware"
	"github.com/bmordt/article-api/src/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	expectedDateFormatString = "2006-01-02"

	apiError = middleware.CustomError{}
)

//ArticleService struct for holding important info to the service
//i.e. db client, env variables etc.
type ArticleService struct {
	DBClient database.DBClient
	Logger   *logrus.Entry
}

//NewArticleService everything we need for the article functions
func NewArticleService(dbClient database.DBClient, logger *logrus.Entry) *ArticleService {
	return &ArticleService{
		DBClient: dbClient,
		Logger:   logger,
	}
}

//CreateArticle gets the fields from the req and creates a new article in the DB
func (a *ArticleService) CreateArticle(w http.ResponseWriter, r *http.Request) {
	a.Logger.Infof("Inside CreateArticle function")

	//parse json request
	newReq := &models.CreateArticleReq{}
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		a.Logger.Errorf("CreateArticle :: Error decoding request: %v", err)
		apiError.ApiError(w, http.StatusBadRequest, "Error decoding request")
		return
	}
	a.Logger.Infof("CreateArticle :: Incoming create article request: %+v", newReq)

	tDate, err := time.Parse(expectedDateFormatString, newReq.Date)
	if err != nil {
		a.Logger.Errorf("CreateArticle :: Error decoding request date: %v", err)
		apiError.ApiError(w, http.StatusBadRequest, fmt.Sprintf("Request date is not expected format \"%s\"", expectedDateFormatString))
		return
	}

	//map request to db article object
	newArticle := mapCreateArticleReqToDBArticle(newReq, tDate)

	//store in db
	newID, err := a.DBClient.CreateArticleRow(newArticle.Title, newArticle.Body, newArticle.Date, newArticle.Tags)
	if err != nil {
		a.Logger.Errorf("CreateArticle :: Error storing article %+v : %v", newArticle, err)
		apiError.ApiError(w, http.StatusInternalServerError, "Internal server error storing article")
		return
	}

	newArticle.ID = strconv.Itoa(newID)
	a.Logger.Infof("CreateArticle :: Successfully created new article ID: %d", newID)

	resp := mapToArticleResponse(newArticle)
	middleware.ModelResponse(w, 201, resp)
	return
}

//GetArticle gets the article from DB that belongs to the ID provided in the path parameter
func (a *ArticleService) GetArticle(w http.ResponseWriter, r *http.Request) {
	a.Logger.Infof("Inside GetArticle function")

	//Make sure path params are okay
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		a.Logger.Warnf("GetArticle :: id is not present in the url path %s", r.URL.Path)
		apiError.ApiError(w, http.StatusBadRequest, "id path parameter is not provided")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		a.Logger.Warnf("GetArticle :: id is not a valid integer %s", id)
		apiError.ApiError(w, http.StatusBadRequest, "id path parameter is not valid")
		return
	}

	//Check to see if it exists first
	article, err := a.DBClient.GetArticleRowByID(idInt)
	if err != nil {
		a.Logger.Errorf("GetArticle :: Error getting article %d from DB : %v", idInt, err)
		apiError.ApiError(w, http.StatusInternalServerError, "Internal server error getting article")
		return
	}

	a.Logger.Infof("GetArticle :: Successfully found article: %+v", article)

	resp := mapToArticleResponse(article)
	middleware.ModelResponse(w, 200, resp)
	return
}

//GetArticlesByTagAndDate gets the article from DB that belongs to the ID provided in the path parameter
func (a *ArticleService) GetArticlesByTagAndDate(w http.ResponseWriter, r *http.Request) {
	a.Logger.Infof("Inside GetArticlesByTagAndDate function")

	//Make sure path params are okay
	vars := mux.Vars(r)
	tagName, ok := vars["tagName"]
	if !ok {
		a.Logger.Warnf("GetArticlesByTagAndDate :: tagName is not present in the url path %s", r.URL.Path)
		apiError.ApiError(w, http.StatusBadRequest, "tagName path parameter is not provided")
		return
	}
	date, ok := vars["date"]
	if !ok {
		a.Logger.Warnf("GetArticlesByTagAndDate :: date is not present in the url path %s", r.URL.Path)
		apiError.ApiError(w, http.StatusBadRequest, "date path parameter is not provided")
		return
	}

	//Validate the date
	_, err := time.Parse(expectedDateFormatString, date)
	if err != nil {
		a.Logger.Errorf("CreateArticle :: Error parsing param date: %v", err)
		apiError.ApiError(w, http.StatusBadRequest, fmt.Sprintf("Path parameter date is not in expected format \"%s\"", expectedDateFormatString))
		return
	}

	//Check to see the tag and date have results
	articles, err := a.DBClient.GetArticleRowByTagAndDate(tagName, date)
	if err != nil {
		a.Logger.Errorf("GetArticlesByTagAndDate :: Error getting articles %s %s from DB : %v", tagName, date, err)
		apiError.ApiError(w, http.StatusInternalServerError, "Internal server error getting article")
		return
	}

	//Return the correct info
	resp := mapToTagGroupArticleResp(*articles, tagName)

	a.Logger.Infof("GetArticlesByTagAndDate :: Successfully found articles and mapped to response: %+v", resp)
	middleware.ModelResponse(w, 200, resp)
	return
}

func mapCreateArticleReqToDBArticle(req *models.CreateArticleReq, reqDate time.Time) *models.Article {
	return &models.Article{
		Title: req.Title,
		Body:  req.Body,
		Tags:  req.Tags,
		Date:  reqDate,
	}
}

func mapToArticleResponse(dbArticle *models.Article) *models.ArticleResp {
	return &models.ArticleResp{
		ID:    dbArticle.ID,
		Title: dbArticle.Title,
		Body:  dbArticle.Body,
		Tags:  dbArticle.Tags,
		Date:  dbArticle.Date.Format(expectedDateFormatString),
	}
}

//mapToTagGroupArticleResp returns the desired response
//Distinct article IDs
//Distinct article tags - needs to return the last 10 entered
//Count of articles
func mapToTagGroupArticleResp(articles []models.Article, tagName string) *models.GroupArticleResp {
	//related_tags field contains a list of tags that are on the articles that the current tag is on for the same day
	resp := &models.GroupArticleResp{
		Tag:   tagName,
		Count: len(articles),
	}

	ids := []string{}
	relatedTagsCount := make(map[string]int)

	//loop through the article to get each ID whilst also appending the tagName to a map to get distinct tag values
	//We know the query returns the articles in createdDate desc order. So the newest articles should be first
	for _, article := range articles {
		if len(ids) < 10 {
			ids = append(ids, article.ID)
		}

		for _, tag := range article.Tags {
			relatedTagsCount[tag]++
		}
	}
	resp.Articles = ids

	//For each unique value append to the tag Name array
	relatedTagsNames := make([]string, len(relatedTagsCount))
	i := 0
	for name := range relatedTagsCount {
		relatedTagsNames[i] = name
		i++
	}
	resp.RelatedTags = relatedTagsNames

	return resp
}
