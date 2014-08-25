package tests

import (
    "encoding/json"
    "github.com/dschlyter/yada/app/models"
    "github.com/revel/revel"
    "net/url"
)

type ApiTest struct {
    revel.TestSuite
}

func (t *ApiTest) Before() {
    models.SetTestMode()
}

func (t ApiTest) TestAddExpense() {
    v := url.Values{}
    v.Set("category", "stuff")
    t.PostForm("/api/add", v)
    t.AssertOk()

    t.Get("/api/list")
    result := []models.Expense{}
    panicOn(json.Unmarshal(t.ResponseBody, &result))

    t.AssertEqual("stuff", result[0].Category)
    t.AssertOk()
}

func (t ApiTest) TestAddTwoExpenses() {
    v := url.Values{}
    v.Set("category", "stuff")
    t.PostForm("/api/add", v)
    t.AssertOk()

    v.Set("category", "stuff2")
    t.PostForm("/api/add", v)
    t.AssertOk()

    t.Get("/api/list")
    result := []models.Expense{}
    panicOn(json.Unmarshal(t.ResponseBody, &result))

    t.AssertEqual(2, len(result))
    //t.AssertEqual("stuff2", result[0].Category)
    //t.AssertEqual("stuff", result[1].Category)
    t.AssertOk()
}

func panicOn(err error) {
    if err != nil {
        panic(err)
    }
}
