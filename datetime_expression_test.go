package timeexpression

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDateTimeExpression(t *testing.T) {

	testDataList := []struct {
		expr         string
		alwaysActive bool
		hasErr       bool
	}{
		{
			expr:         "[*][*][*][*]",
			alwaysActive: true,
			hasErr:       false,
		},
		{
			expr:         "[1991][*][*][*]",
			alwaysActive: false,
			hasErr:       false,
		},
		{
			expr:         "[1991-2021][*][*][*]",
			alwaysActive: false,
			hasErr:       false,
		},
		{
			expr:         "[1991-2021-2022][*][*][*]",
			alwaysActive: false,
			hasErr:       true,
		},
		{
			expr:         "[1991][03][*][*]",
			alwaysActive: false,
			hasErr:       false,
		},
		{
			expr:         "[1991][01-03][*][*]",
			alwaysActive: false,
			hasErr:       false,
		},
		{
			expr:         "[1991][01-13][*][*]",
			alwaysActive: false,
			hasErr:       true,
		},
		{
			expr:         "[1991][01-03][04][*]",
			alwaysActive: false,
			hasErr:       false,
		},
		{
			expr:         "[1991][01-03][04-05][*]",
			alwaysActive: false,
			hasErr:       false,
		},
		{
			expr:         "[1991][01-03][02-32][*]",
			alwaysActive: false,
			hasErr:       true,
		},
		{
			expr:         "[1991][01-03][04-05][11:00:00-12:00:00]",
			alwaysActive: false,
			hasErr:       false,
		},
		{
			expr:         "[1991][01-03][04-05][11:00:00-24:00:01]",
			alwaysActive: false,
			hasErr:       true,
		},
		{
			expr:         "[1991][01-03][04-05][11:00:00-24:00:01][]",
			alwaysActive: false,
			hasErr:       true,
		},
		{
			expr:         "[1991][01-03][04-05][08:00:00-10:00:00,11:00:00-12:00:00]",
			alwaysActive: false,
			hasErr:       false,
		},
	}

	for _, testData := range testDataList {
		expression, err := NewDateTimeExpression(testData.expr)
		assert.Equal(t, testData.hasErr, err != nil)
		if err == nil {
			assert.Equal(t, testData.alwaysActive, expression.alwaysActive)
		}
	}

}

func TestDateTimeExpression_IsIn(t *testing.T) {
	testDataList := []struct {
		exp string
		t   time.Time
		in  bool
	}{
		{
			exp: "[1991-2021][01-03][04-05][08:00:00-10:00:00,11:00:00-12:00:00]",
			t:   time.Date(1991, 01, 03, 8, 0, 1, 2, time.Local),
			in:  false,
		},
		{
			exp: "[1991-2021][01-03][04-05][08:00:00-10:00:00,11:00:00-12:00:00]",
			t:   time.Date(1991, 01, 04, 8, 0, 1, 2, time.Local),
			in:  true,
		},
		{
			exp: "[1991-2021][01-03][04-05][08:00:00-10:00:00,11:00:00-12:00:00]",
			t:   time.Date(2021, 01, 03, 8, 0, 1, 2, time.Local),
			in:  false,
		},
		{
			exp: "[1991-2021][01-03][04-05][08:00:00-10:00:00,11:00:00-12:00:00]",
			t:   time.Date(2021, 01, 04, 8, 0, 1, 2, time.Local),
			in:  true,
		},
		{
			exp: "[1991-2021][01-03][04-05][08:00:00-10:00:00,11:00:00-12:00:00]",
			t:   time.Date(2022, 01, 04, 8, 0, 1, 2, time.Local),
			in:  false,
		},
		{
			exp: "[1991-2021][01-03][04-05][08:00:00-10:00:00,11:00:00-12:00:00]",
			t:   time.Date(2021, 01, 04, 10, 0, 0, 0, time.Local),
			in:  false,
		},
		{
			exp: "[1991-2021][*][*][*]",
			t:   time.Date(2021, 12, 31, 11, 59, 59, 0, time.Local),
			in:  true,
		},
	}

	for _, testData := range testDataList {
		expr, err := NewDateTimeExpression(testData.exp)
		if err != nil {
			panic(err)
		}
		in := expr.IsIn(testData.t)
		ok := assert.Equal(t, testData.in, in)
		if !ok {
			t.Error(testData.t)
		}
	}

}

func TestDateTimeExpression_GetStartTime_Year(t *testing.T) {
	testDatas := []struct {
		index  int
		exp    string
		input  time.Time
		err    error
		result time.Time
	}{
		{
			index:  0,
			exp:    "[*][*][*][*]",
			input:  time.Now(),
			err:    ErrAlwaysActiveNoStartTime,
			result: time.Time{},
		},
		{
			index:  1,
			exp:    "[2000-2002][*][*][*]",
			input:  time.Date(1999, time.January, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local),
		},
		{
			index:  2,
			exp:    "[2000-2002][*][*][*]",
			input:  time.Date(2000, time.January, 1, 1, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local),
		},
		{
			index:  3,
			exp:    "[2000-2002][*][*][*]",
			input:  time.Date(2001, time.January, 1, 1, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local),
		},
		{
			index:  4,
			exp:    "[2000-2002][*][*][*]",
			input:  time.Date(2002, time.January, 1, 1, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.Local),
		},
		{
			index:  5,
			exp:    "[2000-2002][*][*][*]",
			input:  time.Date(2003, time.January, 1, 1, 0, 0, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
	}

	for i, data := range testDatas {
		fmt.Printf("[%d] exp:%s\n", i, data.exp)
		expr, err := NewDateTimeExpression(data.exp)
		if err != nil {
			panic(err)
		}

		startTime, err := expr.GetStartTime(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, startTime)
	}
}

func TestDateTimeExpression_GetStartTime_Month(t *testing.T) {
	testDatas := []struct {
		exp    string
		input  time.Time
		err    error
		result time.Time
	}{
		{
			exp:    "[*][05][*][*]",
			input:  time.Date(2000, 4, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 5, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][05][*][*]",
			input:  time.Date(2000, 6, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 5, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][05-07][*][*]",
			input:  time.Date(2000, 4, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 5, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][05-07][*][*]",
			input:  time.Date(2000, 6, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 5, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][05-07][*][*]",
			input:  time.Date(2000, 8, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 5, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][05-07][*][*]",
			input:  time.Date(1999, 8, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 5, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][05-07][*][*]",
			input:  time.Date(2000, 6, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 5, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2001][05-07][*][*]",
			input:  time.Date(2000, 4, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 5, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2001][05-07][*][*]",
			input:  time.Date(2001, 6, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 5, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2001][05-07][*][*]",
			input:  time.Date(2002, 6, 1, 0, 0, 0, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
	}

	for i, data := range testDatas {
		fmt.Printf("[%d] exp:%s\n", i, data.exp)
		expr, err := NewDateTimeExpression(data.exp)
		if err != nil {
			panic(err)
		}

		startTime, err := expr.GetStartTime(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, startTime)
	}
}

func TestDateTimeExpression_GetStartTime_Day(t *testing.T) {
	testDatas := []struct {
		exp    string
		input  time.Time
		err    error
		result time.Time
	}{
		{
			exp:    "[*][*][05][*]",
			input:  time.Date(2000, 4, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 4, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][*][05][*]",
			input:  time.Date(2000, 4, 6, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 5, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][*][05-10][*]",
			input:  time.Date(2000, 4, 4, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 4, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][*][05-10][*]",
			input:  time.Date(2000, 4, 6, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 4, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][*][05-10][*]",
			input:  time.Date(2000, 4, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 5, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][10][05-10][*]",
			input:  time.Date(2000, 4, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 10, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][10][05-10][*]",
			input:  time.Date(2000, 10, 6, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 10, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][10][05-10][*]",
			input:  time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 10, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][08-10][05-10][*]",
			input:  time.Date(2000, 4, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][08-10][05-10][*]",
			input:  time.Date(2000, 8, 8, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][08-10][05-10][*]",
			input:  time.Date(2000, 8, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 9, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][08-10][05-10][*]",
			input:  time.Date(2000, 9, 7, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 9, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][08-10][05-10][*]",
			input:  time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 8, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][08-10][05-10][*]",
			input:  time.Date(2000, 4, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][08-10][05-10][*]",
			input:  time.Date(2000, 8, 6, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][08-10][05-10][*]",
			input:  time.Date(2000, 8, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 9, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][08-10][05-10][*]",
			input:  time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
		{
			exp:    "[2000-2002][08-10][05-10][*]",
			input:  time.Date(1999, 10, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2002][08-10][05-10][*]",
			input:  time.Date(2000, 9, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 10, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2002][08-10][05-10][*]",
			input:  time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 8, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2002][08-10][05-10][*]",
			input:  time.Date(2001, 9, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 10, 5, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2002][08-10][05-10][*]",
			input:  time.Date(2002, 10, 11, 0, 0, 0, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
	}

	for i, data := range testDatas {
		fmt.Printf("[%d] exp:%s\n", i, data.exp)
		expr, err := NewDateTimeExpression(data.exp)
		if err != nil {
			panic(err)
		}

		startTime, err := expr.GetStartTime(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, startTime)
	}
}

func TestDateTimeExpression_GetStartTime_Hour(t *testing.T) {
	testDatas := []struct {
		exp    string
		input  time.Time
		err    error
		result time.Time
	}{
		{
			exp:    "[*][*][*][8:00:00-10:00:00]",
			input:  time.Date(2000, time.January, 1, 7, 30, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][*][8:00:00-10:00:00]",
			input:  time.Date(2000, time.January, 1, 9, 30, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][*][8:00:00-10:00:00]",
			input:  time.Date(2000, time.January, 1, 10, 30, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 2, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][*][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 1, 10, 30, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 11, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][*][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 1, 11, 30, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 11, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][*][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 1, 12, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 2, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 1, 12, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 8, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 10, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 11, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 11, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 11, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 12, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 0, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 8, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 13, 00, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 6, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 6, 11, 00, 01, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 6, 11, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 7, 13, 00, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 7, 13, 00, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 6, 8, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 6, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 8, 8, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2001, time.February, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 6, 8, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 6, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.March, 3, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.March, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.March, 6, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.March, 6, 11, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.March, 8, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.April, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.April, 8, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2001, time.February, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 8, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 6, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 6, 11, 00, 00, 0, time.Local),
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 8, 0, 0, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.March, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.April, 8, 0, 0, 00, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.April, 7, 12, 30, 30, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
	}
	for i, data := range testDatas {
		fmt.Printf("[%d] exp:%s\n", i, data.exp)
		expr, err := NewDateTimeExpression(data.exp)
		if err != nil {
			panic(err)
		}

		startTime, err := expr.GetStartTime(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, startTime)
	}
}

func TestDateTimeExpression_GetEndTime_Year(t *testing.T) {
	testDatas := []struct {
		index  int
		exp    string
		input  time.Time
		err    error
		result time.Time
	}{
		{
			index:  0,
			exp:    "[*][*][*][*]",
			input:  time.Now(),
			err:    ErrNoEnd,
			result: time.Time{},
		},
		{
			index:  1,
			exp:    "[2000-2002][*][*][*]",
			input:  time.Date(1999, time.January, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2003, time.January, 1, 0, 0, 0, 0, time.Local),
		},
		{
			index:  2,
			exp:    "[2000-2002][*][*][*]",
			input:  time.Date(2000, time.January, 1, 1, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2003, time.January, 1, 0, 0, 0, 0, time.Local),
		},
		{
			index:  3,
			exp:    "[2000-2002][*][*][*]",
			input:  time.Date(2001, time.January, 1, 1, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2003, time.January, 1, 0, 0, 0, 0, time.Local),
		},
		{
			index:  4,
			exp:    "[2000-2002][*][*][*]",
			input:  time.Date(2003, time.January, 1, 0, 0, 0, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
		{
			index:  5,
			exp:    "[2000-2002][*][*][*]",
			input:  time.Date(2003, time.January, 1, 1, 0, 0, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
	}

	for i, data := range testDatas {
		fmt.Printf("[%d] exp:%s\n", i, data.exp)
		expr, err := NewDateTimeExpression(data.exp)
		if err != nil {
			panic(err)
		}

		enTime, err := expr.GetEndTime(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, enTime)
	}

}

func TestDateTimeExpression_GetEndTime_Month(t *testing.T) {
	testDatas := []struct {
		exp    string
		input  time.Time
		err    error
		result time.Time
	}{
		{
			exp:    "[*][05][*][*]",
			input:  time.Date(2000, 4, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 6, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][05][*][*]",
			input:  time.Date(2000, 6, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 6, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][05-07][*][*]",
			input:  time.Date(2000, 4, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][05-07][*][*]",
			input:  time.Date(2000, 6, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][05-07][*][*]",
			input:  time.Date(2000, 8, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 8, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][05-07][*][*]",
			input:  time.Date(1999, 8, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][05-07][*][*]",
			input:  time.Date(2000, 6, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2001][05-07][*][*]",
			input:  time.Date(2000, 4, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2001][05-07][*][*]",
			input:  time.Date(2001, 6, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 8, 1, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2001][05-07][*][*]",
			input:  time.Date(2001, 8, 1, 0, 0, 0, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
	}

	for i, data := range testDatas {
		fmt.Printf("[%d] exp:%s\n", i, data.exp)
		expr, err := NewDateTimeExpression(data.exp)
		if err != nil {
			panic(err)
		}

		endTime, err := expr.GetEndTime(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, endTime)
	}
}

func TestDateTimeExpression_GetEndTime_Day(t *testing.T) {
	testDatas := []struct {
		exp    string
		input  time.Time
		err    error
		result time.Time
	}{
		{
			exp:    "[*][*][05][*]",
			input:  time.Date(2000, 4, 1, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 4, 6, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][*][05][*]",
			input:  time.Date(2000, 4, 6, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 5, 6, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][*][05-10][*]",
			input:  time.Date(2000, 4, 4, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 4, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][*][05-10][*]",
			input:  time.Date(2000, 4, 6, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 4, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][*][05-10][*]",
			input:  time.Date(2000, 4, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 5, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][10][05-10][*]",
			input:  time.Date(2000, 4, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][10][05-10][*]",
			input:  time.Date(2000, 10, 6, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][10][05-10][*]",
			input:  time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 10, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][08-10][05-10][*]",
			input:  time.Date(2000, 4, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][08-10][05-10][*]",
			input:  time.Date(2000, 8, 8, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][08-10][05-10][*]",
			input:  time.Date(2000, 8, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 9, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][08-10][05-10][*]",
			input:  time.Date(2000, 9, 7, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 9, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[*][08-10][05-10][*]",
			input:  time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 8, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][08-10][05-10][*]",
			input:  time.Date(2000, 4, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][08-10][05-10][*]",
			input:  time.Date(2000, 8, 6, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][08-10][05-10][*]",
			input:  time.Date(2000, 8, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 9, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000][08-10][05-10][*]",
			input:  time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
		{
			exp:    "[2000-2002][08-10][05-10][*]",
			input:  time.Date(1999, 10, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 8, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2002][08-10][05-10][*]",
			input:  time.Date(2000, 9, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2002][08-10][05-10][*]",
			input:  time.Date(2000, 10, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 8, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2002][08-10][05-10][*]",
			input:  time.Date(2001, 9, 11, 0, 0, 0, 0, time.Local),
			err:    nil,
			result: time.Date(2001, 10, 11, 0, 0, 0, 0, time.Local),
		},
		{
			exp:    "[2000-2002][08-10][05-10][*]",
			input:  time.Date(2002, 10, 11, 0, 0, 0, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
	}

	for i, data := range testDatas {
		fmt.Printf("[%d] exp:%s\n", i, data.exp)
		expr, err := NewDateTimeExpression(data.exp)
		if err != nil {
			panic(err)
		}

		endTime, err := expr.GetEndTime(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, endTime)
	}
}

func TestDateTimeExpression_GetEndTime_Hour(t *testing.T) {
	testDatas := []struct {
		exp    string
		input  time.Time
		err    error
		result time.Time
	}{
		{
			exp:    "[*][*][*][8:00:00-10:00:00]",
			input:  time.Date(2000, time.January, 1, 7, 30, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][*][8:00:00-10:00:00]",
			input:  time.Date(2000, time.January, 1, 9, 30, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][*][8:00:00-10:00:00]",
			input:  time.Date(2000, time.January, 1, 10, 30, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 2, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][*][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 1, 10, 30, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 12, 30, 30, 0, time.Local),
		},
		{
			exp:    "[*][*][*][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 1, 11, 30, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 1, 12, 30, 30, 0, time.Local),
		},
		{
			exp:    "[*][*][*][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 1, 12, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 2, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 1, 12, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 8, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 10, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 12, 30, 30, 0, time.Local),
		},
		{
			exp:    "[*][*][05][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 11, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 12, 30, 30, 0, time.Local),
		},
		{
			exp:    "[*][*][05][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 12, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 0, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 8, 40, 30, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 5, 13, 00, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 6, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][*][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 6, 11, 00, 01, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.January, 6, 12, 30, 30, 0, time.Local),
		},
		{
			exp:    "[*][*][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 7, 13, 00, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 7, 13, 00, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 6, 8, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 6, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 8, 8, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2001, time.February, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 6, 8, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 6, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.March, 3, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.March, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.March, 6, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.March, 6, 12, 30, 30, 0, time.Local),
		},
		{
			exp:    "[*][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.March, 8, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.April, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[*][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.April, 8, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2001, time.February, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 8, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 6, 11, 30, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 6, 12, 30, 30, 0, time.Local),
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 8, 0, 0, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.March, 5, 10, 00, 00, 0, time.Local),
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.April, 8, 0, 0, 00, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
	}
	for i, data := range testDatas {
		fmt.Printf("[%d] exp:%s\n", i, data.exp)
		expr, err := NewDateTimeExpression(data.exp)
		if err != nil {
			panic(err)
		}

		endTime, err := expr.GetEndTime(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, endTime)
	}
}

func TestDateTimeExpression_GetNextStartTime(t *testing.T) {
	testDatas := []struct {
		exp    string
		input  time.Time
		err    error
		result time.Time
	}{
		{
			exp:    "[*][*][*][*]",
			input:  time.Date(2000, time.January, 8, 0, 0, 00, 0, time.Local),
			err:    ErrAlwaysActiveNoStartTime,
			result: time.Time{},
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.January, 8, 0, 0, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 5, 8, 00, 00, 0, time.Local),
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.February, 5, 9, 0, 00, 0, time.Local),
			err:    nil,
			result: time.Date(2000, time.February, 5, 11, 00, 00, 0, time.Local),
		},
		{
			exp:    "[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]",
			input:  time.Date(2000, time.April, 8, 0, 0, 00, 0, time.Local),
			err:    ErrOutOfDate,
			result: time.Time{},
		},
	}

	for i, data := range testDatas {
		fmt.Printf("[%d] exp:%s\n", i, data.exp)
		expr, err := NewDateTimeExpression(data.exp)
		if err != nil {
			panic(err)
		}

		nextStartTime, err := expr.GetNextStartTime(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, nextStartTime)
	}
}
