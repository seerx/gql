package gqlh

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
	"github.com/seerx/gql/pkg/utils"
)

// ValidatorFn 自定义验证函数
type ValidatorFn func(paramName string, paramValue interface{}, inputParam graphql.ResolveParams) error

// InputValidator 输入参数验证
type InputValidator struct {
	validatorFn          ValidatorFn
	requestObjectManager *RequestObjectManager
	graphqlParam         graphql.ResolveParams
	params               map[string]*paramStatus
}

type paramStatus struct {
	Error string
}

var typeOfInputValidator = reflect.TypeOf(InputValidator{})

// Requires 必须参数
func (v *InputValidator) Requires(params ...string) {
	for _, p := range params {
		sub := strings.Split(p, ".")
		for n := range sub {
			k := strings.Join(sub[:n+1], ".")
			if msg, ok := v.params[k]; ok {
				panic(errors.New(msg.Error))
			}
		}
	}
}

func (v *InputValidator) checkValidate(paramName string, field *Field, val interface{}) {
	//tp := reflect.TypeOf(val)
	//kd := tp.Kind()
	for _, ck := range field.Prop.ValChecker {
		err := ck.Passed(val)
		if err != nil {
			panic(paramName + ":" + err.Error())
			//v.params[paramName] = &paramStatus{
			//	Error: err.Error(),
			//}
			//return
		}
	}
	//if utils.IsStringType(kd) {
	//	// TODO 字符串，可以检测 长度和正则表达式
	//	// return
	//} else if utils.IsIntType(kd) {
	//	// TODO 整形，可以检测最大值和最小值
	//	// return
	//}
	// 使用自定义函数检测参数
	if v.validatorFn != nil {
		err := v.validatorFn(paramName, val, v.graphqlParam)
		if err != nil {
			v.params[paramName] = &paramStatus{
				Error: err.Error(),
			}
		}
	}
}

func (v *InputValidator) parseField(parentParam string, input map[string]interface{}, arg *RequestObject) reflect.Value {
	res := reflect.New(arg.Param.Prop.RealType)
	elem := res.Elem()

	for _, field := range arg.Fields {
		paramKey := field.JSONName
		if parentParam != "" {
			paramKey = parentParam + "." + paramKey
		}

		inputVal, ok := input[field.JSONName]
		if !ok {
			// 未提交的参数
			v.params[paramKey] = &paramStatus{
				Error: fmt.Sprintf("参数 %s 未提交", paramKey),
			}
			continue
		}

		resField := elem.FieldByName(field.Name)
		if utils.IsStructType(field.Prop.Kind) && !utils.IsTimeType(field.Prop.RealType) {
			// 结构类型，且不是时间类型
			inputMap := inputVal.(map[string]interface{})
			argType := v.requestObjectManager.FindOrRegisterObject(field, "")
			val := v.parseField(paramKey, inputMap, argType)
			resField.Set(val)
			v.checkValidate(paramKey, field, inputVal)
			continue
		}

		val := reflect.ValueOf(inputVal)

		if val.Type() != resField.Type() {
			// 类型不匹配
			v.params[paramKey] = &paramStatus{
				Error: fmt.Sprintf("参数 %s 类型不匹配，期望类型：%s, 实际类型：%s",
					paramKey,
					resField.Type().Name(),
					val.Type().Name()),
			}
			continue
		}

		// 赋值
		resField.Set(val)
		// 规则检查
		v.checkValidate(paramKey, field, inputVal)
	}
	return res
}

// ParseInput 解析输入参数
func (v *InputValidator) ParseInput(param *graphql.ResolveParams, arg *RequestObject) (reflect.Value, error) {
	pmap, ok := param.Args[arg.Name].(map[string]interface{})
	if !ok {
		return reflect.ValueOf(nil), fmt.Errorf("Required arguments: %s. ", arg.Name)
	}
	v.params = make(map[string]*paramStatus)

	return v.parseField("", pmap, arg), nil
}

// ParseInput1 解析输入参数
func (v *InputValidator) ParseInput1(param *graphql.ResolveParams, arg *RequestObject) (reflect.Value, error) {
	pmap, ok := param.Args[arg.Name].(map[string]interface{})
	if !ok {
		return reflect.ValueOf(nil), fmt.Errorf("Required arguments: %s. ", arg.Name)
	}
	v.params = make(map[string]*paramStatus)
	// arg.Fields[9].Prop.Key()

	res := reflect.New(arg.Param.Prop.RealType)
	elem := res.Elem()

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
		v.checkValidate("", field, inputVal)
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
