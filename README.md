# timeexpression

[![Issues](https://img.shields.io/github/issues/mingforpc/timeexpression)]()
[![Forks](https://img.shields.io/github/forks/mingforpc/timeexpression)]()
[![Stars](https://img.shields.io/github/stars/mingforpc/timeexpression)]()
[![Forks](https://img.shields.io/github/license/mingforpc/timeexpression)]()
[![Coverage](https://img.shields.io/badge/coverage-93%25-yellowgreen.svg)]()

GitHub forks:
GitHub stars:
GitHub license:


一个时间表达式的工具，方便用于配置活动时间。

## 语法格式

`[*,yyyy,yyyy-yyyy][*,MM,MM-MM][*,dd,dd-dd][*,hh:mm:ss-hh:mm:ss]`

时间遵循左闭右开原则, 开始时间是认为属于周期内，结束时间认为不属于周期内:

etc: `[2001][09][10][18:00:00-19:00:00]`

那么`2001-09-10 18:00:00`是算在范围内的

那么`2001-09-10 19:00:00`是**不**算在范围内的

## 例子

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

0. 更好的参数检验和检查
1. 语法上支持: `w`开头配置按星期几的方式
2. 语法上支持: `r`开头配置相对开始时间的方式
3. 支持计算周期已经开始了多少次，本次周期是第几次等函数