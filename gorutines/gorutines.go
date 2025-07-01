// La palabra reservada para usar una
// Goroutine es "go"

package main

import (
	"fmt"
	"sync"
	"time"
)

func printNumbers() {
	for i := 1; i <= 5; i++ {
		fmt.Println(i)
	}
}

// Sincronizacion de una Goroutine
 // - funcion normal
 func sayHello(name string) {
	for i := 1; i < 3 ; i ++ {
		fmt.Printf("Hola desde %s (iteracion %d)\n", name, i+1)
		time.Sleep(500 * time.Millisecond)
	}
 }

// WaitGroups <-- Forma correcta de sincronizacion
func printNumber(wg *sync.WaitGroup){
	defer wg.Done()
	for i := 1; i <= 5; i++ {
		fmt.Println(i)
	}
}

func main() {
	go printNumbers() // Creamos y ejecutamos una gorutine
	time.Sleep(time.Second) // Esperamos un segundo
	fmt.Println("Main function existing....")


	fmt.Println("=== Ejecucion SECUENCIAL ====")
	sayHello("Carlos")
	sayHello("Rivera")

	fmt.Println("=== Ejecucion CONCURRENTE ====")
	// Lanzar goroutines con la palabra clave "go"
	go sayHello("Charlie")
	go sayHello("Alice")
	 // Sin esto el programa terminaria antes que las gourotines
	 time.Sleep(2 * time.Second)
	 fmt.Println("Main function finished")

	 // WaitGroup <-- 
	 var wg sync.WaitGroup
	 wg.Add(1) // Incrementar el numero de gorutines que se van a esperar
	 go printNumber(&wg)

	 wg.Wait() // Esperar a que todas las gorutines terminen
	 fmt.Println("Main function existing...")
}

