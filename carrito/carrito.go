// Package carrito gestiona el carrito de compras. Cada operación devuelve un carrito NUEVO.
package carrito

import (
	"errors"
	"fmt"

	"tienda/funcional"
	"tienda/productos"
)

// Item es una línea del carrito.
type Item struct {
	ProductoID     int
	Nombre         string
	PrecioCentavos int
	Cantidad       int
}

// Carrito es una lista de items.
type Carrito []Item

var (
	ErrCantidadInvalida = errors.New("la cantidad debe ser mayor que cero")
	ErrNoEnCarrito      = errors.New("el producto no está en el carrito")
)

// Agregar suma cantidad unidades del producto, validando el stock disponible.
func Agregar(c Carrito, p productos.Producto, cantidad int) (Carrito, error) {
	if cantidad <= 0 {
		return c, ErrCantidadInvalida
	}
	yaEnCarrito := funcional.Reducir(c, 0, func(acc int, i Item) int {
		if i.ProductoID == p.ID {
			return acc + i.Cantidad
		}
		return acc
	})
	if yaEnCarrito+cantidad > p.Stock {
		return c, fmt.Errorf("%w: %s (disponible %d, en carrito %d)",
			productos.ErrStockInsuficiente, p.Nombre, p.Stock, yaEnCarrito)
	}
	if yaEnCarrito > 0 {
		return Carrito(funcional.Mapear(c, func(i Item) Item {
			if i.ProductoID == p.ID {
				i.Cantidad += cantidad
			}
			return i
		})), nil
	}
	nuevo := append(Carrito{}, c...)
	return append(nuevo, Item{p.ID, p.Nombre, p.PrecioCentavos, cantidad}), nil
}

// Quitar elimina por completo un producto del carrito.
func Quitar(c Carrito, productoID int) (Carrito, error) {
	resto := funcional.Filtrar(c, func(i Item) bool { return i.ProductoID != productoID })
	if len(resto) == len(c) {
		return c, ErrNoEnCarrito
	}
	return Carrito(resto), nil
}

// Total devuelve el total en centavos.
func Total(c Carrito) int {
	return funcional.Reducir(c, 0, func(acc int, i Item) int { return acc + i.PrecioCentavos*i.Cantidad })
}

// CantidadUnidades devuelve el número total de unidades.
func CantidadUnidades(c Carrito) int {
	return funcional.Reducir(c, 0, func(acc int, i Item) int { return acc + i.Cantidad })
}
