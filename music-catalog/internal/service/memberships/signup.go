package memberships

import (
	"database/sql"
	"errors"
	"music-catalog/internal/models/memberships"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) SignUp(request memberships.SignUpRequest) error {
	existingUser, err := s.repository.GetUser(request.Email, request.Username, 0)
	if err != nil || err != sql.ErrNoRows {
		log.Error().Err(err).Msg("failed to get user")
		return err
	}

	if existingUser != nil {
		log.Error().Msg("user already exists")
		return errors.New("Email or username already exists")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return err
	}

	model := memberships.User{
		Email:     request.Email,
		Username:  request.Username,
		Password:  string(pass),
		CreatedBy: request.Email,
		UpdateBy:  request.Email,
	}
	return s.repository.CreateUser(model)
}
