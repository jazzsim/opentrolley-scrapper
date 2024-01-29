package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

type PageRequest struct {
	StartId  int    `json:"start_id"`
	Genre    string `json:"genre"`
	Url      string `json:"url"`
	FileName string `json:"file_name"`
}

type AllLinks struct {
	Links []string
}

type BookResponse struct {
	Id                 int     `json:"id"`
	ImageUrl           string  `json:"imageUrl"`
	Title              string  `json:"title"`
	Genre              string  `json:"genre"`
	BindingDescription string  `json:"bindingDescription"`
	Language           string  `json:"language"`
	Description        string  `json:"description"`
	Price              float32 `json:"price"`
	DiscountPrice      float32 `json:"discountPrice"`
	ISBN               string  `json:"isbn"`
	Publisher          string  `json:"publisher"`
	PublicationDate    string  `json:"publicationDate"`
	Pages              int     `json:"pages"`
	PostedBy           int     `json:"postedBy"`
}

type BookAuthor struct {
	BookId     int    `json:"bookId"`
	AuthorName string `json:"authorName"`
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.POST("/", scrape)

	router.Run("localhost:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func scrape(c *gin.Context) {
	var response []BookResponse
	pr := PageRequest{}

	if err := c.ShouldBindJSON(&pr); err != nil {
		// Handle error (e.g., invalid JSON format)
		c.JSON(http.StatusBadRequest, gin.H{"400 error": err.Error()})
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new browser session
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	allLinks := pr.getLinks(ctx)

	var authorSlice = []BookAuthor{}
	for index, link := range allLinks.Links {
		// wait for 5 seconds to avoid rate limit
		time.Sleep(5 * time.Second)
		fmt.Println(link)
		response = append(response, getDetails(c, pr.StartId+index, link, &authorSlice, pr.Genre))
	}

	// write authors into json
	authorsData, err := json.MarshalIndent(authorSlice, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	writeFile(pr.FileName, authorsData)

	c.IndentedJSON(http.StatusOK, response)
}

func (pr *PageRequest) getLinks(ctx context.Context) (response AllLinks) {

	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Navigate(pr.Url),
		chromedp.Nodes(`.book-title>a`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		fmt.Println("Images not found")
	}

	for _, node := range nodes {
		fmt.Println(`node =`, node.AttributeValue("href"))
		response.Links = append(response.Links, "https://opentrolley.com.my/"+node.AttributeValue("href"))
	}
	return response
}

func getDetails(c *gin.Context, id int, url string, authorSlice *[]BookAuthor, genre string) (response BookResponse) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
		chromedp.Flag("enable-automation", false),
		// chromedp.Flag("headless", false),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a new browser session
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var ok bool
	var oriPrice string
	var discPrice string
	var language string
	var desc string
	var pages string
	var authorNodes = []string{"", "", ""}

	fmt.Println("url = ", url)
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.AttributeValue(`.book-cover>div>img`, "src", &response.ImageUrl, &ok),
		chromedp.InnerHTML(`.book-title>span`, &response.Title),
		chromedp.InnerHTML(`.book-author>a:nth-child(1)>span`, &authorNodes[0]),
		chromedp.InnerHTML(`.book-author>a:nth-child(3)>span`, &authorNodes[1]),
		chromedp.InnerHTML(`.book-author>a:nth-child(5)>span`, &authorNodes[2]),
		chromedp.InnerHTML(`.book-synopsys`, &desc),
		chromedp.InnerHTML(".price-before-disc>span:nth-child(2)", &oriPrice),
		chromedp.InnerHTML(".price-after-disc>span:nth-child(2)", &discPrice),
		chromedp.InnerHTML(`#ctl00_ContentPlaceHolder1_lblBindingDescription`, &response.BindingDescription),
		chromedp.InnerHTML(`#ctl00_ContentPlaceHolder1_lblLanguage`, &language),
		chromedp.InnerHTML(`.additional-info>div:nth-child(3)>.value>span:nth-child(1)`, &response.ISBN),
		chromedp.InnerHTML(`.additional-info>div:nth-child(5)>.value>a>span`, &response.Publisher),
		chromedp.InnerHTML(`.additional-info>div:nth-child(6)>.value>span:nth-child(1)`, &response.PublicationDate),
		chromedp.InnerHTML(`.additional-info>div:nth-child(7)>.value>span`, &pages),
	)

	if err != nil {
		fmt.Println("err found ", err)
	}

	response.Id = id
	response.Genre = genre
	response.Language = strings.Trim(language[3:], " ")
	response.Publisher = strings.Trim(response.Publisher, " ")
	pagesInt, _ := strconv.ParseInt(pages, 10, 32)
	response.Pages = int(pagesInt)

	// assign original price
	oriPriceFloat, _ := strconv.ParseFloat(oriPrice, 32)
	response.Price = float32(math.Ceil(oriPriceFloat*100) / 100)
	/* at the time of development, there is discount on every products,
	randomly add discount price for variety purposes
	*/
	// random discount event
	randomNumber := rand.Intn(2) + 1
	if randomNumber == 1 {
		discPriceFloat, _ := strconv.ParseFloat(discPrice, 32)
		response.DiscountPrice = float32(math.Ceil(discPriceFloat*100) / 100)
	}

	// handle authors
	// remove "By "
	authorNodes[0] = authorNodes[0][3:]
	authorList := removeEmptyAuthors(authorNodes)
	for _, author := range authorList {
		*authorSlice = append(*authorSlice,
			BookAuthor{
				BookId:     id,
				AuthorName: author,
			},
		)
	}

	// handle description
	response.Description = removeHTMLTags(desc)

	// handle date
	response.PublicationDate = formatDate(response.PublicationDate)

	// randomly assign to users
	randomId := rand.Intn(30) + 1
	response.PostedBy = randomId

	return response
}

func removeHTMLTags(input string) string {
	// Define the regular expression pattern for HTML tags
	htmlTagRegex := regexp.MustCompile("<[^>]+>")

	// Replace all occurrences of HTML tags with an empty string
	textWithoutTags := htmlTagRegex.ReplaceAllString(input, "")

	// Remove newline characters
	result := strings.ReplaceAll(textWithoutTags, "\n", "")

	return strings.Trim(result, " ")
}

func removeEmptyAuthors(s []string) []string {
	var result []string

	for _, author := range s {
		if author != "" {
			result = append(result, author)
		}
	}

	for index, name := range result {
		name = authorName(name)
		result[index] = name
	}

	return result
}

func authorName(name string) string {
	splitName := strings.Split(name, ",")
	name = splitName[0] + splitName[1]
	return name
}

func formatDate(date string) string {
	// Parse the input date string
	inputDate, err := time.Parse("02/01/2006", date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return date
	}

	// Format the date in "YYYY-MM-DD" format
	return inputDate.Format("2006-01-02")
}

func writeFile(fileName string, jsonData []byte) {
	// Write JSON data to a file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
