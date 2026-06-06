package product

import (
	"encoding/json"
	"testing"
)

func TestProductJSON(t *testing.T) {
	p := Product{
		ID:       1,
		Name:     "Digital Thermometer",
		Category: "Diagnostics",
		Price:    499,
		Stock:    25,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("could not marshal product: %v", err)
	}

	var decoded Product
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("could not unmarshal product: %v", err)
	}

	if decoded.Name != "Digital Thermometer" {
		t.Errorf("expected name Digital Thermometer but got %s", decoded.Name)
	}
	if decoded.Price != 499 {
		t.Errorf("expected price 499 but got %v", decoded.Price)
	}
}
