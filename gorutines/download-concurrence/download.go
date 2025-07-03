package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func downloadURL(url string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done() // defer -> sirve para marcar gorutine como terminada
	}
	inicio := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error descarga %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close() // Cerrar respuesta

	// Leer todo el contenido(para simular descarga completa)
	_, err = io.ReadAll(resp.Body)
	
	if err != nil{
		fmt.Printf("Error de lectura %s: %v\n", url, err)
		return
	}

	duracion := time.Since(inicio)
	fmt.Printf("Descargando %s en %v (Estado: %s)\n", url, duracion, resp.Status)
}

func main() {
	urls := []string{
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/2",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/3",
	}
	fmt.Println("=== Descarga SECUENCIAL ===")
	
	inicio := time.Now()

	for _, url := range urls {
		downloadURL(url, nil)
	}

	tiempoSecuencial := time.Since(inicio)
	fmt.Printf("Tiempo secuencial total : %v\n\n", tiempoSecuencial)

	fmt.Println("=== Descarga CONCURRENTE ===")

	inicio = time.Now()
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go downloadURL(url, &wg)
	}

	wg.Wait()

	tiempoConcurrencia := time.Since(inicio)
	fmt.Printf("Tiempo concurrencia total : %v\n\n", tiempoConcurrencia)
	
	fmt.Printf("Mejora: %.2fx mas rapido\n", float64(tiempoSecuencial)/float64(tiempoConcurrencia))
}