package cola

type colaEnlazada[T any] struct {
	primer_nodo *nodoCola[T]
	ultimo_nodo *nodoCola[T]
}

type nodoCola[T any] struct {
	dato T
	prox *nodoCola[T]
}

func crearNodo[T any](dato T) *nodoCola[T] {
	nodo := new(nodoCola[T])
	nodo.dato = dato
	return nodo
}

func CrearColaEnlazada[T any]() Cola[T] {
	cola := new(colaEnlazada[T])
	return cola
}

func (cola *colaEnlazada[T]) EstaVacia() bool {
	return cola.primer_nodo == nil
}

func (cola *colaEnlazada[T]) VerPrimero() T {
	if cola.EstaVacia() {
		panic("La cola esta vacia")
	}
	return cola.primer_nodo.dato
}

func (cola *colaEnlazada[T]) Encolar(dato T) {
	nodo := crearNodo(dato)
	if cola.EstaVacia() {
		cola.primer_nodo = nodo
	} else {
		cola.ultimo_nodo.prox = nodo
	}
	cola.ultimo_nodo = nodo
}

func (cola *colaEnlazada[T]) Desencolar() T {
	if cola.EstaVacia() {
		panic("La cola esta vacia")
	}
	dato := cola.primer_nodo.dato
	cola.primer_nodo = cola.primer_nodo.prox
	if cola.primer_nodo == nil {
		cola.ultimo_nodo = nil
	}
	return dato
}
