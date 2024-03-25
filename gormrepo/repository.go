package gormrepo

import (
	"gorm.io/gorm"
)

type GenericRepository[T interface{}] interface {
	DB() *gorm.DB
	GetModel() T
	Debug() GenericRepository[T]
	Preload(query string, args ...interface{}) GenericRepository[T]
	All() ([]T, error)
	Create(ent T) (*T, error)
	Update(pk uint, ent T) (*T, error)
	Save(ent T) (*T, error)
	Find(pk uint) (*T, error)
	FindByEntity(ent T) (*T, error)
	FindByAttribute(attr string, value interface{}) (*T, error)
	Get(fn func(tx *gorm.DB) (*gorm.DB, error)) ([]T, error)
	GetByEntity(ent T) ([]T, error)
	GetByAttributes(attrs map[string]interface{}) ([]T, error)
	Delete(pk uint) (bool, error)
}

type Preload struct {
	query string
	args  []interface{}
}

type genericRepository[T interface{}] struct {
	model   T
	db      *gorm.DB
	debug   bool
	preload []*Preload
}

func (g *genericRepository[T]) DB() *gorm.DB {
	db := g.db

	if g.debug {
		db = db.Debug()
	}

	if g.preload != nil && len(g.preload) != 0 {
		for _, preload := range g.preload {
			db = db.Preload(preload.query, preload.args...)
		}
	}

	return db
}

func (g *genericRepository[T]) GetModel() T {
	return g.model
}

func (g *genericRepository[T]) Debug() GenericRepository[T] {
	g.debug = true
	return &genericRepository[T]{
		model:   g.model,
		db:      g.db,
		debug:   true,
		preload: g.preload,
	}
}

func (g *genericRepository[T]) Preload(query string, args ...interface{}) GenericRepository[T] {
	preload := &Preload{
		query: query,
		args:  args,
	}

	var preloads []*Preload

	if g.preload != nil {
		preloads = append(preloads, g.preload...)
	}

	preloads = append(preloads, preload)

	return &genericRepository[T]{
		model:   g.model,
		db:      g.db,
		debug:   g.debug,
		preload: preloads,
	}
}

func (g *genericRepository[T]) All() ([]T, error) {
	models := make([]T, 0)

	err := g.DB().Find(&models).Error

	return models, err
}

func (g *genericRepository[T]) Create(ent T) (*T, error) {
	err := g.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(&ent).Error
	})

	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (g *genericRepository[T]) Update(pk uint, ent T) (*T, error) {
	model := g.GetModel()
	err := g.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&model, pk).Error; err != nil {
			return err
		}

		if err := tx.Model(&model).Updates(&ent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	err = g.DB().First(&model, pk).Error

	return &model, err
}

func (g *genericRepository[T]) Save(ent T) (*T, error) {
	err := g.DB().Transaction(func(tx *gorm.DB) error {
		return tx.Save(&ent).Error
	})

	if err != nil {
		return nil, err
	}

	return &ent, nil
}

func (g *genericRepository[T]) Find(pk uint) (*T, error) {
	model := g.GetModel()
	err := g.DB().First(&model, pk).Error

	if err != nil {
		return nil, err
	}

	return &model, err
}

func (g *genericRepository[T]) FindByEntity(ent T) (*T, error) {
	model := g.GetModel()
	if err := g.DB().Where(&ent).First(&model).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (g *genericRepository[T]) FindByAttribute(attr string, value interface{}) (*T, error) {
	model := g.GetModel()
	err := g.DB().Where(map[string]interface{}{attr: value}).First(&model).Error
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (g *genericRepository[T]) Get(fn func(tx *gorm.DB) (*gorm.DB, error)) ([]T, error) {
	models := make([]T, 0)
	err := g.DB().Transaction(func(tx *gorm.DB) error {
		tx, err := fn(tx)
		if err != nil {
			return err
		}
		return tx.Find(&models).Error
	})

	return models, err
}

func (g *genericRepository[T]) GetByEntity(ent T) ([]T, error) {
	models := make([]T, 0)
	err := g.DB().Where(&ent).Find(&models).Error

	return models, err
}

func (g *genericRepository[T]) GetByAttributes(attrs map[string]interface{}) ([]T, error) {
	models := make([]T, 0)
	err := g.DB().Where(attrs).Find(&models).Error

	return models, err
}

func (g *genericRepository[T]) Delete(pk uint) (bool, error) {
	model := g.GetModel()
	err := g.DB().First(&model, pk).Error
	if err != nil {
		return false, err
	}

	err = g.db.Delete(&model).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func NewGenericRepository[T interface{}](db *gorm.DB, model T) GenericRepository[T] {
	return &genericRepository[T]{
		model:   model,
		db:      db,
		debug:   false,
		preload: nil,
	}
}
