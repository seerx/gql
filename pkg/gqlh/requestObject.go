package gqlh

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/seerx/gql/pkg/utils"
)

// RequestObject 请求的参数对象定义
type RequestObject struct {
	Object graphql.Input
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
	//if list {
	//	// 如果是列表，则返回空，不可以使用列表提交
	//	return nil
	//}
	key, keyList := field.Prop.Key(), field.Prop.KeyOfList()
	//obj, ok := objm.objectMap[key]
	var obj *RequestObject
	var ok bool
	if list {
		// 列表
		obj, ok = objm.objectMap[keyList]
		if !ok {
			// 没有列表
			obj, ok = objm.objectMap[key]
			if ok {
				// 有单独的注册对象
				// 注册列表即可
				obj = &RequestObject{
					Object: graphql.NewList(obj.Object),
					Name:   name,
					Param:  field,
				}
				objm.objectMap[keyList] = obj
				//obj = objm.registerList(field, graphql.NewList(obj.Object))
				return obj
			}
		}
	} else {
		obj, ok = objm.objectMap[key]
	}

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

		for n := 0; n < p.RealType.NumField(); n++ {
			field := p.RealType.Field(n)
			prop := utils.ParseTypeProp(field.Type)
			id := utils.ParseStructFieldName(&field)
			utils.ParseValueCheckers(prop, &field)
			typeField := new(graphql.InputObjectFieldConfig)
			ftype, isStruct := utils.StructFieldTypeToGraphType(&field)
			fd := &Field{
				Name:     field.Name,
				JSONName: id,
				Prop:     prop,
			}
			//graphql.Input()
			typeField.Type = ftype
			typeField.Description = prop.Desc

			if isStruct {
				// 是结构类型，递归生成
				typeField.Type = objm.FindOrRegisterObject(fd, id).Object
			}

			//if prop.IsList {
			//	if prop.Kind == reflect.Struct {
			//		// 结构列表
			//		//ftype := objm.FindOrRegisterObject(fd, id).Object
			//	} else {
			//		// 非结构
			//		typeField.Type = graphql.NewList(ftype)
			//	}
			//	// 是列表，不知道怎么支持列表，所以不允许出现列表参数
			//	//return nil
			//}

			// 存储结构字段
			// prop.structField = field
			obj.Fields = append(obj.Fields, fd)

			objFields[id] = typeField
		}

	}
	return obj
}
