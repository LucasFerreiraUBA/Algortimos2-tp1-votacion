package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	votos "rerepolez/votos"

	errores "rerepolez/errores"
)

const INGRESAR = "ingresar"
const DESHACER = "deshacer"
const VOTAR = "votar"
const FIN_VOTO = "fin-votar"
const VALIDACION = "OK"
const CANTIDAD_CARGOS = 3

func main() {
	var lista_partidos []votos.Partido
	args := os.Args
	if len(args) < 3 {
		ErrorParametro := new(errores.ErrorParametros)
		fmt.Fprintf(os.Stdout, "%s\n", ErrorParametro.Error())
		return
	} else {
		archivos := args[1:]
		s1 := abrirArchivo(archivos[0])
		s2 := abrirArchivo(archivos[1])
		votos_blancos := votos.CrearVotosEnBlanco()
		lista_partidos = append(lista_partidos, votos_blancos)
		for s1.Scan() {
			campos := strings.Split(s1.Text(), ",")
			var lista_candidatos [3]string
			for i := 0; i < CANTIDAD_CARGOS; i++ {
				lista_candidatos[i] = campos[i+1]
			}
			partido := votos.CrearPartido(campos[0], lista_candidatos)
			lista_partidos = append(lista_partidos, partido)
			fmt.Fprintf(os.Stdout, "%s\n", s1.Text())
		}
		for s2.Scan() {
			fmt.Fprintf(os.Stdout, "%s\n", s2.Text())
		}
	}
	fila := votos.CrearFilaVotacion()
	for {
		pedir_input(fila, lista_partidos)
	}
}

func ingresar(dni int, fila votos.Fila) {
	votante := votos.CrearVotante(dni)
	fila.Ingresar(votante)
	fmt.Println(VALIDACION)
}

func votar(tipo_voto votos.TipoVoto, alternativa int, fila votos.Fila) error {
	if !fila.HayVotantes() {
		return &errores.FilaVacia{}
	}
	votante := fila.VerActual()
	votante.Votar(tipo_voto, alternativa)
	return nil
}

func finalizar_voto(fila votos.Fila, lista_partidos []votos.Partido) {
	votante := fila.VerActual()
	votos_finales, _ := votante.FinVoto()
	var tipo_de_voto votos.TipoVoto
	for i := 0; i < len(votos_finales.VotoPorTipo); i++ {
		tipo_de_voto = votos.TipoVoto(i)
		lista_partidos[votos_finales.VotoPorTipo[i]].VotadoPara(tipo_de_voto)
	}
	fila.FinalizarVoto()
}

func pedir_input(fila votos.Fila, lista_partidos []votos.Partido) {
	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')
	palabras := strings.Split(input, " ")
	leer_input(palabras, fila, lista_partidos)
}

func leer_input(input []string, fila votos.Fila, lista_partidos []votos.Partido) error {
	switch input[0] {
	case INGRESAR:
		dni_string := strings.Trim(input[1], "\n")
		dni, _ := strconv.Atoi(dni_string)
		if dni <= 0 {
			return &errores.DNIError{}
		}
		ingresar(dni, fila)

	case VOTAR:
		var tipo_voto votos.TipoVoto
		tipo_voto_aux, _ := strconv.Atoi(input[1])
		tipo_voto = votos.TipoVoto(tipo_voto_aux)
		alternativa, _ := strconv.Atoi(input[2])
		votar(tipo_voto, alternativa, fila)

	case DESHACER:
		if fila.HayVotantes() {
			return &errores.FilaVacia{}
		}
		votante := fila.VerActual()
		votante.Deshacer()

	case FIN_VOTO:
		if fila.HayVotantes() {
			return &errores.FilaVacia{}
		}
		finalizar_voto(fila, lista_partidos)
	}
	return nil
}

func abrirArchivo(ruta string) *bufio.Scanner {
	archivo, err := os.Open(ruta)
	if err != nil {
		errorLeerArchivo := new(errores.ErrorLeerArchivo)
		fmt.Fprintf(os.Stderr, "%s", errorLeerArchivo.Error())
	}
	return bufio.NewScanner(archivo)
}
