package timeexpression

import (
	"testing"
	"time"
)

type testData struct {
	t      string
	result bool
}

// TestHourUnitExpression_IsIn 判断日期格式
func TestHourUnitExpression_IsIn(t *testing.T) {
	expression, err := newHourUnitExpression("13:00:03-16:59:59")
	if err != nil {
		t.Fatal(err)
	}

	// 定义并初始化并赋值给 data
	testDatas := []testData{
		{
			t:      "13:00:02",
			result: false,
		},
		{
			t:      "13:00:03",
			result: true,
		},
		{
			t:      "15:00:00",
			result: true,
		},
		{
			t:      "17:00:00",
			result: false,
		},
	}

	testIsInHourUnit(t, expression, testDatas)

	expression, err = newHourUnitExpression("*")
	if err != nil {
		t.Fatal(err)
	}

	// 定义并初始化并赋值给 data
	testDatas = []testData{
		{
			t:      "15:04:05",
			result: true,
		},
		{
			t:      "23:59:59",
			result: true,
		},
		{
			t:      "00:00:00",
			result: true,
		},
	}

	expression, err = newHourUnitExpression("*")
	if err != nil {
		t.Fatal(err)
	}
}

func testIsInHourUnit(t *testing.T, expression *hourUnitExpression, testDatas []testData) {
	for _, data := range testDatas {
		hourTime, err := time.Parse("15:04:05", data.t)
		if err != nil {
			t.Fatal(err)
		}

		in := expression.isIn(hourTime.Hour(), hourTime.Minute(), hourTime.Second())

		if in != data.result {
			t.Errorf("want[%v] but get[%v]", data.result, in)
		}
	}
}
