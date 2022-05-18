package auth_test

import (
	"testing"
	"tupulung/entities"
	web "tupulung/entities/web"
	_userRepository "tupulung/repositories/user"
	_authService "tupulung/services/auth"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	t.Run("invalid-email", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindBy").Return(entities.User{}, web.WebError{})

		authService := _authService.NewAuthService(userRepositoryMock)
		_, err := authService.Login(entities.AuthRequest{
			Email:    userSample.Email + "wrongwrongwrong",
			Password: "password",
		})

		assert.Error(t, err)
	})
	t.Run("invalid-password", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindBy").Return(userSample, nil)

		authService := _authService.NewAuthService(userRepositoryMock)
		_, err := authService.Login(entities.AuthRequest{
			Email:    userSample.Email,
			Password: "invalidpasswordhere",
		})

		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("FindBy").Return(userSample, nil)

		authService := _authService.NewAuthService(userRepositoryMock)
		actual, err := authService.Login(entities.AuthRequest{
			Email:    userSample.Email,
			Password: "password",
		})

		assert.Nil(t, err)
		assert.NotEqual(t, "", actual.Token)
	})
}

func TestMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSample := _userRepository.UserCollection[0]
		userSample.Password = "$2a$12$iu2L7bKpW4Rpe5yPGt3KPOm5N229fSuMlkHYu5l25dIwgvvW6oQYO" // pass: password
		userRepositoryMock := _userRepository.NewUserRepositoryMock(&mock.Mock{})
		userRepositoryMock.Mock.On("Find").Return(userSample, nil)

		jwt := jwt.Token{
			Raw:    "",
			Method: jwt.SigningMethodHS256,
			Claims: jwt.MapClaims{},
		}
		authService := _authService.NewAuthService(userRepositoryMock)
		_, err := authService.Me(int(userSample.ID), &jwt)

		assert.Nil(t, err)
	})
}
