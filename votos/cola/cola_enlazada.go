package cola

type nodoCola[T any] struct {
	dato T
	prox *nodoCola[T]
}

func crearNodo[T any]() *nodoCola[T] {
	nodo := new(nodoCola[T])
	return nodo
}

type colaEnlazada[T any] struct {
	primero *nodoCola[T]
	ultimo  *nodoCola[T]
}

func CrearColaEnlazada[T any]() Cola[T] {
	cola := new(colaEnlazada[T])
	return cola
}

func (cola *colaEnlazada[T]) EstaVacia() bool {
	return cola.primero == nil
}

func (cola *colaEnlazada[T]) Encolar(elem T) {
	nodo := crearNodo[T]()
	nodo.dato = elem
	if cola.EstaVacia() {
		cola.primero = nodo
	} else {
		cola.ultimo.prox = nodo
	}
	cola.ultimo = nodo
}

func (cola *colaEnlazada[T]) Desencolar() T {
	if cola.EstaVacia() {
		panic("La cola esta vacia")
	}
	dato := cola.primero.dato
	cola.primero = cola.primero.prox
	if cola.primero == nil {
		cola.ultimo = nil
	}
	return dato
}

func (cola *colaEnlazada[T]) VerPrimero() T {
	if cola.EstaVacia() {
		panic("La cola esta vacia")
	}
	return cola.primero.dato
}
