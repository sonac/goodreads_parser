package goodreads_parser

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

type Book struct {
	Title     string
	Author    string
	Rating    Rating
	PosterUrl string
}

type Rating struct {
	Count int64
	Avg   float64
}

func (p *Parser) FindBooks(searchString string) (*[]Book, error) {
	html, err := p.fetch(searchString)
	if err != nil {
		log.Println("error occurred during fetch")
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(html)

	var books []Book
	doc.Find("tr[itemtype='http://schema.org/Book']").Each(func(i int, s *goquery.Selection) {

		book, err := p.parseBook(s)
		if err != nil {
			log.Println("error occurred during book parsing")
			return // skip this book
		}

		books = append(books, *book)
	})
	return &books, nil
}

func (p *Parser) parseBook(s *goquery.Selection) (*Book, error) {
	title := s.Find(".bookTitle").Text()
	author := s.Find(".authorName").First().Text()
	rating := s.Find(".minirating").Text()
	posterUrl, _ := s.Find(".bookCover").Attr("src")
	posterUrl = getPosterUrl(posterUrl)

	parsedRating, err := parseRating(strings.TrimSpace(rating))
	if err != nil {
		log.Println("error occurred during rating parsing")
		return nil, err
	}

	book := Book{
		Title:     strings.TrimSpace(title),
		Author:    author,
		Rating:    parsedRating,
		PosterUrl: posterUrl,
	}
	return &book, nil
}

func (p *Parser) fetch(name string) (*string, error) {
	formattedName := strings.ReplaceAll(name, " ", "+")
	baseUrl := "https://www.goodreads.com/search?utf8=%E2%9C%93&query="
	searchUrl := baseUrl + formattedName
	fmt.Println(searchUrl)
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		log.Println("error occurred during request build")
		return nil, err
	}
	respBody, err := Client.Do(req)
	if err != nil {
		log.Println("error occurred during fetching response")
		return nil, err
	}
	html, err := io.ReadAll(respBody.Body)
	htmlString := string(html)
	return &htmlString, nil
}

// parse string like "4.47 avg rating — 9,852,011 ratings" to Rating struct
func parseRating(s string) (Rating, error) {
	avg, err := strconv.ParseFloat(strings.Split(s, " ")[0], 32)
	if err != nil {
		return Rating{}, err
	}
	count, err := strconv.ParseInt(strings.ReplaceAll(strings.Split(s, " ")[4], ",", ""), 10, 32)
	if err != nil {
		return Rating{}, err
	}
	return Rating{
		Count: count,
		Avg:   roundFloat64(avg),
	}, nil
}

// fetching proper url from minimized version
func getPosterUrl(s string) string {
	bookId := strings.Split(s, "/")[7]
	posterId := strings.Split(strings.Split(s, "/")[8], ".")[0]
	return fmt.Sprintf("https://images-na.ssl-images-amazon.com/images/S/compressed.photo.goodreads.com/books/%s/%s.jpg", bookId, posterId)
}

func roundFloat64(f float64) float64 {
	ratio := math.Pow(10, float64(2))
	return math.Round(f*ratio) / ratio

}
