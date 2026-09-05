package user

import (
	"errors"
	"testing"

	"general-agent/internal/errs"

	"gorm.io/gorm"
)

type fakeRepo struct {
	users map[string]*User
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{users: make(map[string]*User)}
}

func (f *fakeRepo) GetByID(id string) (*User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}

func (f *fakeRepo) List(offset, limit int) ([]User, int64, error) {
	users := make([]User, 0, len(f.users))
	for _, u := range f.users {
		users = append(users, *u)
	}
	total := int64(len(users))
	if offset > len(users) {
		return nil, total, nil
	}
	end := offset + limit
	if end > len(users) {
		end = len(users)
	}
	return users[offset:end], total, nil
}

func (f *fakeRepo) Create(user *User) error {
	if _, ok := f.users[user.Username]; ok {
		return gorm.ErrDuplicatedKey
	}
	f.users[user.ID] = user
	return nil
}

func (f *fakeRepo) Update(user *User) error {
	f.users[user.ID] = user
	return nil
}

func (f *fakeRepo) Delete(id string) error {
	delete(f.users, id)
	return nil
}

func TestService_GetByID_NotFound(t *testing.T) {
	svc := NewService(newFakeRepo())
	_, err := svc.GetByID("999")
	if !errors.Is(err, errs.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestService_GetByID_Success(t *testing.T) {
	repo := newFakeRepo()
	repo.users["1"] = &User{ID: "1", Username: "alice"}
	svc := NewService(repo)

	got, err := svc.GetByID("1")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Username != "alice" {
		t.Errorf("expected username=alice, got %s", got.Username)
	}
}

func TestService_List_Defaults(t *testing.T) {
	svc := NewService(newFakeRepo())
	users, total, err := svc.List(0, 0)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if total != 0 {
		t.Errorf("expected total=0, got %d", total)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestService_Create_Success(t *testing.T) {
	svc := NewService(newFakeRepo())
	user, err := svc.Create(&CreateUserReq{Username: "alice", Password: "pass", Name: "Alice"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if user.Username != "alice" {
		t.Errorf("expected username=alice, got %s", user.Username)
	}
}