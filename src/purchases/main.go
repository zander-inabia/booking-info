package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

const (
	defaultPort = "9081"
)

var (
	log    *logrus.Logger
	apikey string
)

type GoodReaderResponse struct {
	XMLName xml.Name `xml:"GoodreadsResponse"`
	Book    []struct {
		XMLName xml.Name `xml:"book"`
		Id      string   `xml:"id"`
		URL     string   `xml:"url"`
	} `xml:"book"`
}

func init() {
	log = logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
}
func mustMapEnv(target *string, envKey string) {
	v := os.Getenv(envKey)
	if v == "" {
		panic(fmt.Sprintf("environment variable %q not set", envKey))
	}
	*target = v
}

func getBookStoreLink(isbn string) (interface{}, error) {
	request := resty.New().R()
	url := fmt.Sprintf("https://www.goodreads.com/book/isbn/%s?key=%s&format=xml", isbn, apikey)

	resp, err := request.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() >= 400 {
		return "", fmt.Errorf("Error while fetching project. %s", resp.Body())
	}
	readerResp := GoodReaderResponse{}
	err = xml.Unmarshal(resp.Body(), &readerResp)

	if err != nil {
		return "", err
	}
	return readerResp.Book[0].URL, nil

}

func main() {
	port := defaultPort
	if value, ok := os.LookupEnv("PORT"); ok {
		port = value
	}

	mustMapEnv(&apikey, "API_KEY")
	r := gin.Default()

	v1 := r.Group("/purchase")
	{

		v1.GET("/goodreads/:bookIsbn", getBookInfo)

	}

	log.Errorln(r.Run(":" + port))
}

func getBookInfo(g *gin.Context) {
	isbn := g.Param("bookIsbn")
	if isbn == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "isbn not found in request"})
		return
	}
	link, err := getBookStoreLink(isbn)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	g.JSON(http.StatusOK, gin.H{"link": link})
	return
}
