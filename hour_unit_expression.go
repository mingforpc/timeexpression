package timeexpression

import (
	"strings"
)

// hourUnitExpression 小时/分钟/秒的最小解析单位
type hourUnitExpression struct {
	start hourUnit
	end   hourUnit
	isAll bool
}

// newHourUnitExpression 格式为 *,hh:mm:ss-hh:mm:ss
func newHourUnitExpression(unitStr string) (*hourUnitExpression, error) {
	unitStr = strings.Trim(unitStr, " ")
	if unitStr == "*" {
		expression := &hourUnitExpression{
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
		}
		expression.isAll = true

		return expression, nil
	}

	expression := &hourUnitExpression{}

	hourSubList := strings.Split(unitStr, "-")
	if len(hourSubList) != 2 {
		// 表达式为00:00:00-24:00:00,必须分为2段
		return nil, ErrHourUnitFormat
	}

	var err error
	// 处理开始时间
	expression.start, err = newHourTimeUnit(hourSubList[0])
	if err != nil {
		return nil, err
	}
	// 处理结束时间
	expression.end, err = newHourTimeUnit(hourSubList[1])
	if err != nil {
		return nil, err
	}

	return expression, nil
}

// isIn 是否在表达式范围内
func (expression *hourUnitExpression) isIn(hour, minute, sec int) bool {

	// 检验小时
	if hour < expression.start.Hour || hour > expression.end.Hour {
		return false
	}
	if hour > expression.start.Hour && hour < expression.end.Hour {
		return true
	}

	if hour == expression.start.Hour {
		// 小时等于开始
		if minute < expression.start.Minute {
			return false
		}
		if minute > expression.start.Minute {
			return true
		}
		// 分钟相等,判断秒
		if sec < expression.start.Sec {
			return false
		}
		return true
	} else {
		// 小时等于结束
		if minute > expression.end.Minute {
			return false
		}
		if minute < expression.end.Minute {
			return true
		}
		// 分钟相等,判断秒
		if sec >= expression.end.Sec {
			return false
		}
		return true
	}
}
