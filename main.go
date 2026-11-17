// Sistema de Gestión de E-commerce (Etapa 1) - programación funcional en Go.
// Este programa demuestra el flujo completo: catálogo -> usuario -> carrito -> pedido.
package main

import (
	"fmt"

	"tienda/carrito"
	"tienda/funcional"
	"tienda/pedidos"
	"tienda/productos"
	"tienda/usuarios"
)

func titulo(texto string) { fmt.Printf("\n=== %s ===\n", texto) }

func mostrarCatalogo(ps []productos.Producto) {
	for _, p := range ps {
		fmt.Printf("  #%d %-22s %-12s %8s  stock: %d\n",
			p.ID, p.Nombre, p.Categoria, productos.FormatearPrecio(p.PrecioCentavos), p.Stock)
	}
}

// revisar imprime el error (si existe) sin detener el programa.
func revisar(etiqueta string, err error) {
	if err != nil {
		fmt.Printf("  [ERROR esperado] %s: %v\n", etiqueta, err)
	}
}

func catalogoInicial() []productos.Producto {
	datos := []struct {
		id        int
		nombre    string
		categoria string
		precio    int
		stock     int
	}{
		{1, "Laptop Pro 14", "Tecnologia", 89999, 5},
		{2, "Mouse inalambrico", "Tecnologia", 1999, 30},
		{3, "Teclado mecanico", "Tecnologia", 5999, 12},
		{4, "Cuaderno universitario", "Papeleria", 350, 100},
		{5, "Mochila antirrobo", "Accesorios", 4500, 8},
		{6, "Audifonos bluetooth", "Tecnologia", 3999, 0},
	}
	catalogo := []productos.Producto{}
	for _, d := range datos {
		p, err := productos.Nuevo(d.id, d.nombre, d.categoria, d.precio, d.stock)
		if err != nil {
			fmt.Println("producto inválido:", err)
			continue
		}
		catalogo = append(catalogo, p)
	}
	return catalogo
}

func main() {
	catalogo := catalogoInicial()

	titulo("1. Catálogo completo")
	mostrarCatalogo(catalogo)

	titulo("2. Búsqueda funcional: Tecnologia, con stock, hasta $60.00, ordenado por precio")
	filtro := productos.Y(productos.PorCategoria("Tecnologia"), productos.ConStock(), productos.PrecioMaximo(6000))
	mostrarCatalogo(productos.OrdenarPorPrecio(productos.Buscar(catalogo, filtro)))

	titulo("3. Usuarios: registro y autenticación")
	var lista []usuarios.Usuario
	sal, _ := usuarios.GenerarSal()
	_, _, err := usuarios.Registrar(lista, 1, "Javier Ortiz", "javier@correo", "Clave1234", sal)
	revisar("correo inválido", err)
	_, _, err = usuarios.Registrar(lista, 1, "Javier Ortiz", "javier@correo.com", "abc", sal)
	revisar("clave débil", err)
	lista, u, err := usuarios.Registrar(lista, 1, "Javier Ortiz", "javier@correo.com", "Clave1234", sal)
	revisar("registro", err)
	fmt.Printf("  Registrado: %s (hash guardado: %s...)\n", u.Nombre, u.Hash[:16])
	_, err = usuarios.Autenticar(lista, "javier@correo.com", "ClaveIncorrecta1")
	revisar("clave incorrecta", err)
	sesion, err := usuarios.Autenticar(lista, "javier@correo.com", "Clave1234")
	revisar("login", err)
	fmt.Printf("  Sesión iniciada: %s\n", sesion.Nombre)

	titulo("4. Carrito (cada operación devuelve un carrito nuevo)")
	var c carrito.Carrito
	laptop, _ := productos.BuscarPorID(catalogo, 1)
	mouse, _ := productos.BuscarPorID(catalogo, 2)
	audifonos, _ := productos.BuscarPorID(catalogo, 6)
	c, err = carrito.Agregar(c, laptop, 1)
	revisar("agregar laptop", err)
	c, err = carrito.Agregar(c, mouse, 2)
	revisar("agregar mouse", err)
	_, err = carrito.Agregar(c, laptop, 10)
	revisar("más unidades que el stock", err)
	_, err = carrito.Agregar(c, audifonos, 1)
	revisar("producto agotado", err)
	fmt.Printf("  Unidades: %d | Total: %s\n", carrito.CantidadUnidades(c), productos.FormatearPrecio(carrito.Total(c)))

	titulo("5. Pedido y ciclo de vida")
	pedido, nuevoCatalogo, err := pedidos.Crear(100, sesion.ID, c, catalogo)
	revisar("crear pedido", err)
	fmt.Printf("  Pedido #%d creado: %s por %s\n", pedido.ID, pedido.Estado, productos.FormatearPrecio(pedido.TotalCentavos))
	fmt.Println("  Stock ORIGINAL de la laptop:", laptop.Stock, "| Stock en el catálogo NUEVO:",
		funcional.Filtrar(nuevoCatalogo, func(p productos.Producto) bool { return p.ID == 1 })[0].Stock)
	for i := 0; i < 3; i++ {
		pedido, err = pedidos.Avanzar(pedido)
		revisar("avanzar", err)
		fmt.Printf("  Estado -> %s\n", pedido.Estado)
	}
	_, err = pedidos.Cancelar(pedido)
	revisar("cancelar pedido entregado", err)
}
