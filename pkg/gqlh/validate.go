package gqlh

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
)

// InputValidator 输入参数验证
type InputValidator struct {
	params map[string]*paramStatus
}

type paramStatus struct {
	Error string
}

var typeOfInputValidator = reflect.TypeOf(InputValidator{})

// Requires 必须参数
func (v *InputValidator) Requires(params ...string) {
	for _, p := range params {
		if msg, ok := v.params[p]; ok {
			panic(errors.New(msg.Error))
		}
	}
}

func (v *InputValidator) checkValidate(field *Field, val interface{}) {
	// tp := reflect.TypeOf(val)
	// if tp.Kind() == reflect.String {
	// 	// 字符串，可以检测 长度和正则表达式
	// } else if (tp.Kind() == reflect.)
}

// 解析参数
func parseInput(param *graphql.ResolveParams, arg *RequestObject) (reflect.Value, *InputValidator, error) {
	pmap, ok := param.Args[arg.Name].(map[string]interface{})
	if !ok {
		return reflect.ValueOf(nil), nil, fmt.Errorf("Required arguments: %s. ", arg.Name)
	}

	res := reflect.New(arg.Param.Prop.RealType)
	elem := res.Elem()

	validator := &InputValidator{
		params: make(map[string]*paramStatus),
	}

	for _, field := range arg.Fields {
		// field.Prop.Type.

		inputVal := pmap[field.JSONName]
		if inputVal == nil {
			// 未提交的参数
			validator.params[field.JSONName] = &paramStatus{
				Error: fmt.Sprintf("参数 %s 未提交", field.JSONName),
			}
			continue
		}
		val := reflect.ValueOf(inputVal)
		resField := elem.FieldByName(field.Name)
		if val.Type() != resField.Type() {
			// 类型不匹配
			validator.params[field.JSONName] = &paramStatus{
				Error: fmt.Sprintf("参数 %s 类型不匹配，期望类型：%s, 实际类型：%s",
					field.JSONName,
					resField.Type().Name(),
					val.Type().Name()),
			}
			continue
		}
		// 赋值
		resField.Set(val)
		// 规则检查
		validator.checkValidate(field, inputVal)
	}

	return res, validator, nil
}
