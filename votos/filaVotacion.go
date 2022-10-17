package votos

import (
	TDACola "rerepolez/votos/cola"
)

type FilaVotacion struct {
	cola_votantes      TDACola.Cola[Votante]
	votantes_repetidos []int
	votante_actual     Votante
}

func CrearFilaVotacion() Fila {
	fila := new(FilaVotacion)
	return fila
}

func (fila FilaVotacion) HayVotantes() bool {
	return !fila.cola_votantes.EstaVacia()
}

func (fila FilaVotacion) VerActual() Votante {
	return fila.votante_actual
}

func (fila *FilaVotacion) Ingresar(votante Votante) {
	fila.cola_votantes.Encolar(votante)
	if !fila.HayVotantes() {
		fila.votante_actual = votante
	}
}

func (fila *FilaVotacion) FinalizarVoto() {
	dni_votante := fila.votante_actual.LeerDNI()
	if dni_votante > fila.votantes_repetidos[len(fila.votantes_repetidos)-1] {
		fila.votantes_repetidos[len(fila.votantes_repetidos)] = dni_votante
	} else {
		for i := 0; i < len(fila.votantes_repetidos); i++ {
			if fila.votantes_repetidos[i] < dni_votante {
				continue
			}
			fila.votantes_repetidos = append(fila.votantes_repetidos[:i+1], fila.votantes_repetidos[i:]...)
			fila.votantes_repetidos[i] = dni_votante
			break
		}
	}
	fila.cola_votantes.Desencolar()
	fila.votante_actual = fila.cola_votantes.VerPrimero()
}

func (fila FilaVotacion) YaVoto(dni int) bool {
	return fila.buscarDNI(fila.votantes_repetidos, dni, 0, len(fila.votantes_repetidos)-1)
}

func (fila FilaVotacion) buscarDNI(lista []int, elemento int, principio int, fin int) bool {
	if fin < principio {
		return false
	}
	medio := (fin + principio) / 2
	if elemento == lista[medio] {
		return true
	}
	if elemento > lista[medio] {
		return fila.buscarDNI(lista, elemento, medio, fin)
	}

	return fila.buscarDNI(lista, elemento, principio, medio)
}
