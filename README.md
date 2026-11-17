# Tienda-go
# Sistema de Gestión de E-commerce (Etapa 1)

Núcleo de una tienda en línea desarrollado en **Go** aplicando los principios de la **programación funcional**, arquitectura modular y controles de seguridad desde el diseño. Este proyecto forma parte de la asignatura *Lenguaje de Programación Go* en la **Universidad Internacional del Ecuador (UIDE)**.

---

## 🛠️ Tecnologías y Librerías
* **Lenguaje:** Go 1.22+ (utilizando exclusivamente la **librería estándar**, sin dependencias de terceros).
* **Control de versiones:** Git y GitHub.

---

## 📂 Estructura del Proyecto

El código está organizado en paquetes desacoplados con responsabilidades únicas, sus respectivos archivos de pruebas y sin dependencias circulares:

```text
tienda/
├── README.md
├── .gitignore
├── go.mod
├── main.go
├── funcional/
│   ├── funcional.go
│   └── funcional_test.go
├── productos/
│   ├── productos.go
│   └── productos_test.go
├── carrito/
│   ├── carrito.go
│   └── carrito_test.go
├── pedidos/
│   ├── pedidos.go
│   └── pedidos_test.go
└── usuarios/
    ├── usuarios.go
    └── usuarios_test.go
```

---

## 🏗️ Arquitectura y Módulos

| Paquete | Responsabilidad Principal |
| :--- | :--- |
| **`funcional`** | Funciones genéricas de orden superior: `Filtrar`, `Mapear` y `Reducir`. Es el núcleo base que no depende de ningún otro paquete. |
| **`productos`** | Catálogo en memoria, validación de datos, búsqueda mediante predicados, ordenamiento por precio y control de stock. |
| **`carrito`** | Operaciones de agregar y quitar productos, y cálculo de totales garantizando inmutabilidad (devuelve un carrito nuevo sin alterar el anterior). |
| **`pedidos`** | Conversión del carrito en pedido, descuento seguro de stock en el catálogo y control de transición de estados (`PENDIENTE` -> `PAGADO` -> `ENVIADO` -> `ENTREGADO` o `CANCELADO`). |
| **`usuarios`** | Registro, validación estricta de correo y claves, hash con sal (SHA-256) y autenticación resistente a ataques de tiempo. |
| **`main.go`** | Programa de demostración que ejecuta el flujo completo de compra de punta a punta con errores controlados. |

---

## 🚀 Programación Funcional Aplicada

El sistema implementa un estilo funcional en Go para garantizar un código predecible, verificable y libre de efectos secundarios:

* **Funciones Puras:** Operaciones como `Total`, `Hashear`, `ValidarClave` y `FormatearPrecio` devuelven siempre el mismo resultado para los mismos datos de entrada. *(Excepción documentada: `GenerarSal` se aísla como la única función impura por requerir aleatoriedad criptográfica).*
* **Inmutabilidad:** Las estructuras de datos no se mutan directamente; operaciones como `Agregar`, `Quitar`, `DescontarStock`, `Avanzar` y `Registrar` generan copias actualizadas del estado.
* **Funciones de Orden Superior y Closures:** Uso de `Filtrar`, `Mapear`, `Reducir`, y closures reutilizables (`PorCategoria`, `PrecioMaximo`, `NombreContiene`, `ConStock`) combinadas mediante composición con `productos.Y(...)`.
* **Errores como Valores:** Todas las funciones devuelven `(resultado, error)` para manejar casos límite sin detener abruptamente la ejecución.

---

## 🔒 Seguridad y Ciberseguridad desde el Diseño

Debido a la naturaleza del manejo de credenciales y transacciones, se integraron controles de seguridad usando la librería estándar de Go (`crypto/sha256`, `crypto/rand`, `crypto/subtle`, `net/mail`):

* **Almacenamiento Seguro de Claves:** Nunca se guarda la contraseña en texto plano; se aplica **SHA-256 con 100,000 iteraciones** para encarecer ataques de fuerza bruta.
* **Protección contra Tablas Arcoíris:** Generación de una sal aleatoria única por usuario mediante `crypto/rand`.
* **Prevención de Ataques de Tiempo y Enumeración:** Las credenciales se verifican mediante comparación en tiempo constante (`crypto/subtle`) y retornando mensajes de error genéricos.
* **Validación Estricta de Entradas:** Política de contraseñas seguras (mínimo 8 caracteres, mayúscula, minúscula y número), validación de formato de email y revalidación de stock al momento de crear el pedido.
* **Precisión Financiera:** Los montos monetarios se manejan en centavos (`int`) para evitar errores de redondeo de punto flotante (`float64`).

*(Nota: En un entorno de producción se recomienda utilizar algoritmos especializados como Argon2id o bcrypt; aquí se empleó SHA-256 iterado para cumplir con la restricción de usar únicamente la librería estándar).*

---

## 🧪 Pruebas Automáticas y Resultados

Cada paquete cuenta con su propia suite de pruebas unitarias (`go test`) para verificar casos exitosos, errores esperados y la preservación de la inmutabilidad:

```bash
# Ejecutar todas las pruebas con reporte de cobertura
go test ./... -cover
```

Resultados de cobertura obtenidos por módulo:
* `funcional`: **100.0%**
* `productos`: **87.8%**
* `carrito`: **95.8%**
* `pedidos`: **100.0%**
* `usuarios`: **84.1%**

---

## 💻 Instrucciones de Ejecución

1. **Clonar el repositorio:**
   ```bash
   git clone https://github.com/tu-usuario/tienda-go.git
   cd tienda
   ```

2. **Ejecutar el flujo de demostración principal:**
   ```bash
   go run .
   ```

---

## 👨‍💻 Autor

* **Bryan Javier Ortiz Cuero**  
* **Carrera:** Ingeniería en Ciberseguridad  
* **Institución:** Universidad Internacional del Ecuador (UIDE)  
* **Asignatura:** Lenguaje de Programación Go  
