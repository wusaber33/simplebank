package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
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
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	user := createRandomUser(t)
	newhashedPassword, err := util.HashedPassword(util.RandomString(6))
	require.NoError(t,err)
    arg := UpdateUserParams{
        Username:       user.Username,
		HashedPassword:	pgtype.Text{
			String: newhashedPassword,
			Valid: true,
		},
    }

    newuser,err := testStore.UpdateUser(context.Background(), arg)
    require.NoError(t, err)
	require.Equal(t,newuser.FullName,user.FullName)
	require.Equal(t,newuser.Email,user.Email)
	require.Equal(t,newuser.Username,user.Username)
	require.Equal(t,newuser.HashedPassword,newhashedPassword)
	require.NotEqual(t,newuser.HashedPassword,user.HashedPassword)
    
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	user := createRandomUser(t)
	newfullname:=util.RandomOwner()
    arg := UpdateUserParams{
        Username:       user.Username,
        FullName:      	pgtype.Text{
			String: newfullname,
            Valid: true,
		},
    }

    newuser,err := testStore.UpdateUser(context.Background(), arg)
    require.NoError(t, err)
	require.NotEqual(t,newuser.FullName,user.FullName)
	require.Equal(t,newuser.FullName,newfullname)
	require.Equal(t,newuser.Username,user.Username)
	require.Equal(t,newuser.Email,user.Email)
	require.Equal(t,newuser.HashedPassword,user.HashedPassword)
    
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	user := createRandomUser(t)
	newemail :=util.RandomEmail()
    arg := UpdateUserParams{
        Username:       user.Username,
        Email:          pgtype.Text{
			String: newemail,
            Valid: true,
		},
    }

    newuser,err := testStore.UpdateUser(context.Background(), arg)
    require.NoError(t, err)
	require.Equal(t,newuser.FullName,user.FullName)
	require.NotEqual(t,newuser.Email,user.Email)
	require.Equal(t,newuser.Username,user.Username)
	require.Equal(t,newuser.Email,newemail)
	require.Equal(t,newuser.HashedPassword,user.HashedPassword)
    
}

func TestUpdateUserAllFields(t *testing.T) {
	user := createRandomUser(t)
	newhashedPassword, err := util.HashedPassword(util.RandomString(6))
	newfullname:=util.RandomOwner()
	newemail :=util.RandomEmail()
	require.NoError(t,err)
    arg := UpdateUserParams{
        Username:       user.Username,
        FullName:      	pgtype.Text{
			String: newfullname,
            Valid: true,
		},
        Email:          pgtype.Text{
			String: newemail,
            Valid: true,
		},
		HashedPassword: pgtype.Text{
			String: newhashedPassword,
			Valid: true,
		},
    }

    newuser,err := testStore.UpdateUser(context.Background(), arg)
    require.NoError(t, err)
	require.NotEqual(t,newuser.FullName,user.FullName)
	require.Equal(t,newuser.FullName,newfullname)
	require.NotEqual(t,newuser.Email,user.Email)
	require.Equal(t,newuser.Username,user.Username)
	require.Equal(t,newuser.Email,newemail)
	require.Equal(t,newuser.HashedPassword,newhashedPassword)
	require.NotEqual(t,newuser.HashedPassword,user.HashedPassword)
    
}