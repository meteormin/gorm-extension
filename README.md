# gorm extensions

## extensions

- generic repository
- hooks

## install
```shell
go get github.com/meteormin/gorm-extension
```

## usage

### generic repository

```go
package main

import (
	"fmt"
	"github.com/meteormin/gorm-extension/gormrepo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

const Host = "localhost"
const Username = "test"
const Password = "test"
const Dbname = "test"
const Port = "5432"
const SslMode = "false"
const TimeZone = "Asia/Seoul"

type TestModel struct {
	gorm.Model
	Name string
}

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		Host, Username, Password, Dbname, Port, SslMode, TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	repo := gormrepo.NewGenericRepository[TestModel](db, TestModel{})

	all, err := repo.All()
	if err != nil {
		panic(err)
	}

	log.Print(all)
}

```

**Interface**
```go
package gormrepo

import "gorm.io/gorm"
// GenericRepository interface
// Rules: Find* > single record, Get* > multiple records
type GenericRepository[T interface{}] interface {
	//DB Get gorm DB
	DB() *gorm.DB

	//Debug use Debug mode, print console query
	Debug() GenericRepository[T]

	//Preload preloading(eager loading)
	Preload(query string, args ...interface{}) GenericRepository[T]

	//All select * from my_table
	All() ([]T, error)

	//Create by input entity
	Create(ent T) (*T, error)

	//Update input pk record by input entity
	Update(pk uint, ent T) (*T, error)

	//Save update by input entity
	Save(ent T) (*T, error)

	//Find by pk
	Find(pk uint) (*T, error)

	//FindByEntity find by entity
	FindByEntity(ent T) (*T, error)

	//FindByAttribute find by single attribute
	FindByAttribute(attr string, value interface{}) (*T, error)

	//Get by callback
	Get(fn func(tx *gorm.DB) (*gorm.DB, error)) ([]T, error)

	//GetByEntity get by entity
	GetByEntity(ent T) ([]T, error)

	//GetByAttributes get by attributes(map)
	GetByAttributes(attrs map[string]interface{}) ([]T, error)

	//Delete by pk
	Delete(pk uint) (bool, error)
}

```

### hooks

```go
package main

import (
	"fmt"
	"github.com/meteormin/gorm-extension/gormhooks"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

const Host = "localhost"
const Username = "test"
const Password = "test"
const Dbname = "test"
const Port = "5432"
const SslMode = "false"
const TimeZone = "Asia/Seoul"

type TestModel struct {
	gorm.Model
	Name string
}

func (tm *TestModel) Hooks() gormhooks.Hooks[*TestModel] {
	return gormhooks.GetHooks(tm)
}

func (tm *TestModel) AfterFind(tx *gorm.DB) error {
	return tm.Hooks().AfterFind(tx)
}

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		Host, Username, Password, Dbname, Port, SslMode, TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	var model *TestModel // must pointer variable

	hooks := gormhooks.New(model)
	hooks.HandleAfterFind(func(q *TestModel, tx *gorm.DB) (err error) {
		log.Print(q)
		return nil
	})

	find := TestModel{}

	db.Find(&find)
}

```
