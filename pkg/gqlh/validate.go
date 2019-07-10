package gqlh

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/seerx/gql/pkg/utils"
)

// ValidatorFn 自定义验证函数
type ValidatorFn func()

// InputValidator 输入参数验证
type InputValidator struct {
	validatorFn ValidatorFn
	params      map[string]*paramStatus
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
	tp := reflect.TypeOf(val)
	kd := tp.Kind()
	if utils.IsStringType(kd) {
		// TODO 字符串，可以检测 长度和正则表达式
		return
	} else if utils.IsIntType(kd) {
		// TODO 整形，可以检测最大值和最小值
		return
	}
	// 使用自定义函数检测参数
	v.validatorFn()
}

func (v *InputValidator) parseField() {

}

// ParseInput 解析输入参数
func (v *InputValidator) ParseInput(param *graphql.ResolveParams, arg *RequestObject) (reflect.Value, error) {
	pmap, ok := param.Args[arg.Name].(map[string]interface{})
	if !ok {
		return reflect.ValueOf(nil), fmt.Errorf("Required arguments: %s. ", arg.Name)
	}

	res := reflect.New(arg.Param.Prop.RealType)
	elem := res.Elem()

	v.params = make(map[string]*paramStatus)

	for _, field := range arg.Fields {
		inputVal, ok := pmap[field.JSONName]
		if !ok {
			// 未提交的参数
			v.params[field.JSONName] = &paramStatus{
				Error: fmt.Sprintf("参数 %s 未提交", field.JSONName),
			}
			continue
		}
		val := reflect.ValueOf(inputVal)
		resField := elem.FieldByName(field.Name)
		if val.Type() != resField.Type() {
			// 类型不匹配
			v.params[field.JSONName] = &paramStatus{
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
		v.checkValidate(field, inputVal)
	}

	return res, nil
}

// 解析参数
// func parseInput(param *graphql.ResolveParams, arg *RequestObject) (reflect.Value, *InputValidator, error) {
// 	pmap, ok := param.Args[arg.Name].(map[string]interface{})
// 	if !ok {
// 		return reflect.ValueOf(nil), nil, fmt.Errorf("Required arguments: %s. ", arg.Name)
// 	}

// 	res := reflect.New(arg.Param.Prop.RealType)
// 	elem := res.Elem()

// 	validator := &InputValidator{
// 		params: make(map[string]*paramStatus),
// 	}

// 	for _, field := range arg.Fields {
// 		// field.Prop.Type.

// 		inputVal, ok := pmap[field.JSONName]
// 		if !ok {
// 			// 未提交的参数
// 			validator.params[field.JSONName] = &paramStatus{
// 				Error: fmt.Sprintf("参数 %s 未提交", field.JSONName),
// 			}
// 			continue
// 		}
// 		val := reflect.ValueOf(inputVal)
// 		resField := elem.FieldByName(field.Name)
// 		if val.Type() != resField.Type() {
// 			// 类型不匹配
// 			validator.params[field.JSONName] = &paramStatus{
// 				Error: fmt.Sprintf("参数 %s 类型不匹配，期望类型：%s, 实际类型：%s",
// 					field.JSONName,
// 					resField.Type().Name(),
// 					val.Type().Name()),
// 			}
// 			continue
// 		}
// 		// 赋值
// 		resField.Set(val)
// 		// 规则检查
// 		validator.checkValidate(field, inputVal)
// 	}

// 	return res, validator, nil
// }
