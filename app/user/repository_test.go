package user

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func seedUser(t *testing.T, db *gorm.DB, u *User) {
	t.Helper()
	if err := db.Create(u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
}

func TestRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	seedUser(t, db, &User{ID: "1", Username: "alice", Name: "Alice"})

	got, err := repo.GetByID("1")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Username != "alice" {
		t.Errorf("expected username=alice, got %s", got.Username)
	}

	_, err = repo.GetByID("999")
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	seedUser(t, db, &User{ID: "1", Username: "alice"})
	seedUser(t, db, &User{ID: "2", Username: "bob"})

	users, total, err := repo.List(0, 10)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total=2, got %d", total)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	err := repo.Create(&User{ID: "1", Username: "alice"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	var count int64
	db.Model(&User{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 user, got %d", count)
	}
}

func TestRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	seedUser(t, db, &User{ID: "1", Username: "alice", Name: "Alice"})

	err := repo.Update(&User{ID: "1", Name: "Alice Updated"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, _ := repo.GetByID("1")
	if got.Name != "Alice Updated" {
		t.Errorf("expected Name='Alice Updated', got %s", got.Name)
	}
}

func TestRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	seedUser(t, db, &User{ID: "1", Username: "alice"})

	err := repo.Delete("1")
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = repo.GetByID("1")
	if err == nil {
		t.Error("expected error after delete")
	}
}