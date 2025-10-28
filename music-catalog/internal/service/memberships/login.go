package memberships

import (
	"errors"
	"music-catalog/internal/models/memberships"
	"music-catalog/pkg/jwt"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *service) Login(request memberships.LoginRequest) (string, error) {
	userDetail, err := s.repository.GetUser(request.Email, "", 0)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Msg("error get user from database")
		return "", err
	}

	if userDetail == nil {
		return "", errors.New("email not exists")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(request.Password))
	if err != nil {
		return "", errors.New("email and password do not match")
	}

	accessToken, err := jwt.CreateToken(int64(userDetail.ID), userDetail.Username, s.cfg.Service.SecretJWT)
	if err != nil {
		log.Error().Err(err).Msg("error create jwt token")
		return "", err
	}

	return accessToken, nil
}
