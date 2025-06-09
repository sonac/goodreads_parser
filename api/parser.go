package goodreads_parser

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var Client HTTPClient

func init() {
	Client = &http.Client{}
}

// Parser provides methods to search and fetch book information from Goodreads
type Parser struct {
	client HTTPClient
}

// NewParser creates and returns a new Parser instance
func NewParser() *Parser {
	return &Parser{
		client: Client,
	}
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

// FindBooks searches for books on Goodreads and returns a limited number of results
// searchString: the query to search for
// limit: maximum number of books to return
func (p *Parser) FindBooks(searchString string, limit int) ([]Book, error) {
	if searchString == "" {
		return nil, fmt.Errorf("search string cannot be empty")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	formattedName := strings.ReplaceAll(searchString, " ", "+")
	baseUrl := "https://www.goodreads.com/search?utf8=%E2%9C%93&query="
	searchUrl := baseUrl + formattedName

	html, err := p.fetch(searchUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search results: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	books := make([]Book, 0, limit)
	bookCount := 0

	selector := doc.Find("tr[itemtype='http://schema.org/Book']")
	slog.Info("Found book elements", "count", selector.Length())

	selector.Each(func(i int, s *goquery.Selection) {
		if bookCount >= limit {
			return
		}

		book, err := parseSearchBook(s)
		if err != nil {
			slog.Error("Failed to parse search book", "index", i, "error", err)
			return // skip this book
		}

		err = p.fetchBookDetails(book)
		if err != nil {
			slog.Error("Failed to fetch book details", "id", book.Id, "title", book.Title, "error", err)
			return // skip this book
		}

		books = append(books, *book)
		bookCount++
	})

	return books, nil
}

func (p *Parser) fetchBookDocument(bookHref string) (*goquery.Document, error) {
	if !strings.HasPrefix(bookHref, "/") {
		return nil, fmt.Errorf("invalid book URL: %s", bookHref)
	}

	baseUrl := "https://www.goodreads.com"
	fullUrl := baseUrl + bookHref

	html, err := p.fetch(fullUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book page: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse book page HTML: %w", err)
	}

	return doc, nil
}

func (p *Parser) fetchBookDetails(b *Book) error {
	if b == nil {
		return fmt.Errorf("book cannot be nil")
	}

	if b.Url == "" {
		return fmt.Errorf("book URL cannot be empty")
	}

	doc, err := p.fetchBookDocument(b.Url)
	if err != nil {
		return err
	}

	titleElem := doc.Find("h1[data-testid='bookTitle']")
	if titleElem.Length() > 0 {
		b.Title = strings.TrimSpace(titleElem.Text())
	}

	authorElem := doc.Find("a.ContributorLink").First()
	if authorElem.Length() > 0 {
		b.Author = strings.TrimSpace(authorElem.Text())
	}

	ratingText := doc.Find("div.RatingStatistics__rating").Text()
	if len(ratingText) > 0 {
		ratingStr := ratingText
		if len(ratingText) > 5 {
			ratingStr = ratingText[:5]
		}

		if rating, err := strconv.ParseFloat(strings.TrimSpace(ratingStr), 64); err == nil {
			b.Rating.Avg = rating
		} else {
			slog.Debug("Failed to parse rating", "text", ratingStr, "error", err)
		}
	}

	ratingsCountElem := doc.Find("span[data-testid='ratingsCount']")
	if ratingsCountElem.Length() > 0 {
		ratingsCountText := ratingsCountElem.Text()
		ratingsCountText = strings.ReplaceAll(ratingsCountText, "ratings", "")
		ratingsCountText = strings.ReplaceAll(ratingsCountText, ",", "")
		parts := strings.Split(ratingsCountText, "\u00a0")
		if len(parts) > 0 {
			countStr := strings.TrimSpace(parts[0])
			if count, err := strconv.ParseInt(countStr, 10, 64); err == nil {
				b.Rating.Count = count
			} else {
				slog.Debug("Failed to parse ratings count", "text", countStr, "error", err)
			}
		}
	}

	pubInfoElem := doc.Find("p[data-testid='publicationInfo']")
	if pubInfoElem.Length() > 0 {
		pubInfoText := pubInfoElem.Text()
		pubInfoText = strings.TrimSpace(strings.Replace(pubInfoText, "First published", "", 1))
		pubInfoSlice := strings.Fields(pubInfoText)

		if len(pubInfoSlice) > 2 {
			yearStr := pubInfoSlice[2]
			if year, err := strconv.Atoi(yearStr); err == nil {
				b.PublisherYear = int16(year)
			} else {
				slog.Debug("Failed to parse publication year", "text", yearStr, "error", err)
			}
		}
	}

	descElem := doc.Find("div.DetailsLayoutRightParagraph__widthConstrained span.Formatted")
	if descElem.Length() > 0 {
		b.Description = strings.TrimSpace(descElem.Text())
	} else {
		slog.Debug("Description not found", "book", b.Title)
	}

	pageCountElem := doc.Find("p[data-testid='pagesFormat']")
	if pageCountElem.Length() > 0 {
		pageCountText := pageCountElem.Text()
		parts := strings.Fields(pageCountText)
		if len(parts) > 0 {
			if pageCount, err := strconv.Atoi(parts[0]); err == nil {
				b.PageCount = int16(pageCount)
			}
		}
	}

	posterElem := doc.Find("img.ResponsiveImage").First()
	if posterUrl, exists := posterElem.Attr("src"); exists {
		b.PosterUrl = posterUrl
	}

	return nil
}

func parseSearchBook(s *goquery.Selection) (*Book, error) {
	if s == nil {
		return nil, fmt.Errorf("selection cannot be nil")
	}

	titleElem := s.Find(".bookTitle")
	if titleElem.Length() == 0 {
		return nil, fmt.Errorf("title element not found")
	}
	title := titleElem.Text()

	url, exists := titleElem.Attr("href")
	if !exists || url == "" {
		return nil, fmt.Errorf("book URL not found")
	}

	authorElem := s.Find(".authorName").First()
	if authorElem.Length() == 0 {
		return nil, fmt.Errorf("author element not found")
	}
	author := authorElem.Text()

	urlParts := strings.Split(url, "/")
	if len(urlParts) < 4 {
		return nil, fmt.Errorf("invalid URL format: %s", url)
	}

	idText := urlParts[3]
	idDigits := trimAfterNonDigit(idText)
	if idDigits == "" {
		return nil, fmt.Errorf("book ID not found in URL: %s", url)
	}

	id, err := strconv.Atoi(idDigits)
	if err != nil {
		return nil, fmt.Errorf("failed to parse book ID: %w", err)
	}

	book := Book{
		Id:     id,
		Title:  strings.TrimSpace(title),
		Author: strings.TrimSpace(author),
		Url:    url,
	}
	return &book, nil
}

func (p *Parser) fetch(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close() // Ensure we close the response body to prevent resource leaks

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(html), nil
}

// trimAfterNonDigit returns the longest prefix of s that contains only digits
func trimAfterNonDigit(s string) string {
	if s == "" {
		return ""
	}

	for i, r := range s {
		if !unicode.IsDigit(r) {
			return s[:i]
		}
	}
	return s
}
