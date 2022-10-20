package votos

import (
	"fmt"
	"os"
	errores"rerepolez/errores"
)

const NOMBRE_PARTIDO_BLANCO = "Votos En Blanco"

type partidoImplementacion struct {
	nombre     string
	candidatos [CANT_VOTACION]string
	votos      [CANT_VOTACION]int
}

type partidoEnBlanco struct {
	nombre string
	votos [CANT_VOTACION]int
}

func CrearPartido(nombre string, candidatos [CANT_VOTACION]string) Partido {
	partido := partidoImplementacion{nombre: nombre, candidatos: candidatos}
	return &partido
}

func CrearVotosEnBlanco() Partido {
	votos_blancos := new(partidoEnBlanco)
	votos_blancos.nombre = NOMBRE_PARTIDO_BLANCO
	return votos_blancos
}

func (partido *partidoImplementacion) VotadoPara(tipo TipoVoto) {
	if tipo < 0 || tipo >= CANT_VOTACION {
		err := new(errores.ErrorTipoVoto)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	partido.votos[tipo]++

}

func (partido partidoImplementacion) Nombre() string {
	return partido.nombre
}

func (blanco partidoEnBlanco) Nombre() string {
	return blanco.nombre
}

func (partido partidoImplementacion) Candidato(cargo int) string {
	return partido.candidatos[cargo]
}

func (blanco partidoEnBlanco) Candidato(cargo int) string {
	return ""
}

func (partido partidoImplementacion) ObtenerResultado(tipo TipoVoto) string {
	if partido.votos[tipo] == 1{
		return fmt.Sprintf("%d voto", partido.votos[tipo])
	}
	return fmt.Sprintf("%d votos", partido.votos[tipo])
}

func (blanco *partidoEnBlanco) VotadoPara(tipo TipoVoto) {
	if tipo < 0 || tipo >= CANT_VOTACION {
		err := new(errores.ErrorTipoVoto)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	blanco.votos[tipo]++
}

func (blanco partidoEnBlanco) ObtenerResultado(tipo TipoVoto) string {
	if blanco.votos[tipo] == 1{
		return fmt.Sprintf("%d voto", blanco.votos[tipo])
	}
	return fmt.Sprintf("%d votos", blanco.votos[tipo])
}