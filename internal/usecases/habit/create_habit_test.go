package habit

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uygardeniz/habit-tracker/internal/apperrors"
	"github.com/uygardeniz/habit-tracker/internal/dto"
	"github.com/uygardeniz/habit-tracker/internal/entity"
)

// MockHabitRepository is a mock implementation of the HabitRepository interface
type MockHabitRepository struct {
	mock.Mock
}

func (m *MockHabitRepository) Create(ctx context.Context, habit *entity.Habit) (*entity.Habit, error) {
	args := m.Called(ctx, habit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Habit), args.Error(1)
}

func (m *MockHabitRepository) FindByID(ctx context.Context, id string) (*entity.Habit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Habit), args.Error(1)
}

func (m *MockHabitRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Habit, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Habit), args.Error(1)
}

func (m *MockHabitRepository) Update(ctx context.Context, habit *entity.Habit) error {
	args := m.Called(ctx, habit)
	return args.Error(0)
}

func (m *MockHabitRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateHabitUsecase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		input          dto.CreateHabitDTO
		setupMock      func(*MockHabitRepository)
		expectedResult func(*entity.Habit) bool
		expectedError  error
	}{
		{
			name:   "successful habit creation with all fields",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "Morning Exercise",
				Description: stringPtr("Daily morning workout"),
				Motivation:  stringPtr("Get fit and healthy"),
				Color:       "#FF5733",
				Category:    stringPtr("Health"),
				Frequency:   "daily",
				TargetCount: 1,
				TargetDays:  nil,
			},
			setupMock: func(m *MockHabitRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(h *entity.Habit) bool {
					return h.Name == "Morning Exercise" &&
						h.UserID == "user-123" &&
						h.Frequency == "daily" &&
						h.TargetCount == 1 &&
						h.Color == "#FF5733" &&
						*h.Description == "Daily morning workout" &&
						*h.Motivation == "Get fit and healthy" &&
						*h.Category == "Health"
				})).Return(&entity.Habit{
					ID:          "habit-123",
					UserID:      "user-123",
					Name:        "Morning Exercise",
					Description: stringPtr("Daily morning workout"),
					Motivation:  stringPtr("Get fit and healthy"),
					Color:       "#FF5733",
					Category:    stringPtr("Health"),
					Frequency:   "daily",
					TargetCount: 1,
					IsActive:    true,
				}, nil)
			},
			expectedResult: func(h *entity.Habit) bool {
				return h != nil &&
					h.Name == "Morning Exercise" &&
					h.UserID == "user-123" &&
					h.Frequency == "daily" &&
					h.TargetCount == 1 &&
					h.Color == "#FF5733" &&
					h.IsActive == true
			},
			expectedError: nil,
		},
		{
			name:   "successful habit creation with weekly target days",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "Gym Workout",
				Color:       "#00FF00",
				Frequency:   "weekly",
				TargetCount: 3,
				TargetDays:  stringPtr(`{"days":["monday","wednesday","friday"]}`),
			},
			setupMock: func(m *MockHabitRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(h *entity.Habit) bool {
					return h.Name == "Gym Workout" &&
						h.TargetDays != nil &&
						len(h.TargetDays.Days) == 3
				})).Return(&entity.Habit{
					ID:          "habit-456",
					UserID:      "user-123",
					Name:        "Gym Workout",
					Color:       "#00FF00",
					Frequency:   "weekly",
					TargetCount: 3,
					TargetDays: &entity.TargetDays{
						Days: []any{"monday", "wednesday", "friday"},
					},
					IsActive: true,
				}, nil)
			},
			expectedResult: func(h *entity.Habit) bool {
				return h != nil &&
					h.Name == "Gym Workout" &&
					h.TargetDays != nil &&
					len(h.TargetDays.Days) == 3
			},
			expectedError: nil,
		},
		{
			name:   "successful habit creation with monthly target days",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "Monthly Review",
				Color:       "#0000FF",
				Frequency:   "monthly",
				TargetCount: 1,
				TargetDays:  stringPtr(`{"days":[1,15]}`),
			},
			setupMock: func(m *MockHabitRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(h *entity.Habit) bool {
					return h.Name == "Monthly Review" &&
						h.TargetDays != nil &&
						len(h.TargetDays.Days) == 2
				})).Return(&entity.Habit{
					ID:          "habit-789",
					UserID:      "user-123",
					Name:        "Monthly Review",
					Color:       "#0000FF",
					Frequency:   "monthly",
					TargetCount: 1,
					TargetDays: &entity.TargetDays{
						Days: []any{float64(1), float64(15)},
					},
					IsActive: true,
				}, nil)
			},
			expectedResult: func(h *entity.Habit) bool {
				return h != nil &&
					h.Name == "Monthly Review" &&
					h.TargetDays != nil &&
					len(h.TargetDays.Days) == 2
			},
			expectedError: nil,
		},
		{
			name:   "invalid target days JSON format",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "Invalid Habit",
				Color:       "#FF0000",
				Frequency:   "weekly",
				TargetCount: 1,
				TargetDays:  stringPtr(`invalid json`),
			},
			setupMock: func(m *MockHabitRepository) {
				// Repository should not be called due to validation error
			},
			expectedResult: func(h *entity.Habit) bool {
				return h == nil
			},
			expectedError: apperrors.ErrInvalidInput,
		},
		{
			name:   "invalid habit validation - empty name",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "",
				Color:       "#FF0000",
				Frequency:   "daily",
				TargetCount: 1,
			},
			setupMock: func(m *MockHabitRepository) {
				// Repository should not be called due to validation error
			},
			expectedResult: func(h *entity.Habit) bool {
				return h == nil
			},
			expectedError: apperrors.ErrInvalidInput,
		},
		{
			name:   "invalid habit validation - invalid frequency",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "Test Habit",
				Color:       "#FF0000",
				Frequency:   "invalid",
				TargetCount: 1,
			},
			setupMock: func(m *MockHabitRepository) {
				// Repository should not be called due to validation error
			},
			expectedResult: func(h *entity.Habit) bool {
				return h == nil
			},
			expectedError: apperrors.ErrInvalidInput,
		},
		{
			name:   "invalid habit validation - invalid color",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "Test Habit",
				Color:       "invalid-color",
				Frequency:   "daily",
				TargetCount: 1,
			},
			setupMock: func(m *MockHabitRepository) {
				// Repository should not be called due to validation error
			},
			expectedResult: func(h *entity.Habit) bool {
				return h == nil
			},
			expectedError: apperrors.ErrInvalidInput,
		},
		{
			name:   "invalid habit validation - zero target count",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "Test Habit",
				Color:       "#FF0000",
				Frequency:   "daily",
				TargetCount: 0,
			},
			setupMock: func(m *MockHabitRepository) {
				// Repository should not be called due to validation error
			},
			expectedResult: func(h *entity.Habit) bool {
				return h == nil
			},
			expectedError: apperrors.ErrInvalidInput,
		},
		{
			name:   "invalid weekly target days",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "Test Habit",
				Color:       "#FF0000",
				Frequency:   "weekly",
				TargetCount: 1,
				TargetDays:  stringPtr(`{"days":["invalidday"]}`),
			},
			setupMock: func(m *MockHabitRepository) {
				// Repository should not be called due to validation error
			},
			expectedResult: func(h *entity.Habit) bool {
				return h == nil
			},
			expectedError: apperrors.ErrInvalidInput,
		},
		{
			name:   "invalid monthly target days",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "Test Habit",
				Color:       "#FF0000",
				Frequency:   "monthly",
				TargetCount: 1,
				TargetDays:  stringPtr(`{"days":[0,32]}`),
			},
			setupMock: func(m *MockHabitRepository) {
				// Repository should not be called due to validation error
			},
			expectedResult: func(h *entity.Habit) bool {
				return h == nil
			},
			expectedError: apperrors.ErrInvalidInput,
		},
		{
			name:   "repository error",
			userID: "user-123",
			input: dto.CreateHabitDTO{
				Name:        "Test Habit",
				Color:       "#FF0000",
				Frequency:   "daily",
				TargetCount: 1,
			},
			setupMock: func(m *MockHabitRepository) {
				m.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedResult: func(h *entity.Habit) bool {
				return h == nil
			},
			expectedError: errors.New("database error"),
		},
		{
			name:   "empty user ID",
			userID: "",
			input: dto.CreateHabitDTO{
				Name:        "Test Habit",
				Color:       "#FF0000",
				Frequency:   "daily",
				TargetCount: 1,
			},
			setupMock: func(m *MockHabitRepository) {
				// Repository should not be called due to validation error
			},
			expectedResult: func(h *entity.Habit) bool {
				return h == nil
			},
			expectedError: apperrors.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockHabitRepository)
			tt.setupMock(mockRepo)

			usecase := NewCreateHabitUsecase(mockRepo)
			ctx := context.Background()

			// Execute
			result, err := usecase.Execute(ctx, tt.userID, tt.input)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.True(t, tt.expectedResult(result))
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateHabitUsecase_Constructor(t *testing.T) {
	mockRepo := new(MockHabitRepository)
	usecase := NewCreateHabitUsecase(mockRepo)

	assert.NotNil(t, usecase)
	assert.Equal(t, mockRepo, usecase.habitRepository)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
