package timeexpression

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testData struct {
	t      string
	result bool
}

func TestHourUnitExpression_New(t *testing.T) {
	testDatas := []struct {
		exp   string
		err   error
		start hourUnit
		end   hourUnit
		isAll bool
	}{
		{
			exp:   "[13]",
			err:   ErrHourUnitFormat,
			start: hourUnit{},
			end:   hourUnit{},
			isAll: false,
		},
		{
			exp: "*",
			err: nil,
			start: hourUnit{
				Hour:   0,
				Minute: 0,
				Sec:    0,
			},
			end: hourUnit{
				Hour:   24,
				Minute: 0,
				Sec:    0,
			},
			isAll: true,
		},
		{
			exp: "12:00:00-13:00:00",
			err: nil,
			start: hourUnit{
				Hour:   12,
				Minute: 0,
				Sec:    0,
			},
			end: hourUnit{
				Hour:   13,
				Minute: 0,
				Sec:    0,
			},
			isAll: false,
		},
		{
			exp:   "13:00:00-12:00:00",
			err:   errors.New("hour error: start after end"),
			start: hourUnit{},
			end:   hourUnit{},
			isAll: false,
		},
	}

	for i, data := range testDatas {
		fmt.Printf("[%d] ex[%s]\n", i, data.exp)

		expression, err := newHourUnitExpression(data.exp)
		if data.err != nil {
			assert.NotNil(t, err)
			assert.EqualError(t, data.err, err.Error())
		} else {
			assert.NotNil(t, expression)
			assert.Equal(t, data.start, expression.start)
			assert.Equal(t, data.end, expression.end)
			assert.Equal(t, data.isAll, expression.isAll)
		}

	}
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
