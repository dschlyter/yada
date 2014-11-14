package tests

import (
	"encoding/json"
	"net/url"
	"os"
	"runtime/debug"

	"github.com/dschlyter/yada/app/models"
	"github.com/revel/revel"
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
	t.PostForm("/api/expenses", v)
	t.AssertOk()

	t.Get("/api/expenses?user=1")
	result := []models.Expense{}
	panicOn(json.Unmarshal(t.ResponseBody, &result))

	t.AssertEqual("stuff", result[0].Category)
	t.AssertOk()
}

func (t ApiTest) TestAddTwoExpenses() {
	v := exampleData()
	t.PostForm("/api/expenses", v)
	t.AssertOk()

	t.PostForm("/api/expenses", v)
	t.AssertOk()

	t.Get("/api/expenses?user=1")
	result := []models.Expense{}
	panicOn(json.Unmarshal(t.ResponseBody, &result))

	t.AssertEqual(2, len(result))
	t.AssertOk()
}

func (t ApiTest) addThreeThings(user string) []models.Expense {
	v := exampleData()
	v.Set("description", "stuff")
	panicOn(models.SetMockTime("2014-08-25T22:00:00"))
	t.PostForm("/api/expenses", v)
	t.AssertOk()

	v.Set("description", "stuff2")
	panicOn(models.SetMockTime("2014-08-25T21:00:00"))
	t.PostForm("/api/expenses", v)
	t.AssertOk()

	v.Set("description", "stuff3")
	v.Set("user", "2")
	v.Set("owedAmount", "20")
	panicOn(models.SetMockTime("2014-08-25T23:00:00"))
	t.PostForm("/api/expenses", v)
	t.AssertOk()

	t.Get("/api/expenses?user=" + user)
	t.AssertOk()
	result := []models.Expense{}
	panicOn(json.Unmarshal(t.ResponseBody, &result))

	return result
}

func (t ApiTest) TestExpensesSortedDescendingByTime() {
	result := t.addThreeThings("1")

	t.AssertEqual(3, len(result))
	t.AssertEqual("stuff3", result[0].Description)
	t.AssertEqual("stuff", result[1].Description)
	t.AssertEqual("stuff2", result[2].Description)
}

func (t ApiTest) TestExpenseSummingForUser1() {
	result := t.addThreeThings("1")

	t.AssertEqual(100, result[0].Balance)
	t.AssertEqual(120, result[1].Balance)
	t.AssertEqual(60, result[2].Balance)
}

func (t ApiTest) TestExpenseSummingForUser2() {
	result := t.addThreeThings("2")

	t.AssertEqual(-100, result[0].Balance)
	t.AssertEqual(-120, result[1].Balance)
	t.AssertEqual(-60, result[2].Balance)
}

func (t ApiTest) TestNegativeAmount_AllowedAndDecreasesSum() {
	v := exampleData()
	v.Set("totalAmount", "-7000")
	v.Set("owedAmount", "-5000")
	t.PostForm("/api/expenses", v)
	t.AssertOk()

	v = exampleData()
	v.Set("totalAmount", "6000")
	v.Set("owedAmount", "4950")
	t.PostForm("/api/expenses", v)
	t.AssertOk()

	t.Get("/api/expenses?user=1")
	t.AssertOk()
	result := []models.Expense{}
	panicOn(json.Unmarshal(t.ResponseBody, &result))
	t.AssertEqual(-50, result[0].Balance)
}

func (t ApiTest) TestValidateOwedAmountOverTotal_NotAllowed() {
	v := exampleData()
	v.Set("totalAmount", "500")
	v.Set("owedAmount", "700")
	t.PostForm("/api/expenses", v)
	t.AssertStatus(400)
}

func (t ApiTest) TestValidateNegativeOwedAmountUnderTotal_NotAllowed() {
	v := exampleData()
	v.Set("totalAmount", "-500")
	v.Set("owedAmount", "-700")
	t.PostForm("/api/expenses", v)
	t.AssertStatus(400)
}

func (t ApiTest) TestValidateDate() {
	v := exampleData()
	v.Set("date", "this is an invalid date")
	t.PostForm("/api/expenses", v)
	t.AssertStatus(400)
}

func (t ApiTest) TestValidateUserExists() {
	v := exampleData()
	v.Del("user")
	t.PostForm("/api/expenses", v)
	t.AssertStatus(400)
}

func (t ApiTest) TestValidateUserValid() {
	v := exampleData()
	v.Set("user", "29")
	t.PostForm("/api/expenses", v)
	t.AssertStatus(400)
}

func (t ApiTest) TestValidateCategoryExists() {
	v := exampleData()
	v.Del("category")
	t.PostForm("/api/expenses", v)
	t.AssertStatus(400)
}

func panicOn(err error) {
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
}
