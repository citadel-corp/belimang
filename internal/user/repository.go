package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/citadel-corp/belimang/internal/common/db"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	Create(ctx context.Context, user *Users) (err error)
	GetByUsername(ctx context.Context, username string) (user *Users, err error)
	// GetByUID(ctx context.Context, uid string) (user *Users, err error)
	// GetByID(ctx context.Context, id uint64) (user *Users, err error)
}

type dbRepository struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &dbRepository{db: db}
}

// Create implements Repository.
func (d *dbRepository) Create(ctx context.Context, user *Users) (err error) {
	createUserQuery := `
		INSERT INTO users (
			uid, username, email, hashed_password, user_type
		) VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING id;
	`
	_, err = d.db.DB().ExecContext(ctx, createUserQuery, user.UID, user.Username, user.Email, user.HashedPassword, user.UserType)
	var pgErr *pgconn.PgError
	if err != nil {
		log.Debug().Msgf("error creating user: %v", err)
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrUserAlreadyExists
			default:
				return
			}
		}
		return
	}

	return
}

func (d *dbRepository) GetByUsername(ctx context.Context, username string) (user *Users, err error) {
	getUserQuery := `
		SELECT id, uid, username, email, hashed_password, user_type, created_at
		FROM users
		WHERE username = $1;
	`
	row := d.db.DB().QueryRowContext(ctx, getUserQuery, username)
	user = &Users{}
	err = row.Scan(&user.ID, &user.UID, &user.Username, &user.Email, &user.HashedPassword, &user.UserType, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrUserNotFound
	}
	if err != nil {
		return
	}

	return
}

func (d *dbRepository) GetByUID(ctx context.Context, uid string) (user *Users, err error) {
	return
}

func (d *dbRepository) GetByID(ctx context.Context, id uint64) (user *Users, err error) {
	return
}
