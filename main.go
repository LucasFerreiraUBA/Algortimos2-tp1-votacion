package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	votos "rerepolez/votos"

	comandos "rerepolez/comandos"
	errores "rerepolez/errores"
)

const (
	CANTIDAD_CARGOS = 3
	PRESIDENTE_STR  = "Presidente"
	INTENDENTE_STR  = "Intendente"
	GOBERNADOR_STR  = "Gobernador"
)

func main() {
	scanners, error := validar_archivos()
	if error != nil {
		fmt.Fprintf(os.Stdout, "%s\n", error.Error())
		return
	}
	lista_partidos := Crear_lista_partidos(scanners[0])
	padron := Crear_Padron(scanners[1])
	fila := votos.CrearFilaVotacion()
	for {
		input := pedir_input(&fila, lista_partidos)
		if input != nil {
			comandos.Leer_input(input, &fila, lista_partidos, &padron)
		}
	}
}

func validar_archivos() ([]bufio.Scanner, error) {

	// Devuelve una lista de archivos que pueden leerse, si alguno no pudo abrir devuelve un error

	var lista_scanners []bufio.Scanner
	args := os.Args
	if len(args) < 3 {
		return nil, &errores.ErrorParametros{}
	}
	archivos := args[1:]
	for i := 0; i < len(archivos); i++ {
		archivo, error := os.Open(archivos[i])
		scanner := bufio.NewScanner(archivo)
		if error != nil {
			return nil, &errores.ErrorLeerArchivo{}
		}
		lista_scanners = append(lista_scanners, *scanner)
	}
	return lista_scanners, nil
}

func Crear_lista_partidos(s1 bufio.Scanner) []votos.Partido {

	// Pre : archivo partidos.csv
	// Post: Crea una lista de elementos type Partido

	var lista_partidos []votos.Partido
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
	return lista_partidos
}

func Crear_Padron(scanner bufio.Scanner) []int {

	// Pre:  archivo padron_txt
	// Post: Crea una lista de elementos type Padron

	padron := make([]int, 0, 10)
	for scanner.Scan() {
		elemento, _ := strconv.Atoi(scanner.Text())
		insertar_elemento(&padron, elemento)
	}
	return padron
}

func insertar_elemento(lista *[]int, elemento int) {

	// Inserta un padron en una lista de padrones de forma ordenada
	// utilizando el algoritmo de inserción.

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

func pedir_input(fila *votos.Fila, lista_partidos []votos.Partido) []string {
	inputReader := bufio.NewReader(os.Stdin)
	input, error := inputReader.ReadString('\n')
	if error == nil {
		return strings.Split(input, " ")
	}
	if error == io.EOF {
		imprimir_resultados(lista_partidos, *fila)
	}
	return nil
}

func imprimir_resultados(lista_partidos []votos.Partido, fila votos.Fila) {
	if fila.HayVotantes() {
		ErrorVotoIncompleto := new(errores.ErrorCiudadanosSinVotar)
		fmt.Fprintf(os.Stdout, "%s\n", ErrorVotoIncompleto.Error())
	}
	lista_cargos := [3]string{PRESIDENTE_STR, GOBERNADOR_STR, INTENDENTE_STR}
	for j := 0; j < CANTIDAD_CARGOS; j++ {
		fmt.Fprintf(os.Stdout, "%s:\n", lista_cargos[j])
		for i := 0; i < len(lista_partidos); i++ {
			if i == 0 {
				fmt.Fprintf(os.Stdout, "%s : %s\n",
					lista_partidos[i].Nombre(),                            // %s
					lista_partidos[i].ObtenerResultado(votos.TipoVoto(i))) // %s
				continue
			}
			fmt.Fprintf(os.Stdout, "%s - %s : %s\n",
				lista_partidos[i].Nombre(),                            // %s
				lista_partidos[i].Candidato(j),                        // %s
				lista_partidos[i].ObtenerResultado(votos.TipoVoto(j))) // %s
		}
	}
	fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d votos", 2) // cambiar luego el 2 por votos impugnados
}
