# timeexpression

[![Issues](https://img.shields.io/github/issues/mingforpc/timeexpression)]()
[![Forks](https://img.shields.io/github/forks/mingforpc/timeexpression)]()
[![Stars](https://img.shields.io/github/stars/mingforpc/timeexpression)]()
[![Forks](https://img.shields.io/github/license/mingforpc/timeexpression)]()
[![Coverage](https://img.shields.io/badge/coverage-94%25-yellowgreen.svg)]()

A time expression tool, maybe useful for get the end time or start time of activity.

一个时间表达式的工具，方便用于配置活动时间。

## Expression(语法格式)

`[*,yyyy,yyyy-yyyy][*,MM,MM-MM][*,dd,dd-dd][*,hh:mm:ss-hh:mm:ss]`

The time expression is follow the principle of left closed and right open, it thinks the start time is in period, but close time not in period.

时间遵循左闭右开原则, 开始时间是认为属于周期内，结束时间认为不属于周期内:

etc: `[2001][09][10][18:00:00-19:00:00]`

So `2001-09-10 18:00:00` is in this period.

那么`2001-09-10 18:00:00`是算在范围内的

But `2001-09-10 19:00:00` is **not** in this period.

那么`2001-09-10 19:00:00`是**不**算在范围内的

## Example(例子)

```go
expr, err := timeexpression.NewDateTimeExpression("[2000][02-04][05-07][8:00:00-10:00:00,11:00:00-12:30:30]")
if err != nil {
    panic(err)
}

now := time.Date(2000, time.February, 8, 0, 0, 00, 0, time.Local)

start, _ := expr.GetStartTime(now) // 结果 2000-03-05 8:00:00 
end, _ := expr.GetEndTime(now)     // 结果 2000-03-05 10:00:00
isIn := expr.IsIn(now)             // 结果 false
```

## TODO:

1. Support `w` to represent week:
    * Expression: `[*,yyyy,yyyy-yyyy][*,MM,MM-MM][wi,wi-j][*,hh:mm:ss-hh:mm:ss]`
    * example 1: `[2000][12][w1][*]` represent all the Monday in December 2000
    * example 2: `[2000][12][w1-2][*]` represent all the Monday to Wednesday in December 2000
2. Support `r` to represent relative time
3. Support how much period had been expired, and get current period cnt.

1. 语法上支持: `w`开头配置按星期几的方式:
    * 语法为: `[*,yyyy,yyyy-yyyy][*,MM,MM-MM][wi,wi-j][*,hh:mm:ss-hh:mm:ss]`
    * 例子1: `[2000][12][w1][*]`表示2000年的12月的所有周1
    * 例子2: `[2000][12][w1-2][*]`表示2000年的12月的所有周1到周3
2. 语法上支持: `r`开头配置相对开始时间的方式
3. 支持计算周期已经开始了多少次，本次周期是第几次等函数