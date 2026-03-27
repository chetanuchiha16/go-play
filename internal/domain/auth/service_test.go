package auth_test

import (
	"testing"

	"github.com/chetanuchiha16/go-play/db"
	"github.com/chetanuchiha16/go-play/internal/domain/auth"
	"github.com/chetanuchiha16/go-play/pkg/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	mockStore := mocks.NewMockUserStore(t)
	password := ";ajdfjaodja"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	expUser := db.User{
		ID:           1,
		Name:         "Chetan Kishor",
		PasswordHash: string(hash),
		Email:        "abcd@gmail.com",
		CreatedAt:    pgtype.Timestamptz{},
	}
	mockStore.On("GetUserByEmail", mock.Anything, "abcd@gmail.com").Return(expUser, nil)

	authService := auth.NewAuthService(mockStore, []byte("sillykey"))
	user, token, err := authService.Login(t.Context(), "abcd@gmail.com", ";ajdfjaodja")

	assert.NoError(t, err)
	assert.Equal(t, expUser.ID, user.ID)
	assert.Equal(t, expUser, user)
	assert.NotNil(t, token)
	// fmt.Print(token)

}
