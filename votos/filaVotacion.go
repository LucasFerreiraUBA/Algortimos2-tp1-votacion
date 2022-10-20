package votos

import (
	TDACola "rerepolez/cola"
)

type FilaVotacion struct {
	cola_votantes TDACola.Cola[Votante]
	votantes_repetidos []int
	votante_actual Votante
}

func CrearFilaVotacion() Fila {
	fila := new(FilaVotacion)
	fila.cola_votantes = TDACola.CrearColaEnlazada[Votante]()
	return fila 
}

func (fila FilaVotacion) HayVotantes() bool {
	return fila.cola_votantes.EstaVacia() == false
}

func (fila FilaVotacion) VerActual() *Votante {
	return &(fila.votante_actual)
}

func (fila *FilaVotacion) Ingresar(votante Votante) {
	if !fila.HayVotantes(){
		fila.votante_actual = votante
	}
	fila.cola_votantes.Encolar(votante)
}

func (fila *FilaVotacion) registrar_voto_finalizado(dni_votante int) {
	if len(fila.votantes_repetidos) == 0{
		fila.votantes_repetidos = append(fila.votantes_repetidos, dni_votante)
		return
	}
	if len(fila.votantes_repetidos) == 1 {
		if (fila.votantes_repetidos)[0] < dni_votante{
			fila.votantes_repetidos = append((fila.votantes_repetidos)[0:], dni_votante)
		} else {
			fila.votantes_repetidos = append((fila.votantes_repetidos)[:1], (fila.votantes_repetidos)[0])
			(fila.votantes_repetidos)[0] = dni_votante
		}
		return
	}
	for i := 0; i < len(fila.votantes_repetidos); i++ {
		if dni_votante > (fila.votantes_repetidos)[i]{
			continue
		}
		fila.votantes_repetidos = append((fila.votantes_repetidos)[:i + 1], (fila.votantes_repetidos)[i:]...)
		(fila.votantes_repetidos)[i] = dni_votante
		return
	}
	fila.votantes_repetidos = append(fila.votantes_repetidos, dni_votante)
}


func (fila *FilaVotacion) FinalizarVoto() {
	dni_votante := fila.votante_actual.LeerDNI()
	fila.registrar_voto_finalizado(dni_votante)
	fila.cola_votantes.Desencolar()
	if fila.cola_votantes.EstaVacia(){
		fila.votante_actual = nil
		return
	}
	fila.votante_actual = fila.cola_votantes.VerPrimero()
}


func (fila FilaVotacion) BuscarDNI(dni int, lista []int) bool {
	return fila.buscarDNI(lista, dni, 0, len(lista) - 1)
}

func (fila FilaVotacion) ValidarDNI(dni int) bool {
	return fila.buscarDNI(fila.votantes_repetidos, dni, 0, len(fila.votantes_repetidos) - 1)
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
		return fila.buscarDNI(lista, elemento, medio + 1, fin)
	}

	return fila.buscarDNI(lista, elemento, principio, medio - 1)
}


