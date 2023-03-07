package gormrepo_test

import (
	"errors"
	"github.com/miniyus/gorm-extension/gormrepo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"testing"
)

type TestEntity struct {
	gorm.Model
	Name string `gorm:"column:name;unique"`
}

func connectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(TestEntity{})
	if err != nil {
		panic(err)
	}

	return db
}

func TestNewGenericRepository(t *testing.T) {
	db := connectDB()

	repo := gormrepo.NewGenericRepository(db, TestEntity{})

	if repo.DB() != db {
		t.Error("fail")
	}
}

func newTestRepo() gormrepo.GenericRepository[TestEntity] {
	db := connectDB()
	return gormrepo.NewGenericRepository(db, TestEntity{})
}

var repo = newTestRepo()

func TestGenericRepository_All(t *testing.T) {
	all, err := repo.Debug().All()
	if err != nil {
		t.Error(err)
	}

	log.Print(all)
}

func TestGenericRepository_Find(t *testing.T) {
	find, err := repo.Debug().Find(1)
	if err != nil {
		t.Error(err)
	}

	if find.ID != 1 {
		t.Errorf("ID Match Fail: req: %d, real: %d", 1, find.ID)
	}

	log.Print(find)
}

func TestGenericRepository_Create(t *testing.T) {
	_, err := repo.Debug().FindByAttribute("name", "TEST")
	if err != nil {
		t.Error(err)
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Error("exists ")
		return
	}

	create, err := repo.Debug().Create(TestEntity{Name: "TEST"})
	if err != nil {
		t.Error(err)
	}

	if create == nil {
		t.Error("failed create... result is nil")
	}

	if create.Name != "TEST" {
		t.Errorf("Not Match Name: req: %s, real: %s", create.Name, "TEST")
	}

	log.Print(create)
}

func TestGenericRepository_FindByAttribute(t *testing.T) {
	findAttr, err := repo.Debug().FindByAttribute("name", "TEST")
	if err != nil {
		t.Error(err)
		return
	}

	log.Print(findAttr)
}

func TestGenericRepository_FindByEntity(t *testing.T) {
	findEnt, err := repo.Debug().FindByEntity(TestEntity{Name: "TEST"})
	if err != nil {
		t.Error(err)
		return
	}

	log.Print(findEnt)
}

func TestGenericRepository_Get(t *testing.T) {
	get, err := repo.Debug().Get(func(tx *gorm.DB) (*gorm.DB, error) {
		tx.Where("name", "TEST")
		return tx, nil
	})

	if err != nil {
		t.Error(err)
		return
	}

	log.Print(get)
}
