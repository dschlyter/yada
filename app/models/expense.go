package models

import (
	"encoding/json"
	"time"
)

type Expense struct {
	User                    int
	Category, Description   string
	TotalAmount, OwedAmount int
	Date, ReportDate        time.Time
	Balance                 int // Not used for input / storage for now
}

func (t Expense) key() string {
	if t.ReportDate.After(t.Date) {
		return t.ReportDate.String()
	}

	return t.Date.String()
}

func (exp Expense) Save() (err error) {
	data, err := json.Marshal(exp)
	if err != nil {
		return
	}

	return save(exp.key(), data)
}

func ListExpenses() (ret []Expense, err error) {
	dataList, err := get("", 10)
	if err != nil {
		return nil, err
	}

	ret = []Expense{}

	for _, data := range dataList {
		element := Expense{}
		err := json.Unmarshal(data, &element)
		if err != nil {
			return nil, err
		}
		ret = append(ret, element)
	}

	return
}
