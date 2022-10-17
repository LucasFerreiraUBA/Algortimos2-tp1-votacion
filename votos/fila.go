package votos

type Fila interface {
	VerActual() Votante

	HayVotantes() bool

	Ingresar(Votante)

	FinalizarVoto()

	
}
