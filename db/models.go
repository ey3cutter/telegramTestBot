package db

type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Callback string `json:"callback"`
}

type SubCategory struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Callback string `json:"callback"`
}

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
