package model

import "time"

type PriceAnalyticsRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Days      int   `json:"days,omitempty"` // По умолчанию 30 дней
}

type PriceAnalyticsResponse struct {
	ProductID    int64               `json:"product_id"`
	ProductName  string              `json:"product_name"`
	CurrentPrice int                 `json:"current_price"`
	MinPrice     int                 `json:"min_price"`
	MaxPrice     int                 `json:"max_price"`
	AvgPrice     float64             `json:"avg_price"`
	PriceHistory []PriceHistoryItem  `json:"price_history"`
	Suppliers    []SupplierAnalytics `json:"suppliers"`
}

type PriceHistoryItem struct {
	Date     time.Time `json:"date"`
	Price    int       `json:"price"`
	Supplier string    `json:"supplier"`
}

type SupplierAnalytics struct {
	SupplierID   int64   `json:"supplier_id"`
	SupplierName string  `json:"supplier_name"`
	CurrentPrice int     `json:"current_price"`
	MinPrice     int     `json:"min_price"`
	MaxPrice     int     `json:"max_price"`
	AvgPrice     float64 `json:"avg_price"`
	PriceChanges int     `json:"price_changes_count"`
}

// Market Analytics Models
type MarketAnalyticsResponse struct {
	Period           int                 `json:"period_days"`
	Category         string              `json:"category,omitempty"`
	TotalProducts    int                 `json:"total_products"`
	AvgPriceChange   float64             `json:"avg_price_change_percent"`
	MarketVolatility float64             `json:"market_volatility_percent"`
	TopGainers       []ProductTrend      `json:"top_gainers"`
	TopLosers        []ProductTrend      `json:"top_losers"`
	TrendData        []MarketTrendItem   `json:"trend_data"`
	Categories       []CategoryAnalytics `json:"categories"`
}

type ProductTrend struct {
	ProductID    int64   `json:"product_id"`
	ProductName  string  `json:"product_name"`
	CurrentPrice int     `json:"current_price"`
	PriceChange  float64 `json:"price_change_percent"`
	Volume       int     `json:"volume"`
}

type MarketTrendItem struct {
	Date        time.Time `json:"date"`
	MarketIndex float64   `json:"market_index"`
	Volume      int       `json:"volume"`
}

type CategoryAnalytics struct {
	CategoryName    string  `json:"category_name"`
	ProductCount    int     `json:"product_count"`
	AvgPriceChange  float64 `json:"avg_price_change_percent"`
	Volatility      float64 `json:"volatility_percent"`
	TopProduct      string  `json:"top_product"`
	TopProductPrice int     `json:"top_product_price"`
}
