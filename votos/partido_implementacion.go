package votos

import (
	"fmt"
	"os"
	"rerepolez/errores"
)

type partidoImplementacion struct {
	nombre     string
	candidatos [CANT_VOTACION]string
	votos      [CANT_VOTACION]int
}

type partidoEnBlanco struct {
	votos [CANT_VOTACION]int
}

func CrearPartido(nombre string, candidatos [CANT_VOTACION]string) Partido {
	partido := partidoImplementacion{nombre: nombre, candidatos: candidatos}
	return &partido
}

func CrearVotosEnBlanco() Partido {
	return new(partidoEnBlanco)
}

func (partido *partidoImplementacion) VotadoPara(tipo TipoVoto) {
	if tipo < 0 || tipo >= CANT_VOTACION {
		err := new(errores.ErrorTipoVoto)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	partido.votos[tipo]++

}

func (partido partidoImplementacion) ObtenerResultado(tipo TipoVoto) string {
	return fmt.Sprintf("%d", partido.votos[tipo])
}

func (blanco *partidoEnBlanco) VotadoPara(tipo TipoVoto) {
	if tipo < 0 || tipo >= CANT_VOTACION {
		err := new(errores.ErrorTipoVoto)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	blanco.votos[tipo]++
}

func (blanco partidoEnBlanco) ObtenerResultado(tipo TipoVoto) string {
	return fmt.Sprintf("%d", blanco.votos[tipo])
}
