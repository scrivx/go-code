package main

import (
	"errors"
	"fmt"
)

// Intefaz comun para todos los procesadores de pago
type PaymentProcessor interface {
	Process(amout float64) error
	GetFree() float64
}

// Procesador de tarjeta de credito
type CreditCardProcessor struct {
	CardNumber string
	FreeRate   float64
}

func (cc CreditCardProcessor) Process(amout float64) error {
	if amout <= 0 {
		return errors.New("Cantidad no valido")
	}
	fmt.Printf("Procesando  $%.2f con tarjeta ****%s\n", amout, cc.CardNumber[len(cc.CardNumber)-4:])
	return nil
}

func (cc CreditCardProcessor) GetFree() float64 {
	return cc.FreeRate
}

// Procesador Paypal
type PaypalProcessor struct {
	Email string
}

func (pp PaypalProcessor) Process(amout float64) error {
	if amout <= 0 {
		return errors.New("Cantidad no valido")
	}
	fmt.Printf("Procesando  $%.2f via Paypal (%s)\n", amout, pp.Email)
	return nil
}

func (pp PaypalProcessor) GetFree() float64 {
	return 0.029 // 2.9%
}

// Procesador de criptomonedas
type CrypoProcessor struct {
	WalletAdress string
	Currency     string
}

func (cp CrypoProcessor) Process(amout float64) error {
	if amout <= 0 {
		return errors.New("Cantidad no valido")
	}
	fmt.Printf("Procesando  $%.2f en %s (wallet %s....)\n", amout, cp.Currency, cp.WalletAdress[:10])
	return nil
}

func (cp CrypoProcessor) GetFree() float64 {
	return 0.01 // 1%
}

// Funcion Polimorfica que funciona con cualquier procesador
func ProcessOrder(processor PaymentProcessor, amout float64) error {
	free := amout * processor.GetFree()
	total := amout + free

	fmt.Printf("Procesando orden por $%.2f ( + $%.2f free) = $%.2f total\n", amout, free, total)
	return processor.Process(total)
}

// Funcion que elige el mejor procesador automaticamente
func ProcessWithBestRate(amout float64, processors []PaymentProcessor) error {
	if len(processors) == 0 {
		return errors.New("No hay procesadores")
	}

	// Encontrar el procesador con menor comision
	bestProcessor := processors[0]
	lowestFree := bestProcessor.GetFree()

	for _, processor := range processors[1:] {
		if processor.GetFree() < lowestFree {
			bestProcessor = processor
			lowestFree = processor.GetFree()
		}
	}
	fmt.Printf("Seleccionando procesador con %.1f%% de comision\n", lowestFree*100)
	return ProcessOrder(bestProcessor, amout)
}

func main() {
	// Creamos diferentes procesadores
	creditCard := CreditCardProcessor{
		CardNumber: "1234-5678-9012-3456",
		FreeRate:   0.035, // 3.5%
	}

	paypal := PaypalProcessor{
		Email: "criv@gmail.com",
	}

	crypto := CrypoProcessor{
		WalletAdress: "0x1234567890123456789012345678901234567890",
		Currency:     "BTC",
	}

	// Polimorfismo en accion
	processors := []PaymentProcessor{creditCard, paypal, crypto}
	fmt.Println("===== Procesando el mejor rate =====")
	ProcessWithBestRate(100.0, processors)

	fmt.Println("\n===== Procesando con cada uno =====")
	for _, processor := range processors {
		ProcessOrder(processor, 50.0)
		fmt.Println()
	}
}
