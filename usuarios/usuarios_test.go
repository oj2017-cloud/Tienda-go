package usuarios

import (
	"errors"
	"testing"
)

func TestValidarClave(t *testing.T) {
	casos := map[string]bool{"Clave1234": true, "clave1234": false, "CLAVE1234": false, "Clave": false, "ClaveClave": false}
	for clave, ok := range casos {
		if (ValidarClave(clave) == nil) != ok {
			t.Errorf("clave %q: esperado ok=%v", clave, ok)
		}
	}
}

func TestValidarEmail(t *testing.T) {
	if ValidarEmail("a@b.com") != nil {
		t.Error("a@b.com debería ser válido")
	}
	for _, malo := range []string{"sinarroba", "a@", "@b.com", "Nombre <a@b.com>", "a@correo"} {
		if ValidarEmail(malo) == nil {
			t.Errorf("%q debería ser inválido", malo)
		}
	}
}

func TestRegistrarYAutenticar(t *testing.T) {
	lista, u, err := Registrar(nil, 1, "Ana", "Ana@Correo.com", "Clave1234", "sal1")
	if err != nil || u.Email != "ana@correo.com" || u.Hash == "Clave1234" {
		t.Fatalf("err=%v usuario=%+v", err, u)
	}
	if _, _, err = Registrar(lista, 2, "Otra", "ana@correo.com", "Clave1234", "sal2"); !errors.Is(err, ErrEmailEnUso) {
		t.Errorf("duplicado: %v", err)
	}
	if _, err = Autenticar(lista, "ana@correo.com", "Clave1234"); err != nil {
		t.Errorf("login válido: %v", err)
	}
	if _, err = Autenticar(lista, "ana@correo.com", "Mala1234"); !errors.Is(err, ErrCredenciales) {
		t.Errorf("clave mala: %v", err)
	}
	if _, err = Autenticar(lista, "noexiste@correo.com", "Clave1234"); !errors.Is(err, ErrCredenciales) {
		t.Errorf("correo inexistente: %v", err)
	}
}

func TestHashEsPuro(t *testing.T) {
	if Hashear("x", "s") != Hashear("x", "s") || Hashear("x", "s") == Hashear("x", "t") {
		t.Error("Hashear debe ser determinista y depender de la sal")
	}
}
