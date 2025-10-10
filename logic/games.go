package logic

import (
	"errors"
	"strings"

	db "github.com/Tomasgithub01/ProgWeb/db/sqlc"
)

func ValidateGame(g db.Game) error {
	if strings.TrimSpace(g.Name) == "" {
		return errors.New("Name cannot be empty")
	}
	if len(g.Name) > 30 {
		return errors.New("Name cannot be longer than 30 characters")
	}

	if strings.TrimSpace(g.Description) == "" {
		return errors.New("Description cannot be empty")
	}
	if len(g.Description) > 200 {
		return errors.New("Description cannot be longer than 200 characters")
	}

	if g.Image.Valid && len(g.Image.String) > 100 {
		return errors.New("Image URL cannot be longer than 100 characters")
	}

	if g.Link.Valid && len(g.Link.String) > 100 {
		return errors.New("Link cannot be empty")
	}
	return nil
}
