package utils

import (
	"encoding/csv"
	"fmt"
	"github.com/drkennetz/azcost/azure"
	"log"
	"os"
)

func WriteCSV(filename string, gb azure.GroupBy) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		log.Fatalln("failed to create file", err)
	}
	w := csv.NewWriter(f)
	defer w.Flush()

	if err = w.Write([]string{"ResourceGroup", "ResourceType", "Cost"}); err != nil {
		log.Fatalln("error writing header", err)
	}
	for key := range gb.Gb {
		for subkey := range gb.Gb[key] {
			stringCost := fmt.Sprintf("%f", gb.Gb[key][subkey])
			if err = w.Write([]string{key, subkey, stringCost}); err != nil {
				log.Fatalln("error while writing record to file", err)
			}
		}
	}
}
