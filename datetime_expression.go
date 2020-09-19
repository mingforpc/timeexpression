package timeexpression

import (
	"strings"
	"time"
)

type DateTimeExpression struct {
	year  *yearExpression
	month *monthExpression
	day   *dayExpression
	hour  *hourExpression

	alwaysActive bool // 表示该表达式是否永远有效
	hasEnd       bool // 表示是否会结束
}

// 时间表达式为[*,yyyy,yyyy-yyyy][*,mm,mm-mm][*,dd,dd-dd][*,h1-h2]
func NewDateTimeExpression(expression string) (*DateTimeExpression, error) {
	dateTimeExpression := &DateTimeExpression{}

	// 去掉最前的'['和最后的']'
	expression = strings.TrimPrefix(expression, "[")
	expression = strings.TrimSuffix(expression, "]")

	expSplits := strings.Split(expression, "][")
	if len(expSplits) != 4 {
		return nil, ErrDateTimeFormat
	}

	var err error
	// 解析年
	dateTimeExpression.year, err = newYearExpression(expSplits[0])
	if err != nil {
		return nil, err
	}
	// 解析月
	dateTimeExpression.month, err = newMonthExpression(expSplits[1])
	if err != nil {
		return nil, err
	}
	// 解析日
	dateTimeExpression.day, err = newDayExpression(expSplits[2])
	if err != nil {
		return nil, err
	}
	// 解析时
	dateTimeExpression.hour, err = newHourExpression(expSplits[3])
	if err != nil {
		return nil, err
	}

	if dateTimeExpression.year.isAll &&
		dateTimeExpression.month.isAll &&
		dateTimeExpression.day.isAll &&
		dateTimeExpression.hour.isAll {
		dateTimeExpression.alwaysActive = true
	}
	if !dateTimeExpression.year.isAll {
		dateTimeExpression.hasEnd = true
	}

	return dateTimeExpression, nil
}

// isIn 判断时间是否在表达式指定范围内
// 实现为左闭右开 etc: [2001][09][10][18:00:00-19:00:00] 那么2001-09-10 19:00:00是不算在范围内的
func (expression *DateTimeExpression) IsIn(t time.Time) bool {

	if expression.alwaysActive {
		return true
	}

	in := expression.year.isIn(t.Year())
	if !in {
		return false
	}

	in = expression.month.isIn(int(t.Month()))
	if !in {
		return false
	}

	in = expression.day.isIn(t.Day())
	if !in {
		return false
	}

	in = expression.hour.isIn(t.Hour(), t.Minute(), t.Second())
	return in
}

// GetStartTime 获取开始时间
// 1. 如果在周期内,则返回本次周期的开始时间
// 2. 如果在周期外,则返回下次周期的开始时间
func (expression *DateTimeExpression) GetStartTime(t time.Time) (time.Time, error) {
	if expression.alwaysActive {
		return time.Time{}, ErrAlwaysActiveNoStartTime
	}

	// 0.0 计算开始年
	startTime, err := expression.calculateStartYear(t)
	if err != nil {
		return time.Time{}, err
	}

	// 1.0 计算开始月
	startTime, err = expression.calculateStartMonth(startTime)
	if err != nil {
		return time.Time{}, err
	}

	// 2.0 计算开始天
	startTime, err = expression.calculateStartDay(startTime)
	if err != nil {
		return time.Time{}, err
	}

	// 3.0 计算时分秒
	startTime, err = expression.calculateStartHourUnit(startTime)
	if err != nil {
		return time.Time{}, err
	}

	return startTime, nil
}

// calculateStartYear 计算开始年
func (expression *DateTimeExpression) calculateStartYear(t time.Time) (time.Time, error) {
	startYear, err := expression.year.getStart(t.Year())
	if err != nil {
		return time.Time{}, err
	}
	if expression.month.isAll && expression.day.isAll && expression.hour.isAll ||
		t.Year() < startYear {
		// 如果只设了年，则返回开始的年
		// OR
		// 如果时间在年之前，则将时间直接挪到开始的年第一时刻
		t = time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.Local)
	}
	// 其他一般情况，都是当年开始的

	return t, nil
}

// calculateStartMonth 计算开始月
func (expression *DateTimeExpression) calculateStartMonth(t time.Time) (time.Time, error) {
	startMonth, addYear, err := expression.month.getStart(int(t.Month()))
	if err != nil {
		return time.Time{}, err
	}

	if (!expression.day.isAll || !expression.hour.isAll) && expression.month.isIn(int(t.Month())) {
		// 特殊情况 [*][*][05][*] 或者 [*][04-12][05][*]，这种月份需求是跟着当月，但开始月会获取到1的
		//  expression.month.isIn(int(t.month())) 保证了时间t是在范围内
		// starttime的月份不变
		return t, nil
	}

	if addYear {
		// 时间跨越1年，需要从新的一年的第一个时刻开始
		t = time.Date(t.Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
		if t.Year() > expression.year.end {
			return time.Time{}, ErrOutOfDate
		}
	}
	if int(t.Month()) < startMonth {
		// 当月份向后推进时，防止日时分秒跨越
		t = time.Date(t.Year(), time.Month(startMonth), 1, 0, 0, 0, 0, time.Local)
	} else {
		t = time.Date(t.Year(), time.Month(startMonth), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	}

	return t, nil
}

// calculateStartDay 计算开始天
func (expression *DateTimeExpression) calculateStartDay(t time.Time) (time.Time, error) {
	startDay, addMonth, err := expression.day.getStart(t.Year(), t.Month(), t.Day())
	if err != nil {
		return time.Time{}, err
	}

	if (!expression.hour.isAll) && expression.day.isIn(t.Day()) {
		// 特殊情况 [*][*][*][01:00:00-12:00:00] 或者 [*][*][03-09][01:00:00-12:00:00]，类似这种天需求是跟着当天，但开始天会获取到1或者配置的开始天
		//  expression.day.isIn(t.day()) 保证了时间t是在范围内
		// starttime的天不变
		return t, nil
	}

	if addMonth {
		// 时间跨越1个月，需要从新的一个月的第一个时刻开始
		t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.Local)
		if t.Year() > expression.year.end {
			return time.Time{}, ErrOutOfDate
		}
		if int(t.Month()) > expression.month.end {
			// 有可能下一年的,所以挪到下一年的进行递归
			t = time.Date(t.Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
			t, err = expression.calculateStartYear(t)
			if err != nil {
				return time.Time{}, err
			}
			t, err = expression.calculateStartMonth(t)
			if err != nil {
				return time.Time{}, err
			}

			return expression.calculateStartDay(t)
		}
	}
	if t.Day() < startDay {
		// 当日向后推进时，防止时分秒跨越
		t = time.Date(t.Year(), t.Month(), startDay, 0, 0, 0, 0, time.Local)
	} else {
		t = time.Date(t.Year(), t.Month(), startDay, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	}

	return t, nil
}

// calculateStartHourUnit 计算开始的时分秒
func (expression *DateTimeExpression) calculateStartHourUnit(t time.Time) (time.Time, error) {
	startHourUnit, addDay, err := expression.hour.getStart(t.Hour(), t.Minute(), t.Second())
	if err != nil {
		return time.Time{}, err
	}
	if addDay {
		// 推后1天，就从新的一天的第一刻开始
		t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, time.Local)
		if t.Year() > expression.year.end {
			return time.Time{}, ErrOutOfDate
		}
		if int(t.Month()) > expression.month.end {
			return time.Time{}, ErrOutOfDate
		}
		if t.Day() > expression.day.end {
			// 有可能下个月的, 递归处理
			t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.Local)
			return expression.GetStartTime(t)
		}
	}
	t = time.Date(t.Year(), t.Month(), t.Day(), startHourUnit.Hour, startHourUnit.Minute, startHourUnit.Sec,
		0, time.Local)

	return t, nil
}

//GetNextStartTime 获取下次开始时间,不管是否在周期内，都获取下次的时间
func (expression *DateTimeExpression) GetNextStartTime(t time.Time) (time.Time, error) {
	if expression.alwaysActive {
		return time.Time{}, ErrAlwaysActiveNoStartTime
	}

	in := expression.IsIn(t)
	var err error
	if in {
		// 获取当前周期的结束时间
		t, err = expression.GetEndTime(t)
		if err != nil {
			return time.Time{}, err
		}
	}

	return expression.GetStartTime(t)
}

// GetEndTime 获取结束时间,仅在周期内有效
// 实现为左闭右开 etc: [2001][09][10][18:00:00-19:00:00] 那么结束时间为2001-09-10 19:00:00(因为结束时间,已经结束了)
func (expression *DateTimeExpression) GetEndTime(t time.Time) (time.Time, error) {
	if expression.alwaysActive {
		return time.Time{}, ErrNoEnd
	}

	endTime, err := expression.calculateEndYear(t)
	if err != nil {
		return time.Time{}, err
	}

	endTime, err = expression.calculateEndMonth(endTime)
	if err != nil {
		return time.Time{}, err
	}

	endTime, err = expression.calculateEndDay(endTime)
	if err != nil {
		return time.Time{}, err
	}

	endTime, err = expression.calculateEndHourUnit(endTime)
	if err != nil {
		return time.Time{}, err
	}

	return endTime, nil
}

// calculateEndYear 计算结束年
func (expression *DateTimeExpression) calculateEndYear(t time.Time) (time.Time, error) {
	// 1.0 处理年
	endYear, err := expression.year.getEnd(t.Year())
	if err != nil {
		return time.Time{}, err
	}
	if expression.month.isAll && expression.day.isAll && expression.hour.isAll {
		// 如果只设了年，则返回结束的年
		t = time.Date(endYear, time.January, 1, 0, 0, 0, 0, time.Local)
	}
	startYear, err := expression.year.getStart(t.Year())
	if err != nil {
		return time.Time{}, err
	}
	if t.Year() < startYear {
		// 如果当前年都比开始年小，则偏移到开始年
		t = time.Date(startYear, time.January, 1, 0, 0, 0, 0, time.Local)
	}

	// 其他一般情况，都是当年结束的

	return t, nil
}

// calculateEndMonth 计算结束月
func (expression *DateTimeExpression) calculateEndMonth(t time.Time) (time.Time, error) {
	// 2.0  处理月
	endMonth, addYear, err := expression.month.getEnd(int(t.Month()))
	if err != nil {
		return time.Time{}, err
	}

	if int(t.Month()) < expression.month.start {
		t = time.Date(t.Year(), time.Month(expression.month.start), 1, 0, 0, 0, 0, time.Local)

	}

	if (!expression.day.isAll || !expression.hour.isAll) && expression.month.isIn(int(t.Month())) {
		// 特殊情况 [*][*][05][*] 或者 [*][04-12][05][*]，这种月份需求是跟着当月，但结束月会获取到12的
		//  expression.month.isIn(int(t.month())) 保证了时间t是在范围内
		// endtime的月份不变
		return t, nil
	}

	if addYear {
		// 时间跨越1年，需要从新的一年的第一个时刻开始
		t = time.Date(t.Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
		if t.Year() > expression.year.end {
			return time.Time{}, ErrOutOfDate
		}
	}

	if int(t.Month()) < endMonth {
		// 当月份向后推进时，防止日时分秒跨越
		t = time.Date(t.Year(), time.Month(endMonth), 1, 0, 0, 0, 0, time.Local)
	} else {
		t = time.Date(t.Year(), time.Month(endMonth), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	}

	return t, nil
}

// calculateEndDay 计算结束天
func (expression *DateTimeExpression) calculateEndDay(t time.Time) (time.Time, error) {
	endDay, addMonth, err := expression.day.getEnd(t.Year(), t.Month(), t.Day())
	if err != nil {
		return time.Time{}, err
	}

	if t.Day() < expression.day.start {
		t = time.Date(t.Year(), t.Month(), expression.day.start, 0, 0, 0, 0, time.Local)
	}

	if (!expression.hour.isAll) && expression.day.isIn(t.Day()) {
		// 特殊情况 [*][*][*][01:00:00-12:00:00] 或者 [*][*][03-09][01:00:00-12:00:00]，类似这种天需求是跟着当天，但结束天会获取到12或者配置的结束天
		//  expression.day.isIn(t.day()) 保证了时间t是在范围内
		// starttime的天不变
		return t, nil
	}

	if addMonth {
		// 时间跨越1个月，需要从新的一个月的第一个时刻开始
		t = time.Date(t.Year(), t.Month()+1, expression.day.start, 0, 0, 0, 0, time.Local)
		if t.Year() > expression.year.end {
			return time.Time{}, ErrOutOfDate
		}
		if int(t.Month()) > expression.month.end {
			// 有可能下一年的,所以挪到下一年的进行递归
			t = time.Date(t.Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
			t, err = expression.calculateEndYear(t)
			if err != nil {
				return time.Time{}, err
			}
			t, err = expression.calculateEndMonth(t)
			if err != nil {
				return time.Time{}, err
			}
			return expression.calculateEndDay(t)
		}
	}
	if t.Day() < endDay && expression.hour.isAll {
		// 当日向后推进时，防止时分秒跨越
		t = time.Date(t.Year(), t.Month(), endDay, 0, 0, 0, 0, time.Local)
	}

	return t, nil
}

// calculateEndHourUnit 计算结束的时分秒
func (expression *DateTimeExpression) calculateEndHourUnit(t time.Time) (time.Time, error) {
	endHourUnit, addDay, err := expression.hour.getEnd(t.Hour(), t.Minute(), t.Second())
	if err != nil {
		return time.Time{}, err
	}
	if addDay {
		// 推后1天，就从新的一天的第一刻开始
		t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, time.Local)
		if t.Year() > expression.year.end {
			return time.Time{}, ErrOutOfDate
		}
		if int(t.Month()) > expression.month.end {
			return time.Time{}, ErrOutOfDate
		}
		if t.Day() > expression.day.end {
			// 有可能下个月的, 递归处理
			t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, time.Local)
			return expression.GetEndTime(t)
		}
	}
	t = time.Date(t.Year(), t.Month(), t.Day(), endHourUnit.Hour, endHourUnit.Minute, endHourUnit.Sec,
		0, time.Local)

	return t, nil
}
