package carrito

import (
	"errors"
	"testing"

	"tienda/productos"
)

var laptop = productos.Producto{ID: 1, Nombre: "Laptop", Categoria: "T", PrecioCentavos: 80000, Stock: 3}

func TestAgregarYTotal(t *testing.T) {
	var c Carrito
	c, _ = Agregar(c, laptop, 1)
	c, _ = Agregar(c, laptop, 2)
	if len(c) != 1 || c[0].Cantidad != 3 || Total(c) != 240000 || CantidadUnidades(c) != 3 {
		t.Errorf("carrito inesperado: %+v", c)
	}
}

func TestAgregarErrores(t *testing.T) {
	var c Carrito
	if _, err := Agregar(c, laptop, 0); !errors.Is(err, ErrCantidadInvalida) {
		t.Errorf("cantidad 0: %v", err)
	}
	if _, err := Agregar(c, laptop, 4); !errors.Is(err, productos.ErrStockInsuficiente) {
		t.Errorf("sin stock: %v", err)
	}
}

func TestInmutabilidad(t *testing.T) {
	original, _ := Agregar(nil, laptop, 1)
	_, _ = Agregar(original, laptop, 1)
	if original[0].Cantidad != 1 {
		t.Errorf("el carrito original fue modificado")
	}
}

func TestQuitar(t *testing.T) {
	c, _ := Agregar(nil, laptop, 1)
	c, err := Quitar(c, 1)
	if err != nil || len(c) != 0 {
		t.Errorf("err=%v c=%v", err, c)
	}
	if _, err = Quitar(c, 1); !errors.Is(err, ErrNoEnCarrito) {
		t.Errorf("se esperaba ErrNoEnCarrito: %v", err)
	}
}
