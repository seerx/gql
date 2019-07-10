package gqlh

import (
	"github.com/graphql-go/graphql"
	"github.com/seerx/gql/utils"
)

// ResponseObject 返回对象说明定义
type ResponseObject struct {
	Object graphql.Output
	Param  *Field // 对应的参数
}

// ResponseObjectManager  请求对象管理
type ResponseObjectManager struct {
	objectMap map[string]*ResponseObject
}

// NewResponseObjectManager 创建管理器
func NewResponseObjectManager() *ResponseObjectManager {
	return &ResponseObjectManager{
		objectMap: make(map[string]*ResponseObject),
	}
}

// FindOrRegisterObject 查找查询对象，如果找不到则注册
func (objm *ResponseObjectManager) FindOrRegisterObject(field *Field) *ResponseObject {
	list := field.Prop.IsList
	key, keyList := field.Prop.Key(), field.Prop.KeyOfList()
	var obj *ResponseObject
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
				obj = objm.registerList(field, graphql.NewList(obj.Object))
				return obj
			}
		}
	} else {
		obj, ok = objm.objectMap[key]
	}
	if !ok {
		// 没有找到，注册
		objFields := graphql.Fields{}

		p := field.Prop
		name := p.TypeName
		// typeField := new(graphql.Field)
		for n := 0; n < p.RealType.NumField(); n++ {
			field := p.RealType.Field(n)
			typeField := new(graphql.Field)
			ftype, isStruct := utils.StructFieldTypeToGraphType(&field)
			if isStruct {
				structFieldProp := utils.ParseTypeProp(field.Type)
				structFieldParam := &Field{Prop: structFieldProp}
				// 递归创建下属对象
				fobj := objm.FindOrRegisterObject(structFieldParam)
				typeField.Type = fobj.Object
			} else {
				typeField.Type = ftype
			}
			id := utils.ParseStructFieldName(&field)
			objFields[id] = typeField
		}

		// 注册单个查询对象
		gobj := graphql.NewObject(graphql.ObjectConfig{
			Name:   name,
			Fields: objFields,
		})

		obj = objm.registerObject(field, gobj)
		if list {
			// 是数组
			obj = objm.registerList(field, graphql.NewList(gobj))
		}
	}
	return obj
}

func (objm *ResponseObjectManager) registerObject(param *Field, obj *graphql.Object) *ResponseObject {
	qobj := &ResponseObject{
		Object: obj,
		Param:  param,
	}
	objm.objectMap[param.Prop.Key()] = qobj
	return qobj
}

func (objm *ResponseObjectManager) registerList(param *Field, list *graphql.List) *ResponseObject {
	qobj := &ResponseObject{
		Object: list,
		Param:  param,
	}
	objm.objectMap[param.Prop.KeyOfList()] = qobj
	return qobj
}
