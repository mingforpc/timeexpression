package timeexpression

import "errors"

var (
	// ErrDateTimeFormat 年月日时表达格式不对
	ErrDateTimeFormat = errors.New("date time format not math")
	// ErrYearFormat 年的表达式格式不对
	ErrYearFormat = errors.New("year format not math")
	// ErrMonthFormat 月的表达式格式不对
	ErrMonthFormat = errors.New("month format not math")
	// ErrDayFormat 日的表达式格式不对
	ErrDayFormat = errors.New("day format not math")
	// ErrHourUnitFormat 时分秒的表达式格式不对
	ErrHourUnitFormat = errors.New("hour unit format not math")
	// ErrAlwaysActiveNoStartTime 表达式总是有效,所以没有开始时间
	ErrAlwaysActiveNoStartTime = errors.New("expression always active, no start time")
	// ErrOutOfDate 超过了表达式的时间范围
	ErrOutOfDate = errors.New("expression is out of date")
	// ErrNoEnd 没有结束范围
	ErrNoEnd = errors.New("expression is no end time")
)

const (
	// 支持的最大年份
	MaxYear = 9999
)
