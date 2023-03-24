package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {

	idField := flag.String("id", "", "which column to use as key/id")

	flag.Parse()

	reader := csv.NewReader(os.Stdin)
	reader.FieldsPerRecord = -1

	var headers []string
	var data []map[string]string

	for i := 0; ; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read csv: %s", err)
		}

		if i == 0 {
			// first line. gather headers
			headers = record
		} else {
			vals := map[string]string{}
			for j := 0; j < len(record) && j < len(headers); j++ {
				vals[headers[j]] = record[j]
			}
			if len(vals) > 0 {
				data = append(data, vals)
			}
		}
	}

	if len(flag.Args()) > 0 {
		// only include the named columns
		include := flag.Args()
		if *idField != "" {
			include = append(include, *idField)
		}

		cut := []map[string]string{}

		for _, row := range data {
			cutRow := map[string]string{}
			for _, fld := range include {
				cutRow[fld] = row[fld]
			}
			cut = append(cut, cutRow)
		}

		data = cut
	}

	if *idField != "" {
		byID := map[string]map[string]string{}

		for _, row := range data {
			idColVal := row[*idField]
			delete(row, *idField)
			byID[idColVal] = row
		}

		writeData(byID)
	} else {
		writeData(data)
	}

}

func writeData(data any) {
	out, err := yaml.Marshal(data)
	if err != nil {
		log.Fatalf("failed to marshal data to yaml: %v", err)
	}

	os.Stdout.Write(out)
	os.Stdout.WriteString("\n")
}
