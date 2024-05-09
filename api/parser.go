package goodreads_parser

import (
	"io"
	"log/slog"
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
	Id            int
	Title         string
	Author        string
	Rating        Rating
	PosterUrl     string
	PublisherYear int16
	Description   string
	PageCount     int16
	Url           string
}

type Rating struct {
	Count int64
	Avg   float64
}

func (p *Parser) FindBooks(searchString string, lim int) (*[]Book, error) {
	formattedName := strings.ReplaceAll(searchString, " ", "+")
	baseUrl := "https://www.goodreads.com/search?utf8=%E2%9C%93&query="
	searchUrl := baseUrl + formattedName
	html, err := fetch(searchUrl)
	if err != nil {
		slog.Error("error occurred during fetch")
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*html))
	if err != nil {
		slog.Error("error occured when parsing html")
		return nil, err
	}

	var books []Book
	iter := 0
	doc.Find("tr[itemtype='http://schema.org/Book']").Each(func(i int, s *goquery.Selection) {
		if iter == lim {
			return
		}
		book, err := parseSearchBook(s)
		if err != nil {
			slog.Error("error occurred during search book parsing", err)
			return // skip this book
		}
		err = book.parseBookDetails()
		if err != nil {
			slog.Error("error occurred during book details parsing", err)
			return // skip this book
		}
		books = append(books, *book)
		iter += 1
	})
	return &books, nil
}

func fetchBook(bookHref string) (*goquery.Document, error) {
	baseUrl := "https://www.goodreads.com"
	fullUrl := baseUrl + bookHref
	html, err := fetch(fullUrl)
	if err != nil {
		slog.Error("error occured during fetch")
		return nil, err
	}

	return goquery.NewDocumentFromReader(strings.NewReader(*html))
}

func (b *Book) parseBookDetails() error {
	d, err := fetchBook(b.Url)
	if err != nil {
		return nil
	}
	// Parsing base
	b.Title = d.Find("h1[data-testid='bookTitle']").Text()
	b.Author = d.Find("a.ContributorLink").First().Text()

	// Parsing rating and count
	ratingText := d.Find("div.RatingStatistics__rating").Text()
	if len(ratingText) > 5 {
		if rating, err := strconv.ParseFloat(ratingText[:5], 64); err == nil {
			b.Rating.Avg = rating
		}
	}
	ratingsCountText := d.Find("span[data-testid='ratingsCount']").Text()
	ratingsCountText = strings.ReplaceAll(ratingsCountText, "ratings", "")
	ratingsCountText = strings.ReplaceAll(ratingsCountText, ",", "")
	ratingsCountText = strings.Split(ratingsCountText, "\u00a0")[0]
	if count, err := strconv.ParseInt(strings.TrimSpace(ratingsCountText), 10, 64); err == nil {
		b.Rating.Count = count
	}

	// Rest of the data
	pubInfoText := d.Find("p[data-testid='publicationInfo']").Text()
	pubInfoText = strings.TrimSpace(strings.Replace(pubInfoText, "First published", "", 1))
	year, _ := strconv.Atoi(strings.Split(pubInfoText, " ")[2])
	b.PublisherYear = int16(year)

	b.Description = d.Find("div.DetailsLayoutRightParagraph__widthConstrained").Find("span.Formatted").Nodes[0].LastChild.Data

	pageCountText := d.Find("p[data-testid='pagesFormat']").Text()
	pageCountText = strings.Split(pageCountText, " ")[0]
	pageCount, _ := strconv.Atoi(pageCountText)
	b.PageCount = int16(pageCount)

	b.PosterUrl, _ = d.Find("img.ResponsiveImage").First().Attr("src")

	return nil
}

func parseSearchBook(s *goquery.Selection) (*Book, error) {
	title := s.Find(".bookTitle").Text()
	url, _ := s.Find(".bookTitle").Attr("href")
	author := s.Find(".authorName").First().Text()

	idText := strings.Split(strings.Split(url, "/")[3], "-")
	var id int
	var err error
	if len(idText) > 0 {
		if id, err = strconv.Atoi(idText[0]); err != nil {
			return nil, err
		}
	}

	book := Book{
		Id:     id,
		Title:  strings.TrimSpace(title),
		Author: author,
		Url:    url,
	}
	return &book, nil
}

func fetch(url string) (*string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("error occurred during request build")
		return nil, err
	}
	respBody, err := Client.Do(req)
	if err != nil {
		slog.Error("error occurred during fetching response")
		return nil, err
	}
	html, err := io.ReadAll(respBody.Body)
	if err != nil {
		slog.Error("error occurred during reading response body")
		return nil, err
	}
	htmlString := string(html)
	return &htmlString, nil
}
