package comandos

import (
	"fmt"
	"os"
	"rerepolez/errores"
	"rerepolez/votos"
	"strconv"
	"strings"
)

const (
	// input del usuario:
	INGRESAR = "ingresar"
	DESHACER = "deshacer"
	VOTAR    = "votar"
	FIN_VOTO = "fin-votar"
	// ------------------
	VALIDACION     = "OK"
	PRESIDENTE_STR = "Presidente"
	INTENDENTE_STR = "Intendente"
	GOBERNADOR_STR = "Gobernador"
	PRESIDENTE_INT = 0
	GOBERNADOR_INT = 1
	INTENDENTE_INT = 2
)

func Leer_input(input []string, fila *votos.Fila, lista_partidos []votos.Partido, padron *[]int) {
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

func atrapar_errores_dni(dni int, fila *votos.Fila, padron *[]int) error {
	if dni <= 0 {
		return &errores.DNIError{}
	}
	if !(*fila).BuscarDNI(dni, *padron) {
		return &errores.DNIFueraPadron{}
	}
	return nil
}

func ingresar(dni int, fila *votos.Fila) {
	votante := votos.CrearVotante(dni)
	(*fila).Ingresar(votante)
	fmt.Fprintf(os.Stdout, "%s\n", VALIDACION)
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

func atrapar_errores_deshacer(fila votos.Fila, dni int) error {
	if !fila.HayVotantes() {
		return &errores.FilaVacia{}
	}
	if fila.ValidarDNI(dni) {
		return &errores.ErrorVotanteFraudulento{Dni: dni}
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
