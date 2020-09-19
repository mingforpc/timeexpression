package timeexpression

import (
	"errors"
	"strings"
)

//hourExpression 解析小时的表达式
type hourExpression struct {
	hourUnits []*hourUnitExpression
	isAll     bool
}

// newHourExpression 格式为 *,hh:mm:ss-hh:mm:ss[,hh:mm:ss-hh:mm:ss]...
func newHourExpression(hourStr string) (*hourExpression, error) {
	hourStr = strings.Trim(hourStr, " ")

	var hourUnits []*hourUnitExpression
	isAll := false
	splitHourStr := strings.Split(hourStr, ",")
	for _, splitStr := range splitHourStr {
		unitExpression, err := newHourUnitExpression(splitStr)
		if err != nil {
			return nil, err
		}
		if unitExpression.isAll {
			isAll = true
		}

		hourUnits = append(hourUnits, unitExpression)
	}

	expression := &hourExpression{
		hourUnits: hourUnits,
		isAll:     isAll,
	}

	return expression, nil
}

func (expression *hourExpression) isIn(hour, minute, sec int) bool {
	for _, unitExpression := range expression.hourUnits {
		if unitExpression.isIn(hour, minute, sec) {
			return true
		}
	}

	return false
}

// getStart 获取开始的时分秒
// 1. 如果在周期内,则返回本次周期的时分秒
// 2. 如果在周期外
//    在开始前,返回开始的时分秒
//    在范围后，则返回错误
func (expression *hourExpression) getStart(hour, min, sec int) (hourUint hourUnit, addDay bool, err error) {

	// 尝试是否有范围内的
	for _, unit := range expression.hourUnits {
		if unit.isIn(hour, min, sec) {
			return unit.start, false, nil
		}
	}

	// 都在范围外了, 则找一个将来最接近的
	paramUnit := hourUnit{
		Hour:   hour,
		Minute: min,
		Sec:    sec,
	}
	paramSec := paramUnit.toSec()
	var targetUnit *hourUnit
	var minUnit *hourUnit
	for _, unit := range expression.hourUnits {
		if minUnit == nil || minUnit.toSec() > unit.start.toSec() {
			minUnit = &unit.start
		}

		if paramSec >= unit.end.toSec() {
			continue
		}
		if targetUnit == nil || targetUnit.toSec() > unit.start.toSec() {
			targetUnit = &unit.start
		}

	}

	if targetUnit != nil {
		return *targetUnit, false, nil
	}

	if minUnit == nil {
		return hourUnit{}, false, errors.New("hourExpression getStart get unreachable error")
	}

	return *minUnit, true, nil
}

// getEnd 获取结束的时分秒
// 1. 如果在周期内,则返回本次周期的结束时分秒
// 2. 如果在周期外
//    在结束前,返回结束的时分秒
//    在范围后，则返回错误
func (expression *hourExpression) getEnd(hour, min, sec int) (hourUint hourUnit, addDay bool, err error) {
	// 尝试是否有范围内的
	for _, unit := range expression.hourUnits {
		if unit.isIn(hour, min, sec) {
			return unit.end, false, nil
		}
	}

	// 都在范围外了, 则找一个将来最接近的
	paramUnit := hourUnit{
		Hour:   hour,
		Minute: min,
		Sec:    sec,
	}
	paramSec := paramUnit.toSec()
	var targetUnit *hourUnit
	var minUnit *hourUnit
	for _, unit := range expression.hourUnits {
		if minUnit == nil || minUnit.toSec() > unit.end.toSec() {
			minUnit = &unit.end
		}

		if paramSec > unit.end.toSec() {
			continue
		}
		if targetUnit == nil || targetUnit.toSec() > unit.end.toSec() {
			targetUnit = &unit.end
		}

	}

	if targetUnit != nil {
		return *targetUnit, false, nil
	}

	if minUnit == nil {
		return hourUnit{}, false, errors.New("hourExpression getEnd get unreachable error")
	}

	return *minUnit, true, nil
}
