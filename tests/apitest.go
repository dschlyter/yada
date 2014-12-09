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

func timeIs(time string) {
	panicOn(models.SetMockTime(time))
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

func (t ApiTest) post(v url.Values) {
	t.PostForm("/api/expenses", v)
	t.AssertOk()
}

func (t ApiTest) get(user string) []models.Expense {
	t.Get("/api/expenses?user=" + user)
	t.AssertOk()
	result := []models.Expense{}
	panicOn(json.Unmarshal(t.ResponseBody, &result))

	return result
}

func (t ApiTest) TestAddExpense() {
	t.post(exampleData())

	t.Get("/api/expenses?user=1")
	result := []models.Expense{}
	panicOn(json.Unmarshal(t.ResponseBody, &result))

	t.AssertEqual("stuff", result[0].Category)
	t.AssertOk()
}

func (t ApiTest) TestAddTwoExpenses() {
	t.post(exampleData())
	t.post(exampleData())

	result := t.get("1")

	t.AssertEqual(2, len(result))
}

func (t ApiTest) addThreeThings(user string) []models.Expense {
	v := exampleData()
	v.Set("description", "stuff")
	timeIs("2014-08-25T22:00:00")
	t.post(v)

	v.Set("description", "stuff2")
	timeIs("2014-08-25T21:00:00")
	t.post(v)

	v.Set("description", "stuff3")
	v.Set("user", "2")
	v.Set("owedAmount", "20")
	timeIs("2014-08-25T23:00:00")
	t.post(v)

	return t.get(user)
}

func (t ApiTest) Test_ReportDateAfterDate_SortedByReportDate() {
	result := t.addThreeThings("1")

	t.AssertEqual(3, len(result))
	t.AssertEqual("stuff3", result[0].Description)
	t.AssertEqual("stuff", result[1].Description)
	t.AssertEqual("stuff2", result[2].Description)
}

func (t ApiTest) Test_ReportDateBeforeDate_SortedByDate() {
	v := exampleData()

	timeIs("2014-08-25T21:00:01")
	v.Set("description", "second")
	v.Set("date", "2015-02-02T20:03:04Z")
	t.post(v)

	timeIs("2014-08-25T21:00:02")
	v.Set("description", "third")
	v.Set("date", "2015-03-02T20:03:04Z")
	t.post(v)

	timeIs("2014-08-25T21:00:04")
	v.Set("description", "a long time ago")
	v.Set("date", "1970-01-02T20:03:04Z")
	t.post(v)

	timeIs("2014-08-25T21:00:03")
	v.Set("description", "first")
	v.Set("date", "2015-01-02T20:03:04Z")
	t.post(v)

	result := t.get("1")

	t.AssertEqual(4, len(result))
	t.AssertEqual("third", result[0].Description)
	t.AssertEqual("second", result[1].Description)
	t.AssertEqual("first", result[2].Description)
	t.AssertEqual("a long time ago", result[3].Description)
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
	t.post(v)

	v = exampleData()
	v.Set("totalAmount", "6000")
	v.Set("owedAmount", "4950")
	t.post(v)

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
