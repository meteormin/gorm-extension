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
	Name              string `gorm:"column:name;unique"`
	TestRelationModel TestRelationModel
}

type TestRelationModel struct {
	gorm.Model
	TestEntityId uint `gorm:"column:test_entity_id"`
	Seq          int
}

func connectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(TestEntity{}, TestRelationModel{})
	if err != nil {
		panic(err)
	}

	db.Debug().Exec("DELETE FROM `test_entities`;")
	db.Debug().Exec("UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'test_entities';")
	db.Debug().Exec("DELETE FROM `test_relation_models`;")
	db.Debug().Exec("UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'test_relation_models';")
	return db
}

func TestNewGenericRepository(t *testing.T) {
	db := connectDB()

	repo = gormrepo.NewGenericRepository(db, TestEntity{})

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
		return
	}

	if find.ID != 1 {
		t.Errorf("ID Match Fail: req: %d, real: %d", 1, find.ID)
		return
	}

	log.Print(find)
}

func TestGenericRepository_Create(t *testing.T) {
	_, err := repo.Debug().FindByAttribute("name", "TEST")
	if err != nil {
		t.Error(err)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Error("exists ")
			return
		}
	}

	create, err := repo.Debug().Create(TestEntity{Name: "TEST", TestRelationModel: TestRelationModel{Seq: 1}})
	if err != nil {
		t.Error(err)
		return
	}

	if create == nil {
		t.Error("failed create... result is nil")
		return
	}

	if create.Name != "TEST" {
		t.Errorf("Not Match Name: req: %s, real: %s", create.Name, "TEST")
		return
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

func TestGenericRepository_GetByAttributes(t *testing.T) {
	getByAttr, err := repo.Debug().GetByAttributes(map[string]interface{}{"name": "TEST"})
	if err != nil {
		t.Error(err)
		return
	}

	log.Print(getByAttr)
}

func TestGenericRepository_GetByEntity(t *testing.T) {
	getByEnt, err := repo.Debug().GetByEntity(TestEntity{
		Name: "TEST",
	})
	if err != nil {
		t.Error(err)
		return
	}

	log.Print(getByEnt)
}

func TestGenericRepository_Save(t *testing.T) {
	find, err := repo.Debug().Find(1)
	if err != nil {
		t.Error(err)
		return
	}

	find.Name = "TEST2"

	save, err := repo.Debug().Save(*find)
	if err != nil {
		t.Error(err)
		return
	}

	log.Print(save)
}

func TestGenericRepository_Update(t *testing.T) {
	update, err := repo.Debug().Update(1, TestEntity{Name: "TEST"})
	if err != nil {
		t.Error(err)
		return
	}

	log.Print(update)
}

func TestGenericRepository_Preload(t *testing.T) {
	preload, err := repo.Debug().Preload("TestRelationModel").Find(1)
	if err != nil {
		t.Error(err)
		return
	}

	log.Print(preload)
}

func TestGenericRepository_Delete(t *testing.T) {
	d, err := repo.Debug().Delete(1)
	if err != nil {
		t.Error(err)
		return
	}
	log.Print(d)
}
