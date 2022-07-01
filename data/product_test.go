package data

import "testing"

func TestCheckValidation(t *testing.T) {

	p := &Product{
		Name:  "Coffee",
		Price: 1.0,
		SKU:   "ddd-ddd-dsjbdns",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
