package timeexpression

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDayExpression(t *testing.T) {
	testDatas := []struct {
		exp   string
		err   error
		start int
		end   int
		isAll bool
	}{
		{
			exp:   "*",
			err:   nil,
			start: 1,
			end:   31,
			isAll: true,
		},
		{
			exp:   "02",
			err:   nil,
			start: 2,
			end:   2,
			isAll: false,
		},
		{
			exp:   "02-03",
			err:   nil,
			start: 2,
			end:   3,
			isAll: false,
		},
		{
			exp:   "02-03",
			err:   nil,
			start: 2,
			end:   3,
			isAll: false,
		},
		{
			exp:   "02-03-04",
			err:   ErrDayFormat,
			start: 0,
			end:   0,
			isAll: false,
		},
		{
			exp: "a-03",
			err: &time.ParseError{
				Layout:     "02",
				Value:      "a",
				LayoutElem: "02",
				ValueElem:  "a",
				Message:    "",
			},
			start: 0,
			end:   0,
			isAll: false,
		},
		{
			exp:   "03-02",
			err:   errors.New("day error: start after end"),
			start: 0,
			end:   0,
			isAll: false,
		},
	}

	for _, data := range testDatas {
		exp, err := newDayExpression(data.exp)
		assert.Equal(t, data.err, err)
		if err != nil {
			continue
		}
		assert.Equal(t, data.start, exp.start)
		assert.Equal(t, data.end, exp.end)
		assert.Equal(t, data.isAll, exp.isAll)
	}
}

func TestDayExpression_IsIn(t *testing.T) {
	expression, err := newDayExpression("02-03")
	if err != nil {
		panic(err)
	}

	in := expression.isIn(2)
	assert.True(t, in)
	in = expression.isIn(1)
	assert.False(t, in)
}

func TestDayExpression_GetStart(t *testing.T) {
	testDatas := []struct {
		exp            string
		inputYear      int
		inputMonth     time.Month
		inputDay       int
		err            error
		resultStart    int
		resultAddMonth bool
	}{
		{
			exp:            "*",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       20,
			err:            nil,
			resultStart:    1,
			resultAddMonth: false,
		},
		{
			exp:            "04",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       1,
			err:            nil,
			resultStart:    4,
			resultAddMonth: false,
		},
		{
			exp:            "04",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       4,
			err:            nil,
			resultStart:    4,
			resultAddMonth: false,
		},
		{
			exp:            "04",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       7,
			err:            nil,
			resultStart:    4,
			resultAddMonth: true,
		},
		{
			exp:            "04-08",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       7,
			err:            nil,
			resultStart:    4,
			resultAddMonth: false,
		},
		{
			exp:            "04-08",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       9,
			err:            nil,
			resultStart:    4,
			resultAddMonth: true,
		},
		{
			exp:            "26-27",
			inputYear:      2000,
			inputMonth:     time.February,
			inputDay:       28,
			err:            nil,
			resultStart:    26,
			resultAddMonth: true,
		},
		{
			exp:        "29",
			inputYear:  2000,
			inputMonth: time.February,
			inputDay:   31,
			err:        errors.New("dayExpression getStart get unreachable error"),
		},
	}

	for i, data := range testDatas {
		fmt.Printf("[%d] exp[%s]\n", i, data.exp)
		exp, err := newDayExpression(data.exp)
		if err != nil {
			panic(err)
		}

		start, addMonth, err := exp.getStart(data.inputYear, data.inputMonth, data.inputDay)
		if data.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, data.err.Error())
		}
		assert.Equal(t, data.resultStart, start)
		assert.Equal(t, data.resultAddMonth, addMonth)
	}
}

func TestDayExpression_GetEnd(t *testing.T) {
	testDatas := []struct {
		exp            string
		inputYear      int
		inputMonth     time.Month
		inputDay       int
		err            error
		resultEnd      int
		resultAddMonth bool
	}{
		{
			exp:            "*",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       20,
			err:            nil,
			resultEnd:      31,
			resultAddMonth: false,
		},
		{
			exp:            "04",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       1,
			err:            nil,
			resultEnd:      4,
			resultAddMonth: false,
		},
		{
			exp:            "04",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       4,
			err:            nil,
			resultEnd:      4,
			resultAddMonth: false,
		},
		{
			exp:            "04",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       7,
			err:            nil,
			resultEnd:      4,
			resultAddMonth: true,
		},
		{
			exp:            "04-08",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       7,
			err:            nil,
			resultEnd:      8,
			resultAddMonth: false,
		},
		{
			exp:            "04-08",
			inputYear:      2000,
			inputMonth:     time.January,
			inputDay:       9,
			err:            nil,
			resultEnd:      8,
			resultAddMonth: true,
		},
		{
			exp:            "26-27",
			inputYear:      2000,
			inputMonth:     time.February,
			inputDay:       28,
			err:            nil,
			resultEnd:      27,
			resultAddMonth: true,
		},
		{
			exp:        "29",
			inputYear:  2000,
			inputMonth: time.February,
			inputDay:   31,
			err:        errors.New("dayExpression getEnd get unreachable error"),
		},
	}

	for i, data := range testDatas {
		fmt.Printf("[%d] exp[%s]\n", i, data.exp)
		exp, err := newDayExpression(data.exp)
		if err != nil {
			panic(err)
		}

		end, addMonth, err := exp.getEnd(data.inputYear, data.inputMonth, data.inputDay)
		if data.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, data.err.Error())
		}
		assert.Equal(t, data.resultEnd, end)
		assert.Equal(t, data.resultAddMonth, addMonth)
	}
}
