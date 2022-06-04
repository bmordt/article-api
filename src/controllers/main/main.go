package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/bmordt/article-api/src/database"
	"github.com/bmordt/article-api/src/services"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	portNum, dbName, dbUser, dbPassword, dbHost, dbPort string

	logger *logrus.Entry
)

func init() {
	logger = newLogger()
}

func newLogger() *logrus.Entry {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
	})
	log := logrus.NewEntry(l).WithFields(logrus.Fields{
		"Product": "Aritcle-api",
	})
	return log
}

func main() {
	//put here instead of in init function so test can run
	initEnvVariables()

	muxrouter := mux.NewRouter()

	dbClient := database.NewArticleDBClient(dbUser, dbName, dbPassword, dbHost, dbPort, logger)
	articleService := services.NewArticleService(dbClient, logger)

	// -- article routes
	muxrouter.HandleFunc("/articles", articleService.CreateArticle).Methods("POST")
	muxrouter.HandleFunc("/articles/{id}", articleService.GetArticle).Methods("GET")
	muxrouter.HandleFunc("/tags/{tagName}/{date}", articleService.GetArticlesByTagAndDate).Methods("GET")

	//Router end
	logger.Fatalf("%v", http.ListenAndServe(":"+portNum, muxrouter))
}

//initEnvVariables gets required env variables
func initEnvVariables() {
	portNum = GetAPIPort()
	if strings.Compare(portNum, "") == 0 {
		logger.Fatalf("Server Port env \"APIPORT\" variable is not set: %s", portNum)
	}
	dbUser = GetDBUser()
	if strings.Compare(dbUser, "") == 0 {
		logger.Fatalf("Database user env \"DBUSER\" variable is not set: %s", dbUser)
	}
	dbName = GetDBName()
	if strings.Compare(dbName, "") == 0 {
		logger.Fatalf("Database name env \"DBUSER\" variable is not set: %s", dbName)
	}
	dbPassword = GetDBPassword()
	if strings.Compare(dbPassword, "") == 0 {
		logger.Fatalf("Database password env \"DBPASSWORD\" variable is not set: %s", dbPassword)
	}
	dbHost = GetDBHost()
	if strings.Compare(dbHost, "") == 0 {
		logger.Fatalf("Database host env \"DBHOST\" variable is not set: %s", dbHost)
	}
	dbPort = GetDBPort()
	if strings.Compare(dbPort, "") == 0 {
		logger.Fatalf("Database port env \"DBPORT\" variable is not set: %s", dbPort)
	}
}

//GetAPIPort gets the api port from env
func GetAPIPort() string {
	return os.Getenv("APIPORT")
}

//GetDBUser gets the user from env
func GetDBUser() string {
	return os.Getenv("DBUSER")
}

//GetDBName gets the db name from env
func GetDBName() string {
	return os.Getenv("DBNAME")
}

//GetDBPassword gets the db password from env
func GetDBPassword() string {
	return os.Getenv("DBPASSWORD")
}

//GetDBHost gets the db host from env
func GetDBHost() string {
	return os.Getenv("DBHOST")
}

//GetDBPort gets the db port from env
func GetDBPort() string {
	return os.Getenv("DBPORT")
}
