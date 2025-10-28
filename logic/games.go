package logic

import (
	"errors"
	"strings"

	db "ProgWeb/db/sqlc"
)

func ValidateGame(g db.Game) error {
	if strings.TrimSpace(g.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if len(g.Name) > 30 {
		return errors.New("name cannot be longer than 30 characters")
	}

	if strings.TrimSpace(g.Description) == "" {
		return errors.New("description cannot be empty")
	}
	if len(g.Description) > 200 {
		return errors.New("description cannot be longer than 200 characters")
	}

	if len(g.Image) > 200 {
		return errors.New("image URL cannot be longer than 200 characters")
	}

	if len(g.Link) > 200 {
		return errors.New("link cannot be longer than 200 characters")
	}
	return nil
}

func ValidateUser(u db.User) error {
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name cannot be empty")
	}
	if len(u.Name) > 30 {
		return errors.New("name cannot be longer than 30 characters")
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password cannot be empty")
	}
	if len(u.Password) > 255 {
		return errors.New("password cannot be longer than 255 characters")
	}
	return nil
}
