package pila

/* Definición del struct pila proporcionado por la cátedra. */

type pilaDinamica[T any] struct {
	datos    []T
	cantidad int
}

func CrearPilaDinamica[T any]() Pila[T] {
	pila := new(pilaDinamica[T]) // ACA esta el cambio
	// hago lo que deba hacer
	pila.redimensionar(5)

	return pila
}
func (pila *pilaDinamica[T]) Apilar(dato T) {
	capacidadPila := len(pila.datos)
	pila.datos[pila.cantidad] = dato
	pila.cantidad++
	if pila.cantidad == capacidadPila {
		pila.redimensionar(capacidadPila * 2)
	}
}

func (pila *pilaDinamica[T]) Desapilar() T {
	capacidadPila := len(pila.datos)
	if pila.EstaVacia() {
		panic("La pila esta vacia")
	}
	dato := pila.datos[pila.cantidad-1]
	pila.cantidad--
	if pila.cantidad <= capacidadPila/4 {
		pila.redimensionar(capacidadPila / 2)
	}
	return dato

}

func (pila pilaDinamica[T]) EstaVacia() bool {
	return pila.cantidad == 0
}

func (pila *pilaDinamica[T]) VerTope() T {
	if pila.EstaVacia() {
		panic("La pila esta vacia")
	}
	return pila.datos[pila.cantidad-1]
}

func (pila *pilaDinamica[T]) redimensionar(capacidad int) {
	nueva_pila := make([]T, capacidad)
	copy(nueva_pila, pila.datos)
	pila.datos = nueva_pila
}
