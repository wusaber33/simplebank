package gapi

import (
	"context"
	mockdb "simplebank/db/mock"
	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/token"
	"simplebank/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func TestUpdateAPI(t *testing.T) {
	user, _ := randomUser(t)

	newName := util.RandomOwner()
	newEmail := util.RandomEmail()
	invalidEmail := "invalid-email"
	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContexts func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
						Username: user.Username,
						FullName: pgtype.Text{
							String: newName,
							Valid: true,
						},
						Email:   pgtype.Text{
							String: newEmail,
							Valid: true,
						},
					}
				updatedUser := db.User{
					Username: user.Username,
					HashedPassword: user.HashedPassword,
					FullName: newName,
					Email: newEmail,
					PasswordChangedAt: user.CreatedAt,
					IsEmailVerified: user.IsEmailVerified,
				}
				store.EXPECT().
					UpdateUser(gomock.Any(),gomock.Eq(arg)).
					Times(1).
					Return(updatedUser,nil)
			},
			buildContexts: func(t *testing.T,tokenMaker token.Maker)context.Context{
				return newContextWithBearerToken(t,tokenMaker,user.Username,user.Role,time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updatedUser := res.GetUser()
				require.Equal(t, user.Username, updatedUser.Username)
				require.Equal(t, newName, updatedUser.FullName)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
		{name: "UserNotFound",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(),gomock.Any()).
					Times(1).
					Return(db.User{},db.ErrRecordNotFound)
			},
			buildContexts: func(t *testing.T,tokenMaker token.Maker)context.Context{
				return newContextWithBearerToken(t,tokenMaker,user.Username,user.Role,time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st,ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{name: "ExpiredToken",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(),gomock.Any()).
					Times(0)
			},
			buildContexts: func(t *testing.T,tokenMaker token.Maker)context.Context{
				return newContextWithBearerToken(t,tokenMaker,user.Username,user.Role,-time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st,ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{name: "NoAuthorization",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(),gomock.Any()).
					Times(0)
			},
			buildContexts: func(t *testing.T,tokenMaker token.Maker)context.Context{
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st,ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{name: "InvalidEmail",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &invalidEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(),gomock.Any()).
					Times(0)
			},
			buildContexts: func(t *testing.T,tokenMaker token.Maker)context.Context{
				return newContextWithBearerToken(t,tokenMaker,user.Username,user.Role,time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st,ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storectrl := gomock.NewController(t)
			defer storectrl.Finish()
			store := mockdb.NewMockStore(storectrl)

			tc.buildStubs(store)
			// start test server and send request
			server := newTestServer(t, store,nil)

			ctx := tc.buildContexts(t,server.tokenMaker)
			res, err := server.UpdateUser(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}

}
