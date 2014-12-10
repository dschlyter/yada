package models

import (
	"encoding/json"
	"time"
)

type Expense struct {
	Id                      string
	User                    int
	Category, Description   string
	TotalAmount, OwedAmount int
	Date, ReportDate        time.Time
	Balance                 int // Not used for input / storage for now
}

func (t Expense) key() string {
	date := t.Date

	if t.ReportDate.After(t.Date) {
		date = t.ReportDate
	}

	return date.Format("2006-01-02T15:04:05Z")
}

func (exp Expense) Save() (err error) {
	exp.Id = exp.key()
	data, err := json.Marshal(exp)
	if err != nil {
		return
	}

	return save(exp.Id, data)
}

func ListExpenses() (ret []Expense, err error) {
	blobs, err := getBlobs()
	if err != nil {
		return nil, err
	}

	ret = []Expense{}

	for _, data := range blobs {
		exp := Expense{}
		err := json.Unmarshal(data.Value, &exp)
		if err != nil {
			return nil, err
		}
		exp.Id = data.Key
		ret = append(ret, exp)
	}

	return
}
