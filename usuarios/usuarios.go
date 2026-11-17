// Package usuarios gestiona el registro y la autenticación.
// Seguridad: la clave nunca se guarda en texto plano; se guarda un hash con sal.
package usuarios

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode"

	"tienda/funcional"
)

// iteraciones encarece los ataques de fuerza bruta. En producción se usaría bcrypt o argon2
// (paquetes de terceros); aquí se usa solo la librería estándar.
const iteraciones = 100_000

// Usuario guarda el hash y la sal, nunca la clave.
type Usuario struct {
	ID     int
	Nombre string
	Email  string
	Sal    string
	Hash   string
}

var (
	ErrDatoInvalido = errors.New("dato de usuario inválido")
	ErrEmailEnUso   = errors.New("el correo ya está registrado")
	ErrClaveDebil   = errors.New("la clave no cumple la política de seguridad")
	// Mensaje único para no revelar si el correo existe o si falló la clave.
	ErrCredenciales = errors.New("correo o clave incorrectos")
)

// GenerarSal crea 16 bytes aleatorios criptográficamente seguros (única función impura).
func GenerarSal() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Hashear es una función pura: mismos datos, mismo resultado.
func Hashear(clave, sal string) string {
	h := sha256.Sum256([]byte(sal + clave))
	for i := 1; i < iteraciones; i++ {
		h = sha256.Sum256(append(h[:], []byte(sal)...))
	}
	return hex.EncodeToString(h[:])
}

// ValidarClave exige mínimo 8 caracteres, mayúscula, minúscula y número.
func ValidarClave(clave string) error {
	var may, min, num bool
	for _, r := range clave {
		switch {
		case unicode.IsUpper(r):
			may = true
		case unicode.IsLower(r):
			min = true
		case unicode.IsDigit(r):
			num = true
		}
	}
	if len([]rune(clave)) < 8 || !may || !min || !num {
		return fmt.Errorf("%w: mínimo 8 caracteres con mayúscula, minúscula y número", ErrClaveDebil)
	}
	return nil
}

// ValidarEmail acepta solo direcciones con formato válido.
func ValidarEmail(email string) error {
	dir, err := mail.ParseAddress(email)
	if err != nil || dir.Address != email || !strings.Contains(email[strings.LastIndex(email, "@"):], ".") {
		return fmt.Errorf("%w: correo con formato inválido", ErrDatoInvalido)
	}
	return nil
}

// Registrar devuelve la lista nueva con el usuario agregado.
func Registrar(lista []Usuario, id int, nombre, email, clave, sal string) ([]Usuario, Usuario, error) {
	nombre = strings.TrimSpace(nombre)
	email = strings.ToLower(strings.TrimSpace(email))
	if nombre == "" {
		return lista, Usuario{}, fmt.Errorf("%w: el nombre no puede estar vacío", ErrDatoInvalido)
	}
	if err := ValidarEmail(email); err != nil {
		return lista, Usuario{}, err
	}
	if err := ValidarClave(clave); err != nil {
		return lista, Usuario{}, err
	}
	if len(funcional.Filtrar(lista, func(u Usuario) bool { return u.Email == email })) > 0 {
		return lista, Usuario{}, ErrEmailEnUso
	}
	u := Usuario{ID: id, Nombre: nombre, Email: email, Sal: sal, Hash: Hashear(clave, sal)}
	return append(append([]Usuario{}, lista...), u), u, nil
}

// Autenticar verifica correo y clave comparando en tiempo constante.
func Autenticar(lista []Usuario, email, clave string) (Usuario, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	encontrados := funcional.Filtrar(lista, func(u Usuario) bool { return u.Email == email })
	if len(encontrados) == 0 {
		Hashear(clave, "sal-falsa") // gasta el mismo tiempo para no delatar si el correo existe
		return Usuario{}, ErrCredenciales
	}
	u := encontrados[0]
	if subtle.ConstantTimeCompare([]byte(Hashear(clave, u.Sal)), []byte(u.Hash)) != 1 {
		return Usuario{}, ErrCredenciales
	}
	return u, nil
}
