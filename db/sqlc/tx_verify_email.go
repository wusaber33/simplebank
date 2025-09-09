package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

// VerifyEmailParams contains the input parameters of the transfer transaction
type VerifyEmailTxParams struct {
	EmailId int64
	SecretCode string
}

// VerifyEmailResult is the result of the transfer transaction
type VerifyEmailTxResult struct {
	User User
	VerifyEmail VerifyEmail 
}

//var txKey = struct{}{}


// VerifyEmail performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update account' balance within a single database transcation
func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail,err =q.UpdateVerifyEmail(ctx,UpdateVerifyEmailParams{
			ID: arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
            return err
        }
		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{
				Bool: true,
				Valid: true,
			},
		})     
        return err
	})

	return result, err
}