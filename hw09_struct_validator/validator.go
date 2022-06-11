package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	validationErrors := make([]string, 0)
	for _, err := range v {
		validationErrors = append(validationErrors, fmt.Sprintf("Field: %v, Error: %v", err.Field, err.Err))
	}

	return strings.Join(validationErrors, "\n")
}

var ErrIntOutOfList = errors.New("the field value is out of list")
var ErrIntOutOfRange = errors.New("the field value is out of range")
var ErrIntTooSmall = errors.New("the field value is too small")
var ErrIntTooLarge = errors.New("the field value is too big")
var ErrStringOutOfList = errors.New("the field value is out of list")
var ErrStringLen = errors.New("the filed value length is invalid")
var ErrStringRegexp = errors.New("the field value does not apply to regexp pattern")

func Validate(v interface{}) error {
	reflectedStruct := reflect.ValueOf(v)

	if reflectedStruct.Kind() != reflect.Struct {
		return errors.New("type must be struct")
	}

	reflectedStructType := reflectedStruct.Type()

	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < reflectedStructType.NumField(); i++ {
		field := reflectedStructType.Field(i)
		validationRule := field.Tag.Get("validate")

		if len(validationRule) == 0 {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Int:
			err := validateInt(field.Name, int(reflectedStruct.Field(i).Int()), validationRule, &validationErrors)
			if err != nil {
				return err
			}
		case reflect.String:
			err := validateString(field.Name, reflectedStruct.Field(i).String(), validationRule, &validationErrors)
			if err != nil {
				return err
			}
		case reflect.Slice:
			err := validateSlice(field.Name, reflectedStruct.Field(i), validationRule, &validationErrors)
			if err != nil {
				return err
			}
		default:
			return errors.New("unsupported file type")
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateSlice(fieldName string, values reflect.Value, rule string, errors *ValidationErrors) error {
	stringSlice, isStrings := values.Interface().([]string)
	intSlice, isIntS := values.Interface().([]int)

	if isStrings {
		for _, value := range stringSlice {
			err := validateString(fieldName, value, rule, errors)
			if err != nil {
				return nil
			}
		}
	} else if isIntS {
		for _, value := range intSlice {
			err := validateInt(fieldName, value, rule, errors)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}

func validateInt(fieldName string, value int, rule string, errors *ValidationErrors) error {
	isValid := false

	isListValidation, _ := regexp.Match("^in:", []byte(rule))
	if isListValidation {
		isValid = false
		rule = strings.TrimPrefix(rule, "in:")
		values := strings.Split(rule, ",")
		for _, v := range values {
			intValue, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			if value == intValue {
				isValid = true
				break
			}
		}
		if !isValid {
			validationError := ValidationError{Field: fieldName, Err: ErrIntOutOfList}
			*errors = append(*errors, validationError)
		}
	}

	isMinValidation, _ := regexp.Match("min:(\\d+)", []byte(rule))
	isMaxValidation, _ := regexp.Match("max:(\\d+)", []byte(rule))
	isRangeLengthValidation := isMinValidation && isMaxValidation

	if isRangeLengthValidation {
		pattern := regexp.MustCompile(`\\d+`)
		lengthValues := pattern.FindAllString(rule, 2)
		if len(lengthValues) == 2 {
			minLength, err := strconv.Atoi(lengthValues[0])
			if err != nil {
				return err
			}
			maxLength, err := strconv.Atoi(lengthValues[1])
			if err != nil {
				return err
			}
			isValid = value >= minLength && value <= maxLength
			if !isValid {
				validationError := ValidationError{Field: fieldName, Err: ErrIntOutOfRange}
				*errors = append(*errors, validationError)
			}
		}
	} else if isMinValidation {
		pattern := regexp.MustCompile(`\\d+`)
		minValue := pattern.FindAllString(rule, 1)
		if len(minValue) == 1 {
			minLength, err := strconv.Atoi(minValue[0])
			if err != nil {
				return err
			}
			isValid = value >= minLength
			if !isValid {
				validationError := ValidationError{Field: fieldName, Err: ErrIntTooSmall}
				*errors = append(*errors, validationError)
			}
		}
	} else if isMaxValidation {
		pattern := regexp.MustCompile(`\\d+`)
		maxValue := pattern.FindAllString(rule, 1)
		if len(maxValue) == 1 {
			minLength, err := strconv.Atoi(maxValue[0])
			if err != nil {
				return err
			}
			isValid = value <= minLength
			if !isValid {
				validationError := ValidationError{Field: fieldName, Err: ErrIntTooLarge}
				*errors = append(*errors, validationError)
			}
		}
	}

	return nil
}

func validateString(fieldName string, value string, rule string, errors *ValidationErrors) error {
	isValid := false

	isRegexpValidation, _ := regexp.Match("^regexp:", []byte(rule))
	if isRegexpValidation {
		validationPattern := strings.TrimPrefix(rule, "regexp:")
		isValid, err := regexp.Match(validationPattern, []byte(value))
		if err != nil {
			return err
		}
		if !isValid {
			validationError := ValidationError{Field: fieldName, Err: ErrStringRegexp}
			*errors = append(*errors, validationError)
		}
	}

	isListValidation, _ := regexp.Match("^in:", []byte(rule))
	if isListValidation {
		rule = strings.TrimPrefix(rule, "in:")
		values := strings.Split(rule, ",")
		isValid = false
		for _, v := range values {
			if value == v {
				isValid = true
				break
			}
		}
		if !isValid {
			validationError := ValidationError{Field: fieldName, Err: ErrStringOutOfList}
			*errors = append(*errors, validationError)
		}
	}

	isLengthValidation, _ := regexp.Match("^len:(\\d+)", []byte(rule))
	if isLengthValidation {
		fieldLength, err := strconv.Atoi(strings.TrimPrefix(rule, "len:"))
		if err != nil {
			return err
		}
		isValid = len(value) == fieldLength
		if !isValid {
			validationError := ValidationError{Field: fieldName, Err: ErrStringLen}
			*errors = append(*errors, validationError)
		}
	}

	return nil
}
