package handler

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	modelApi "diploma/modules/product/handler/model"

	"github.com/gin-gonic/gin"
)

// GetMarketAnalytics godoc
// @Summary      Get market analytics
// @Description  Retrieve general market price trends and analytics (mock data)
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        days         query     int     false "Number of days for history (default 30)"
// @Param        category     query     string  false "Category filter"
// @Success      200  {object}  modelApi.MarketAnalyticsResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/product/market/analytics [get]
func (h *CatalogHandler) GetMarketAnalytics(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	category := c.DefaultQuery("category", "")

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	if days > 365 {
		days = 365
	}

	analytics := generateMockMarketAnalytics(days, category)
	c.JSON(http.StatusOK, analytics)
}

func generateMockMarketAnalytics(days int, category string) *modelApi.MarketAnalyticsResponse {
	rand.Seed(time.Now().UnixNano())

	categories := []struct {
		name       string
		avgChange  float64
		volatility float64
	}{
		{"Молочные продукты", 2.5, 15.0},
		{"Мясные продукты", 5.2, 25.0},
		{"Хлебобулочные изделия", 1.8, 10.0},
		{"Овощи и фрукты", 8.7, 35.0},
		{"Крупы и макароны", 3.1, 12.0},
	}

	var topGainers []modelApi.ProductTrend
	var topLosers []modelApi.ProductTrend
	var trendData []modelApi.MarketTrendItem

	// Генерируем данные тренда рынка
	for i := days; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)

		// Базовый индекс с небольшой волатильностью
		baseIndex := 100.0
		trend := float64(days-i) * 0.1             // Небольшой общий тренд
		volatility := (rand.Float64() - 0.5) * 5.0 // ±2.5% волатильность

		marketIndex := baseIndex + trend + volatility

		trendData = append(trendData, modelApi.MarketTrendItem{
			Date:        date,
			MarketIndex: marketIndex,
			Volume:      1000 + rand.Intn(5000), // Объем торгов
		})
	}

	// Генерируем топ товаров по росту/падению
	products := []string{
		"Молоко 3.2%", "Хлеб белый", "Мясо говядина", "Картофель", "Рис",
		"Масло подсолнечное", "Сахар", "Гречка", "Сыр твердый", "Курица",
	}

	for i, product := range products {
		change := (rand.Float64() - 0.5) * 20.0 // ±10% изменение
		currentPrice := 1000 + rand.Intn(3000)

		trend := modelApi.ProductTrend{
			ProductID:    int64(i + 1),
			ProductName:  product,
			CurrentPrice: currentPrice,
			PriceChange:  change,
			Volume:       100 + rand.Intn(900),
		}

		if change > 0 {
			topGainers = append(topGainers, trend)
		} else {
			topLosers = append(topLosers, trend)
		}
	}

	// Ограничиваем до топ 5
	if len(topGainers) > 5 {
		topGainers = topGainers[:5]
	}
	if len(topLosers) > 5 {
		topLosers = topLosers[:5]
	}

	// Общая статистика
	totalProducts := 1000 + rand.Intn(500)
	avgPriceChange := (rand.Float64() - 0.5) * 6.0 // ±3%
	marketVolatility := 5.0 + rand.Float64()*10.0  // 5-15%

	return &modelApi.MarketAnalyticsResponse{
		Period:           days,
		Category:         category,
		TotalProducts:    totalProducts,
		AvgPriceChange:   avgPriceChange,
		MarketVolatility: marketVolatility,
		TopGainers:       topGainers,
		TopLosers:        topLosers,
		TrendData:        trendData,
		Categories:       generateCategoryAnalytics(categories),
	}
}

func generateCategoryAnalytics(categories []struct {
	name       string
	avgChange  float64
	volatility float64
}) []modelApi.CategoryAnalytics {
	var result []modelApi.CategoryAnalytics

	for _, cat := range categories {
		// Добавляем небольшую случайность к базовым значениям
		change := cat.avgChange + (rand.Float64()-0.5)*2.0
		volatility := cat.volatility + (rand.Float64()-0.5)*5.0

		result = append(result, modelApi.CategoryAnalytics{
			CategoryName:    cat.name,
			ProductCount:    50 + rand.Intn(200),
			AvgPriceChange:  change,
			Volatility:      volatility,
			TopProduct:      generateRandomProductName(),
			TopProductPrice: 500 + rand.Intn(2000),
		})
	}

	return result
}

func generateRandomProductName() string {
	products := []string{
		"Молоко пастеризованное", "Хлеб ржаной", "Говядина высший сорт",
		"Картофель молодой", "Рис длиннозерный", "Масло сливочное",
		"Сахар-рафинад", "Гречка отборная", "Сыр голландский", "Филе куриное",
	}
	return products[rand.Intn(len(products))]
}
