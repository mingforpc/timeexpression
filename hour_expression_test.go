package timeexpression

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHourExpression_IsIn(t *testing.T) {
	expression, err := newHourExpression("13:00:03-16:59:59,23:00:00-23:30:00")
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
		{
			t:      "23:30:00",
			result: true,
		},
	}

	testIsInHour(t, expression, testDatas)

}

func testIsInHour(t *testing.T, expression *hourExpression, testDatas []testData) {
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

func TestHourExpression_GetStart(t *testing.T) {

	testDatas := []struct {
		exp          string
		inputHour    int
		inputMin     int
		inputSec     int
		resultHour   int
		resultMin    int
		resultSec    int
		resultAddDay bool
		err          error
	}{
		{
			exp:          "*",
			inputHour:    12,
			inputMin:     30,
			inputSec:     30,
			resultHour:   0,
			resultMin:    0,
			resultSec:    0,
			resultAddDay: false,
			err:          nil,
		},
		{
			exp:          "03:30:30-05:00:30",
			inputHour:    2,
			inputMin:     30,
			inputSec:     30,
			resultHour:   3,
			resultMin:    30,
			resultSec:    30,
			resultAddDay: false,
			err:          nil,
		},
		{
			exp:          "03:30:30-05:00:30",
			inputHour:    4,
			inputMin:     30,
			inputSec:     30,
			resultHour:   3,
			resultMin:    30,
			resultSec:    30,
			resultAddDay: false,
			err:          nil,
		},
		{
			exp:          "03:30:30-05:00:30",
			inputHour:    5,
			inputMin:     30,
			inputSec:     30,
			resultHour:   3,
			resultMin:    30,
			resultSec:    30,
			resultAddDay: true,
			err:          nil,
		},
	}

	for _, data := range testDatas {
		exp, err := newHourExpression(data.exp)
		if err != nil {
			panic(err)
		}

		startHourUnit, addDay, err := exp.getStart(data.inputHour, data.inputMin, data.inputSec)
		if data.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, data.err.Error())
		}
		assert.Equal(t, data.resultHour, startHourUnit.Hour)
		assert.Equal(t, data.resultMin, startHourUnit.Minute)
		assert.Equal(t, data.resultSec, startHourUnit.Sec)
		assert.Equal(t, data.resultAddDay, addDay)
	}
}

func TestHourExpression_GetEnd(t *testing.T) {
	testDatas := []struct {
		exp          string
		inputHour    int
		inputMin     int
		inputSec     int
		resultHour   int
		resultMin    int
		resultSec    int
		resultAddDay bool
		err          error
	}{
		{
			exp:          "*",
			inputHour:    12,
			inputMin:     30,
			inputSec:     30,
			resultHour:   24,
			resultMin:    0,
			resultSec:    0,
			resultAddDay: false,
			err:          nil,
		},
		{
			exp:          "03:30:30-05:00:30",
			inputHour:    2,
			inputMin:     30,
			inputSec:     30,
			resultHour:   5,
			resultMin:    00,
			resultSec:    30,
			resultAddDay: false,
			err:          nil,
		},
		{
			exp:          "03:30:30-05:00:30",
			inputHour:    4,
			inputMin:     30,
			inputSec:     30,
			resultHour:   5,
			resultMin:    00,
			resultSec:    30,
			resultAddDay: false,
			err:          nil,
		},
		{
			exp:          "03:30:30-05:00:30",
			inputHour:    5,
			inputMin:     30,
			inputSec:     30,
			resultHour:   5,
			resultMin:    00,
			resultSec:    30,
			resultAddDay: true,
			err:          nil,
		},
	}

	for _, data := range testDatas {
		exp, err := newHourExpression(data.exp)
		if err != nil {
			panic(err)
		}

		endHourUnit, addDay, err := exp.getEnd(data.inputHour, data.inputMin, data.inputSec)
		if data.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, data.err.Error())
		}
		assert.Equal(t, data.resultHour, endHourUnit.Hour)
		assert.Equal(t, data.resultMin, endHourUnit.Minute)
		assert.Equal(t, data.resultSec, endHourUnit.Sec)
		assert.Equal(t, data.resultAddDay, addDay)
	}
}
