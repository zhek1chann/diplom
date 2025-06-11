package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// ========================================
// Структуры данных
// ========================================

// DeliveryCondition structure for delivery_conditions.json
type DeliveryCondition struct {
	ConditionID               int     `json:"condition_id"`
	MinimumFreeDeliveryAmount float64 `json:"minimum_free_delivery_amount"`
	DeliveryFee               float64 `json:"delivery_fee"`
}

// Supplier structure for suppliers.json (базовые поставщики)
type Supplier struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	ConditionID int    `json:"condition_id"`
}

// ExtendedSupplier structure для расширенных поставщиков
type ExtendedSupplier struct {
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	ConditionID int    `json:"condition_id"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Address     string `json:"address,omitempty"`
	TIN         string `json:"tin,omitempty"`
}

// Product structure for products.json (базовые продукты)
type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ImageURL   string `json:"image_url"`
	GTIN       int64  `json:"gtin"`
	MinPrice   int    `json:"min_price"`
	SupplierID int    `json:"supplier_id"`
}

// EnhancedProduct structure для расширенных продуктов
type EnhancedProduct struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	ImageURL          string `json:"image_url"`
	GTIN              string `json:"gtin"`
	MinPrice          int    `json:"min_price"`
	CategoryID        int    `json:"category_id"`
	SubcategoryID     int    `json:"subcategory_id"`
	Brand             string `json:"brand"`
	Weight            string `json:"weight"`
	CountryOrigin     string `json:"country_origin"`
	ShelfLifeDays     int    `json:"shelf_life_days"`
	StorageConditions string `json:"storage_conditions"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	SupplierID        int    `json:"supplier_id"`
}

// ProductSupplier structure for products_supplier.json
type ProductSupplier struct {
	ProductID     int `json:"product_id"`
	SupplierID    int `json:"supplier_id"`
	Price         int `json:"price"`
	MinSellAmount int `json:"min_sell_amount"`
}

// Category structure для категорий
type Category struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Subcategories []Subcategory `json:"subcategories"`
}

// Subcategory structure для подкатегорий
type Subcategory struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CategoryID int    `json:"category_id"`
}

// PriceHistory structure для истории цен
type PriceHistory struct {
	ID           int    `json:"id"`
	ProductID    int    `json:"product_id"`
	SupplierID   int    `json:"supplier_id"`
	Price        int    `json:"price"`
	Date         string `json:"date"`
	ChangeReason string `json:"change_reason"`
}

// ========================================
// Вспомогательные функции
// ========================================

func generatePhoneNumber(userID int) string {
	// Генерируем телефон казахстанского формата +7XXXXXXXXX (12 символов максимум)
	return fmt.Sprintf("+7%09d", userID%1000000000)
}

func generateHashedPassword() (string, error) {
	const password = "password1@"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func generateOrderAmountByIndex(index int) int {
	orderAmounts := []int{30000, 50000, 70000, 100000, 150000}
	return orderAmounts[index%len(orderAmounts)]
}

func loadJSON(fileName string, data interface{}) error {
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %s: %v", fileName, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("не удалось декодировать JSON из %s: %v", fileName, err)
	}
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// ========================================
// Функции загрузки данных
// ========================================

func loadDeliveryConditions(db *sqlx.DB) error {
	fmt.Println("📦 Загружаем условия доставки...")

	deliveryConditions := []DeliveryCondition{}
	if err := loadJSON("delivery_conditions.json", &deliveryConditions); err != nil {
		return err
	}

	for _, condition := range deliveryConditions {
		_, err := db.Exec(`
			INSERT INTO delivery_conditions (condition_id, minimum_free_delivery_amount, delivery_fee)
			VALUES ($1, $2, $3)
			ON CONFLICT (condition_id) DO UPDATE SET 
				minimum_free_delivery_amount = $2, delivery_fee = $3`,
			condition.ConditionID, condition.MinimumFreeDeliveryAmount, condition.DeliveryFee)
		if err != nil {
			return fmt.Errorf("ошибка загрузки условий доставки: %v", err)
		}
	}

	fmt.Printf("✅ Загружено %d условий доставки\n", len(deliveryConditions))
	return nil
}

func loadCategories(db *sqlx.DB) error {
	if !fileExists("categories.json") {
		fmt.Println("⚠️  Файл categories.json не найден, пропускаем...")
		return nil
	}

	fmt.Println("📂 Загружаем категории...")

	categories := []Category{}
	if err := loadJSON("categories.json", &categories); err != nil {
		return err
	}

	for _, category := range categories {
		_, err := db.Exec(`
			INSERT INTO categories (id, name, description) 
			VALUES ($1, $2, $3) 
			ON CONFLICT (id) DO UPDATE SET name = $2, description = $3`,
			category.ID, category.Name, category.Description)
		if err != nil {
			return fmt.Errorf("ошибка добавления категории %d: %v", category.ID, err)
		}

		for _, subcategory := range category.Subcategories {
			_, err := db.Exec(`
				INSERT INTO subcategories (id, category_id, name) 
				VALUES ($1, $2, $3) 
				ON CONFLICT (id) DO UPDATE SET category_id = $2, name = $3`,
				subcategory.ID, subcategory.CategoryID, subcategory.Name)
			if err != nil {
				return fmt.Errorf("ошибка добавления подкатегории %d: %v", subcategory.ID, err)
			}
		}
	}

	fmt.Printf("✅ Загружено %d категорий\n", len(categories))
	return nil
}

func loadBasicSuppliers(db *sqlx.DB) error {
	fmt.Println("🏪 Загружаем базовых поставщиков...")

	suppliers := []Supplier{}
	if err := loadJSON("suppliers.json", &suppliers); err != nil {
		return err
	}

	for i, supplier := range suppliers {
		phoneNumber := generatePhoneNumber(supplier.UserID)
		hashedPassword, err := generateHashedPassword()
		if err != nil {
			return fmt.Errorf("ошибка генерации пароля: %v", err)
		}

		_, err = db.Exec(`
			INSERT INTO users (id, name, phone_number, hashed_password, role) 
			VALUES ($1, $2, $3, $4, 1)
			ON CONFLICT (id) DO UPDATE SET name = $2, phone_number = $3`,
			supplier.UserID, supplier.Name, phoneNumber, hashedPassword)
		if err != nil {
			return fmt.Errorf("ошибка добавления пользователя %d: %v", supplier.UserID, err)
		}

		orderAmount := generateOrderAmountByIndex(i)
		_, err = db.Exec(`
			INSERT INTO suppliers (user_id, condition_id, name, order_amount) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE SET 
				condition_id = $2, name = $3, order_amount = $4`,
			supplier.UserID, supplier.ConditionID, supplier.Name, orderAmount)
		if err != nil {
			return fmt.Errorf("ошибка добавления поставщика %d: %v", supplier.UserID, err)
		}
	}

	fmt.Printf("✅ Загружено %d базовых поставщиков\n", len(suppliers))
	return nil
}

func loadExtendedSuppliers(db *sqlx.DB) error {
	if !fileExists("extended_suppliers.json") {
		fmt.Println("⚠️  Файл extended_suppliers.json не найден, пропускаем...")
		return nil
	}

	fmt.Println("🏪 Загружаем расширенных поставщиков...")

	extendedSuppliers := []ExtendedSupplier{}
	if err := loadJSON("extended_suppliers.json", &extendedSuppliers); err != nil {
		return err
	}

	for i, supplier := range extendedSuppliers {
		phoneNumber := generatePhoneNumber(supplier.UserID)
		hashedPassword, err := generateHashedPassword()
		if err != nil {
			return fmt.Errorf("ошибка генерации пароля: %v", err)
		}

		_, err = db.Exec(`
			INSERT INTO users (id, name, phone_number, hashed_password, role) 
			VALUES ($1, $2, $3, $4, 1)
			ON CONFLICT (id) DO UPDATE SET name = $2, phone_number = $3`,
			supplier.UserID, supplier.Name, phoneNumber, hashedPassword)
		if err != nil {
			return fmt.Errorf("ошибка добавления пользователя %d: %v", supplier.UserID, err)
		}

		orderAmount := generateOrderAmountByIndex(i)
		_, err = db.Exec(`
			INSERT INTO suppliers (user_id, condition_id, name, order_amount) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE SET 
				condition_id = $2, name = $3, order_amount = $4`,
			supplier.UserID, supplier.ConditionID, supplier.Name, orderAmount)
		if err != nil {
			return fmt.Errorf("ошибка добавления поставщика %d: %v", supplier.UserID, err)
		}
	}

	fmt.Printf("✅ Загружено %d расширенных поставщиков\n", len(extendedSuppliers))
	return nil
}

func loadBasicProducts(db *sqlx.DB) error {
	fmt.Println("📦 Загружаем базовые продукты...")

	products := []Product{}
	if err := loadJSON("products.json", &products); err != nil {
		return err
	}

	for _, product := range products {
		_, err := db.Exec(`
			INSERT INTO products (id, name, image_url, gtin)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET 
				name = $2, image_url = $3, gtin = $4`,
			product.ID, product.Name, product.ImageURL, product.GTIN)
		if err != nil {
			return fmt.Errorf("ошибка добавления продукта %d: %v", product.ID, err)
		}
	}

	fmt.Printf("✅ Загружено %d базовых продуктов\n", len(products))
	return nil
}

func loadEnhancedProducts(db *sqlx.DB) error {
	fmt.Println("📦 Обновляем продукты расширенными данными...")

	// Загружаем products_enhanced.json для обновления существующих
	if fileExists("products_enhanced.json") {
		enhancedProducts := []EnhancedProduct{}
		if err := loadJSON("products_enhanced.json", &enhancedProducts); err != nil {
			return err
		}

		for _, product := range enhancedProducts {
			createdAt, _ := time.Parse("2006-01-02T15:04:05Z", product.CreatedAt)
			updatedAt, _ := time.Parse("2006-01-02T15:04:05Z", product.UpdatedAt)

			_, err := db.Exec(`
				UPDATE products SET 
					category_id = $2, 
					subcategory_id = $3,
					created_at = $4,
					updated_at = $5
				WHERE id = $1`,
				product.ID, product.CategoryID, product.SubcategoryID,
				createdAt, updatedAt)
			if err != nil {
				fmt.Printf("⚠️  Предупреждение: не удалось обновить продукт %d: %v\n", product.ID, err)
			}
		}
		fmt.Printf("✅ Обновлено %d продуктов категориями\n", len(enhancedProducts))
	}

	// Загружаем extended_products.json для добавления новых
	if fileExists("extended_products.json") {
		extendedProducts := []EnhancedProduct{}
		if err := loadJSON("extended_products.json", &extendedProducts); err != nil {
			return err
		}

		for _, product := range extendedProducts {
			createdAt, _ := time.Parse("2006-01-02T15:04:05Z", product.CreatedAt)
			updatedAt, _ := time.Parse("2006-01-02T15:04:05Z", product.UpdatedAt)

			_, err := db.Exec(`
				INSERT INTO products (id, name, image_url, gtin, category_id, subcategory_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				ON CONFLICT (id) DO UPDATE SET 
					name = $2, image_url = $3, gtin = $4, 
					category_id = $5, subcategory_id = $6, updated_at = $8`,
				product.ID, product.Name, product.ImageURL, product.GTIN,
				product.CategoryID, product.SubcategoryID, createdAt, updatedAt)
			if err != nil {
				return fmt.Errorf("ошибка добавления расширенного продукта %d: %v", product.ID, err)
			}
		}
		fmt.Printf("✅ Добавлено %d новых продуктов\n", len(extendedProducts))
	}

	return nil
}

func loadProductSuppliers(db *sqlx.DB) error {
	fmt.Println("💰 Загружаем связи продукт-поставщик...")

	productSuppliers := []ProductSupplier{}
	if err := loadJSON("products_supplier.json", &productSuppliers); err != nil {
		return err
	}

	successCount := 0
	for _, ps := range productSuppliers {
		_, err := db.Exec(`
			INSERT INTO products_supplier (product_id, supplier_id, price, sell_amount)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (product_id, supplier_id) DO UPDATE SET 
				price = $3, sell_amount = $4`,
			ps.ProductID, ps.SupplierID, ps.Price, ps.MinSellAmount)
		if err != nil {
			fmt.Printf("⚠️  Предупреждение: не удалось добавить связь продукт %d - поставщик %d: %v\n",
				ps.ProductID, ps.SupplierID, err)
		} else {
			successCount++
		}
	}

	fmt.Printf("✅ Загружено %d связей продукт-поставщик\n", successCount)
	return nil
}

func loadPriceHistory(db *sqlx.DB) error {
	fmt.Println("📈 Загружаем историю цен...")

	totalLoaded := 0

	// Загружаем базовую историю цен
	if fileExists("price_history.json") {
		priceHistory := []PriceHistory{}
		if err := loadJSON("price_history.json", &priceHistory); err != nil {
			return err
		}

		for _, history := range priceHistory {
			date, _ := time.Parse("2006-01-02T15:04:05Z", history.Date)

			_, err := db.Exec(`
				INSERT INTO price_history (id, product_id, supplier_id, price, date, change_reason)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (id) DO UPDATE SET 
					product_id = $2, supplier_id = $3, price = $4, date = $5, change_reason = $6`,
				history.ID, history.ProductID, history.SupplierID,
				history.Price, date, history.ChangeReason)
			if err != nil {
				fmt.Printf("⚠️  Предупреждение: не удалось добавить запись истории %d: %v\n", history.ID, err)
			} else {
				totalLoaded++
			}
		}
	}

	// Загружаем расширенную историю цен
	if fileExists("extended_price_history.json") {
		extendedPriceHistory := []PriceHistory{}
		if err := loadJSON("extended_price_history.json", &extendedPriceHistory); err != nil {
			return err
		}

		for _, history := range extendedPriceHistory {
			date, _ := time.Parse("2006-01-02T15:04:05Z", history.Date)

			_, err := db.Exec(`
				INSERT INTO price_history (id, product_id, supplier_id, price, date, change_reason)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (id) DO UPDATE SET 
					product_id = $2, supplier_id = $3, price = $4, date = $5, change_reason = $6`,
				history.ID, history.ProductID, history.SupplierID,
				history.Price, date, history.ChangeReason)
			if err != nil {
				fmt.Printf("⚠️  Предупреждение: не удалось добавить запись истории %d: %v\n", history.ID, err)
			} else {
				totalLoaded++
			}
		}
	}

	if totalLoaded == 0 {
		fmt.Println("⚠️  Файлы истории цен не найдены, пропускаем...")
	} else {
		fmt.Printf("✅ Загружено %d записей истории цен\n", totalLoaded)
	}
	return nil
}

func updateProductPrices(db *sqlx.DB) error {
	fmt.Println("🔄 Обновляем минимальные цены продуктов...")

	rows, err := db.Query(`
		SELECT ps.product_id, ps.price AS lowest_price, ps.supplier_id
		FROM products_supplier ps
		WHERE ps.price = (
			SELECT MIN(price)
			FROM products_supplier
			WHERE product_id = ps.product_id
		)
		GROUP BY ps.product_id, ps.supplier_id, ps.price;
	`)
	if err != nil {
		return fmt.Errorf("ошибка запроса минимальных цен: %v", err)
	}
	defer rows.Close()

	updateCount := 0
	updateQuery := `UPDATE products SET lowest_price = $1, lowest_supplier_id = $2 WHERE id = $3`

	for rows.Next() {
		var productID, lowestPrice, lowestSupplierID int

		if err := rows.Scan(&productID, &lowestPrice, &lowestSupplierID); err != nil {
			return fmt.Errorf("ошибка сканирования строки: %v", err)
		}

		_, err := db.Exec(updateQuery, lowestPrice, lowestSupplierID, productID)
		if err != nil {
			fmt.Printf("⚠️  Предупреждение: не удалось обновить цену продукта %d: %v\n", productID, err)
		} else {
			updateCount++
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("ошибка итерации строк: %v", err)
	}

	fmt.Printf("✅ Обновлено %d минимальных цен\n", updateCount)
	return nil
}

// ========================================
// Главная функция
// ========================================

func main() {
	fmt.Println("🚀 Универсальная загрузка данных каталога в PostgreSQL")
	fmt.Println(strings.Repeat("=", 60))

	// Подключение к базе данных
	connStr := "host=localhost port=5432 dbname=catalog user=note-user password=note-password sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Ошибка ping БД: %v", err)
	}
	fmt.Println("✅ Подключение к PostgreSQL установлено")

	// Этапы загрузки данных
	steps := []struct {
		name string
		fn   func(*sqlx.DB) error
	}{
		{"Условия доставки", loadDeliveryConditions},
		{"Категории и подкатегории", loadCategories},
		{"Базовые поставщики", loadBasicSuppliers},
		{"Расширенные поставщики", loadExtendedSuppliers},
		{"Базовые продукты", loadBasicProducts},
		{"Расширенные продукты", loadEnhancedProducts},
		{"Связи продукт-поставщик", loadProductSuppliers},
		{"История цен", loadPriceHistory},
		{"Обновление минимальных цен", updateProductPrices},
	}

	// Выполняем каждый этап
	for i, step := range steps {
		fmt.Printf("\n📋 Этап %d/%d: %s\n", i+1, len(steps), step.name)
		if err := step.fn(db); err != nil {
			log.Fatalf("❌ Ошибка на этапе '%s': %v", step.name, err)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 Загрузка данных завершена успешно!")
	fmt.Println("\n📊 Что загружено:")
	fmt.Println("   ✅ Условия доставки")
	fmt.Println("   ✅ Категории и подкатегории (если доступны)")
	fmt.Println("   ✅ Поставщики (базовые + расширенные)")
	fmt.Println("   ✅ Продукты (базовые + расширенные)")
	fmt.Println("   ✅ Связи продукт-поставщик с ценами")
	fmt.Println("   ✅ История изменений цен (если доступна)")
	fmt.Println("   ✅ Автоматический расчет минимальных цен")

	fmt.Println("\n🚀 База данных готова к использованию!")
}
