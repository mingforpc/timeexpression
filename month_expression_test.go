package timeexpression

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMonthExpression(t *testing.T) {
	testDataList := []struct {
		expr  string
		isAll bool
		start int
		end   int
		err   error
	}{
		{
			expr:  "*",
			isAll: true,
			start: 1,
			end:   12,
		},
		{
			expr:  "08",
			isAll: false,
			start: 8,
			end:   8,
		},
		{
			expr:  "01-11",
			isAll: false,
			start: 1,
			end:   11,
		},
		{
			expr: "01-11-1",
			err:  ErrMonthFormat,
		},
		{
			expr: "a-11",
			err: &time.ParseError{
				Layout:     "01",
				Value:      "a",
				LayoutElem: "01",
				ValueElem:  "a",
				Message:    "",
			},
		},
		{
			expr: "01-b",
			err: &time.ParseError{
				Layout:     "01",
				Value:      "b",
				LayoutElem: "01",
				ValueElem:  "b",
				Message:    "",
			},
		},
		{
			expr: "06-03",
			err:  errors.New("month error: start after end"),
		},
	}

	for _, testData := range testDataList {
		expression, err := newMonthExpression(testData.expr)
		if err != nil {
			assert.EqualError(t, testData.err, err.Error())
			continue
		}
		assert.Equal(t, testData.isAll, expression.isAll)
		assert.Equal(t, testData.start, expression.start)
		assert.Equal(t, testData.end, expression.end)
	}

}

func TestMonthExpression_IsIn(t *testing.T) {
	testExpr := "03-07"
	testDataList := []struct {
		val  int
		isIn bool
	}{
		{
			val:  3,
			isIn: true,
		},
		{
			val:  5,
			isIn: true,
		},
		{
			val:  7,
			isIn: true,
		},
		{
			val:  1,
			isIn: false,
		},
		{
			val:  12,
			isIn: false,
		},
	}

	expression, err := newMonthExpression(testExpr)
	if err != nil {
		panic(err)
	}
	for _, testData := range testDataList {
		assert.Equal(t, testData.isIn, expression.isIn(testData.val))
	}

}

func TestMonthExpression_GetStart(t *testing.T) {
	testDatas := []struct {
		exp           string
		input         int
		err           error
		resultStart   int
		resultAddYear bool
	}{
		{
			exp:           "*",
			input:         2,
			err:           nil,
			resultStart:   1,
			resultAddYear: false,
		},
		{
			exp:           "03",
			input:         2,
			err:           nil,
			resultStart:   3,
			resultAddYear: false,
		},
		{
			exp:           "03",
			input:         3,
			err:           nil,
			resultStart:   3,
			resultAddYear: false,
		},
		{
			exp:           "03",
			input:         4,
			err:           nil,
			resultStart:   3,
			resultAddYear: true,
		},
		{
			exp:           "03-05",
			input:         1,
			err:           nil,
			resultStart:   3,
			resultAddYear: false,
		},
		{
			exp:           "03-05",
			input:         4,
			err:           nil,
			resultStart:   3,
			resultAddYear: false,
		},
		{
			exp:           "03-05",
			input:         6,
			err:           nil,
			resultStart:   3,
			resultAddYear: true,
		},
		{
			exp:   "12",
			input: 13,
			err:   errors.New("monthExpression getStart get unreachable error"),
		},
	}

	for _, data := range testDatas {
		exp, err := newMonthExpression(data.exp)
		if err != nil {
			panic(err)
		}

		start, addYear, err := exp.getStart(data.input)
		if data.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, data.err.Error())
		}
		assert.Equal(t, data.resultStart, start)
		assert.Equal(t, data.resultAddYear, addYear)
	}
}

func TestMonthExpression_GetEnd(t *testing.T) {
	testDatas := []struct {
		exp           string
		input         int
		err           error
		resultEnd     int
		resultAddYear bool
	}{
		{
			exp:           "*",
			input:         2,
			err:           nil,
			resultEnd:     12,
			resultAddYear: false,
		},
		{
			exp:           "03",
			input:         2,
			err:           nil,
			resultEnd:     3,
			resultAddYear: false,
		},
		{
			exp:           "03",
			input:         3,
			err:           nil,
			resultEnd:     3,
			resultAddYear: false,
		},
		{
			exp:           "03",
			input:         4,
			err:           nil,
			resultEnd:     3,
			resultAddYear: true,
		},
		{
			exp:           "03-05",
			input:         1,
			err:           nil,
			resultEnd:     5,
			resultAddYear: false,
		},
		{
			exp:           "03-05",
			input:         4,
			err:           nil,
			resultEnd:     5,
			resultAddYear: false,
		},
		{
			exp:           "03-05",
			input:         6,
			err:           nil,
			resultEnd:     5,
			resultAddYear: true,
		},
		{
			exp:   "12",
			input: 13,
			err:   errors.New("monthExpression getEnd get unreachable error"),
		},
	}

	for _, data := range testDatas {
		exp, err := newMonthExpression(data.exp)
		if err != nil {
			panic(err)
		}

		end, addYear, err := exp.getEnd(data.input)
		if data.err == nil {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, data.err.Error())
		}
		assert.Equal(t, data.resultEnd, end)
		assert.Equal(t, data.resultAddYear, addYear)
	}
}
