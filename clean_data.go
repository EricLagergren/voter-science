// +build ignore

package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("tabula-aadt.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
	r.FieldsPerRecord = -1

	out, err := os.Create("fixed-data.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	w := csv.NewWriter(out)
	w.Comma = '\t'

	header, err := r.Read()
	if err != nil {
		log.Fatalln(err)
	}

	rec, err := r.Read()
	if err != nil {
		log.Fatalln(err)
	}
	street := rec[0]

	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln(err)
		}
		// First record is 'street name', but some fields are missing
		// values (because of how the original PDF was laid out).
		if rec[0] == "" {
			rec[0] = street
		} else {
			street = rec[0]
		}

		// Some rows lack full fields, so pad them.
		if diff := len(header) - len(rec); diff > 0 {
			rec = append(rec, strings.Repeat("", diff))
		}

		if err := w.Write(rec); err != nil {
			log.Fatalln(err)
		}
	}
}