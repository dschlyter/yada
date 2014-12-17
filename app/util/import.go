package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/dschlyter/yada/app/models"
)

// Import expenses from cvs (google spreadsheet format)

func main() {
	file, err := os.Open("import.csv")
	panicOn(err)
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	parsed, err := reader.ReadAll()
	panicOn(err)

	models.InitDB(models.DB)

	for _, v := range parsed {
		date, err := time.Parse("1/2/2006 15:04:05", v[0])
		panicOn(err)
		description := v[5]

		user := -1
		category := ""
		amount := 0
		split := 0.0

		for i := 1; i < 5; i++ {
			if v[i] != "" {
				if user != -1 {
					panic("Multiple input amounts, please scrub data more")
				}

				amount, err = strconv.Atoi(v[i])
				panicOn(err)

				switch i {
				case 1:
					user = 2
					split = 0.6
					category = "Mat"
				case 2:
					user = 2
					split = 1
					category = "Betalning"
				case 3:
					user = 1
					split = 0.4
					category = "Mat"
				case 4:
					user = 1
					split = 1
					category = "Betalning"
				default:
					panic("Amount index out of range")
				}
			}
		}

		expense := models.Expense{
			User:        user,
			Category:    category,
			Description: description,
			TotalAmount: amount,
			OwedAmount:  int(float64(amount)*split + 0.5),
			Date:        date,
			ReportDate:  date,
		}
		expense.Save()
	}
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
