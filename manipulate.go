// +build ignore

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("fixed-data.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	rows, err := buildRows(r)
	if err != nil {
		log.Fatalln(err)
	}
	sort.Sort(rows)
	fmt.Println(rows)
}

// parseInt parses s into an integer, accounting for quotes and commas.
// It'll panic on invalid inputs.
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	s = strings.Replace(s, `"`, "", -1)
	s = strings.Replace(s, ",", "", -1)
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return int(n)
}

func buildRows(r *csv.Reader) (rs rows, err error) {
	// Header.
	_, err = r.Read()
	if err != nil {
		return nil, err
	}
	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return rs, nil
			}
			return nil, err
		}

		rw := row{
			name:     rec[0],
			location: rec[1],
			aadt09:   parseInt(rec[2]),
			aadt10:   parseInt(rec[3]),
			aadt11:   parseInt(rec[4]),
			aadt12:   parseInt(rec[5]),
			aadt13:   parseInt(rec[6]),
		}

		applyGF := func(n int, old ...int) int {
			for _, v := range old {
				if v != 0 {
					// From http://www.wsdot.wa.gov/mapsdata/travel/pdf/ShortCountFactoringGuide2016.pdf
					return int(math.Ceil(float64(n) * 1.0261))
				}
			}
			return 0 // no data
		}

		rw.aadt14 = applyGF(rw.aadt13, rw.aadt12, rw.aadt11, rw.aadt10, rw.aadt09)
		rw.aadt15 = applyGF(rw.aadt14, rw.aadt13, rw.aadt12, rw.aadt11, rw.aadt10, rw.aadt09)

		rs = append(rs, rw)
	}
	return rs, nil
}

type row struct {
	name     string
	location string

	aadt09 int
	aadt10 int
	aadt11 int
	aadt12 int
	aadt13 int

	// Extrapolated.
	aadt14 int
	aadt15 int
}

type rows []row

func (r rows) Len() int {
	return len(r)
}

func (r rows) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r rows) Less(i, j int) bool {
	return r[i].aadt15 < r[j].aadt15
}

func (r rows) String() string {
	if r == nil {
		return "rows(nil)"
	}
	if len(r) == 0 {
		return "[]"
	}

	var buf bytes.Buffer
	buf.WriteString(r[0].location)
	buf.WriteString(" @ ")
	buf.WriteString(r[0].name)
	buf.WriteString(": ")
	buf.WriteString(strconv.Itoa(r[0].aadt15))

	if len(r) > 0 {
		for _, rw := range r[1:] {
			buf.WriteByte('\n')
			buf.WriteString(rw.location)
			buf.WriteString(" @ ")
			buf.WriteString(rw.name)
			buf.WriteString(": ")
			buf.WriteString(strconv.Itoa(rw.aadt15))
		}
	}
	return buf.String()
}
