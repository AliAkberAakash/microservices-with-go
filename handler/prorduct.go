package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/AliAkberAakash/microservices-with-go/data"
	"github.com/gorilla/mux"
)

type Product struct {
	logger *log.Logger
}

func NewProduct(logger *log.Logger) *Product {
	return &Product{logger: logger}
}

func (p *Product) GetProducts(rw http.ResponseWriter, r *http.Request) {

	p.logger.Println("Gettinig all products")

	productList := data.GetProducts()

	err := productList.ToJson(rw)

	if err != nil {
		http.Error(rw, "Unable to parse data", http.StatusInternalServerError)
	}
}
func (p *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {

	p.logger.Println("Adding product")

	product := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(product)
}

func (p *Product) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
		return
	}

	p.logger.Println("Updating product")

	product := r.Context().Value(KeyProduct{}).(*data.Product)
	err = data.UpdateProduct(id, product)

	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Product) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := &data.Product{}
		err := product.FromJson(r.Body)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
