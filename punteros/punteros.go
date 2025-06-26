package punteros

// Ejemplo: sencillo del uso de un puntero

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	
	// inicializamos la estructura
	p := Person{ 
		Name: "Carlos", 
		Age: 24,
	}
	fmt.Println("Antes de usar puntero:", p)
	
	// usamos el puntero
	ConPuntero(&p) // "&" sirve para indicar que se trata de un puntero
	fmt.Println("Despues de usar puntero:", p)

}

// creamos una función que recibe un puntero
func ConPuntero(ptr *Person) { // "*" indica que se trata de un puntero
 	ptr.Name = "criv"
	ptr.Age = 25
}