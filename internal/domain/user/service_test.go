package user_test

import (
	"context"
	"errors"

	// "fmt"
	"testing"

	"github.com/chetanuchiha16/go-play/db"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/pkg/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Table(t *testing.T) {
	// 1. Define the table structure
	testcases := []struct {
		name           string               // Name of the test case
		inputArgs      user.CreateUserShema // Input payload
		mockReturnUser db.User              // What the mock should return
		mockReturnErr  error                // What error the mock should return
		wantErr        bool                 // Do we expect an error?
	}{
		{
			name: "Success",
			inputArgs: user.CreateUserShema{
				Name:     "Chetan Kishor",
				Email:    "chetan16ck@gmail.com",
				Password: "password",
			},
			mockReturnUser: db.User{
				ID:           1,
				Name:         "Chetan Kishor",
				PasswordHash: ";ajdfjaodja", // Simulated hashed password
				Email:        "chetan16ck@gmail.com",
				CreatedAt:    pgtype.Timestamptz{},
			},
			mockReturnErr: nil,
			wantErr:       false,
		},
		{
			name: "Database Error - Duplicate Email",
			inputArgs: user.CreateUserShema{
				Name:     "Duplicate User",
				Email:    "chetan16ck@gmail.com",
				Password: "password123",
			},
			mockReturnUser: db.User{},
			mockReturnErr:  errors.New("unique constraint violation: email already exists"),
			wantErr:        true,
		},
	}

	for _, testcase := range testcases {
		// 2. Run as a sub-test
		t.Run(testcase.name, func(t *testing.T) {
			mockStore := mocks.NewMockUserStore(t)

			// 3. Setup the mock based on the table row
			// Using mock.Anything for context and the DB user model passed to the store
			mockStore.On("CreateUser", mock.Anything, mock.Anything).
				Return(testcase.mockReturnUser, testcase.mockReturnErr)

			service := user.NewUserService(mockStore)

			// 4. Execute
			resUser, err := service.CreateUser(context.Background(), testcase.inputArgs)

			// 5. Assertions based on the table row
			if testcase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testcase.inputArgs.Name, resUser.Name)
				assert.Equal(t, testcase.inputArgs.Email, resUser.Email)
				// Ensure the password returned is the hash, not the plain text input
				assert.NotEqual(t, testcase.inputArgs.Password, resUser.PasswordHash)
			}

			// 6. Assert that the mocked method was actually called
			mockStore.AssertExpectations(t)
		})
	}
}

func TestGetUser_Table(t *testing.T) {
	// 1. Define the table structure
	testcases := []struct {
		name          string  // Name of the test case
		userID        int64   // Input
		mockReturn    db.User // What the mock should return
		mockErr       error   // What error the mock should return
		wantErr       bool    // Do we expect an error?
		expectedEmail string  // Expected result field
	}{
		{
			name:   "Success",
			userID: 1,
			mockReturn: db.User{
				ID:    1,
				Name:  "Chetan",
				Email: "chetan@example.com",
			},
			mockErr:       nil,
			wantErr:       false,
			expectedEmail: "chetan@example.com",
		},
		{
			name:          "User Not Found",
			userID:        999,
			mockReturn:    db.User{},
			mockErr:       pgx.ErrNoRows,
			wantErr:       true,
			expectedEmail: "",
		},
		{
			name:          "Database Connection Error",
			userID:        500,
			mockReturn:    db.User{},
			mockErr:       errors.New("connection pool exhausted"),
			wantErr:       true,
			expectedEmail: "",
		},
	}

	for _, testcase := range testcases {
		// 2. Run as a sub-test
		t.Run(testcase.name, func(t *testing.T) {
			mockStore := mocks.NewMockUserStore(t)

			// 3. Setup the mock based on the table row
			mockStore.On("GetUser", mock.Anything, testcase.userID).
				Return(testcase.mockReturn, testcase.mockErr)

			service := user.NewUserService(mockStore)

			// 4. Execute
			resUser, err := service.GetUser(context.Background(), testcase.userID)

			// 5. Assertions based on the table row
			if testcase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testcase.expectedEmail, resUser.Email)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	mockStore := mocks.NewMockUserStore(t)
	mockStore.On("DeleteUser", mock.Anything, int64(1)).Return(nil)

	userService := user.NewUserService(mockStore)

	err := userService.DeleteUser(context.Background(), 1)

	assert.NoError(t, err)
}

func TestListUsers(t *testing.T) {
	mockStore := mocks.NewMockUserStore(t)
	dbUsers := []db.ListUsersRow{{
		ID:   1,
		Name: "Chetan Kishor",
		// PasswordHash: ";ajdfjaodja",
		Email:     "chetan16ck@gmail.com",
		CreatedAt: pgtype.Timestamptz{},
	},
		{
			ID:   2,
			Name: "Chetan Kishor",
			// PasswordHash: ";ajwwfjaodja",
			Email:     "ck1234@gmail.com",
			CreatedAt: pgtype.Timestamptz{},
		},
	}
	mockStore.On("ListUsers", mock.Anything, int32(2)).Return(dbUsers, nil)

	userService := user.NewUserService(mockStore)
	users, err := userService.ListUsers(t.Context(), int32(2))

	assert.NoError(t, err)
	assert.Equal(t, dbUsers, users)
}
