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

var minValueValidationRegexp = regexp.MustCompile("min:(\\d+)")
var maxValueValidationRegexp = regexp.MustCompile("max:(\\d+)")

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

		var validErrors ValidationErrors
		var err error

		switch field.Type.Kind() {
		case reflect.Int:
			validErrors, err = validateInt(field.Name, int(reflectedStruct.Field(i).Int()), validationRule)
			if err != nil {
				return err
			}
		case reflect.String:
			validErrors, err = validateString(field.Name, reflectedStruct.Field(i).String(), validationRule)
			if err != nil {
				return err
			}
		case reflect.Slice:
			validErrors, err = validateSlice(field.Name, reflectedStruct.Field(i), validationRule)
			if err != nil {
				return err
			}
		default:
			return errors.New("unsupported file type")
		}

		validationErrors = append(validationErrors, validErrors...)
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateSlice(fieldName string, values reflect.Value, rule string) (ValidationErrors, error) {
	validationErrors := make(ValidationErrors, 0)

	switch values.Interface().(type) {
	case []int:
		for _, value := range values.Interface().([]int) {
			validationErrors, err := validateInt(fieldName, value, rule)
			if err != nil {
				return validationErrors, err
			}
		}
	case []string:
		for _, value := range values.Interface().([]string) {
			validationErrors, err := validateString(fieldName, value, rule)
			if err != nil {
				return validationErrors, err
			}
		}
	}

	return validationErrors, nil
}

func validateInt(fieldName string, value int, rule string) (ValidationErrors, error) {
	validationErrors := make(ValidationErrors, 0)
	isValid := false

	isListValidation := strings.HasPrefix(rule, "in:")
	if isListValidation {
		isValid = false
		rule = strings.TrimPrefix(rule, "in:")
		values := strings.Split(rule, ",")
		for _, v := range values {
			intValue, err := strconv.Atoi(v)
			if err != nil {
				return validationErrors, err
			}
			if value == intValue {
				isValid = true
				break
			}
		}
		if !isValid {
			validationError := ValidationError{Field: fieldName, Err: ErrIntOutOfList}
			validationErrors = append(validationErrors, validationError)
		}
	}

	isMinValidation := minValueValidationRegexp.Match([]byte(rule))
	isMaxValidation := maxValueValidationRegexp.Match([]byte(rule))
	isRangeLengthValidation := isMinValidation && isMaxValidation

	if isRangeLengthValidation {
		pattern := regexp.MustCompile(`\\d+`)
		lengthValues := pattern.FindAllString(rule, 2)
		if len(lengthValues) == 2 {
			minLength, err := strconv.Atoi(lengthValues[0])
			if err != nil {
				return validationErrors, err
			}
			maxLength, err := strconv.Atoi(lengthValues[1])
			if err != nil {
				return validationErrors, err
			}
			isValid = value >= minLength && value <= maxLength
			if !isValid {
				validationError := ValidationError{Field: fieldName, Err: ErrIntOutOfRange}
				validationErrors = append(validationErrors, validationError)
			}
		}
	} else if isMinValidation {
		pattern := regexp.MustCompile(`\\d+`)
		minValue := pattern.FindAllString(rule, 1)
		if len(minValue) == 1 {
			minLength, err := strconv.Atoi(minValue[0])
			if err != nil {
				return validationErrors, err
			}
			isValid = value >= minLength
			if !isValid {
				validationError := ValidationError{Field: fieldName, Err: ErrIntTooSmall}
				validationErrors = append(validationErrors, validationError)
			}
		}
	} else if isMaxValidation {
		pattern := regexp.MustCompile(`\\d+`)
		maxValue := pattern.FindAllString(rule, 1)
		if len(maxValue) == 1 {
			minLength, err := strconv.Atoi(maxValue[0])
			if err != nil {
				return validationErrors, err
			}
			isValid = value <= minLength
			if !isValid {
				validationError := ValidationError{Field: fieldName, Err: ErrIntTooLarge}
				validationErrors = append(validationErrors, validationError)
			}
		}
	}

	return validationErrors, nil
}

func validateString(fieldName string, value string, rule string) (ValidationErrors, error) {
	validationErrors := make(ValidationErrors, 0)
	isValid := false

	isRegexpValidation := strings.HasPrefix(rule, "regexp:")
	if isRegexpValidation {
		validationPattern := strings.TrimPrefix(rule, "regexp:")
		isValid, err := regexp.Match(validationPattern, []byte(value))
		if err != nil {
			return validationErrors, err
		}
		if !isValid {
			validationError := ValidationError{Field: fieldName, Err: ErrStringRegexp}
			validationErrors = append(validationErrors, validationError)
		}
	}

	isListValidation := strings.HasPrefix(rule, "in:")
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
			validationErrors = append(validationErrors, validationError)
		}
	}

	isLengthValidation, err := regexp.Match("^len:(\\d+)", []byte(rule))
	if err != nil {
		return validationErrors, err
	}
	if isLengthValidation {
		fieldLength, err := strconv.Atoi(strings.TrimPrefix(rule, "len:"))
		if err != nil {
			return validationErrors, err
		}
		isValid = len(value) == fieldLength
		if !isValid {
			validationError := ValidationError{Field: fieldName, Err: ErrStringLen}
			validationErrors = append(validationErrors, validationError)
		}
	}

	return validationErrors, nil
}
