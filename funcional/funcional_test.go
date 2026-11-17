package funcional

import (
	"reflect"
	"testing"
)

func TestFiltrarMapearReducir(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	if got := Filtrar(nums, func(n int) bool { return n%2 == 0 }); !reflect.DeepEqual(got, []int{2, 4}) {
		t.Errorf("Filtrar = %v", got)
	}
	if got := Mapear(nums, func(n int) int { return n * n }); !reflect.DeepEqual(got, []int{1, 4, 9, 16, 25}) {
		t.Errorf("Mapear = %v", got)
	}
	if got := Reducir(nums, 0, func(a, n int) int { return a + n }); got != 15 {
		t.Errorf("Reducir = %d", got)
	}
}

func TestNoModificaEntrada(t *testing.T) {
	nums := []int{3, 1, 2}
	Mapear(nums, func(n int) int { return n * 10 })
	if !reflect.DeepEqual(nums, []int{3, 1, 2}) {
		t.Errorf("la entrada fue modificada: %v", nums)
	}
}
