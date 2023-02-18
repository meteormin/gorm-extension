package gormhooks

import (
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"reflect"
)

var hooksMap = make(map[reflect.Type]interface{})

type HasHooks interface {
	callbacks.BeforeSaveInterface
	callbacks.AfterSaveInterface
	callbacks.BeforeCreateInterface
	callbacks.AfterCreateInterface
	callbacks.BeforeUpdateInterface
	callbacks.AfterUpdateInterface
	callbacks.BeforeDeleteInterface
	callbacks.AfterDeleteInterface
	callbacks.AfterFindInterface
}

type Hooks[T interface{}] struct {
	model T
	saving[T]
	creating[T]
	updating[T]
	deleting[T]
	querying[T]
	modify[T]
}

type saving[T interface{}] struct {
	beforeSave func(s T, tx *gorm.DB) (err error)
	afterSave  func(s T, tx *gorm.DB) (err error)
}
type creating[T interface{}] struct {
	beforeCreate func(c T, tx *gorm.DB) (err error)
	afterCreate  func(c T, tx *gorm.DB) (err error)
}
type updating[T interface{}] struct {
	beforeUpdate func(u T, tx *gorm.DB) (err error)
	afterUpdate  func(u T, tx *gorm.DB) (err error)
}
type deleting[T interface{}] struct {
	beforeDelete func(d T, tx *gorm.DB) (err error)
	afterDelete  func(d T, tx *gorm.DB) (err error)
}
type querying[T interface{}] struct {
	afterFind func(q T, tx *gorm.DB) (err error)
}
type modify[T interface{}] struct {
	before func(m T, tx *gorm.DB) error
}

func New[T interface{}](model T) *Hooks[T] {
	h := &Hooks[T]{
		model:    model,
		saving:   saving[T]{beforeSave: nil, afterSave: nil},
		creating: creating[T]{beforeCreate: nil, afterCreate: nil},
		updating: updating[T]{beforeUpdate: nil, afterUpdate: nil},
		deleting: deleting[T]{beforeDelete: nil, afterDelete: nil},
		querying: querying[T]{afterFind: nil},
		modify:   modify[T]{before: nil},
	}

	hooksMap[reflect.TypeOf(model)] = h
	return h
}

func GetHooks[T interface{}](model T) *Hooks[T] {
	if hooksMap[reflect.TypeOf(model)] == nil {
		return New(model)
	}

	get := hooksMap[reflect.TypeOf(model)]

	h, ok := get.(*Hooks[T])
	if !ok {
		return nil
	}

	h.model = model

	return h
}

func (h *Hooks[T]) BeforeSave(tx *gorm.DB) (err error) {
	if h.saving.beforeSave == nil {
		return nil
	}
	return h.saving.beforeSave(h.model, tx)
}

func (h *Hooks[T]) HandleBeforeSave(beforeSave func(s T, tx *gorm.DB) (err error)) {
	h.saving.beforeSave = beforeSave
}

func (h *Hooks[T]) AfterSave(tx *gorm.DB) (err error) {
	if h.saving.afterSave == nil {
		return nil
	}
	return h.saving.afterSave(h.model, tx)
}

func (h *Hooks[T]) HandleAfterSave(afterSave func(s T, tx *gorm.DB) (err error)) {
	h.saving.afterSave = afterSave
}

func (h *Hooks[T]) BeforeCreate(tx *gorm.DB) (err error) {
	if h.creating.beforeCreate == nil {
		return nil
	}
	return h.creating.beforeCreate(h.model, tx)
}

func (h *Hooks[T]) HandleBeforeCreate(beforeCreate func(c T, tx *gorm.DB) (err error)) {
	h.creating.beforeCreate = beforeCreate
}

func (h *Hooks[T]) AfterCreate(tx *gorm.DB) (err error) {
	if h.creating.afterCreate == nil {
		return nil
	}
	return h.creating.afterCreate(h.model, tx)
}

func (h *Hooks[T]) HandleAfterCreate(afterCreate func(c T, tx *gorm.DB) (err error)) {
	h.creating.afterCreate = afterCreate
}

func (h *Hooks[T]) BeforeUpdate(tx *gorm.DB) (err error) {
	if h.updating.beforeUpdate == nil {
		return nil
	}
	return h.updating.beforeUpdate(h.model, tx)
}

func (h *Hooks[T]) HandleBeforeUpdate(beforeUpdate func(u T, tx *gorm.DB) (err error)) {
	h.updating.beforeUpdate = beforeUpdate
}

func (h *Hooks[T]) AfterUpdate(tx *gorm.DB) (err error) {
	if h.updating.afterUpdate == nil {
		return nil
	}
	return h.updating.afterUpdate(h.model, tx)
}

func (h *Hooks[T]) HandleAfterUpdate(afterUpdate func(u T, tx *gorm.DB) (err error)) {
	h.updating.afterUpdate = afterUpdate
}

func (h *Hooks[T]) BeforeDelete(tx *gorm.DB) (err error) {
	if h.deleting.beforeDelete == nil {
		return nil
	}
	return h.deleting.beforeDelete(h.model, tx)
}

func (h *Hooks[T]) HandleBeforeDelete(beforeDelete func(d T, tx *gorm.DB) (err error)) {
	h.deleting.beforeDelete = beforeDelete
}

func (h *Hooks[T]) AfterDelete(tx *gorm.DB) (err error) {
	if h.deleting.afterDelete == nil {
		return nil
	}
	return h.deleting.afterDelete(h.model, tx)
}

func (h *Hooks[T]) HandleAfterDelete(afterDelete func(d T, tx *gorm.DB) (err error)) {
	h.deleting.afterDelete = afterDelete
}

func (h *Hooks[T]) AfterFind(tx *gorm.DB) (err error) {
	if h.querying.afterFind == nil {
		return nil
	}
	return h.querying.afterFind(h.model, tx)
}

func (h *Hooks[T]) HandleAfterFind(afterFind func(q T, tx *gorm.DB) (err error)) {
	h.querying.afterFind = afterFind
}

func (h *Hooks[T]) Before(tx *gorm.DB) error {
	if h.modify.before == nil {
		return nil
	}
	return h.modify.before(h.model, tx)
}

func (h *Hooks[T]) HandleBefore(before func(m T, tx *gorm.DB) error) {
	h.modify.before = before
}
