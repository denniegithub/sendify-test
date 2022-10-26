package main

import (
	"encoding/csv"
	"log"
	"os"
)

type Country struct {
	Name      string
	AlphaCode string
	Region    string
}

// contains list of countries and corresponding regions/ISO alpha codes
var countries []Country

func main() {

	// open file
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

	// init countries list
	countries = createCountriesList(data)

	// print the array
	// for _, v := range countries {
	// 	fmt.Printf("%+v\n", v)
	// }

}

func createCountriesList(data [][]string) []Country {
	var countries []Country
	for i, line := range data {
		if i > 0 { // omit header line
			var rec Country
			// omit some fields that we are not interested in
			for j, field := range line {
				if j == 0 {
					rec.Name = field
				} else if j == 1 {
					rec.AlphaCode = field
				} else if j == 5 {
					rec.Region = field
				}
			}
			countries = append(countries, rec)
		}
	}
	return countries
}
