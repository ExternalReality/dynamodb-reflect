package dynamodb_reflect

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"reflect"
	"strconv"
	"strings"
)

type Name struct {
	Name []string `dynamodb:"struct_name"`
}

func encode(ptr interface{}) (map[string]*dynamodb.AttributeValue, error) {
	values := make(map[string]*dynamodb.AttributeValue)
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("dynamodb")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}

		value, err := toDynamodbAttributeValue(v.Field(i))

		if err != nil {
			return nil, err
		}
		values[name] = value
	}

	return values, nil
}

func toDynamodbAttributeValue(value reflect.Value) (*dynamodb.AttributeValue, error) {
	var valueType reflect.Kind

	switch value.Kind() {
	case reflect.Map, reflect.Ptr, reflect.Array, reflect.Slice:
		valueType = value.Type().Elem().Kind()
		fmt.Printf("type: %s", valueType.String())
	default:
		valueType = value.Kind()
	}

	switch valueType {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:

		if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
			var stringSlice []string
			for i := 0; i < value.Len(); i++ {
				stringSlice = append(stringSlice, strconv.FormatInt(value.Index(i).Int(), 10))
			}
			return &dynamodb.AttributeValue{NS: aws.StringSlice(stringSlice)}, nil
		}

		stringValue := strconv.FormatInt(value.Int(), 10)
		return &dynamodb.AttributeValue{N: aws.String(stringValue)}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:

		if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
			var stringSlice []string
			for i := 0; i < value.Len(); i++ {
				stringSlice = append(stringSlice, strconv.FormatUint(value.Index(i).Uint(), 10))
			}
			return &dynamodb.AttributeValue{NS: aws.StringSlice(stringSlice)}, nil
		}

		stringValue := strconv.FormatUint(value.Uint(), 10)
		return &dynamodb.AttributeValue{N: aws.String(stringValue)}, nil

	case reflect.String:
		if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
			var stringSlice []string
			for i := 0; i < value.Len(); i++ {
				stringSlice = append(stringSlice, value.Index(i).String())
			}
			return &dynamodb.AttributeValue{SS: aws.StringSlice(stringSlice)}, nil
		}

		return &dynamodb.AttributeValue{S: aws.String(value.String())}, nil

	case reflect.Bool:
		return &dynamodb.AttributeValue{BOOL: aws.Bool(value.Bool())}, nil
	}

	return nil, fmt.Errorf("cannot convert type to dynamodb attribute value %s", value.Kind().String())
}
