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
			leer_input(input, &fila, lista_partidos, &padron)
		}
	}
}

func atrapar_errores_deshacer(fila votos.Fila, dni int) error {
	if !fila.HayVotantes() {
		return &errores.FilaVacia{}
	}
	if fila.ValidarDNI(dni) {
		return &errores.ErrorVotanteFraudulento{Dni: dni}
	}
	return nil
}

func insertar_elemento(lista *[]int, elemento int) {
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

func Crear_lista_partidos(s1 bufio.Scanner) []votos.Partido {
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
	padron := make([]int, 0, 10)
	for scanner.Scan() {
		elemento, _ := strconv.Atoi(scanner.Text())
		insertar_elemento(&padron, elemento)
	}
	return padron
}

func ingresar(dni int, fila *votos.Fila) {
	votante := votos.CrearVotante(dni)
	(*fila).Ingresar(votante)
	fmt.Fprintf(os.Stdout, "%s\n", VALIDACION)
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
				fmt.Fprintf(os.Stdout, "%s : %s\n", lista_partidos[i].Nombre(), lista_partidos[i].ObtenerResultado(votos.TipoVoto(i)))
				continue
			}
			fmt.Fprintf(os.Stdout, "%s - %s : %s\n", lista_partidos[i].Nombre(), lista_partidos[i].Candidato(j), lista_partidos[i].ObtenerResultado(votos.TipoVoto(j)))
		}
	}
	fmt.Fprintf(os.Stdout, "\nVotos Impugnados: %d votos", 2) // cambiar luego el 2 por votos impugnados
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

func atrapar_errores_dni(dni int, fila *votos.Fila, padron *[]int) error {
	if dni <= 0 {
		return &errores.DNIError{}
	}
	if !(*fila).BuscarDNI(dni, *padron) {
		return &errores.DNIFueraPadron{}
	}
	return nil
}

func atrapar_errores_finalizar(fila votos.Fila) error {
	if !fila.HayVotantes() {
		return &errores.FilaVacia{}
	}
	dni := (*fila.VerActual()).LeerDNI()
	if fila.ValidarDNI(dni) {
		return &errores.ErrorVotanteFraudulento{Dni: dni}
	}
	return nil

}

func atrapar_errores_votar(alternativa int, lista_partidos []votos.Partido, dni int, fila votos.Fila, comando string) (votos.TipoVoto, error) {
	if alternativa > len(lista_partidos)-1 {
		return votos.TipoVoto(0), &errores.ErrorAlternativaInvalida{}
	}
	if fila.ValidarDNI(dni) {
		return votos.TipoVoto(0), &errores.ErrorVotanteFraudulento{Dni: dni}
	}
	if !fila.HayVotantes() {
		return votos.TipoVoto(0), &errores.FilaVacia{}
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

func leer_input(input []string, fila *votos.Fila, lista_partidos []votos.Partido, padron *[]int) {
	comando := strings.Trim(input[0], "\n")
	switch comando {

	case INGRESAR:
		if len(input) < 2 {
			return
		}
		dni_string := strings.Trim(input[1], "\n")
		dni, _ := strconv.Atoi(dni_string)
		error := atrapar_errores_dni(dni, fila, padron)
		if error != nil {
			fmt.Fprintf(os.Stdout, "%s\n", error.Error())
			return
		}
		ingresar(dni, fila)

	case VOTAR:
		if len(input) < 3 {
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
		error := atrapar_errores_finalizar(*fila)
		if error != nil {
			fmt.Fprintf(os.Stdout, "%s\n", error.Error())
			return
		}
		finalizar_voto(*fila, lista_partidos)
	}
}

func validar_archivos() ([]bufio.Scanner, error) {
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
