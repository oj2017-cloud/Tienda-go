// Package funcional reúne las funciones de orden superior que usa todo el sistema.
// Ninguna modifica los datos de entrada: siempre devuelven una colección nueva o un valor.
package funcional

// Filtrar devuelve un slice nuevo con los elementos que cumplen el predicado.
func Filtrar[T any](xs []T, pred func(T) bool) []T {
	resultado := make([]T, 0, len(xs))
	for _, x := range xs {
		if pred(x) {
			resultado = append(resultado, x)
		}
	}
	return resultado
}

// Mapear aplica f a cada elemento y devuelve un slice nuevo con los resultados.
func Mapear[T, U any](xs []T, f func(T) U) []U {
	resultado := make([]U, 0, len(xs))
	for _, x := range xs {
		resultado = append(resultado, f(x))
	}
	return resultado
}

// Reducir combina todos los elementos en un solo valor, empezando desde inicial.
func Reducir[T, A any](xs []T, inicial A, f func(A, T) A) A {
	acumulado := inicial
	for _, x := range xs {
		acumulado = f(acumulado, x)
	}
	return acumulado
}
