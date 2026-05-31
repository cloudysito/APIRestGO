package repository

import "github.com/cloudysito/apirestgo/models"

type JugadorRepository interface {
	Guardar(jugador models.Jugador) error
	ObtenerTodos() ([]models.Jugador, error)
	ObtenerPorNombre(nombre string) (models.Jugador, error)
}

type MemoriaRepository struct {
	usuarios []models.Jugador
}

func NewMemoriaRepository() *MemoriaRepository {
	return &MemoriaRepository{
		usuarios: []models.Jugador{},
	}
}

func (r *MemoriaRepository) Guardar(jugador models.Jugador) error {
	r.usuarios = append(r.usuarios, jugador)
	return nil
}

func (r *MemoriaRepository) ObtenerTodos() ([]models.Jugador, error) {
	return r.usuarios, nil
}
