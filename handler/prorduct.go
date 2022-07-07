// Package Classification Product API
//
// Documentation for Product API
//
//	Schemes http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/AliAkberAakash/microservices-with-go/data"
	"github.com/gorilla/mux"
)

// A list of prducts returns in the response
// swagger:response productsResponse
type productsResponse struct {
	// All products in the system
	// in: body
	Body []data.Product
}

type Product struct {
	logger *log.Logger
}

func NewProduct(logger *log.Logger) *Product {
	return &Product{logger: logger}
}

// swagger:route GET /products products listProducts
//
//	Returns a list of all products
//
// responses:
//	200: productsResponse
//
//
//
func (p *Product) GetProducts(rw http.ResponseWriter, r *http.Request) {

	p.logger.Println("Gettinig all products")

	productList := data.GetProducts()

	err := productList.ToJson(rw)

	if err != nil {
		http.Error(rw, "Unable to parse data", http.StatusInternalServerError)
	}
}

func (p *Product) DeleteProduct(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Invalid id", http.StatusBadRequest)
		return
	}

	p.logger.Println("Deleting product")

	err = data.DeleteProduct(id)

	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
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
			p.logger.Println("[Error] deserializing product")
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		err = product.Validate()
		if err != nil {
			p.logger.Println("[Error] deserializing product")
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err.Error()),
				http.StatusBadRequest,
			)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
