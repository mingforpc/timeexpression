package timeexpression

import (
	"errors"
	"strings"
	"time"
)

func parseMonthInt(monthStr string) (int, error) {
	month, err := time.Parse("01", monthStr)
	if err != nil {
		return 0, err
	}
	return int(month.Month()), nil
}

type monthExpression struct {
	start int
	end   int
	isAll bool
}

// newMonthExpression 创建月的时间表达式,支持格式为 [*,mm,mm-mm]
func newMonthExpression(expression string) (*monthExpression, error) {
	monthExpression := &monthExpression{}

	expression = strings.Trim(expression, " ")
	if expression == "*" {
		// *的情况
		monthExpression.start = 1
		monthExpression.end = 12
		monthExpression.isAll = true
		return monthExpression, nil
	}
	splitMonthStr := strings.Split(expression, "-")
	if len(splitMonthStr) > 2 {
		return nil, ErrMonthFormat
	}

	var err error
	monthExpression.start, err = parseMonthInt(splitMonthStr[0])
	if err != nil {
		return nil, err
	}

	if len(splitMonthStr) == 2 {
		monthExpression.end, err = parseMonthInt(splitMonthStr[1])
		if err != nil {
			return nil, err
		}
	} else {
		monthExpression.end = monthExpression.start
	}

	return monthExpression, nil
}

// isIn 月份是否在周期内
func (expression *monthExpression) isIn(month int) bool {
	if expression.start <= month && expression.end >= month {
		return true
	}

	return false
}

// getStart 获取开始月
func (expression *monthExpression) getStart(month int) (start int, addYear bool, err error) {

	idx := 1
	addYear = false
	for idx <= int(time.December) {
		if expression.isIn(month) {
			return expression.start, addYear, nil
		}
		month += 1
		if month > int(time.December) {
			month = 1
			addYear = true
		}
		idx += 1
	}

	return 0, false, errors.New("monthExpression getStart get unreachable error")
}

// getEnd 获取结束月, 仅支持周期内
// 如果是"*"则结束月份返回12
// PS:如果年没有超出期限，则结束月为配置周期的结束月
func (expression *monthExpression) getEnd(month int) (end int, addYear bool, err error) {
	idx := 1
	addYear = false
	for idx <= int(time.December) {
		if expression.isIn(month) {
			return expression.end, addYear, nil
		}
		month += 1
		if month > int(time.December) {
			month = 1
			addYear = true
		}
		idx += 1
	}

	return 0, false, errors.New("monthExpression getEnd get unreachable error")
}
