package timeexpression

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestHourUnit_new(t *testing.T) {
	testDatas := []struct {
		exp  string
		err  error
		hour int
		min  int
		sec  int
	}{
		{
			exp:  "22:33:44",
			err:  nil,
			hour: 22,
			min:  33,
			sec:  44,
		},
		{
			exp: "2233:44",
			err: ErrHourUnitFormat,
		},
		{
			exp: "22:60:44",
			err: ErrHourUnitFormat,
		},
		{
			exp: "22:33:60",
			err: ErrHourUnitFormat,
		},
		{
			exp: "25:33:44",
			err: ErrHourUnitFormat,
		},
		{
			exp: "a:33:44",
			err: &strconv.NumError{
				Func: "Atoi",
				Num:  "a",
				Err:  errors.New("invalid syntax"),
			},
		},
		{
			exp: "23:a:44",
			err: &strconv.NumError{
				Func: "Atoi",
				Num:  "a",
				Err:  errors.New("invalid syntax"),
			},
		},
		{
			exp: "23:33:a",
			err: &strconv.NumError{
				Func: "Atoi",
				Num:  "a",
				Err:  errors.New("invalid syntax"),
			},
		},
	}

	for _, data := range testDatas {
		unit, err := newHourTimeUnit(data.exp)
		if err != nil {
			assert.EqualError(t, data.err, err.Error())
			continue
		}

		assert.Equal(t, data.hour, unit.Hour)
		assert.Equal(t, data.min, unit.Minute)
		assert.Equal(t, data.sec, unit.Sec)
	}
}

func TestHourUnit_before(t *testing.T) {
	unitBefore, err := newHourTimeUnit("22:33:44")
	if err != nil {
		panic(err)
	}
	unitAfter1, err := newHourTimeUnit("22:33:55")
	if err != nil {
		panic(err)
	}
	unitAfter2, err := newHourTimeUnit("22:33:44")
	if err != nil {
		panic(err)
	}
	unitAfter3, err := newHourTimeUnit("22:33:33")
	if err != nil {
		panic(err)
	}

	assert.True(t, unitBefore.before(unitAfter1))
	assert.True(t, unitBefore.before(unitAfter2))
	assert.False(t, unitBefore.before(unitAfter3))
}

func TestHourUnit_after(t *testing.T) {
	unitBefore, err := newHourTimeUnit("22:33:44")
	if err != nil {
		panic(err)
	}
	unitAfter1, err := newHourTimeUnit("22:33:55")
	if err != nil {
		panic(err)
	}
	unitAfter2, err := newHourTimeUnit("22:33:44")
	if err != nil {
		panic(err)
	}
	unitAfter3, err := newHourTimeUnit("22:33:33")
	if err != nil {
		panic(err)
	}

	assert.False(t, unitBefore.after(unitAfter1))
	assert.False(t, unitBefore.after(unitAfter2))
	assert.True(t, unitBefore.after(unitAfter3))
}
