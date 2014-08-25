package controllers

import (
    "github.com/dschlyter/yada/app/models"
    "github.com/revel/revel"
)

type Api struct {
    *revel.Controller
}

func (c Api) Add(user int, category, description string, totalAmount, owedAmount float64) revel.Result {
    newData := models.Expense{User: user, Category: category, Description: description, TotalAmount: totalAmount, OwedAmount: owedAmount}
    newData.ReportDate = models.GetTime()
    err := newData.Save()

    if err == nil {
        return c.RenderJson("success")
    } else {
        return c.RenderError(err)
    }
}

func (c Api) List(nextKey string) revel.Result {
    list, err := models.ListExpenses()

    if err == nil {
        return c.RenderJson(list)
    } else {
        return c.RenderError(err)
    }
}
