package def

import (
	"reflect"

	"github.com/seerx/gql/pkg/utils"
)

// FuncInfo 函数定义
type FuncInfo struct {
	Pkg    string        // 所在包
	Name   string        // 函数名称
	Func   reflect.Value // 函数
	Type   reflect.Type  // 函数类型
	Struct interface{}   // 函数所属结构,可能是 nil
}

// GetStructName 获取函数所在结构的名称
func (fi *FuncInfo) GetStructName() string {
	if fi.Struct != nil {
		structType := reflect.TypeOf(fi.Struct)
		return structType.Name()
	}

	return ""
}

// GetDescribe 获取函数的描述信息
// 函数名称以 Desc 结尾的，多时描述信息提供函数
func (fi *FuncInfo) GetDescribe() string {
	if fi.Struct != nil {
		descFnName := fi.Name + "Desc"
		structType := reflect.TypeOf(fi.Struct)
		method, ok := structType.MethodByName(descFnName)
		if ok {
			tp := method.Type
			if tp.NumIn() == 1 && tp.NumOut() == 1 {
				// 该函数没有输入参数，且只有一个输出参数
				op := tp.Out(0)
				if utils.IsStringType(op.Kind()) {
					// 输出参数是 string 类型
					res := method.Func.Call([]reflect.Value{reflect.ValueOf(fi.Struct)})
					v, o := res[0].Interface().(string)
					if o {
						return v
					}
				}
			}
		}
	}

	return ""
}
