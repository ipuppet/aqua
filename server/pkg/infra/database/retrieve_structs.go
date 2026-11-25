package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/lib/pq"
)

// FieldInfo 存储单一对应的结构体及相关的数据库映射信息
type FieldInfo struct {
	Index        []int // 配合 reflect.Value.FieldByIndex 寻找嵌套深层的字段
	NeedsPQArray bool  // 缓存是否需要包裹 pq.Array() 以解析数组
}

var (
	// typeCache 用于缓存结构体列名和内存索引的映射，显著提升查询性能
	typeCache sync.Map // map[reflect.Type]map[string]FieldInfo

	// scannerInterface 用于比对某结构体字段是否自动实现了 sql.Scanner
	scannerInterface = reflect.TypeFor[sql.Scanner]()
)

// QueryContextExecutor is the minimal query interface required by Retrieve.
type QueryContextExecutor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// Retrieve performs a database query and maps the result to the provided struct or slice.
func Retrieve(ctx context.Context, db QueryContextExecutor, target any, query string, args ...any) error {
	reflectValue := reflect.ValueOf(target)
	if reflectValue.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	isSlice := false
	var structType reflect.Type    // 始终为 struct 类型，用于 Scan
	var sliceElemType reflect.Type // 用户期望的 slice 元素类型：*Struct 或 Struct
	elemType := reflectValue.Elem().Type()

	switch elemType.Kind() {
	case reflect.Slice:
		isSlice = true
		sliceElemType = elemType.Elem()
		switch sliceElemType.Kind() {
		case reflect.Ptr:
			if sliceElemType.Elem().Kind() != reflect.Struct {
				return errors.New("slice elements must be pointers to structs")
			}
			structType = sliceElemType.Elem()
		case reflect.Struct:
			structType = sliceElemType
		default:
			return errors.New("slice elements must be structs or pointers to structs")
		}
	case reflect.Struct:
		structType = elemType
	default:
		return errors.New("target must be a pointer to struct or slice of structs/struct pointers")
	}

	if isSlice {
		result, err := queryMultipleRows(ctx, db, structType, sliceElemType, query, args...)
		if err != nil {
			return err
		}
		reflectValue.Elem().Set(reflect.ValueOf(result))
	} else {
		result, err := querySingleRow(ctx, db, target, structType, query, args...)
		if err != nil {
			return err
		}
		// querySingleRow 内部已处理指针赋值，无需再手动 Set。
		_ = result
	}

	return nil
}

// querySingleRow queries a single row.
func querySingleRow(ctx context.Context, db QueryContextExecutor, target any, structType reflect.Type, query string, args ...any) (any, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	fieldMap := getCachedFieldMap(structType)
	params := make([]any, len(columns))
	targetElem := reflect.ValueOf(target).Elem()

	for i, col := range columns {
		if info, ok := fieldMap[col]; ok {
			fieldVal := targetElem.FieldByIndex(info.Index)
			if info.NeedsPQArray {
				params[i] = pq.Array(fieldVal.Addr().Interface())
			} else {
				params[i] = fieldVal.Addr().Interface()
			}
		} else {
			return nil, fmt.Errorf("no matching field found for column: %s in struct type %s", col, structType.Name())
		}
	}

	if err := rows.Scan(params...); err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return target, nil
}

// queryMultipleRows queries multiple rows.
func queryMultipleRows(ctx context.Context, db QueryContextExecutor, structType, sliceElemType reflect.Type, query string, args ...any) (any, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	fieldMap := getCachedFieldMap(structType)

	// Pre-compute scan plan: 对每列缓存对应的 field 信息，避免对每行重复查找 map
	scanPlan := make([]FieldInfo, len(columns))
	for i, col := range columns {
		if info, ok := fieldMap[col]; ok {
			scanPlan[i] = info
		} else {
			return nil, fmt.Errorf("no matching field found for column: %s in struct type %s", col, structType.Name())
		}
	}

	resultSlice := reflect.MakeSlice(reflect.SliceOf(sliceElemType), 0, 0)
	wantPtr := sliceElemType.Kind() == reflect.Ptr

	for rows.Next() {
		newVal := reflect.New(structType)
		elem := newVal.Elem()
		params := make([]any, len(columns))

		for i, plan := range scanPlan {
			fieldVal := elem.FieldByIndex(plan.Index)
			if plan.NeedsPQArray {
				params[i] = pq.Array(fieldVal.Addr().Interface())
			} else {
				params[i] = fieldVal.Addr().Interface()
			}
		}

		if err := rows.Scan(params...); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		if wantPtr {
			resultSlice = reflect.Append(resultSlice, newVal)
		} else {
			resultSlice = reflect.Append(resultSlice, elem)
		}
	}

	if resultSlice.Len() == 0 {
		return nil, sql.ErrNoRows
	}

	return resultSlice.Interface(), nil
}

// getCachedFieldMap 获取（如果不存在则先创建并缓存）指定 struct 类型的字段映射
func getCachedFieldMap(typ reflect.Type) map[string]FieldInfo {
	if cache, ok := typeCache.Load(typ); ok {
		return cache.(map[string]FieldInfo)
	}

	fieldMap := make(map[string]FieldInfo)
	buildFieldInfoRecursive(typ, nil, "", fieldMap)

	typeCache.Store(typ, fieldMap)
	return fieldMap
}

// buildFieldInfoRecursive 根据 reflect.Type 彻底遍历出所有字段对应列名，同时生成 Field Index 以加速访问
func buildFieldInfoRecursive(structType reflect.Type, baseIndex []int, prefix string, fieldMap map[string]FieldInfo) {
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	for i := 0; i < structType.NumField(); i++ {
		fieldDesc := structType.Field(i)

		// 跳过非嵌入的私有字段
		if fieldDesc.PkgPath != "" && !fieldDesc.Anonymous {
			continue
		}

		dbTag := fieldDesc.Tag.Get("db")
		if dbTag == "-" {
			continue
		}

		currentFieldName := dbTag
		if currentFieldName == "" {
			if !fieldDesc.Anonymous {
				currentFieldName = strings.ToLower(fieldDesc.Name)
			}
		}

		isScanner := fieldDesc.Type.Implements(scannerInterface) || reflect.PointerTo(fieldDesc.Type).Implements(scannerInterface)

		// 合并之前的深度路径构造当前深度 Index
		currentIndex := make([]int, len(baseIndex)+1)
		copy(currentIndex, baseIndex)
		currentIndex[len(baseIndex)] = i

		// 如果它是个嵌入子 struct（并没有自己声明 db:"" 参数也没有被覆盖 sql.Scanner），深入解析
		if fieldDesc.Anonymous && !isScanner && fieldDesc.Type.Kind() == reflect.Struct {
			buildFieldInfoRecursive(fieldDesc.Type, currentIndex, prefix, fieldMap)
		} else if fieldDesc.Type.Kind() == reflect.Struct && !isScanner {
			newPrefix := prefix
			if currentFieldName != "" {
				if newPrefix != "" {
					newPrefix += "."
				}
				newPrefix += currentFieldName
			}
			buildFieldInfoRecursive(fieldDesc.Type, currentIndex, newPrefix, fieldMap)
		} else {
			finalMapKey := prefix
			if currentFieldName != "" {
				if finalMapKey != "" {
					finalMapKey += "."
				}
				finalMapKey += currentFieldName
			} else if fieldDesc.Anonymous {
				defaultAnonName := strings.ToLower(fieldDesc.Name)
				if finalMapKey != "" {
					finalMapKey += "."
				}
				finalMapKey += defaultAnonName
			}

			needsPQArray := false
			if fieldDesc.Type.Kind() == reflect.Slice && fieldDesc.Type.Elem().Kind() != reflect.Uint8 {
				if !isScanner {
					needsPQArray = true
				}
			}

			if finalMapKey != "" && (prefix != "" || currentFieldName != "" || (fieldDesc.Anonymous && strings.ToLower(fieldDesc.Name) != "")) {
				fieldMap[finalMapKey] = FieldInfo{Index: currentIndex, NeedsPQArray: needsPQArray}
			} else if prefix == "" && currentFieldName == "" && fieldDesc.Anonymous {
				fieldMap[strings.ToLower(fieldDesc.Name)] = FieldInfo{Index: currentIndex, NeedsPQArray: needsPQArray}
			}
		}
	}
}
