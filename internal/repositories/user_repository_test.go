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

func TestUserRepository_GetOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	tests := []struct {
		name    string
		params  GetOneParams
		mock    func()
		want    *models.User
		wantErr bool
	}{
		{
			name: "Get by ID",
			params: GetOneParams{
				Id: func() *uuid.UUID { id := uuid.New(); return &id }(),
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "email", "password", "provider", "role", "created_at", "updated_at"}).
					AddRow(uuid.New(), "John Doe", "johndoe", "john@example.com", "password", "local", "user", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, username, email, password, provider, role, created_at, updated_at FROM users WHERE id = ?").
					WithArgs(sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			want:    &models.User{},
			wantErr: false,
		},
		{
			name: "Get by username",
			params: GetOneParams{
				Username: func() *string { s := "johndoe"; return &s }(),
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "email", "password", "provider", "role", "created_at", "updated_at"}).
					AddRow(uuid.New(), "John Doe", "johndoe", "john@example.com", "password", "local", "user", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, username, email, password, provider, role, created_at, updated_at FROM users WHERE username = ?").
					WithArgs("johndoe").
					WillReturnRows(rows)
			},
			want:    &models.User{},
			wantErr: false,
		},
		{
			name: "Get by email",
			params: GetOneParams{
				Email: func() *string { s := "john@example.com"; return &s }(),
			},
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "username", "email", "password", "provider", "role", "created_at", "updated_at"}).
					AddRow(uuid.New(), "John Doe", "johndoe", "john@example.com", "password", "local", "user", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, username, email, password, provider, role, created_at, updated_at FROM users WHERE email = ?").
					WithArgs("john@example.com").
					WillReturnRows(rows)
			},
			want:    &models.User{},
			wantErr: false,
		},
		{
			name:   "No valid query field",
			params: GetOneParams{
				// No fields set
			},
			mock:    func() {},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Not found",
			params: GetOneParams{
				Username: func() *string { s := "nonexistent"; return &s }(),
			},
			mock: func() {
				mock.ExpectQuery("SELECT id, name, username, email, password, provider, role, created_at, updated_at FROM users WHERE username = ?").
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.GetOne(context.Background(), tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepository.GetOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got == nil {
				t.Errorf("userRepository.GetOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		expectedUsers := []models.User{
			{
				ID:        uuid.New(),
				Name:      "John Doe",
				Username:  "johndoe",
				Email:     "john@example.com",
				Provider:  "local",
				Role:      "user",
				CreatedAt: now.String(),
				UpdatedAt: now.String(),
			},
			{
				ID:        uuid.New(),
				Name:      "Jane Doe",
				Username:  "janedoe",
				Email:     "jane@example.com",
				Provider:  "local",
				Role:      "admin",
				CreatedAt: now.String(),
				UpdatedAt: now.String(),
			},
		}

		rows := sqlmock.NewRows([]string{"id", "name", "username", "email", "provider", "role", "created_at", "updated_at"})
		for _, user := range expectedUsers {
			rows.AddRow(user.ID, user.Name, user.Username, user.Email, user.Provider, user.Role, user.CreatedAt, user.UpdatedAt)
		}

		mock.ExpectQuery("SELECT id, name, username, email, provider, role, created_at, updated_at FROM users").
			WillReturnRows(rows)

		users, err := repo.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, len(expectedUsers), len(users))
		assert.Equal(t, expectedUsers[0].Username, users[0].Username)
		assert.Equal(t, expectedUsers[1].Username, users[1].Username)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, username, email, provider, role, created_at, updated_at FROM users").
			WillReturnError(errors.New("some error"))

		users, err := repo.GetAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, users)
	})
}

func TestUserRepository_CreateOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("Success", func(t *testing.T) {
		now := time.Now()
		params := CreateOneParams{
			Name:       "John Doe",
			Username:   "johndoe",
			Email:      "john@example.com",
			Password:   "hashedpassword",
			JWTVersion: "1",
		}

		expectedUser := &models.User{
			ID:        uuid.New(),
			Name:      params.Name,
			Username:  params.Username,
			Email:     params.Email,
			Password:  params.Password,
			Provider:  "local",
			Role:      "user",
			CreatedAt: now.String(),
			UpdatedAt: now.String(),
		}

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(params.Name, params.Username, params.Email, params.Password, params.JWTVersion).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "username", "provider", "role", "updated_at", "created_at"}).
				AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Email, expectedUser.Password, expectedUser.Username, expectedUser.Provider, expectedUser.Role, expectedUser.UpdatedAt, expectedUser.CreatedAt))

		user, err := repo.CreateOne(context.Background(), params)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Username, user.Username)
		assert.Equal(t, expectedUser.Email, user.Email)
	})

	t.Run("Error", func(t *testing.T) {
		params := CreateOneParams{
			Name:       "John Doe",
			Username:   "johndoe",
			Email:      "john@example.com",
			Password:   "hashedpassword",
			JWTVersion: "1",
		}

		mock.ExpectQuery("INSERT INTO users").
			WithArgs(params.Name, params.Username, params.Email, params.Password, params.JWTVersion).
			WillReturnError(errors.New("some error"))

		user, err := repo.CreateOne(context.Background(), params)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("Success", func(t *testing.T) {
		userId := uuid.New()
		expectedUser := &models.User{
			ID:        userId,
			Name:      "John Doe",
			Username:  "johndoe",
			Email:     "john@example.com",
			Password:  "hashedpassword",
			Provider:  "local",
			Role:      "user",
			CreatedAt: time.Now().String(),
			UpdatedAt: time.Now().String(),
		}

		mock.ExpectQuery("SELECT id, name, username, email, password, provider, role, created_at, updated_at FROM users WHERE id = ?").
			WithArgs(userId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "email", "password", "provider", "role", "created_at", "updated_at"}).
				AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Username, expectedUser.Email, expectedUser.Password, expectedUser.Provider, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt))

		user, err := repo.GetById(context.Background(), userId)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Username, user.Username)
	})

	t.Run("Not found", func(t *testing.T) {
		userId := uuid.New()
		mock.ExpectQuery("SELECT id, name, username, email, password, provider, role, created_at, updated_at FROM users WHERE id = ?").
			WithArgs(userId).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetById(context.Background(), userId)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("Success", func(t *testing.T) {
		username := "johndoe"
		expectedUser := &models.User{
			ID:        uuid.New(),
			Name:      "John Doe",
			Username:  username,
			Email:     "john@example.com",
			Password:  "hashedpassword",
			Provider:  "local",
			Role:      "user",
			CreatedAt: time.Now().String(),
			UpdatedAt: time.Now().String(),
		}

		mock.ExpectQuery("SELECT id, name, username, email, password, provider, role, created_at, updated_at FROM users WHERE username = ?").
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "email", "password", "provider", "role", "created_at", "updated_at"}).
				AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Username, expectedUser.Email, expectedUser.Password, expectedUser.Provider, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt))

		user, err := repo.GetByUsername(context.Background(), username)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Username, user.Username)
	})

	t.Run("Not found", func(t *testing.T) {
		username := "nonexistent"
		mock.ExpectQuery("SELECT id, name, username, email, password, provider, role, created_at, updated_at FROM users WHERE username = ?").
			WithArgs(username).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByUsername(context.Background(), username)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("Success", func(t *testing.T) {
		email := "john@example.com"
		expectedUser := &models.User{
			ID:        uuid.New(),
			Name:      "John Doe",
			Username:  "johndoe",
			Email:     email,
			Password:  "hashedpassword",
			Provider:  "local",
			Role:      "user",
			CreatedAt: time.Now().String(),
			UpdatedAt: time.Now().String(),
		}

		mock.ExpectQuery("SELECT id, name, username, email, password, provider, role, created_at, updated_at FROM users WHERE email = ?").
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "email", "password", "provider", "role", "created_at", "updated_at"}).
				AddRow(expectedUser.ID, expectedUser.Name, expectedUser.Username, expectedUser.Email, expectedUser.Password, expectedUser.Provider, expectedUser.Role, expectedUser.CreatedAt, expectedUser.UpdatedAt))

		user, err := repo.GetByEmail(context.Background(), email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Email, user.Email)
	})

	t.Run("Not found", func(t *testing.T) {
		email := "nonexistent@example.com"
		mock.ExpectQuery("SELECT id, name, username, email, password, provider, role, created_at, updated_at FROM users WHERE email = ?").
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetByEmail(context.Background(), email)
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestUserRepository_UpdateOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{
			ID:        uuid.New(),
			Name:      "John Doe Updated",
			Username:  "johndoe",
			Email:     "john.updated@example.com",
			Password:  "newhashedpassword",
			Role:      "admin",
			Provider:  "local",
			CreatedAt: time.Now().String(),
		}

		mock.ExpectQuery("UPDATE users").
			WithArgs(user.Username, user.Email, user.Name, user.Password, user.Role, user.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "email", "password", "provider", "role", "created_at", "updated_at"}).
				AddRow(user.ID, user.Name, user.Username, user.Email, user.Password, user.Provider, user.Role, user.CreatedAt, time.Now()))

		updatedUser, err := repo.UpdateOne(context.Background(), user)
		assert.NoError(t, err)
		assert.Equal(t, user.Name, updatedUser.Name)
		assert.Equal(t, user.Email, updatedUser.Email)
	})

	t.Run("Error", func(t *testing.T) {
		user := &models.User{
			ID:        uuid.New(),
			Name:      "John Doe Updated",
			Username:  "johndoe",
			Email:     "john.updated@example.com",
			Password:  "newhashedpassword",
			Role:      "admin",
			Provider:  "local",
			CreatedAt: time.Now().String(),
		}

		mock.ExpectQuery("UPDATE users").
			WithArgs(user.Username, user.Email, user.Name, user.Password, user.Role, user.ID).
			WillReturnError(errors.New("some error"))

		updatedUser, err := repo.UpdateOne(context.Background(), user)
		assert.Error(t, err)
		assert.Nil(t, updatedUser)
	})
}
