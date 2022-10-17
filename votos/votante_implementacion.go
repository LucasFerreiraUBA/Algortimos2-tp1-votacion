package votos

import (
	"rerepolez/errores"
	TDAPila "rerepolez/votos/pila"
)

type votoProvisorio struct {
	alternativas TDAPila.Pila[int]
	puestos      TDAPila.Pila[TipoVoto]
}

type votanteImplementacion struct {
	dni             int
	voto            Voto
	ya_voto         bool
	voto_provisorio *votoProvisorio
}

func CrearVotante(dni int) Votante {
	votante := new(votanteImplementacion)
	return votante
}

func (votante votanteImplementacion) YaVoto() bool {
	return votante.ya_voto
}

func (votante votanteImplementacion) LeerDNI() int {
	return votante.dni
}

func (votante *votanteImplementacion) Votar(tipo TipoVoto, alternativa int) error {
	if votante.voto.Impugnado || votante.YaVoto() {
		return &errores.ErrorVotanteFraudulento{votante.dni}
	}
	votante.voto_provisorio.alternativas.Apilar(votante.voto.VotoPorTipo[tipo])
	votante.voto_provisorio.puestos.Apilar(tipo)
	votante.voto.VotoPorTipo[tipo] = alternativa
	return nil
}

func (votante *votanteImplementacion) Deshacer() error {
	if votante.voto_provisorio.puestos.EstaVacia() {
		return &errores.ErrorNoHayVotosAnteriores{}
	}
	votante.voto.VotoPorTipo[votante.voto_provisorio.puestos.Desapilar()] = votante.voto_provisorio.alternativas.Desapilar()
	return nil
}

func (votante *votanteImplementacion) Invalidar() {
	votante.ya_voto = true
}

func (votante *votanteImplementacion) FinVoto() (Voto, error) {
	if votante.ya_voto {
		error := &errores.ErrorVotanteFraudulento{votante.dni}
		return votante.voto, error
	}
	return votante.voto, nil
}
