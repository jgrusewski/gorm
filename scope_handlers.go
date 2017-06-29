package gorm

import (
	"database/sql/driver"
	"reflect"
)

// custom handlers

//var handlers []*customHandler

type CustomHandler interface {
	Scan(in interface{}) error
	Value() (driver.Value, error)

	GetValue() interface{}
	SetValue(interface{})
}

type customHandler struct {
	name string
	pkg  string
	typ  reflect.Type
	intf interface{}
	h    CustomHandler
}

// scan to customHandler Scan interface
func (c customHandler) Scan(intf interface{}) error {
	return c.h.Scan(intf)
}

// return custom handler driver.Value
func (c customHandler) Value() (driver.Value, error) {
	return (c.h.Value())
}

func (c *customHandler) GetValue() interface{} {
	return c.intf
}
func (c *customHandler) SetValue(intf interface{}) {
	c.h.SetValue(intf)
}

func (d *DB) lookupCustomHandler(field *Field) CustomHandler {
	if field.Field.Kind() == reflect.Ptr {
		for _, hnd := range d.handlers {
			if hnd.typ == field.Field.Type() {
				return hnd.h
			}
		}
	}
	return nil
}

func (d *DB) AddCustomHandler(intf interface{}, h CustomHandler) {
	rt := reflect.TypeOf(intf)
	d.handlers = append(d.handlers, &customHandler{
		pkg:  rt.PkgPath(),
		name: rt.Name(),
		intf: intf,
		h:    h,
		typ:  reflect.TypeOf(intf),
	})
}
