// Package productos gestiona el catálogo de la tienda.
// El dinero se guarda en centavos (int) para evitar errores de redondeo de float64.
package productos

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"tienda/funcional"
)

// Producto es un artículo del catálogo.
type Producto struct {
	ID             int
	Nombre         string
	Categoria      string
	PrecioCentavos int
	Stock          int
}

// Predicado es una función que decide si un producto cumple una condición.
type Predicado func(Producto) bool

var (
	ErrNoEncontrado      = errors.New("producto no encontrado")
	ErrStockInsuficiente = errors.New("stock insuficiente")
	ErrDatoInvalido      = errors.New("dato de producto inválido")
)

// Nuevo valida los datos y devuelve un Producto (o un error).
func Nuevo(id int, nombre, categoria string, precioCentavos, stock int) (Producto, error) {
	nombre = strings.TrimSpace(nombre)
	switch {
	case id <= 0:
		return Producto{}, fmt.Errorf("%w: el ID debe ser positivo", ErrDatoInvalido)
	case nombre == "":
		return Producto{}, fmt.Errorf("%w: el nombre no puede estar vacío", ErrDatoInvalido)
	case precioCentavos <= 0:
		return Producto{}, fmt.Errorf("%w: el precio debe ser mayor que cero", ErrDatoInvalido)
	case stock < 0:
		return Producto{}, fmt.Errorf("%w: el stock no puede ser negativo", ErrDatoInvalido)
	}
	return Producto{id, nombre, strings.TrimSpace(categoria), precioCentavos, stock}, nil
}

// --- Predicados (funciones que devuelven funciones: closures) ---

func PorCategoria(categoria string) Predicado {
	return func(p Producto) bool { return strings.EqualFold(p.Categoria, categoria) }
}

func PrecioMaximo(maxCentavos int) Predicado {
	return func(p Producto) bool { return p.PrecioCentavos <= maxCentavos }
}

func NombreContiene(texto string) Predicado {
	t := strings.ToLower(texto)
	return func(p Producto) bool { return strings.Contains(strings.ToLower(p.Nombre), t) }
}

func ConStock() Predicado {
	return func(p Producto) bool { return p.Stock > 0 }
}

// Y combina varios predicados: el producto debe cumplirlos todos.
func Y(predicados ...Predicado) Predicado {
	return func(p Producto) bool {
		for _, pred := range predicados {
			if !pred(p) {
				return false
			}
		}
		return true
	}
}

// --- Consultas ---

// Buscar devuelve los productos que cumplen el predicado.
func Buscar(catalogo []Producto, pred Predicado) []Producto {
	return funcional.Filtrar(catalogo, pred)
}

// BuscarPorID devuelve el producto con ese ID o ErrNoEncontrado.
func BuscarPorID(catalogo []Producto, id int) (Producto, error) {
	encontrados := funcional.Filtrar(catalogo, func(p Producto) bool { return p.ID == id })
	if len(encontrados) == 0 {
		return Producto{}, fmt.Errorf("%w (ID %d)", ErrNoEncontrado, id)
	}
	return encontrados[0], nil
}

// OrdenarPorPrecio devuelve una COPIA ordenada de menor a mayor; no altera el original.
func OrdenarPorPrecio(catalogo []Producto) []Producto {
	copia := append([]Producto{}, catalogo...)
	sort.Slice(copia, func(i, j int) bool { return copia[i].PrecioCentavos < copia[j].PrecioCentavos })
	return copia
}

// DescontarStock devuelve un catálogo nuevo con el stock reducido.
func DescontarStock(catalogo []Producto, id, cantidad int) ([]Producto, error) {
	p, err := BuscarPorID(catalogo, id)
	if err != nil {
		return nil, err
	}
	if cantidad <= 0 || cantidad > p.Stock {
		return nil, fmt.Errorf("%w: %s (disponible %d, pedido %d)", ErrStockInsuficiente, p.Nombre, p.Stock, cantidad)
	}
	return funcional.Mapear(catalogo, func(x Producto) Producto {
		if x.ID == id {
			x.Stock -= cantidad // x es una copia: el original no cambia
		}
		return x
	}), nil
}

// FormatearPrecio convierte centavos a texto: 1999 -> "$19.99".
func FormatearPrecio(centavos int) string {
	return fmt.Sprintf("$%d.%02d", centavos/100, centavos%100)
}
