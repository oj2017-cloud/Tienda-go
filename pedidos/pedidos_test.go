package pedidos

import (
	"errors"
	"testing"

	"tienda/carrito"
	"tienda/productos"
)

func escenario() (carrito.Carrito, []productos.Producto) {
	cat := []productos.Producto{{ID: 1, Nombre: "Laptop", Categoria: "T", PrecioCentavos: 80000, Stock: 2}}
	c, _ := carrito.Agregar(nil, cat[0], 2)
	return c, cat
}

func TestCrearPedido(t *testing.T) {
	c, cat := escenario()
	p, nuevo, err := Crear(1, 7, c, cat)
	if err != nil || p.Estado != Pendiente || p.TotalCentavos != 160000 || nuevo[0].Stock != 0 || cat[0].Stock != 2 {
		t.Errorf("err=%v pedido=%+v", err, p)
	}
}

func TestCrearCarritoVacio(t *testing.T) {
	_, cat := escenario()
	if _, _, err := Crear(1, 7, nil, cat); !errors.Is(err, ErrCarritoVacio) {
		t.Errorf("err=%v", err)
	}
}

func TestStockCambioDespuesDeAgregar(t *testing.T) {
	c, cat := escenario()
	cat[0].Stock = 1 // alguien más compró antes
	if _, _, err := Crear(1, 7, c, cat); !errors.Is(err, productos.ErrStockInsuficiente) {
		t.Errorf("err=%v", err)
	}
}

func TestCicloDeVida(t *testing.T) {
	p := Pedido{Estado: Pendiente}
	esperados := []Estado{Pagado, Enviado, Entregado}
	for _, e := range esperados {
		var err error
		p, err = Avanzar(p)
		if err != nil || p.Estado != e {
			t.Fatalf("esperado %s, obtuvo %s (%v)", e, p.Estado, err)
		}
	}
	if _, err := Avanzar(p); !errors.Is(err, ErrTransicionInvalida) {
		t.Errorf("err=%v", err)
	}
	if _, err := Cancelar(p); !errors.Is(err, ErrTransicionInvalida) {
		t.Errorf("err=%v", err)
	}
	if c, err := Cancelar(Pedido{Estado: Pagado}); err != nil || c.Estado != Cancelado {
		t.Errorf("err=%v estado=%s", err, c.Estado)
	}
}
