package db

import (
	"database/sql"
	_ "database/sql"
	_ "fmt"
	"log"
	_ "log"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetCategories() ([]Category, error) {
	rows, err := Db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Panic(err)
		}
	}(rows)

	categories := make([]Category, 0)
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *Repository) GetSubCategories(categoryID int) ([]SubCategory, error) {
	rows, err := Db.Query("SELECT id, name, callback FROM subcategories WHERE category_id = $1", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subCategories []SubCategory
	for rows.Next() {
		var subCategory SubCategory // Используем SubCategory вместо Category
		err := rows.Scan(&subCategory.ID, &subCategory.Name, &subCategory.Callback)
		if err != nil {
			return nil, err
		}
		subCategories = append(subCategories, subCategory) // Добавляем заполненную структуру
	}
	return subCategories, nil
}

func (r *Repository) GetProducts(subCategoryID int) ([]Product, error) {
	rows, err := Db.Query("SELECT id, name, description FROM products WHERE subcategory_id = $1", subCategoryID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Panic(err)
		}
	}(rows)

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
