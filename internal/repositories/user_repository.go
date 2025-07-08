package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"my-go-api/internal/models"

	"github.com/google/uuid"
)

type GetOneParams struct {
	Username *string
	Email    *string
	Id       *uuid.UUID
}

type CreateOneParams struct {
	Name       string
	Username   string
	Email      string
	Password   string
	JWTVersion string
}

type IUserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	CreateOne(ctx context.Context, params CreateOneParams) (*models.User, error)
	GetById(ctx context.Context, userId uuid.UUID) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateOne(ctx context.Context, user *models.User) (*models.User, error)
	GetOne(ctx context.Context, params GetOneParams) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) IUserRepository {
	return &userRepository{db: db}
}

func (s *userRepository) GetOne(ctx context.Context, params GetOneParams) (*models.User, error) {
	user := &models.User{}

	var whereClause string
	var value any

	switch {
	case params.Id != nil:
		whereClause = "id = $1"
		value = params.Id
	case params.Username != nil:
		whereClause = "username = $1"
		value = *params.Username
	case params.Email != nil:
		whereClause = "email = $1"
		value = *params.Email
	default:
		return nil, errors.New("no valid query field provided")
	}

	sqlQuery := fmt.Sprintf(`
	SELECT %s 
	FROM users 
	WHERE %s`,
		userSelectedFields, whereClause,
	)

	if err := s.db.QueryRowContext(ctx, sqlQuery, value).Scan(scanUser(user)...); err != nil {
		return nil, err
	}

	return user, nil

}

func (s *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := fmt.Sprintf(`SELECT %s FROM users`, userSelectedFields)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(scanUser(&user)...)
		if err != nil {
			log.Printf("Failed to scan user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *userRepository) CreateOne(
	ctx context.Context,
	params CreateOneParams,
) (*models.User, error) {
	user := &models.User{}

	query := fmt.Sprintf(`
		INSERT INTO users (
			name, 
			username, 
			email, 
			password, 
			jwt_version
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING %s
	`, userSelectedFields)

	if err := s.db.QueryRowContext(ctx, query,
		params.Name,
		params.Username,
		params.Email,
		params.Password,
		params.JWTVersion,
	).
		Scan(scanUser(user)...); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userRepository) GetById(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	user := &models.User{}

	query := fmt.Sprintf(`SELECT %s FROM users WHERE id = $1`, userSelectedFields)

	if err := s.db.QueryRowContext(ctx, query, userId).
		Scan(scanUser(user)...); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*models.User, error) {
	user := &models.User{}

	query := fmt.Sprintf(`
	SELECT %s 
	FROM users 
	WHERE username = $1`,
		userSelectedFields,
	)

	if err := s.db.QueryRowContext(ctx, query, username).
		Scan(scanUser(user)...); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := fmt.Sprintf(`
		SELECT %s
		FROM users 
		WHERE email = $1`,
		userSelectedFields,
	)

	if err := s.db.QueryRowContext(ctx, query, email).
		Scan(scanUser(user)...); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userRepository) UpdateOne(ctx context.Context, user *models.User) (*models.User, error) {
	log.Println(user)
	query := fmt.Sprintf(`
		UPDATE users
		SET username=$1, 
				email=$2, 
				name=$3, 
				password=$4, 
				role=$5, 
				jwt_version=$6, 
				is_verified=$7, 
				updated_at=NOW()
		WHERE id=$8 
		RETURNING %s`,
		userSelectedFields,
	)

	if err := s.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Name,
		user.Password,
		user.Role,
		user.JwtVersion,
		user.IsVerified,
		user.ID,
	).
		Scan(scanUser(user)...); err != nil {
		return nil, err
	}

	return user, nil
}

func scanUser(user *models.User) []any {
	return []any{
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.JwtVersion,
		&user.Provider,
		&user.IsVerified,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	}
}

const userSelectedFields = `
		id, 
		name, 
		email, 
		username, 
		password, 
		jwt_version,
		provider, 
		is_verified,
		role,
		created_at, 
		updated_at 
`
