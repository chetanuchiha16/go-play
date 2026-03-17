package user_test

import (
	"context"
	"testing"

	"github.com/chetanuchiha16/go-play/db"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/pkg/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	mockUserStore := mocks.NewMockUserStore(t)
	mockUserStore.On("CreateUser", mock.Anything, mock.Anything).Return(db.User{
		ID:           1,
		Name:         "Chetan Kishor",
		PasswordHash: ";ajdfjaodja",
		Email:        "chetan16ck@gmail.com",
		CreatedAt:    pgtype.Timestamptz{},
	}, nil)
	userService := user.NewService(mockUserStore)
	args := user.CreateUserShema{
		Name:     "Chetan Kishor",
		Email:    "chetan16ck@gmail.com",
		Password: "password",
	}

	user, err := userService.CreateUser(context.Background(), args)

	assert.NoError(t, err)
	assert.Equal(t, args.Name, user.Name)
	assert.NotEqual(t, args.Password, user.PasswordHash)

	mockUserStore.AssertExpectations(t)
}

func TestGetUser(t *testing.T) {
	mockUserStore := mocks.NewMockUserStore(t)
	mockUserStore.On("GetUser", mock.Anything, int64(1)).Return(db.User{
		ID:           1,
		Name:         "Chetan Kishor",
		PasswordHash: ";ajdfjaodja",
		Email:        "chetan16ck@gmail.com",
		CreatedAt:    pgtype.Timestamptz{},
	}, nil)

	userService := user.NewService(mockUserStore)
	resultUser := db.User{
		ID:           1,
		Name:         "Chetan Kishor",
		PasswordHash: ";ajdfjaodja",
		Email:        "chetan16ck@gmail.com",
		CreatedAt:    pgtype.Timestamptz{},
	}
	user, err := userService.GetUser(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, user, resultUser)
}

func TestDeleteUser(t *testing.T) {
	mockStore := mocks.NewMockUserStore(t)
	mockStore.On("DeleteUser", mock.Anything, int64(1)).Return(nil)

	userService := user.NewService(mockStore)

	err := userService.DeleteUser(context.Background(), 1)

	assert.NoError(t, err)
}

func TestListUsers(t *testing.T) {
	mockStore := mocks.NewMockUserStore(t)
	dbUsers := []db.User{{
		ID:           1,
		Name:         "Chetan Kishor",
		PasswordHash: ";ajdfjaodja",
		Email:        "chetan16ck@gmail.com",
		CreatedAt:    pgtype.Timestamptz{},
	},
	{
		ID:           2,
		Name:         "Chetan Kishor",
		PasswordHash: ";ajwwfjaodja",
		Email:        "ck1234@gmail.com",
		CreatedAt:    pgtype.Timestamptz{},
	},
}
	mockStore.On("ListUsers", mock.Anything).Return(dbUsers, nil)

	userService := user.NewService(mockStore)
	users, err := userService.ListUsers(t.Context())

	assert.NoError(t, err)
	assert.Equal(t, dbUsers, users)
}