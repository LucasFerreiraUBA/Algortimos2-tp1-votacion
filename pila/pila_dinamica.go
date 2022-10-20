package pila

/* Definición del struct pila proporcionado por la cátedra. */

const CAPACIDAD_INICIAL = 10
const CRITERIO_DISMINUICION = 4
const CANTIDAD_INCREMENTACION = 2
const CANTIDAD_REDUCCION = 2

type pilaDinamica[T any] struct {
	datos     []T
	cantidad  int
	capacidad int
}

func CrearPilaDinamica[T any]() Pila[T] {
	pila := new(pilaDinamica[T])
	pila.cantidad = 0
	pila.datos = make([]T, CAPACIDAD_INICIAL, CAPACIDAD_INICIAL)
	pila.capacidad = CAPACIDAD_INICIAL
	return pila
}

func (pila *pilaDinamica[T]) EstaVacia() bool {
	return pila.cantidad == 0
}

func (pila *pilaDinamica[T]) VerTope() T {
	if pila.EstaVacia() {
		panic("La pila esta vacia")
	}
	return pila.datos[pila.cantidad-1]
}

func (pila *pilaDinamica[T]) Apilar(elemento T) {
	if pila.cantidad == pila.capacidad {
		pila.redimenzionar(pila.capacidad * CANTIDAD_INCREMENTACION)
	}
	pila.datos[pila.cantidad] = elemento
	pila.cantidad++
}

func (pila *pilaDinamica[T]) Desapilar() T {
	if pila.EstaVacia() {
		panic("La pila esta vacia")
	}
	if pila.cantidad*CRITERIO_DISMINUICION <= pila.capacidad {
		pila.redimenzionar(pila.capacidad / CANTIDAD_REDUCCION)
	}
	elemento_tope := pila.datos[pila.cantidad-1]
	pila.cantidad--
	return elemento_tope
}

func (pila *pilaDinamica[T]) redimenzionar(nuevo_tamanio int) {
	datos_nuevo := make([]T, nuevo_tamanio)
	copy(datos_nuevo, pila.datos)
	(*pila).datos = datos_nuevo
	(*pila).capacidad = nuevo_tamanio
}
