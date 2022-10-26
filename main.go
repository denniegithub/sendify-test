package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type RegionClass int64

const (
	Domestic RegionClass = iota
	Europe
	International
)

func (t RegionClass) multiplier() float32 {
	switch t {
	case Domestic:
		return 1
	case Europe:
		return 1.5
	case International:
		return 2.5
	default:
		return 1
	}
}

type WeightClass int64

const (
	Small WeightClass = iota
	Medium
	Large
	Huge
)

func (w WeightClass) price() float32 {
	switch w {
	case Small:
		return 100
	case Medium:
		return 300
	case Large:
		return 500
	case Huge:
		return 2000
	default:
		return 100
	}
}

// contains map corresponding regions/ISO alpha codes
var countries map[string]string

type Shipment struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Weight   float32 `json:"weight"`
}

func main() {

	// open file containing all countries
	f, err := os.Open("countries.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// init countries map
	countries = createCountriesMap(data)

	r := chi.NewRouter()
	r.Post("/shipments", createShipment)
	http.ListenAndServe(":8080", r)

}

func createShipment(w http.ResponseWriter, r *http.Request) {

	// Parse JSON request
	var shipment Shipment
	err := json.NewDecoder(r.Body).Decode(&shipment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	price, err := shipment.calculatePrice()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf("price: %d kr", int(*price))))
}

func (shipment *Shipment) calculatePrice() (*float32, error) {
	var regionA, regionB string

	if region, ok := countries[shipment.Sender]; ok {
		regionA = region
	} else {
		return nil, fmt.Errorf("sender country %s doesn't exist", shipment.Sender)
	}

	if region, ok := countries[shipment.Receiver]; ok {
		regionB = region
	} else {
		return nil, fmt.Errorf("receiver country %s doesn't exist", shipment.Receiver)
	}

	var regionClass RegionClass
	if shipment.Sender == shipment.Receiver {
		// Domestic
		regionClass = Domestic
	} else if regionA == "Europe" && regionB == "Europe" {
		// Europe
		regionClass = Europe
	} else if regionA != regionB {
		// International
		regionClass = International
	}

	price := calculateWeightClassPrice(shipment.Weight) * regionClass.multiplier()

	return &price, nil
}

func calculateWeightClassPrice(weight float32) float32 {
	var weightClass WeightClass
	if weight > 0.0 && weight < 10 {
		weightClass = Small
	} else if weight >= 10 && weight < 25 {
		weightClass = Medium
	} else if weight >= 25 && weight < 50 {
		weightClass = Large
	} else if weight >= 50 && weight <= 1000 {
		weightClass = Huge
	}

	return weightClass.price()
}

func createCountriesMap(data [][]string) map[string]string {
	countries := make(map[string]string, 0)
	for i, line := range data {
		if i > 0 { // omit header line
			var alphaCode, region string
			// omit some fields that we are not interested in
			for j, field := range line {
				if j == 1 {
					alphaCode = field
				} else if j == 5 {
					region = field
				}
			}
			countries[alphaCode] = region
		}
	}

	return countries
}
