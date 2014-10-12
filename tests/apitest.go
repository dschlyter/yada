package tests

import (
    "encoding/json"
    "github.com/dschlyter/yada/app/models"
    "github.com/revel/revel"
    "net/url"
    "os"
)

type ApiTest struct {
    revel.TestSuite
}

func (t *ApiTest) Before() {
    models.CloseDB()
    os.RemoveAll(models.DB_TEST) // Wipe previous data
    models.InitDB(models.DB_TEST)
}

func (t *ApiTest) After() {
    models.InitDB(models.DB)
    models.ClearMockTime()
}

func exampleData() url.Values {
    ret := url.Values{}
    ret.Set("user", "1")
    ret.Set("category", "stuff")
    ret.Set("date", "2010-01-02T20:03:04Z")
    ret.Set("description", "some stuff")
    ret.Set("totalAmount", "100")
    ret.Set("owedAmount", "60")
    return ret
}

func (t ApiTest) TestAddExpense() {
    v := exampleData()
    t.PostForm("/api/add", v)
    t.AssertOk()

    t.Get("/api/list")
    result := []models.Expense{}
    panicOn(json.Unmarshal(t.ResponseBody, &result))

    t.AssertEqual("stuff", result[0].Category)
    t.AssertOk()
}

func (t ApiTest) TestAddTwoExpenses() {
    v := exampleData()
    t.PostForm("/api/add", v)
    t.AssertOk()

    t.PostForm("/api/add", v)
    t.AssertOk()

    t.Get("/api/list")
    result := []models.Expense{}
    panicOn(json.Unmarshal(t.ResponseBody, &result))

    t.AssertEqual(2, len(result))
    t.AssertOk()
}

func (t ApiTest) TestExpensesSortedByTime() {
    v := exampleData()
    v.Set("description", "stuff")
    panicOn(models.SetMockTime("2014-08-25T22:00:00"))
    t.PostForm("/api/add", v)
    t.AssertOk()

    v.Set("description", "stuff2")
    panicOn(models.SetMockTime("2014-08-25T21:00:00"))
    t.PostForm("/api/add", v)
    t.AssertOk()

    v.Set("description", "stuff3")
    panicOn(models.SetMockTime("2014-08-25T23:00:00"))
    t.PostForm("/api/add", v)
    t.AssertOk()

    t.Get("/api/list")
    result := []models.Expense{}
    panicOn(json.Unmarshal(t.ResponseBody, &result))

    t.AssertEqual(3, len(result))
    t.AssertEqual("stuff2", result[0].Description)
    t.AssertEqual("stuff", result[1].Description)
    t.AssertEqual("stuff3", result[2].Description)
    t.AssertOk()
}

func (t ApiTest) TestValidateTotalAmount() {
    v := exampleData()
    v.Set("totalAmount", "0")
    t.PostForm("/api/add", v)
    t.AssertStatus(400)
}

func (t ApiTest) TestValidateOwedAmount() {
    v := exampleData()
    v.Set("owedAmount", "-9000")
    t.PostForm("/api/add", v)
    t.AssertStatus(400)
}

func (t ApiTest) TestValidateDate() {
    v := exampleData()
    v.Set("date", "this is an invalid date")
    t.PostForm("/api/add", v)
    t.AssertStatus(400)
}

func (t ApiTest) TestValidateUserExists() {
    v := exampleData()
    v.Del("user")
    t.PostForm("/api/add", v)
    t.AssertStatus(400)
}

func (t ApiTest) TestValidateUserValid() {
    v := exampleData()
    v.Set("user", "29")
    t.PostForm("/api/add", v)
    t.AssertStatus(400)
}

func (t ApiTest) TestValidateCategoryExists() {
    v := exampleData()
    v.Del("category")
    t.PostForm("/api/add", v)
    t.AssertStatus(400)
}

func panicOn(err error) {
    if err != nil {
        panic(err)
    }
}
