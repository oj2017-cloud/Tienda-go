package productos

import (
	"errors"
	"testing"
)

func catalogoPrueba() []Producto {
	return []Producto{
		{1, "Laptop", "Tecnologia", 80000, 2},
		{2, "Mouse", "Tecnologia", 2000, 0},
		{3, "Cuaderno", "Papeleria", 350, 10},
	}
}

func TestNuevoValida(t *testing.T) {
	casos := []struct {
		nombre         string
		id, precio, st int
		texto          string
		ok             bool
	}{
		{"valido", 1, 100, 1, "Mouse", true},
		{"id cero", 0, 100, 1, "Mouse", false},
		{"nombre vacio", 1, 100, 1, "  ", false},
		{"precio cero", 1, 0, 1, "Mouse", false},
		{"stock negativo", 1, 100, -1, "Mouse", false},
	}
	for _, c := range casos {
		_, err := Nuevo(c.id, c.texto, "X", c.precio, c.st)
		if (err == nil) != c.ok {
			t.Errorf("%s: err = %v", c.nombre, err)
		}
	}
}

func TestBusquedaCombinada(t *testing.T) {
	r := Buscar(catalogoPrueba(), Y(PorCategoria("tecnologia"), ConStock()))
	if len(r) != 1 || r[0].ID != 1 {
		t.Errorf("resultado inesperado: %v", r)
	}
}

func TestDescontarStockNoModificaOriginal(t *testing.T) {
	orig := catalogoPrueba()
	nuevo, err := DescontarStock(orig, 1, 2)
	if err != nil || nuevo[0].Stock != 0 || orig[0].Stock != 2 {
		t.Errorf("err=%v nuevo=%d orig=%d", err, nuevo[0].Stock, orig[0].Stock)
	}
	if _, err = DescontarStock(orig, 1, 3); !errors.Is(err, ErrStockInsuficiente) {
		t.Errorf("se esperaba ErrStockInsuficiente, hubo %v", err)
	}
	if _, err = DescontarStock(orig, 99, 1); !errors.Is(err, ErrNoEncontrado) {
		t.Errorf("se esperaba ErrNoEncontrado, hubo %v", err)
	}
}

func TestOrdenarNoAlteraOriginal(t *testing.T) {
	orig := catalogoPrueba()
	o := OrdenarPorPrecio(orig)
	if o[0].ID != 3 || orig[0].ID != 1 {
		t.Errorf("orden incorrecto u original alterado")
	}
}

func TestFormatearPrecio(t *testing.T) {
	if got := FormatearPrecio(1999); got != "$19.99" {
		t.Errorf("got %s", got)
	}
	if got := FormatearPrecio(5); got != "$0.05" {
		t.Errorf("got %s", got)
	}
}
