package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidateWithErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: App{"v1.1"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: ErrStringLen},
			},
		},
		{
			in: User{
				ID:     "1",
				Name:   "Nick",
				Age:    20,
				Email:  "test",
				Role:   "user",
				Phones: []string{"89999999999"},
				meta:   json.RawMessage("does not matter"),
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrStringLen},
				ValidationError{Field: "Email", Err: ErrStringRegexp},
				ValidationError{Field: "Role", Err: ErrStringOutOfList},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			assert.True(t, errors.As(err, &tt.expectedErr))
			assert.Equal(t, tt.expectedErr.Error(), err.Error())
		})
	}
}

func TestValidateWithoutErrors(t *testing.T) {
	tests := []interface{}{
		App{Version: "v1234"},
		Token{},
		User{
			ID:     "fbac1f75-cb30-4048-9519-5b2156602762",
			Name:   "Nick",
			Age:    20,
			Email:  "test@test.com",
			Role:   "admin",
			Phones: []string{"89999999999"},
			meta:   json.RawMessage("does not matter"),
		},
		Response{
			Code: 200,
			Body: "does not matter",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)

			assert.Nil(t, err)
		})
	}
}
