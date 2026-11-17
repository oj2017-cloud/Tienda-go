// Package pedidos convierte un carrito en un pedido y controla su ciclo de vida.
package pedidos

import (
	"errors"
	"fmt"

	"tienda/carrito"
	"tienda/productos"
)

// Estado del pedido.
type Estado string

const (
	Pendiente Estado = "PENDIENTE"
	Pagado    Estado = "PAGADO"
	Enviado   Estado = "ENVIADO"
	Entregado Estado = "ENTREGADO"
	Cancelado Estado = "CANCELADO"
)

// Pedido es una compra confirmada.
type Pedido struct {
	ID            int
	UsuarioID     int
	Items         carrito.Carrito
	TotalCentavos int
	Estado        Estado
}

var (
	ErrCarritoVacio       = errors.New("el carrito está vacío")
	ErrTransicionInvalida = errors.New("transición de estado no permitida")
)

// siguiente define el flujo normal: cada estado tiene a lo sumo un siguiente.
var siguiente = map[Estado]Estado{Pendiente: Pagado, Pagado: Enviado, Enviado: Entregado}

// Crear genera un pedido y devuelve el catálogo con el stock ya descontado.
// Si algún producto no tiene stock suficiente, no cambia nada y devuelve el error.
func Crear(id, usuarioID int, c carrito.Carrito, catalogo []productos.Producto) (Pedido, []productos.Producto, error) {
	if len(c) == 0 {
		return Pedido{}, catalogo, ErrCarritoVacio
	}
	nuevoCatalogo := catalogo
	for _, item := range c {
		var err error
		nuevoCatalogo, err = productos.DescontarStock(nuevoCatalogo, item.ProductoID, item.Cantidad)
		if err != nil {
			return Pedido{}, catalogo, err // el catálogo original queda intacto
		}
	}
	return Pedido{
		ID:            id,
		UsuarioID:     usuarioID,
		Items:         append(carrito.Carrito{}, c...),
		TotalCentavos: carrito.Total(c),
		Estado:        Pendiente,
	}, nuevoCatalogo, nil
}

// Avanzar mueve el pedido al siguiente estado del flujo normal.
func Avanzar(p Pedido) (Pedido, error) {
	sig, ok := siguiente[p.Estado]
	if !ok {
		return p, fmt.Errorf("%w: %s no tiene estado siguiente", ErrTransicionInvalida, p.Estado)
	}
	p.Estado = sig // p es una copia
	return p, nil
}

// Cancelar solo es posible antes del envío.
func Cancelar(p Pedido) (Pedido, error) {
	if p.Estado != Pendiente && p.Estado != Pagado {
		return p, fmt.Errorf("%w: no se puede cancelar un pedido %s", ErrTransicionInvalida, p.Estado)
	}
	p.Estado = Cancelado
	return p, nil
}
