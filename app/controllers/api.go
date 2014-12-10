package controllers

import (
	"errors"
	"math"
	"time"

	"github.com/dschlyter/yada/app/models"
	"github.com/revel/revel"
)

type Api struct {
	*revel.Controller
}

func (c Api) Add(user int, category, description, date string, totalAmount, owedAmount int) revel.Result {
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

func parse(user int, category, description, date string, totalAmount, owedAmount int) (models.Expense, error) {
	newData := models.Expense{User: user, Category: category, Description: description, TotalAmount: totalAmount, OwedAmount: owedAmount}
	newData.ReportDate = models.GetTime()

	parsedDate, validationError := time.Parse("2006-01-02T15:04:05Z", date)
	newData.Date = parsedDate

	if validationError == nil {
		validationError = validate(newData)
	}

	return newData, validationError
}

func validate(data models.Expense) error {
	if math.Abs(float64(data.TotalAmount)) < math.Abs(float64(data.OwedAmount)) {
		return errors.New("Owed amount cannot be larger than total")
	}
	if len(data.Category) <= 0 {
		return errors.New("Category required")
	}
	err := validateUser(data.User)
	if err != nil {
		return err
	}

	return nil
}

func validateUser(user int) error {
	if user != 1 && user != 2 {
		return errors.New("User must be 1 or 2")
	}

	return nil
}

func (c Api) List(afterKey string, user, limit int) revel.Result {
	err := validateUser(user)

	if afterKey == "" {
		afterKey = "LAST" // Default, higher than any timestamp
	}

	if limit <= 0 {
		limit = 10 // Default
	}

	var list []models.Expense
	if err == nil {
		list, err = models.ListExpenses()
	}

	if err != nil {
		return c.RenderError(err)
	}

	// TODO more optimized balance calculation, without iterating trough all values, if needed
	balance := 0
	for i := len(list) - 1; i >= 0; i-- {
		expense := &list[i]
		if user == expense.User {
			balance += expense.OwedAmount
		} else {
			balance -= expense.OwedAmount
		}
		expense.Balance = balance
	}

	start := 0
	for i := 0; i < len(list); i++ {
		if list[i].Id < afterKey {
			start = i
			break
		}
	}

	end := start + limit

	if end > len(list) {
		end = len(list)
	}

	return c.RenderJson(list[start:end])
}
