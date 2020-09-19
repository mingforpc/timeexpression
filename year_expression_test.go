package timeexpression

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewYearExpression(t *testing.T) {
	// 定义并初始化并赋值给 data
	testDatas := []struct {
		exp   string
		err   error
		start int
		end   int
		IsAll bool
	}{
		{
			exp:   "*",
			err:   nil,
			start: 0,
			end:   MaxYear,
			IsAll: true,
		},
		{
			exp:   "1991",
			err:   nil,
			start: 1991,
			end:   1991,
			IsAll: false,
		},
		{
			exp:   "1991-2005",
			err:   nil,
			start: 1991,
			end:   2005,
			IsAll: false,
		},
		{
			exp: "a",
			err: &time.ParseError{
				Layout:     "2006",
				Value:      "a",
				LayoutElem: "2006",
				ValueElem:  "a",
				Message:    "",
			},
		},
		{
			exp: "1999-a",
			err: &time.ParseError{
				Layout:     "2006",
				Value:      "a",
				LayoutElem: "2006",
				ValueElem:  "a",
				Message:    "",
			},
		},
	}

	for _, data := range testDatas {

		exp, err := newYearExpression(data.exp)
		if err != nil {
			assert.EqualError(t, data.err, err.Error())
			continue
		}
		assert.Equal(t, data.start, exp.start)
		assert.Equal(t, data.end, exp.end)
		assert.Equal(t, data.IsAll, exp.isAll)
	}
}

func TestYearExpression_IsIn(t *testing.T) {
	// 定义并初始化并赋值给 data
	testDatas := []testData{
		{
			t:      "1991",
			result: true,
		},
		{
			t:      "0001",
			result: true,
		},
		{
			t:      "2099",
			result: true,
		},
	}
	expression, err := newYearExpression("*")
	if err != nil {
		t.Fatal(err)
	}
	testIsInYear(t, expression, testDatas)

	testDatas = []testData{
		{
			t:      "1991",
			result: true,
		},
		{
			t:      "0001",
			result: false,
		},
		{
			t:      "2099",
			result: false,
		},
	}
	expression, err = newYearExpression("1991-2098")
	if err != nil {
		t.Fatal(err)
	}
	testIsInYear(t, expression, testDatas)

	testDatas = []testData{
		{
			t:      "1991",
			result: true,
		},
		{
			t:      "0001",
			result: false,
		},
		{
			t:      "2099",
			result: false,
		},
	}
	expression, err = newYearExpression("1991")
	if err != nil {
		t.Fatal(err)
	}
	testIsInYear(t, expression, testDatas)
}

func testIsInYear(t *testing.T, expression *yearExpression, testDatas []testData) {
	for _, data := range testDatas {
		yearTime, err := time.Parse("2006", data.t)
		if err != nil {
			t.Fatal(err)
		}

		in := expression.isIn(yearTime.Year())

		if in != data.result {
			t.Errorf("data[%s] want[%v] but get[%v]", data.t, data.result, in)
		}
	}
}

func TestYearExpression_GetStart(t *testing.T) {
	// 定义并初始化并赋值给 data
	testDatas := []struct {
		exp    string
		input  int
		err    error
		result int
	}{
		{
			exp:    "*",
			input:  1991,
			err:    nil,
			result: 1991,
		},
		{
			exp:    "2000",
			input:  1999,
			err:    nil,
			result: 2000,
		},
		{
			exp:    "2000",
			input:  2000,
			err:    nil,
			result: 2000,
		},
		{
			exp:    "2000",
			input:  2001,
			err:    ErrOutOfDate,
			result: 0,
		},
		{
			exp:    "2000-2003",
			input:  2001,
			err:    nil,
			result: 2000,
		},
		{
			exp:    "2000-2003",
			input:  1999,
			err:    nil,
			result: 2000,
		},
		{
			exp:    "2000-2003",
			input:  2004,
			err:    ErrOutOfDate,
			result: 0,
		},
	}

	for _, data := range testDatas {
		exp, err := newYearExpression(data.exp)
		if err != nil {
			panic(err)
		}

		start, err := exp.getStart(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, start)
	}
}

func TestYearExpression_GetEnd(t *testing.T) {
	// 定义并初始化并赋值给 data
	testDatas := []struct {
		exp    string
		input  int
		err    error
		result int
	}{
		{
			exp:    "*",
			input:  1991,
			err:    nil,
			result: 1991,
		},
		{
			exp:    "2000",
			input:  1999,
			err:    nil,
			result: 2000,
		},
		{
			exp:    "2000",
			input:  2000,
			err:    nil,
			result: 2000,
		},
		{
			exp:    "2000",
			input:  2001,
			err:    ErrOutOfDate,
			result: 0,
		},
		{
			exp:    "2000-2003",
			input:  2001,
			err:    nil,
			result: 2003,
		},
		{
			exp:    "2000-2003",
			input:  1999,
			err:    nil,
			result: 2003,
		},
		{
			exp:    "2000-2003",
			input:  2004,
			err:    ErrOutOfDate,
			result: 0,
		},
	}

	for _, data := range testDatas {
		exp, err := newYearExpression(data.exp)
		if err != nil {
			panic(err)
		}

		end, err := exp.getEnd(data.input)
		assert.Equal(t, data.err, err)
		assert.Equal(t, data.result, end)
	}
}
