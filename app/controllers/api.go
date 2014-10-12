package controllers

import (
    "errors"
    "github.com/dschlyter/yada/app/models"
    "github.com/revel/revel"
    "time"
)

type Api struct {
    *revel.Controller
}

func (c Api) Add(user int, category, description, date string, totalAmount, owedAmount float64) revel.Result {
    newData, validationError := parse(user, category, description, date, totalAmount, owedAmount)

    if validationError != nil {
        c.Response.Status = 400
        return c.RenderJson(models.ReturnError(validationError.Error()))
    }

    saveError := newData.Save()

    if saveError != nil {
        c.Response.Status = 500
        return c.RenderJson(models.ReturnError(saveError.Error()))
    }

    return c.RenderJson(models.ReturnSuccess())
}

func parse(user int, category, description, date string, totalAmount, owedAmount float64) (models.Expense, error) {
    newData := models.Expense{User: user, Category: category, Description: description, TotalAmount: totalAmount, OwedAmount: owedAmount}
    newData.ReportDate = models.GetTime()

    parsedDate, validationError := time.Parse("2006-01-02T15:03:04Z", date)
    newData.Date = parsedDate

    if validationError == nil {
        validationError = validate(newData)
    }

    return newData, validationError
}

func validate(data models.Expense) error {
    if data.TotalAmount <= 0.0 {
        return errors.New("Positive totalAmount required")
    }
    if data.OwedAmount <= 0.0 {
        return errors.New("Positive owedAmount required")
    }
    if len(data.Category) <= 0 {
        return errors.New("Category required")
    }
    if data.User != 1 && data.User != 2 {
        return errors.New("User must be 1 or 2")
    }

    return nil
}

func (c Api) List(nextKey string) revel.Result {
    list, err := models.ListExpenses()

    if err == nil {
        return c.RenderJson(list)
    } else {
        return c.RenderError(err)
    }
}
