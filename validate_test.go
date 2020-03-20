package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type validateCustom struct{ mock.Mock }

func (m *validateCustom) Validate() error { return m.Called().Error(0) }

type validateStruct struct {
	F1 bool
	F2 validateCustom
	F3 *validateCustom
}

func (sv *validateStruct) Validate() error {
	if sv.F1 == true {
		return errors.New("f1 is true and should not")
	}
	return nil
}

func TestValidate(t *testing.T) {
	var sv *validateStruct

	t.Run("nil value", func(t *testing.T) {
		require.Error(t, Validate(nil))
	})

	t.Run("one ptr is nil", func(t *testing.T) {
		sv = new(validateStruct)

		sv.F1 = false
		sv.F2.On("Validate").Return(nil).Once()

		err := Validate(&sv)
		assert.NoError(t, err)

		sv.F2.AssertExpectations(t)
	})

	t.Run("no validation error", func(t *testing.T) {
		sv = new(validateStruct)
		sv.F3 = new(validateCustom)

		sv.F1 = false
		sv.F2.On("Validate").Return(nil).Once()
		sv.F3.On("Validate").Return(nil).Once()

		err := Validate(&sv)
		assert.NoError(t, err)

		sv.F2.AssertExpectations(t)
		sv.F3.AssertExpectations(t)
	})

	t.Run("one reflect validation error", func(t *testing.T) {
		sv = new(validateStruct)
		sv.F3 = new(validateCustom)

		sv.F1 = false
		sv.F2.On("Validate").Return(nil).Once()
		sv.F3.On("Validate").Return(errors.New("should be true")).Once()

		err := Validate(&sv)
		assert.Error(t, err)
		assert.Equal(t, "validation error: field f3 should be true", err.Error())

		sv.F2.AssertExpectations(t)
		sv.F3.AssertExpectations(t)
	})

	t.Run("multiple reflect validation error", func(t *testing.T) {
		sv = new(validateStruct)
		sv.F3 = new(validateCustom)

		sv.F1 = false
		sv.F2.On("Validate").Return(errors.New("should be nil")).Once()
		sv.F3.On("Validate").Return(errors.New("should be non nil")).Once()

		err := Validate(&sv)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field f2 should be nil")
		assert.Contains(t, err.Error(), "field f3 should be non nil")

		sv.F2.AssertExpectations(t)
		sv.F3.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		sv = new(validateStruct)
		sv.F3 = new(validateCustom)

		sv.F1 = true
		sv.F2.On("Validate").Return(nil).Once()
		sv.F3.On("Validate").Return(nil).Once()

		err := Validate(&sv)
		assert.Error(t, err)
		assert.Equal(t, "validation error: f1 is true and should not", err.Error())

		sv.F2.AssertExpectations(t)
		sv.F3.AssertExpectations(t)
	})
}
