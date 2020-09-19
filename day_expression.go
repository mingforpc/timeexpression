package timeexpression

import (
	"errors"
	"strings"
	"time"
)

func parseDayInt(dayStr string) (int, error) {
	day, err := time.Parse("02", dayStr)
	if err != nil {
		return 0, err
	}
	return day.Day(), nil
}

type dayExpression struct {
	start int
	end   int
	isAll bool
}

// newDayExpression 创建日的时间表达式,支持格式为 [*,dd,dd-dd]
func newDayExpression(expression string) (*dayExpression, error) {
	dayExpression := &dayExpression{}

	expression = strings.Trim(expression, " ")
	if expression == "*" {
		// *的情况
		dayExpression.start = 1
		dayExpression.end = 31
		dayExpression.isAll = true
		return dayExpression, nil
	}
	splitDayStr := strings.Split(expression, "-")
	if len(splitDayStr) > 2 {
		return nil, ErrDayFormat
	}

	var err error
	dayExpression.start, err = parseDayInt(splitDayStr[0])
	if err != nil {
		return nil, err
	}

	if len(splitDayStr) == 2 {
		dayExpression.end, err = parseDayInt(splitDayStr[1])
		if err != nil {
			return nil, err
		}
	} else {
		dayExpression.end = dayExpression.start
	}

	err = dayExpression.check()
	if err != nil {
		return nil, err
	}

	return dayExpression, nil
}

// check 检查参数
func (expression *dayExpression) check() error {
	if expression.start > expression.end {
		return errors.New("day error: start after end")
	}

	return nil
}

func (expression *dayExpression) isIn(day int) bool {
	if expression.start <= day && expression.end >= day {
		return true
	}

	return false
}

// getStart 获取开始日期
// 1. 如果在周期内,则返回本次周期的日期
// 2. 如果在周期外,则返回下次的开始日期
func (expression *dayExpression) getStart(year int, month time.Month, day int) (start int, addMonth bool, err error) {
	idx := 1
	addMonth = false
	days := daysIn(month, year)

	for idx <= days {
		if expression.isIn(day) {
			return expression.start, addMonth, nil
		}
		day += 1
		if day > days {
			day = 1
			addMonth = true
		}
		idx += 1
	}
	return 0, false, errors.New("dayExpression getStart get unreachable error")
}

// getEnd 获取结束日期
// 1. 如果在周期内,则返回本次周期的结束日期
// 2. 如果在周期外,则返回下次的结束日期
func (expression *dayExpression) getEnd(year int, month time.Month, day int) (end int, addMonth bool, err error) {
	idx := 1
	addMonth = false
	days := daysIn(month, year)

	for idx <= days {
		if expression.isIn(day) {
			return expression.end, addMonth, nil
		}
		day += 1
		if day > days {
			day = 1
			addMonth = true
		}
		idx += 1
	}
	return 0, false, errors.New("dayExpression getEnd get unreachable error")
}
