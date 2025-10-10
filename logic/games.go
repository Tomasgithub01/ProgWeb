package logic

import (
	"errors"
	"strings"

	db "github.com/Tomasgithub01/ProgWeb/db/sqlc"
)

func ValidateGame(g db.Game) error {
	if strings.TrimSpace(g.Name) == "" {
		return errors.New("El nombre no puede estar vacío")
	}
	if len(g.Name) > 30 {
		return errors.New("El nombre no puede tener más de 30 caracteres")
	}

	if strings.TrimSpace(g.Description) == "" {
		return errors.New("La descripción no puede estar vacía")
	}
	if len(g.Description) > 200 {
		return errors.New("la descripción no puede tener más de 200 caracteres")
	}

	if g.Image.Valid && len(g.Image.String) > 100 {
		return errors.New("la URL de la imagen no puede tener más de 100 caracteres")
	}

}
