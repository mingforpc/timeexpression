package timeexpression

import (
	"strconv"
	"strings"
	"time"
)

// hourUnit 时分秒
type hourUnit struct {
	Hour   int
	Minute int
	Sec    int
}

// newHourTimeUnit 表达式要为hh:mm:ss
func newHourTimeUnit(unitStr string) (hourUnit, error) {
	var err error
	unit := hourUnit{}
	vals := strings.Split(unitStr, ":")
	if len(vals) != 3 {
		// 表达式为 00:00:00 必须分为3段
		return hourUnit{}, ErrHourUnitFormat
	}
	unit.Hour, err = strconv.Atoi(vals[0])
	if err != nil {
		return hourUnit{}, err
	}
	unit.Minute, err = strconv.Atoi(vals[1])
	if err != nil {
		return hourUnit{}, err
	}
	unit.Sec, err = strconv.Atoi(vals[2])
	if err != nil {
		return hourUnit{}, err
	}

	// 校验解析出来的值
	err = unit.check()
	if err != nil {
		return hourUnit{}, err
	}

	return unit, nil
}

// check 检验值
func (unit *hourUnit) check() error {
	// 检验时
	if unit.Hour < 0 || unit.Hour > 24 {
		return ErrHourUnitFormat
	}
	// 检验分
	if unit.Minute < 0 || unit.Minute > 59 {
		return ErrHourUnitFormat
	}
	// 检验秒
	if unit.Sec < 0 || unit.Sec > 59 {
		return ErrHourUnitFormat
	}

	// 只能24:00:00
	if unit.Hour == 24 && (unit.Minute > 0 || unit.Sec > 0) {
		return ErrHourUnitFormat
	}

	return nil
}

// before 判断是否在目标时间之前，相等也返回true
func (unit *hourUnit) before(target hourUnit) bool {

	if unit.Hour < target.Hour {
		return true
	} else if unit.Hour > target.Hour {
		return false
	}

	if unit.Minute < target.Minute {
		return true
	} else if unit.Minute > target.Minute {
		return false
	}
	if unit.Sec <= target.Sec {
		return true
	}

	return false
}

// before 判断是否在目标时间之后，相等也返回false
func (unit *hourUnit) after(target hourUnit) bool {
	return !unit.before(target)
}

// toSec 转换为秒
func (unit *hourUnit) toSec() int {
	duration := time.Duration(unit.Hour)*time.Hour + time.Duration(unit.Minute)*time.Minute + time.Duration(unit.Sec)*time.Second
	return int(duration.Seconds())
}
