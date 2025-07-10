package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"my-go-api/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {
	t.Run("GetAll", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		now := time.Now().String()
		expectedUsers := []models.User{
			{
				ID:         uuid.New(),
				Name:       "User 1",
				Email:      "user1@example.com",
				Username:   "user1",
				Password:   "hashed1",
				JwtVersion: "1",
				Provider:   "local",
				IsVerified: true,
				Role:       "user",
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				ID:         uuid.New(),
				Name:       "User 2",
				Email:      "user2@example.com",
				Username:   "user2",
				Password:   "hashed2",
				JwtVersion: "1",
				Provider:   "local",
				IsVerified: true,
				Role:       "admin",
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "username", "password", "jwt_version",
			"provider", "is_verified", "role", "created_at", "updated_at",
		})
		for _, user := range expectedUsers {
			rows.AddRow(
				user.ID, user.Name, user.Email, user.Username, user.Password, user.JwtVersion,
				user.Provider, user.IsVerified, user.Role, user.CreatedAt, user.UpdatedAt,
			)
		}

		mock.ExpectQuery(`SELECT (.+) FROM users`).WillReturnRows(rows)

		users, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("CreateOne", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		now := time.Now().String()
		expectedUser := &models.User{
			ID:         uuid.New(),
			Name:       "New User",
			Email:      "new@example.com",
			Username:   "newuser",
			Password:   "hashed",
			JwtVersion: "1",
			Provider:   "local",
			IsVerified: false,
			Role:       "user",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		params := CreateOneParams{
			Name:       "New User",
			Username:   "newuser",
			Email:      "new@example.com",
			Password:   "hashed",
			JWTVersion: "1",
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "username", "password", "jwt_version",
			"provider", "is_verified", "role", "created_at", "updated_at",
		}).AddRow(
			expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Username,
			expectedUser.Password, expectedUser.JwtVersion, expectedUser.Provider,
			expectedUser.IsVerified, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt,
		)

		mock.ExpectQuery(`INSERT INTO users`).
			WithArgs(params.Name, params.Username, params.Email, params.Password, params.JWTVersion).
			WillReturnRows(rows)

		user, err := repo.CreateOne(context.Background(), params)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetById", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		userID := uuid.New()
		now := time.Now().String()
		expectedUser := &models.User{
			ID:         userID,
			Name:       "Test User",
			Email:      "test@example.com",
			Username:   "testuser",
			Password:   "hashed",
			JwtVersion: "1",
			Provider:   "local",
			IsVerified: true,
			Role:       "user",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "username", "password", "jwt_version",
			"provider", "is_verified", "role", "created_at", "updated_at",
		}).AddRow(
			expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Username,
			expectedUser.Password, expectedUser.JwtVersion, expectedUser.Provider,
			expectedUser.IsVerified, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt,
		)

		mock.ExpectQuery(`SELECT (.+) FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnRows(rows)

		user, err := repo.GetById(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetByUsername", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		username := "testuser"
		now := time.Now().String()
		expectedUser := &models.User{
			ID:         uuid.New(),
			Name:       "Test User",
			Email:      "test@example.com",
			Username:   username,
			Password:   "hashed",
			JwtVersion: "1",
			Provider:   "local",
			IsVerified: true,
			Role:       "user",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "username", "password", "jwt_version",
			"provider", "is_verified", "role", "created_at", "updated_at",
		}).AddRow(
			expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Username,
			expectedUser.Password, expectedUser.JwtVersion, expectedUser.Provider,
			expectedUser.IsVerified, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt,
		)

		mock.ExpectQuery(`SELECT (.+) FROM users WHERE username = \$1`).
			WithArgs(username).
			WillReturnRows(rows)

		user, err := repo.GetByUsername(context.Background(), username)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetByEmail", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		email := "test@example.com"
		now := time.Now().String()
		expectedUser := &models.User{
			ID:         uuid.New(),
			Name:       "Test User",
			Email:      email,
			Username:   "testuser",
			Password:   "hashed",
			JwtVersion: "1",
			Provider:   "local",
			IsVerified: true,
			Role:       "user",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "username", "password", "jwt_version",
			"provider", "is_verified", "role", "created_at", "updated_at",
		}).AddRow(
			expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Username,
			expectedUser.Password, expectedUser.JwtVersion, expectedUser.Provider,
			expectedUser.IsVerified, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt,
		)

		mock.ExpectQuery(`SELECT (.+) FROM users WHERE email = \$1`).
			WithArgs(email).
			WillReturnRows(rows)

		user, err := repo.GetByEmail(context.Background(), email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("UpdateOne", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		now := time.Now()
		userToUpdate := &models.User{
			ID:         uuid.New(),
			Name:       "Updated User",
			Email:      "updated@example.com",
			Username:   "updateduser",
			Password:   "newhashed",
			JwtVersion: "2",
			Provider:   "local",
			IsVerified: true,
			Role:       "admin",
			CreatedAt:  now.Add(-24 * time.Hour).String(),
			UpdatedAt:  now.String(),
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "username", "password", "jwt_version",
			"provider", "is_verified", "role", "created_at", "updated_at",
		}).AddRow(
			userToUpdate.ID, userToUpdate.Name, userToUpdate.Email, userToUpdate.Username,
			userToUpdate.Password, userToUpdate.JwtVersion, userToUpdate.Provider,
			userToUpdate.IsVerified, userToUpdate.Role, userToUpdate.CreatedAt, userToUpdate.UpdatedAt,
		)

		mock.ExpectQuery(`UPDATE users`).
			WithArgs(
				userToUpdate.Username, userToUpdate.Email, userToUpdate.Name,
				userToUpdate.Password, userToUpdate.Role, userToUpdate.JwtVersion,
				userToUpdate.IsVerified, userToUpdate.ID,
			).
			WillReturnRows(rows)

		updatedUser, err := repo.UpdateOne(context.Background(), userToUpdate)
		assert.NoError(t, err)
		assert.Equal(t, userToUpdate, updatedUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOne with ID", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		userID := uuid.New()
		now := time.Now().String()
		expectedUser := &models.User{
			ID:         userID,
			Name:       "Test User",
			Email:      "test@example.com",
			Username:   "testuser",
			Password:   "hashed",
			JwtVersion: "1",
			Provider:   "local",
			IsVerified: true,
			Role:       "user",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "username", "password", "jwt_version",
			"provider", "is_verified", "role", "created_at", "updated_at",
		}).AddRow(
			expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Username,
			expectedUser.Password, expectedUser.JwtVersion, expectedUser.Provider,
			expectedUser.IsVerified, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt,
		)

		mock.ExpectQuery(`SELECT (.+) FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnRows(rows)

		params := GetOneParams{Id: &userID}
		user, err := repo.GetOne(context.Background(), params)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOne with Username", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		username := "testuser"
		now := time.Now().String()
		expectedUser := &models.User{
			ID:         uuid.New(),
			Name:       "Test User",
			Email:      "test@example.com",
			Username:   username,
			Password:   "hashed",
			JwtVersion: "1",
			Provider:   "local",
			IsVerified: true,
			Role:       "user",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "username", "password", "jwt_version",
			"provider", "is_verified", "role", "created_at", "updated_at",
		}).AddRow(
			expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Username,
			expectedUser.Password, expectedUser.JwtVersion, expectedUser.Provider,
			expectedUser.IsVerified, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt,
		)

		mock.ExpectQuery(`SELECT (.+) FROM users WHERE username = \$1`).
			WithArgs(username).
			WillReturnRows(rows)

		params := GetOneParams{Username: &username}
		user, err := repo.GetOne(context.Background(), params)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOne with Email", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		email := "test@example.com"
		now := time.Now().String()
		expectedUser := &models.User{
			ID:         uuid.New(),
			Name:       "Test User",
			Email:      email,
			Username:   "testuser",
			Password:   "hashed",
			JwtVersion: "1",
			Provider:   "local",
			IsVerified: true,
			Role:       "user",
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		rows := sqlmock.NewRows([]string{
			"id", "name", "email", "username", "password", "jwt_version",
			"provider", "is_verified", "role", "created_at", "updated_at",
		}).AddRow(
			expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Username,
			expectedUser.Password, expectedUser.JwtVersion, expectedUser.Provider,
			expectedUser.IsVerified, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt,
		)

		mock.ExpectQuery(`SELECT (.+) FROM users WHERE email = \$1`).
			WithArgs(email).
			WillReturnRows(rows)

		params := GetOneParams{Email: &email}
		user, err := repo.GetOne(context.Background(), params)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetOne with no params", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		params := GetOneParams{}
		user, err := repo.GetOne(context.Background(), params)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, "no valid query field provided")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetAll error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		mock.ExpectQuery(`SELECT (.+) FROM users`).WillReturnError(errors.New("query error"))

		users, err := repo.GetAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, users)
		assert.EqualError(t, err, "query error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("GetById not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewUserRepository(db)

		userID := uuid.New()

		mock.ExpectQuery(`SELECT (.+) FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetById(context.Background(), userID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
