package models

import (
    "encoding/json"
    "time"
)

type Expense struct {
    User                    int
    Category, Description   string
    TotalAmount, OwedAmount float64
    Date, ReportDate        time.Time
}

func (e Expense) Save() (err error) {
    data, err := json.Marshal(e)
    if err != nil {
        return
    }

    return save(data)
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
