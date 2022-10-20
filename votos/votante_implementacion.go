package votos

import (
	"rerepolez/errores"
	TDAPila "rerepolez/pila"
)


type votanteImplementacion struct {
	dni int
	voto Voto
	ya_voto bool
	alternativas TDAPila.Pila[int]
	puestos TDAPila.Pila[TipoVoto]
}


func CrearVotante(dni int) Votante {
	votante := new(votanteImplementacion)
	votante.dni = dni
	votante.alternativas = TDAPila.CrearPilaDinamica[int]()
	votante.puestos = TDAPila.CrearPilaDinamica[TipoVoto]()
	return votante
}

func (votante votanteImplementacion) YaVoto() bool {
	return votante.ya_voto
}

func (votante votanteImplementacion) LeerDNI() int {
	return votante.dni
}

func (votante *votanteImplementacion) Votar(tipo TipoVoto, alternativa int) error {
	if votante.voto.Impugnado {
		return &errores.ErrorVotanteFraudulento{votante.dni}
	}
	votante.alternativas.Apilar(votante.voto.VotoPorTipo[tipo])
	votante.puestos.Apilar(tipo)
	votante.voto.VotoPorTipo[tipo] = alternativa
	return nil
}

func (votante *votanteImplementacion) Deshacer() error {
	if votante.puestos.EstaVacia() {
		return &errores.ErrorNoHayVotosAnteriores{}
	}
	votante.voto.VotoPorTipo[votante.puestos.Desapilar()] = votante.alternativas.Desapilar()
	return nil
}

func (votante *votanteImplementacion) Invalidar() {
	votante.ya_voto = true
}

func (votante *votanteImplementacion) FinVoto() (Voto, error) {
	if votante.ya_voto  {
		error := &errores.ErrorVotanteFraudulento{votante.dni}
		return votante.voto, error
	}
	return votante.voto, nil
}
