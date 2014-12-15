package tests

import (
	"github.com/dschlyter/yada/app/models"
	"github.com/revel/revel"
)

type TimeTest struct {
	revel.TestSuite
}

func (t TimeTest) TestTime_IsRoundedToSecond() {
	now := models.GetTime()

	t.AssertEqual(0, now.Nanosecond())
}
