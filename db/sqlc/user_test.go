package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/julysNICK/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	user1 := createRandomUser(t)

		newFullName := utils.RandomOwner()
	arg := UpdateUserParams{
		Username: user1.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.NotEqual(t, user1.FullName, user2.FullName)
	require.Equal(t, newFullName, user2.FullName)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)


}


func TestUpdateUserOnlyEmail(t *testing.T) {
	user1 := createRandomUser(t)

		newEmail := utils.RandomEmail()
	arg := UpdateUserParams{
		Username: user1.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.NotEqual(t, user1.Email, user2.Email)
	require.Equal(t, newEmail, user2.Email)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	user1 := createRandomUser(t)

		newPassword := utils.RandomString(6)
		hashedPassword, err := utils.HashPassword(newPassword)
		require.NoError(t, err)
	arg := UpdateUserParams{
		Username: user1.Username,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.NotEqual(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, hashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
}


func TestUpdateAllFields(t *testing.T) {
	user1 := createRandomUser(t)

		newFullName := utils.RandomOwner()
		newEmail := utils.RandomEmail()
		newPassword := utils.RandomString(6)
		hashedPassword, err := utils.HashPassword(newPassword)
		require.NoError(t, err)
	arg := UpdateUserParams{
		Username: user1.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.NotEqual(t, user1.HashedPassword, user2.HashedPassword)
	require.NotEqual(t, user1.FullName, user2.FullName)
	require.NotEqual(t, user1.Email, user2.Email)
	require.Equal(t, hashedPassword, user2.HashedPassword)
	require.Equal(t, newEmail, user2.Email)
	require.Equal(t, newFullName, user2.FullName)
	require.Equal(t, user1.Username, user2.Username)
}
