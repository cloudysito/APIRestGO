package models

type Jugador struct {
	Nombre string `json:"nombre"`
	Rango  string `json:"rango"`
}

var BaseDatos []Jugador
