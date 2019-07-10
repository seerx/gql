package gqlh

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/seerx/gql/utils"
)

// RequestObject 请求的参数对象定义
type RequestObject struct {
	Object *graphql.InputObject
	Param  *Field   // 对应的参数
	Name   string   // 结构名称
	Fields []*Field // 结构中的字段列表
}

// RequestObjectManager  请求对象管理
type RequestObjectManager struct {
	objectMap map[string]*RequestObject
}

// NewRequestObjectManager 创建管理器
func NewRequestObjectManager() *RequestObjectManager {
	return &RequestObjectManager{
		objectMap: make(map[string]*RequestObject),
	}
}

// FindOrRegisterObject 查找查询对象，如果找不到则注册
func (objm *RequestObjectManager) FindOrRegisterObject(field *Field, name string) *RequestObject {
	list := field.Prop.IsList
	if list {
		// 如果是列表，则返回空，不可以使用列表提交
		return nil
	}
	key := field.Prop.Key()

	obj, ok := objm.objectMap[key]

	if !ok {
		if name == "" {
			// name 为空时，必须找到对象
			panic(fmt.Errorf("Cann't find input object %s", field))
		}
		// 没有找到，注册
		objFields := graphql.InputObjectConfigFieldMap{}

		p := field.Prop
		// name := p.TypeName
		var fields []*Field
		for n := 0; n < p.RealType.NumField(); n++ {
			field := p.RealType.Field(n)
			typeField := new(graphql.InputObjectFieldConfig)
			ftype, isStruct := utils.StructFieldTypeToGraphType(&field)
			if isStruct {
				// 是结构类型，不支持
				return nil
			}
			id := utils.ParseStructFieldName(&field)
			prop := utils.ParseTypeProp(field.Type)
			// 存储结构字段
			// prop.structField = field
			fields = append(fields, &Field{
				Name:     field.Name,
				JSONName: id,
				Prop:     prop,
			})

			typeField.Type = ftype

			objFields[id] = typeField
		}

		// 注册单个查询对象
		gobj := graphql.NewInputObject(
			graphql.InputObjectConfig{
				Name:   "input" + name,
				Fields: objFields,
			})

		obj = &RequestObject{
			Object: gobj,
			Param:  field,
			Name:   name,
			Fields: fields,
		}
		objm.objectMap[key] = obj
	}
	return obj
}