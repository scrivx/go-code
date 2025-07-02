package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	// IMPORTANTE: Marcar oomo terminado al salir
	defer wg.Done()
	fmt.Printf("Trabajo %d iniciando\n", id)

	// Simular trabajo
	time.Sleep(time.Duration(id) * time.Second)
	fmt.Printf("Trabajo %d finalizando\n", id)
}

func main(){
	var wg sync.WaitGroup
	numbWorkers := 3

	for i := 1; i <= numbWorkers; i++ {
		wg.Add(1) // Incrementar contador
		go worker(i, &wg) // Lanzar gorutine
	}

	fmt.Println("Esperando a que todos los trabajos terminen...")
	wg.Wait() // Esperar que todas terminen
	fmt.Println("Todos los trabajos terminados")
}