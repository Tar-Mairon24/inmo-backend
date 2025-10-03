package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/usecase"
	"inmo-backend/test/mocks/middleware"
	"inmo-backend/test/mocks/repository"
	"inmo-backend/test/mocks/service"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		logrus.Error("Could not load .env: ", err)
	}
	m.Run()
}
func TestAuthUseCase_Login(t *testing.T) {
    testUser := &models.User{
        ID:       1,
        Email:    "test@example.com",
        Password: "hashedpassword",
        Username: "testuser",
    }
    userResponse := models.UserResponse{
        ID:       testUser.ID,
        Email:    testUser.Email,
        Username: testUser.Username,
    }

    t.Run("empty email or password", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()
        
        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        resp, err := uc.Login("", "password")
        assert.Nil(t, resp)
        assert.EqualError(t, err, "email and password cannot be empty")

        resp, err = uc.Login("email@example.com", "")
        assert.Nil(t, resp)
        assert.EqualError(t, err, "email and password cannot be empty")
    })

    t.Run("user not found", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()
        
        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        userRepo.On("GetByEmail", "notfound@example.com").Return(nil, errors.New("not found"))
        
        resp, err := uc.Login("notfound@example.com", "password")
        assert.Nil(t, resp)
        assert.EqualError(t, err, "user not found")
        
        userRepo.AssertExpectations(t)
    })

    t.Run("password verification failed", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()
        
        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        userRepo.On("GetByEmail", testUser.Email).Return(testUser, nil)
        hashing.On("VerifyPassword", "hashedpassword", "wrongpassword").Return(errors.New("invalid password"))
        
        resp, err := uc.Login(testUser.Email, "wrongpassword")
        assert.Nil(t, resp)
        assert.EqualError(t, err, "invalid password")
        
        userRepo.AssertExpectations(t)
        hashing.AssertExpectations(t)
    })

    t.Run("error getting old refresh token - but login continues", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()
        
        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        userRepo.On("GetByEmail", testUser.Email).Return(testUser, nil)
        hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
        tokenRepo.On("GetTokenByUserID", testUser.ID).Return(nil, errors.New("db error"))
        
        // Login should continue despite the error
        authMiddleware.On("GenerateRefreshToken").Return("refreshToken", "idToken", nil)
        tokenRepo.On("SaveToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
        jwtService.On("GenerateToken", testUser).Return("jwtToken", nil)

        resp, err := uc.Login(testUser.Email, "password")
        assert.NoError(t, err)
        assert.NotNil(t, resp)
        assert.Equal(t, "jwtToken", resp.Token)
        assert.Equal(t, "refreshToken", resp.RefreshToken)

        userRepo.AssertExpectations(t)
        hashing.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
        authMiddleware.AssertExpectations(t)
        jwtService.AssertExpectations(t)
    })

    t.Run("delete old refresh token error", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()
        
        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        oldToken := &models.RefreshToken{ID: "oldtokenid"}
        userRepo.On("GetByEmail", testUser.Email).Return(testUser, nil)
        hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
        tokenRepo.On("GetTokenByUserID", testUser.ID).Return(oldToken, nil)
        tokenRepo.On("DeleteToken", oldToken.ID).Return(errors.New("delete error"))

        resp, err := uc.Login(testUser.Email, "password")
        assert.Nil(t, resp)
        assert.EqualError(t, err, "delete error")

        userRepo.AssertExpectations(t)
        hashing.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
    })

    t.Run("generate refresh token error", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()
        
        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        userRepo.On("GetByEmail", testUser.Email).Return(testUser, nil)
        hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
        tokenRepo.On("GetTokenByUserID", testUser.ID).Return(nil, nil)
        authMiddleware.On("GenerateRefreshToken").Return("", "", errors.New("refresh error"))

        resp, err := uc.Login(testUser.Email, "password")
        assert.Nil(t, resp)
        assert.EqualError(t, err, "refresh error")

        userRepo.AssertExpectations(t)
        hashing.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
        authMiddleware.AssertExpectations(t)
    })

    t.Run("save token error", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()
        
        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        userRepo.On("GetByEmail", testUser.Email).Return(testUser, nil)
        hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
        tokenRepo.On("GetTokenByUserID", testUser.ID).Return(nil, nil)
        authMiddleware.On("GenerateRefreshToken").Return("refreshToken", "idToken", nil)
        tokenRepo.On("SaveToken", mock.AnythingOfType("*models.RefreshToken")).Return(errors.New("save error"))

        resp, err := uc.Login(testUser.Email, "password")
        assert.Nil(t, resp)
        assert.EqualError(t, err, "save error")

        userRepo.AssertExpectations(t)
        hashing.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
        authMiddleware.AssertExpectations(t)
    })

    t.Run("generate jwt token error", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()
        
        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        userRepo.On("GetByEmail", testUser.Email).Return(testUser, nil)
        hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
        tokenRepo.On("GetTokenByUserID", testUser.ID).Return(nil, nil)
        authMiddleware.On("GenerateRefreshToken").Return("refreshToken", "idToken", nil)
        tokenRepo.On("SaveToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
        jwtService.On("GenerateToken", testUser).Return("", errors.New("jwt error"))

        resp, err := uc.Login(testUser.Email, "password")
        assert.Nil(t, resp)
        assert.EqualError(t, err, "jwt error")

        userRepo.AssertExpectations(t)
        hashing.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
        authMiddleware.AssertExpectations(t)
        jwtService.AssertExpectations(t)
    })

    t.Run("success", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()
        
        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        userRepo.On("GetByEmail", testUser.Email).Return(testUser, nil)
        hashing.On("VerifyPassword", "hashedpassword", "password").Return(nil)
        tokenRepo.On("GetTokenByUserID", testUser.ID).Return(nil, nil)
        authMiddleware.On("GenerateRefreshToken").Return("refreshToken", "idToken", nil)
        tokenRepo.On("SaveToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
        jwtService.On("GenerateToken", testUser).Return("jwtToken", nil)

        resp, err := uc.Login(testUser.Email, "password")
        assert.NoError(t, err)
        assert.NotNil(t, resp)
        assert.NotNil(t, resp.User, "User should not be nil")
        
        assert.Equal(t, userResponse.ID, resp.User.ID)
        assert.Equal(t, userResponse.Email, resp.User.Email)
        assert.Equal(t, userResponse.Username, resp.User.Username)
        
        assert.Equal(t, "jwtToken", resp.Token)
        assert.Equal(t, "refreshToken", resp.RefreshToken)

        userRepo.AssertExpectations(t)
        hashing.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
        authMiddleware.AssertExpectations(t)
        jwtService.AssertExpectations(t)
    })
}
func TestAuthUseCase_Logout(t *testing.T) {
    t.Run("userID is zero", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        err := uc.Logout(0)
        assert.EqualError(t, err, "user ID cannot be empty")
    })

    t.Run("error getting token ID by user ID", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        tokenRepo.On("GetTokenIDByUserID", uint(1)).Return("", errors.New("db error"))

        err := uc.Logout(1)
        assert.EqualError(t, err, "db error")
        tokenRepo.AssertExpectations(t)
    })

    t.Run("no token found for user", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        tokenRepo.On("GetTokenIDByUserID", uint(2)).Return("", nil)

        err := uc.Logout(2)
        assert.EqualError(t, err, "no token found for the given user ID, user was not logged in")
        tokenRepo.AssertExpectations(t)
    })

    t.Run("error deleting token", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        tokenRepo.On("GetTokenIDByUserID", uint(3)).Return("tokenid123", nil)
        tokenRepo.On("DeleteToken", "tokenid123").Return(errors.New("delete error"))

        err := uc.Logout(3)
        assert.EqualError(t, err, "delete error")
        tokenRepo.AssertExpectations(t)
    })

    t.Run("success", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        tokenRepo.On("GetTokenIDByUserID", uint(4)).Return("tokenid456", nil)
        tokenRepo.On("DeleteToken", "tokenid456").Return(nil)

        err := uc.Logout(4)
        assert.NoError(t, err)
        tokenRepo.AssertExpectations(t)
    })
}
func TestAuthUseCase_RefreshToken(t *testing.T) {
    testUserID := uint(1)
    validJwt := "valid.jwt.token"
    validRefresh := "valid-refresh-token"
    expiredRefresh := &models.RefreshToken{
        ID:        "tokenid",
        UserID:    testUserID,
        Token:     validRefresh,
        ExpiresAt: time.Now().Add(-time.Hour).Unix(),
        CreatedAt: time.Now().Add(-24 * time.Hour),
    }
    validRefreshToken := &models.RefreshToken{
        ID:        "tokenid",
        UserID:    testUserID,
        Token:     validRefresh,
        ExpiresAt: time.Now().Add(time.Hour).Unix(),
        CreatedAt: time.Now().Add(-24 * time.Hour),
    }

    t.Run("empty jwt or refresh token", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        resp, err := uc.RefreshToken(models.RefreshTokenData{JwtToken: "", RefreshToken: ""})
        assert.Nil(t, resp)
        assert.EqualError(t, err, "JWT token and refresh token cannot be empty")
    })

    t.Run("error getting user id from claims", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        jwtService.On("GetUserIDFromClaims", validJwt).Return(uint(0), errors.New("claims error"))

        resp, err := uc.RefreshToken(models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
        assert.Nil(t, resp)
        assert.EqualError(t, err, "claims error")
        jwtService.AssertExpectations(t)
    })

    t.Run("error getting refresh token by user id", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        jwtService.On("GetUserIDFromClaims", validJwt).Return(testUserID, nil)
        tokenRepo.On("GetTokenByUserID", testUserID).Return(nil, errors.New("db error"))

        resp, err := uc.RefreshToken(models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
        assert.Nil(t, resp)
        assert.EqualError(t, err, "db error")
        jwtService.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
    })

    t.Run("refresh token not found", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        jwtService.On("GetUserIDFromClaims", validJwt).Return(testUserID, nil)
        tokenRepo.On("GetTokenByUserID", testUserID).Return(nil, nil)

        resp, err := uc.RefreshToken(models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
        assert.Nil(t, resp)
        assert.EqualError(t, err, "refresh token not found")
        jwtService.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
    })

    t.Run("refresh token expired", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        jwtService.On("GetUserIDFromClaims", validJwt).Return(testUserID, nil)
        tokenRepo.On("GetTokenByUserID", testUserID).Return(expiredRefresh, nil)

        resp, err := uc.RefreshToken(models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
        assert.Nil(t, resp)
        assert.EqualError(t, err, "refresh token expired")
        jwtService.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
    })

    t.Run("invalid refresh token", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        jwtService.On("GetUserIDFromClaims", validJwt).Return(testUserID, nil)
        tokenRepo.On("GetTokenByUserID", testUserID).Return(&models.RefreshToken{
            ID:        "tokenid",
            UserID:    testUserID,
            Token:     "other-refresh-token",
            ExpiresAt: time.Now().Add(time.Hour).Unix(),
            CreatedAt: time.Now().Add(-24 * time.Hour),
        }, nil)

        resp, err := uc.RefreshToken(models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
        assert.Nil(t, resp)
        assert.EqualError(t, err, "invalid refresh token")
        jwtService.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
    })

    t.Run("error refreshing jwt token", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        jwtService.On("GetUserIDFromClaims", validJwt).Return(testUserID, nil)
        tokenRepo.On("GetTokenByUserID", testUserID).Return(validRefreshToken, nil)
        jwtService.On("RefreshToken", validJwt).Return("", errors.New("jwt refresh error"))

        resp, err := uc.RefreshToken(models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
        assert.Nil(t, resp)
        assert.EqualError(t, err, "jwt refresh error")
        jwtService.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
    })

    t.Run("success", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        jwtService.On("GetUserIDFromClaims", validJwt).Return(testUserID, nil)
        tokenRepo.On("GetTokenByUserID", testUserID).Return(validRefreshToken, nil)
        jwtService.On("RefreshToken", validJwt).Return("new.jwt.token", nil)

        resp, err := uc.RefreshToken(models.RefreshTokenData{JwtToken: validJwt, RefreshToken: validRefresh})
        assert.NoError(t, err)
        assert.NotNil(t, resp)
        assert.Equal(t, "new.jwt.token", resp.JwtToken)
        assert.Equal(t, validRefresh, resp.RefreshToken)
        jwtService.AssertExpectations(t)
        tokenRepo.AssertExpectations(t)
    })
}
func TestAuthUseCase_GetStatus(t *testing.T) {
    testUserID := uint(1)
    validRefresh := "valid-refresh-token"
    expiredRefresh := &models.RefreshToken{
        ID:        "tokenid",
        UserID:    testUserID,
        Token:     validRefresh,
        ExpiresAt: time.Now().Add(-time.Hour).Unix(),
        CreatedAt: time.Now().Add(-24 * time.Hour),
    }
    validRefreshToken := &models.RefreshToken{
        ID:        "tokenid",
        UserID:    testUserID,
        Token:     validRefresh,
        ExpiresAt: time.Now().Add(time.Hour).Unix(),
        CreatedAt: time.Now().Add(-24 * time.Hour),
    }

    t.Run("empty userID or refreshToken", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        err := uc.GetStatus(0, "")
        assert.EqualError(t, err, "user ID and refresh token cannot be empty")
        err = uc.GetStatus(testUserID, "")
        assert.EqualError(t, err, "user ID and refresh token cannot be empty")
        err = uc.GetStatus(0, validRefresh)
        assert.EqualError(t, err, "user ID and refresh token cannot be empty")
    })

    t.Run("error getting refresh token by user ID", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        tokenRepo.On("GetTokenByUserID", testUserID).Return(nil, errors.New("db error"))

        err := uc.GetStatus(testUserID, validRefresh)
        assert.EqualError(t, err, "db error")
        tokenRepo.AssertExpectations(t)
    })

    t.Run("refresh token not found", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        tokenRepo.On("GetTokenByUserID", testUserID).Return(nil, nil)

        err := uc.GetStatus(testUserID, validRefresh)
        assert.EqualError(t, err, "refresh token not found")
        tokenRepo.AssertExpectations(t)
    })

    t.Run("refresh token expired", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        tokenRepo.On("GetTokenByUserID", testUserID).Return(expiredRefresh, nil)

        err := uc.GetStatus(testUserID, validRefresh)
        assert.EqualError(t, err, "refresh token expired")
        tokenRepo.AssertExpectations(t)
    })

    t.Run("invalid refresh token", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        // Token mismatch
        tokenRepo.On("GetTokenByUserID", testUserID).Return(&models.RefreshToken{
            ID:        "tokenid",
            UserID:    testUserID,
            Token:     "other-refresh-token",
            ExpiresAt: time.Now().Add(time.Hour).Unix(),
            CreatedAt: time.Now().Add(-24 * time.Hour),
        }, nil)

        err := uc.GetStatus(testUserID, validRefresh)
        assert.EqualError(t, err, "invalid refresh token")
        tokenRepo.AssertExpectations(t)

        // Token empty
        tokenRepo = repositoryMock.NewMockTokenRepo()
        uc = usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)
        tokenRepo.On("GetTokenByUserID", testUserID).Return(&models.RefreshToken{
            ID:        "tokenid",
            UserID:    testUserID,
            Token:     "",
            ExpiresAt: time.Now().Add(time.Hour).Unix(),
            CreatedAt: time.Now().Add(-24 * time.Hour),
        }, nil)

        err = uc.GetStatus(testUserID, validRefresh)
        assert.EqualError(t, err, "invalid refresh token")
        tokenRepo.AssertExpectations(t)
    })

    t.Run("success", func(t *testing.T) {
        userRepo := repositoryMock.NewMockUserRepo()
        tokenRepo := repositoryMock.NewMockTokenRepo()
        jwtService := serviceMocks.NewMockJWTService()
        hashing := middlewareMock.NewMockHashing()
        authMiddleware := middlewareMock.NewMockMiddleware()

        uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService, hashing, authMiddleware)

        tokenRepo.On("GetTokenByUserID", testUserID).Return(validRefreshToken, nil)

        err := uc.GetStatus(testUserID, validRefresh)
        assert.NoError(t, err)
        tokenRepo.AssertExpectations(t)
    })
}


