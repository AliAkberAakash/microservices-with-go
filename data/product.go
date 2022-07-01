package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedAt   string  `json:"-"`
	UpdatedAt   string  `json:"-"`
	DeletedAt   string  `json:"-"`
}

type ProductList []*Product

func (p *ProductList) ToJson(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(p)
}

func (p *Product) FromJson(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Product) Validate() error {
	validator := validator.New()
	validator.RegisterValidation("sku", validateSKU)

	return validator.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := reg.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

func GetProducts() ProductList {
	return productList
}

func AddProduct(p *Product) {
	p.ID = getNextId()
	productList = append(productList, p)
}

func UpdateProduct(id int, p *Product) error {

	_, index, err := findProduct(id)

	if err != nil {
		return err
	}

	p.ID = id
	productList[index] = p

	return nil
}

func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

var ErrorProductNotFound = fmt.Errorf("Prodct not found")

func findProduct(id int) (*Product, int, error) {
	for index, product := range productList {
		if product.ID == id {
			return product, index, nil
		}
	}

	return nil, -1, ErrorProductNotFound
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "aj20-bA12Ls2-sNx1KG902eds",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "aj20-dub87db92-sNx1KG3723",
		CreatedAt:   time.Now().UTC().String(),
		UpdatedAt:   time.Now().UTC().String(),
	},
}
