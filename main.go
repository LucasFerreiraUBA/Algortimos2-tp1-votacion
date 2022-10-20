package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	votos "rerepolez/votos"

	errores "rerepolez/errores"
)

const (
	INGRESAR        = "ingresar"
	DESHACER        = "deshacer"
	VOTAR           = "votar"
	FIN_VOTO        = "fin-votar"
	VALIDACION      = "OK"
	CANTIDAD_CARGOS = 3
	PRESIDENTE_STR  = "Presidente"
	INTENDENTE_STR  = "Intendente"
	GOBERNADOR_STR  = "Gobernador"
	PRESIDENTE_INT  = 0
	GOBERNADOR_INT  = 1
	INTENDENTE_INT  = 2
)

func main() {
	var lista_partidos []votos.Partido
	padron := make([]int, 0, 10)
	args := os.Args
	if len(args) < 3 {
		// ERROR: FALTAN PARÁMETROS
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
		}
		for s2.Scan() {
			elemento, _ := strconv.Atoi(s2.Text())
			insertar_elemento(&padron, elemento)
		}
	}
	fila := votos.CrearFilaVotacion()
	for {
		pedir_input(&fila, lista_partidos, &padron)
	}
}

func atrapar_errores_deshacer(fila votos.Fila, dni int) error {
	if fila.ValidarDNI(dni) {
		return &errores.ErrorVotanteFraudulento{dni}
	}
	return nil
}

func insertar_elemento(lista *[]int, elemento int) {
	// utiliza el algoritmo de inserción para insertar ordenadamente los elementos en la lista
	if len(*lista) == 0 {
		*lista = append(*lista, elemento)
		return
	}
	if len(*lista) == 1 {
		if (*lista)[0] < elemento {
			*lista = append((*lista)[0:], elemento)
		} else {
			*lista = append((*lista)[:1], (*lista)[0])
			(*lista)[0] = elemento
		}
		return
	}
	for i := 0; i < len(*lista); i++ {
		if elemento > (*lista)[i] {
			continue
		}
		*lista = append((*lista)[:i+1], (*lista)[i:]...)
		(*lista)[i] = elemento
		return
	}
	*lista = append(*lista, elemento)
}

func ingresar(dni int, fila *votos.Fila) {
	votante := votos.CrearVotante(dni)
	(*fila).Ingresar(votante)
	fmt.Println(VALIDACION)
}

func votar(tipo_voto votos.TipoVoto, alternativa int, fila *votos.Fila) {
	if !(*fila).HayVotantes() {
		ErrorFila := new(errores.FilaVacia)
		fmt.Fprintf(os.Stdout, "%s\n", ErrorFila.Error())
		return
	}
	votante := (*fila).VerActual()
	(*votante).Votar(tipo_voto, alternativa)
	fmt.Println(VALIDACION)
}

func finalizar_voto(fila votos.Fila, lista_partidos []votos.Partido) {
	votante := fila.VerActual()
	votos_finales, _ := (*votante).FinVoto()
	var tipo_de_voto votos.TipoVoto
	for i := 0; i < len(votos_finales.VotoPorTipo); i++ {
		tipo_de_voto = votos.TipoVoto(i)
		lista_partidos[votos_finales.VotoPorTipo[i]].VotadoPara(tipo_de_voto)
	}
	fila.FinalizarVoto()
	fmt.Println(VALIDACION)
}

func imprimir_resultados(lista_partidos []votos.Partido) {
	fmt.Fprintf(os.Stdout, "%s:\n", PRESIDENTE_STR)
	for i := 0; i < len(lista_partidos); i++ {
		fmt.Fprintf(os.Stdout, "%s", lista_partidos[i].ObtenerResultado(0))
	}
}

func pedir_input(fila *votos.Fila, lista_partidos []votos.Partido, padron *[]int) {
	inputReader := bufio.NewReader(os.Stdin) // El usuario ingresa el comando
	input, error := inputReader.ReadString('\n')
	if error == nil {
		palabras := strings.Split(input, " ")
		leer_input(palabras, fila, lista_partidos, padron)
	}
	if error == io.EOF {
		imprimir_resultados(lista_partidos)
	}
}

func atrapar_errores_dni(dni int, fila *votos.Fila, padron *[]int) bool {
	if dni <= 0 {
		ErrorDNI := new(errores.DNIError)
		fmt.Fprintf(os.Stdout, "%s\n", ErrorDNI.Error())
		return true
	}
	if (*fila).BuscarDNI(dni, *padron) == false {
		ErrorPadron := new(errores.DNIFueraPadron)
		fmt.Fprintf(os.Stdout, "%s\n", ErrorPadron.Error())
		fmt.Fprintf(os.Stdout, "%d\n", dni)
		return true
	}
	return false
}

func atrapar_errores_votar(alternativa int, lista_partidos []votos.Partido, dni int, fila votos.Fila, comando string) (votos.TipoVoto, error) {
	if alternativa > len(lista_partidos)-1 {
		return votos.TipoVoto(0), &errores.ErrorAlternativaInvalida{}
	}
	if fila.ValidarDNI(dni) {
		return votos.TipoVoto(0), &errores.ErrorVotanteFraudulento{dni}
	}
	tipo_voto, error := mappear_tipos_voto(comando)
	return tipo_voto, error
}

func mappear_tipos_voto(cargo string) (votos.TipoVoto, error) {
	if cargo != PRESIDENTE_STR && cargo != GOBERNADOR_STR && cargo != INTENDENTE_STR {
		return votos.TipoVoto(0), &errores.ErrorTipoVoto{}
	}
	m := make(map[string]int)
	m[PRESIDENTE_STR] = PRESIDENTE_INT
	m[GOBERNADOR_STR] = GOBERNADOR_INT
	m[INTENDENTE_STR] = INTENDENTE_INT
	return votos.TipoVoto(m[cargo]), nil
}

func leer_input(
	input []string,
	fila *votos.Fila,
	lista_partidos []votos.Partido,
	padron *[]int) {

	comando := strings.Trim(input[0], "\n")
	switch comando {

	case INGRESAR:
		if len(input) != 2 {
			return
		}
		dni_string := strings.Trim(input[1], "\n")
		dni, _ := strconv.Atoi(dni_string)
		if !atrapar_errores_dni(dni, fila, padron) {
			ingresar(dni, fila)
		}

	case VOTAR:
		if len(input) != 3 {
			return
		}
		if !(*fila).HayVotantes() {
			ErrorFila := new(errores.FilaVacia)
			fmt.Fprintf(os.Stdout, "%s\n", ErrorFila.Error())
			return
		}
		alternativa_string := strings.Trim(input[2], "\n")
		alternativa, _ := strconv.Atoi(alternativa_string)
		tipo_voto, error := atrapar_errores_votar(alternativa, lista_partidos, (*(*fila).VerActual()).LeerDNI(), *fila, input[1])
		if error != nil {
			fmt.Fprintf(os.Stdout, "%s\n", error.Error())
			return
		}
		votar(tipo_voto, alternativa, fila)

	case DESHACER:
		if !(*fila).HayVotantes() {
			ErrorFila := new(errores.FilaVacia)
			fmt.Fprintf(os.Stdout, "%s\n", ErrorFila.Error())
			return
		}
		votante := (*fila).VerActual()
		error := (*votante).Deshacer()
		if error != nil {
			fmt.Fprintf(os.Stdout, "%s\n", error.Error())
			return
		}
		error = atrapar_errores_deshacer(*fila, (*votante).LeerDNI())
		if error != nil {
			fmt.Fprintf(os.Stdout, "%s\n", error.Error())
			return
		}
		fmt.Println(VALIDACION)

	case FIN_VOTO:
		if !(*fila).HayVotantes() {
			ErrorFila := new(errores.FilaVacia)
			fmt.Fprintf(os.Stdout, "%s\n", ErrorFila.Error())
			return
		}
		finalizar_voto(*fila, lista_partidos)
	}
}

func abrirArchivo(ruta string) *bufio.Scanner {
	// Devuelve un archivo abierto.

	archivo, err := os.Open(ruta)
	if err != nil {
		errorLeerArchivo := new(errores.ErrorLeerArchivo)
		fmt.Fprintf(os.Stderr, "%s", errorLeerArchivo.Error())
	}
	return bufio.NewScanner(archivo)
}
