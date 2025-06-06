package nct

import (
	"diploma/modules/product/model"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type NCTParser struct {
	baseURL string
}

func NewNCTParser(baseURL string) *NCTParser {
	return &NCTParser{
		baseURL: baseURL,
	}
}

// ParseProductByGTIN searches for a product by its GTIN and returns found products.
func (p *NCTParser) ParseProductByGTIN(gtin string) ([]model.Product, error) {
	fmt.Println("Starting NCT parser...")
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("error starting Playwright: %w", err)
	}
	defer pw.Stop()
	fmt.Println("Playwright started successfully")
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("error launching Firefox: %w", err)
	}

	fmt.Println("Browser launched successfully")
	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("error creating page: %w", err)
	}

	// Use the baseURL from the parser
	if _, err := page.Goto(p.baseURL); err != nil {
		return nil, fmt.Errorf("error navigating to URL: %w", err)
	}

	fmt.Println("Navigated to NCT successfully")
	inputSelector := `input[placeholder="Введите наименование или номер штрихкода"]`
	if _, err = page.WaitForSelector(inputSelector); err != nil {
		return nil, fmt.Errorf("search input not found: %w", err)
	}

	if err = page.Fill(inputSelector, gtin); err != nil {
		return nil, fmt.Errorf("error filling search input: %w", err)
	}
	fmt.Println("Search input filled successfully")

	if err = page.Click("button:has-text(\"Поиск\")"); err != nil {
		return nil, fmt.Errorf("error clicking search button: %w", err)
	}

	// Waiting for the table to load completely.
	time.Sleep(5 * time.Second)

	html, err := page.InnerHTML("table")
	if err != nil {
		return nil, fmt.Errorf("error reading table: %w", err)
	}
	return ConvertHtmlToProduct(html), nil
}

// ConvertHtmlToProduct parses the HTML and returns a slice of Product.
func ConvertHtmlToProduct(html string) []model.Product {
	reTr := regexp.MustCompile(`(?s)<tr.*?>(.*?)</tr>`)
	reTd := regexp.MustCompile(`(?s)<td.*?>(.*?)</td>`)
	reTags := regexp.MustCompile(`<.*?>`)

	cleanText := func(html string) string {
		text := reTags.ReplaceAllString(html, "")
		return strings.TrimSpace(text)
	}

	var products []model.Product

	trMatches := reTr.FindAllStringSubmatch(html, -1)
	for _, tr := range trMatches {
		tdMatches := reTd.FindAllStringSubmatch(tr[1], -1)
		if len(tdMatches) < 5 {
			continue
		}
		p := model.Product{
			Name:        cleanText(tdMatches[0][1]),
			Category:    cleanText(tdMatches[3][1]),
			Subcategory: cleanText(tdMatches[4][1]),
		}
		products = append(products, p)
	}

	return products
}
