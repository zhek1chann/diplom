package handler

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	modelApi "diploma/modules/product/handler/model"

	"github.com/gin-gonic/gin"
)

// GetPriceAnalytics godoc
// @Summary      Get price analytics for a product
// @Description  Retrieve price history and analytics for a specific product (mock data)
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id   path      int     true  "Product ID"
// @Param        days         query     int     false "Number of days for history (default 30)"
// @Success      200  {object}  modelApi.PriceAnalyticsResponse
// @Failure      400  {object}  modelApi.ErrorResponse
// @Router       /api/product/{id}/analytics [get]
func (h *CatalogHandler) GetPriceAnalytics(c *gin.Context) {
	productIDStr := c.Param("id")
	daysStr := c.DefaultQuery("days", "30")

	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, modelApi.ErrorResponse{Err: "Invalid product ID"})
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30
	}

	if days > 365 {
		days = 365 // Максимум год
	}

	// Генерируем mock данные
	analytics := generateMockPriceAnalytics(productID, days)

	c.JSON(http.StatusOK, analytics)
}

func generateMockPriceAnalytics(productID int64, days int) *modelApi.PriceAnalyticsResponse {
	rand.Seed(time.Now().UnixNano() + productID) // Seed для консистентности данных

	// Базовые цены для разных поставщиков
	basePrice := 1000 + rand.Intn(5000) // От 1000 до 6000 тенге

	suppliers := []struct {
		id   int64
		name string
	}{
		{1, "ТОО \"Алматы Продукт\""},
		{2, "ОАО \"Караганда Снаб\""},
		{3, "ИП Петров И.И."},
		{4, "ТОО \"Астана Дистрибьюшн\""},
	}

	var priceHistory []modelApi.PriceHistoryItem
	var supplierAnalytics []modelApi.SupplierAnalytics

	minPrice := basePrice
	maxPrice := basePrice
	totalPrice := 0
	count := 0

	// Генерируем историю цен
	for i := days; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)

		// Генерируем изменения цен для случайных поставщиков
		numSuppliers := 1 + rand.Intn(3) // 1-3 поставщика в день
		usedSuppliers := make(map[int]bool)

		for j := 0; j < numSuppliers; j++ {
			supplierIdx := rand.Intn(len(suppliers))
			if usedSuppliers[supplierIdx] {
				continue
			}
			usedSuppliers[supplierIdx] = true

			// Цена с небольшой вариацией
			variation := float64(basePrice) * (0.8 + rand.Float64()*0.4) // ±20%
			price := int(variation)

			// Добавляем тренд (небольшое изменение со временем)
			trend := float64(i-days/2) * (float64(basePrice) * 0.001) // Небольшой тренд
			price += int(trend)

			if price < minPrice {
				minPrice = price
			}
			if price > maxPrice {
				maxPrice = price
			}
			totalPrice += price
			count++

			priceHistory = append(priceHistory, modelApi.PriceHistoryItem{
				Date:     date,
				Price:    price,
				Supplier: suppliers[supplierIdx].name,
			})
		}
	}

	// Генерируем аналитику по поставщикам
	for _, supplier := range suppliers {
		supplierMinPrice := basePrice + rand.Intn(500) - 250
		supplierMaxPrice := supplierMinPrice + rand.Intn(1000)
		supplierAvgPrice := float64(supplierMinPrice+supplierMaxPrice) / 2
		currentPrice := supplierMinPrice + rand.Intn(supplierMaxPrice-supplierMinPrice)

		supplierAnalytics = append(supplierAnalytics, modelApi.SupplierAnalytics{
			SupplierID:   supplier.id,
			SupplierName: supplier.name,
			CurrentPrice: currentPrice,
			MinPrice:     supplierMinPrice,
			MaxPrice:     supplierMaxPrice,
			AvgPrice:     supplierAvgPrice,
			PriceChanges: 5 + rand.Intn(20), // 5-25 изменений
		})
	}

	avgPrice := float64(totalPrice) / float64(count)
	currentPrice := basePrice + rand.Intn(1000) - 500

	return &modelApi.PriceAnalyticsResponse{
		ProductID:    productID,
		ProductName:  generateMockProductName(productID),
		CurrentPrice: currentPrice,
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
		AvgPrice:     avgPrice,
		PriceHistory: priceHistory,
		Suppliers:    supplierAnalytics,
	}
}

func generateMockProductName(productID int64) string {
	products := []string{
		"Молоко «Пискаревский» 3,2% 900 мл",
		"Хлеб «Дарницкий» нарезной 700 г",
		"Масло подсолнечное «Злато» 1 л",
		"Сахар-песок 1 кг",
		"Мука пшеничная высший сорт 2 кг",
		"Рис шлифованный 1 кг",
		"Гречка ядрица 800 г",
		"Макароны «Барилла» спагетти 500 г",
		"Колбаса вареная «Докторская» 400 г",
		"Сыр «Российский» 45% 200 г",
	}

	index := int(productID) % len(products)
	return products[index]
}
