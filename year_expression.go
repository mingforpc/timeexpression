package timeexpression

import (
	"errors"
	"strings"
	"time"
)

func parseYearInt(yearStr string) (int, error) {
	year, err := time.Parse("2006", yearStr)
	if err != nil {
		return 0, err
	}
	return year.Year(), nil
}

type yearExpression struct {
	start int
	end   int
	isAll bool
}

//newYearExpression 创建年的时间表达式,支持格式为 [*,yyyy,yyyy-yyyy]
func newYearExpression(expression string) (*yearExpression, error) {
	yearExpression := &yearExpression{}

	expression = strings.Trim(expression, " ")
	if expression == "*" {
		// *的情况
		yearExpression.start = 0
		yearExpression.end = MaxYear
		yearExpression.isAll = true
		return yearExpression, nil
	}

	splitYearStr := strings.Split(expression, "-")
	if len(splitYearStr) > 2 {
		return nil, ErrYearFormat
	}

	var err error
	yearExpression.start, err = parseYearInt(splitYearStr[0])
	if err != nil {
		return nil, err
	}

	if len(splitYearStr) == 2 {
		yearExpression.end, err = parseYearInt(splitYearStr[1])
		if err != nil {
			return nil, err
		}
	} else {
		yearExpression.end = yearExpression.start
	}

	err = yearExpression.check()
	if err != nil {
		return nil, err
	}

	return yearExpression, nil
}

// check 检查参数
func (expression *yearExpression) check() error {
	if expression.start > expression.end {
		return errors.New("year error: start after end")
	}

	return nil
}

// isIn 是否在范围内
func (expression *yearExpression) isIn(year int) bool {

	if year >= expression.start && year <= expression.end {
		return true
	}

	return false
}

// getStart 获取开始年
// 1. 如果在周期内,则返回本次周期的年
// 2. 如果在周期外
//    在开始前,返回开始的年
//    在范围后，则返回错误
func (expression *yearExpression) getStart(year int) (int, error) {
	if expression.isAll {
		return year, nil
	}

	if expression.end < year {
		return 0, ErrOutOfDate
	}

	return expression.start, nil
}

// getEnd 获取结束年, 仅支持周期内
// PS:如果年没有超出期限，则结束年为配置周期的结束年
func (expression *yearExpression) getEnd(year int) (int, error) {
	if expression.isAll {
		return year, nil
	}

	if expression.end < year {
		return 0, ErrOutOfDate
	}

	return expression.end, nil
}
