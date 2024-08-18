package database

import (
	"fmt"
	"log/slog"
	"time"
)

func InsertToken(token string, rawRefreshToken string, userGUID string) error {
	tx := DB.MustBegin()

	_, err := tx.Exec("INSERT INTO refresh_tokens (rawRefreshToken, token_hash, user_guid) VALUES ($1, $2, $3)", rawRefreshToken, token, userGUID)

	if err != nil {
		slog.Error("error insertToken")
		return err
	}

	tx.Commit()

	return nil
}

func GetToken(rawRefreshToken string) (Token, error) {
	var userGUID string
    var createdAt time.Time
	var tokenHash string

    query := `SELECT user_guid, created_at, token_hash FROM refresh_tokens WHERE rawRefreshToken = $1`

    err := DB.QueryRow(query, rawRefreshToken).Scan(&userGUID, &createdAt, &tokenHash)
	if err != nil {
		slog.Error(fmt.Sprintf("GetToken() error = %v, can't get value from DB", err))
		return Token{}, err
	}

	token := Token{
		RawRefreshToken: rawRefreshToken,
		TokenHash: tokenHash,
		UserGUID: userGUID,
		CreatedAt: createdAt,
	}
	return token, nil
}

func DeleteToken(rawRefreshToken string) (error) {
	tx := DB.MustBegin()

	_, err := tx.Exec("DELETE FROM refresh_tokens WHERE rawRefreshToken=$1", rawRefreshToken)
	if err != nil {
		slog.Error(fmt.Sprintf("PopToken() error = %v, can't delete token", err))
		return err
	}
	tx.Commit()
	return nil
}

type Token struct {
	RawRefreshToken string `db:"rawRefreshToken"`
	TokenHash       string `db:"token_hash"`
	UserGUID        string `db:"user_guid"`
	CreatedAt       time.Time `db:"created_at"`
}
