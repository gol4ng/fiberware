// Code generated by mockery v2.4.0. DO NOT EDIT.

package mocks

import (
	logger "github.com/gol4ng/logger"
	mock "github.com/stretchr/testify/mock"
)

// LoggerInterface is an autogenerated mock type for the LoggerInterface type
type LoggerInterface struct {
	mock.Mock
}

// Alert provides a mock function with given fields: message, field
func (_m *LoggerInterface) Alert(message string, field ...logger.Field) {
	_va := make([]interface{}, len(field))
	for _i := range field {
		_va[_i] = field[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Critical provides a mock function with given fields: message, field
func (_m *LoggerInterface) Critical(message string, field ...logger.Field) {
	_va := make([]interface{}, len(field))
	for _i := range field {
		_va[_i] = field[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Debug provides a mock function with given fields: message, field
func (_m *LoggerInterface) Debug(message string, field ...logger.Field) {
	_va := make([]interface{}, len(field))
	for _i := range field {
		_va[_i] = field[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Emergency provides a mock function with given fields: message, field
func (_m *LoggerInterface) Emergency(message string, field ...logger.Field) {
	_va := make([]interface{}, len(field))
	for _i := range field {
		_va[_i] = field[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Error provides a mock function with given fields: message, field
func (_m *LoggerInterface) Error(message string, field ...logger.Field) {
	_va := make([]interface{}, len(field))
	for _i := range field {
		_va[_i] = field[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Info provides a mock function with given fields: message, field
func (_m *LoggerInterface) Info(message string, field ...logger.Field) {
	_va := make([]interface{}, len(field))
	for _i := range field {
		_va[_i] = field[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Log provides a mock function with given fields: message, level, field
func (_m *LoggerInterface) Log(message string, level logger.Level, field ...logger.Field) {
	_va := make([]interface{}, len(field))
	for _i := range field {
		_va[_i] = field[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, message, level)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Notice provides a mock function with given fields: message, field
func (_m *LoggerInterface) Notice(message string, field ...logger.Field) {
	_va := make([]interface{}, len(field))
	for _i := range field {
		_va[_i] = field[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Warning provides a mock function with given fields: message, field
func (_m *LoggerInterface) Warning(message string, field ...logger.Field) {
	_va := make([]interface{}, len(field))
	for _i := range field {
		_va[_i] = field[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, message)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}
