package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Service define las operaciones del servicio de contraseñas
type Service interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

// service implementa el servicio de contraseñas
type service struct {
	cost int
}

// NewService crea una nueva instancia del servicio de contraseñas
func NewService(cost int) Service {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &service{cost: cost}
}

// Hash genera un hash de la contraseña
func (s *service) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify verifica si la contraseña coincide con el hash
func (s *service) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
