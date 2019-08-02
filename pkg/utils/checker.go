package utils

import (
	"fmt"
)

// ValueChecker 数据检查
type ValueChecker interface {
	Passed(val interface{}) error
}

// IntegerMax 监测整形最大值
type IntegerMax struct {
	Max   int
	Equal bool
}

func typeError(expected string) error {
	return fmt.Errorf("expect type %s", expected)
}

func (imax *IntegerMax) Passed(val interface{}) error {
	n, ok := val.(int)
	if !ok {
		return typeError("int")
	}
	if imax.Equal {
		if n > imax.Max {
			return fmt.Errorf("the value must less then(or equal): %d", imax.Max)
		}
	} else {
		if n >= imax.Max {
			return fmt.Errorf("the value must less then: %d", imax.Max)
		}
	}
	return nil
}

// IntegerMax 监测整形最小值
type IntegerMin struct {
	Min   int
	Equal bool
}

func (imin *IntegerMin) Passed(val interface{}) error {
	n, ok := val.(int)
	if !ok {
		return typeError("int")
	}

	if imin.Equal {
		if n < imin.Min {
			return fmt.Errorf("the value must great then (or equal): %d", imin.Min)
		}
	} else {
		if n <= imin.Min {
			return fmt.Errorf("the value must great then : %d", imin.Min)
		}
	}
	return nil

	//return (imin.Equal && n >= imin.Min) || n > imin.Min
}

// FloatMax 监测整形最大值
type FloatMax struct {
	Max   float64
	Equal bool
}

func (fmax *FloatMax) Passed(val interface{}) error {
	n, ok := val.(float64)
	if !ok {
		return typeError("float")
	}
	if fmax.Equal {
		if n > fmax.Max {
			return fmt.Errorf("the value must less then (or equal): %f", fmax.Max)
		}
	} else {
		if n >= fmax.Max {
			return fmt.Errorf("the value must less then: %f", fmax.Max)
		}
	}
	return nil
}

// FloatMin 监测整形最小值
type FloatMin struct {
	Min   float64
	Equal bool
}

func (fmin *FloatMin) Passed(val interface{}) error {
	n, ok := val.(float64)
	if !ok {
		return typeError("float")
	}
	if fmin.Equal {
		if n < fmin.Min {
			return fmt.Errorf("the value must great then (or equal): %f", fmin.Min)
		}
	} else {
		if n <= fmin.Min {
			return fmt.Errorf("the value must great then: %f", fmin.Min)
		}
	}
	return nil
	//return (fmin.Equal && n >= fmin.Min) || n > fmin.Min
}

// StringMax 监测整形最大值
type StringMax struct {
	Max   int
	Equal bool
}

func (smax *StringMax) Passed(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return typeError("string")
	}
	n := len(str)
	if smax.Equal {
		if n > smax.Max {
			return fmt.Errorf("the string's length must less then(or equal): %d", smax.Max)
		}
	} else {
		if n >= smax.Max {
			return fmt.Errorf("the string's length must less then: %d", smax.Max)
		}
	}
	return nil
}

// StringMin 监测整形最小值
type StringMin struct {
	Min   int
	Equal bool
}

func (smin *StringMin) Passed(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return typeError("string")
	}
	n := len(str)

	if smin.Equal {
		if n < smin.Min {
			return fmt.Errorf("the string's length must great then(or equal): %d", smin.Min)
		}
	} else {
		if n <= smin.Min {
			return fmt.Errorf("the string's length must great then: %d", smin.Min)
		}
	}
	return nil
}
