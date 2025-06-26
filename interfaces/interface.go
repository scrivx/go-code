package interfaces

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ==========================================
// PASO 1: INTERFACES BÁSICAS
// ==========================================

// Notificador define la funcionalidad básica de envío
type Notificador interface {
	EnviarNotificacion(destinatario, mensaje string) error
}

// ValidadorMensaje valida contenido antes del envío
type ValidadorMensaje interface {
	ValidarMensaje(mensaje string) error
	ValidarDestinantario(destinatario string) error
}

//Rastreador permite hacer seguimiento de notificaciones
type Rastreador interface {
	ObtenerEstado(id string) (string, error)
	ObtenerEstadisticas() map[string]int
}

// Logger registra eventos del sistema
type Logger interface {
	Log(nivel, mensaje string)
	LogError(error)
	LogInfo(string)
}

// ==========================================
// PASO 2: INTERFACES COMPUESTAS
// ==========================================

// NotificadorCompleto combina funcionalidades básicas
type NotificadorCompleto interface {
	Notificador
	ValidadorMensaje
}

// Notificador Avanzado incluye todas las funcionalidades
type NotificadorAvanzado interface {
	Notificador
	ValidadorMensaje
	Rastreador
	Logger
}

// ==========================================
// PASO 3: STRUCTS Y TIPOS DE DATOS
// ==========================================

type TipoNotificacion string

const (
	Email TipoNotificacion = "email"
	SMS   TipoNotificacion = "sms"
	Push  TipoNotificacion = "push"
	Slack TipoNotificacion = "slack"
)

type EstadoNotificacion string

const (
	Pendiente EstadoNotificacion = "pendiente"
	Enviado   EstadoNotificacion = "enviado"
	Fallida   EstadoNotificacion = "fallida"
	Entregada EstadoNotificacion = "entregada"
)

type RegistroNotificacion struct {
	ID           string
	Tipo         TipoNotificacion
	Destinatario string
	Mensaje      string
	Estado       EstadoNotificacion
	Timestamp    time.Time
	Intentos     int
	Error        error
}

type ConfiguracionNotificacion struct {
	MaxIntentos     int
	TimeoutSegundos int
	ReintentoAuto   bool
}

// ==========================================
// PASO 4: IMPLEMENTACIONES CONCRETAS
// ==========================================

// EmailNotificador - Implementa múltiples interfaces
type EmailNotificador struct {
	servidor      string
	puerto        int
	usuario       string
	password      string
	configuracion ConfiguracionNotificacion
	registros     map[string]*RegistroNotificacion
}

// Constructor para EmailNotificador
func NuevoEmailNotificador(servidor string, puerto int, usuario string, password string, configuracion ConfiguracionNotificacion) *EmailNotificador {
	return &EmailNotificador{
		servidor: servidor,
		puerto:   puerto,
		usuario:  usuario,
		password: password,
		configuracion: ConfiguracionNotificacion{
			MaxIntentos:     3,
			TimeoutSegundos: 30,
			ReintentoAuto:   true,
		},
		registros: make(map[string]*RegistroNotificacion),
	}
}


// Implementa Notificador
func(e *EmailNotificador) EnviarNotificacion(destinatario, mensaje string) error {
	// Validar antes de enviar
	if err := e.ValidarDestinantario(destinatario); err != nil {
		return err
	}
	if err := e.ValidarMensaje(mensaje); err != nil {
		return err
	}

	// Crear registro
	id := fmt.Sprintf("email_%d", time.Now().UnixNano())
	registro := &RegistroNotificacion{
		ID:           id,
		Tipo:         Email,
		Destinatario: destinatario,
		Mensaje:      mensaje,
		Estado:       Pendiente,
		Timestamp:    time.Now(),
		Intentos:     1,
	}
	e.registros[id] = registro
	// Simular envio de email
	e.LogInfo(fmt.Sprintf("Enviando email a %s", destinatario))
	time.Sleep(100 * time.Millisecond) // Simular latencia

	//Simular exito/fallo (90% de exito)
	if time.Now().UnixNano()%10 == 0 {
		registro.Estado = Fallida
		registro.Error = errors.New("Servidor SMTP no disponible")
		e.LogError(registro.Error)
		return errors.New("fallo al enviar email")
	} 

	registro.Estado = Enviado
	e.LogInfo(fmt.Sprintf("Email exitosamente enviado a %s", id))
	return nil
}

// Implementa ValidadorMensaje
func(e *EmailNotificador) ValidarMensaje(mensaje string) error {
	if len(mensaje) == 0 {
		return errors.New("mensaje no puede estar vacío")
	}
	if len(mensaje) > 1000 {
		return  errors.New("mensaje demasiado largo")
	}
	return nil
}

func(e *EmailNotificador) ValidarDestinantario(destinatario string) error {
	if !strings.Contains(destinatario, "@") {
		return errors.New("Email invalido: debe contener @")
	}
	if !strings.Contains(destinatario, ".") {
		return errors.New("Email invalido: debe contener .")
	}
	return nil
}

// Implementa Rastreador
func (e *EmailNotificador) ObtenerEstado(id string) (string, error) {
	if registro, existe := e.registros[id]; existe {
		return string(registro.Estado), nil
	}
	return "", errors.New("Notificacion no encontrada")
}

func (e *EmailNotificador) ObtenerEstadisticas() map[string]int {
	stats := map[string]int {
		"total": 0,
		"enviados": 0,
		"fallidos": 0,
		"pendientes": 0,
	}
	for _, registro := range e.registros {
		stats["total"]++
		switch registro.Estado {
		case Enviado:
			stats["enviados"]++
		case Fallida:
			stats["fallidos"]++
		case Pendiente:
			stats["pendientes"]++
		}
	}
	return stats
}

// Implementa Logger
func (e *EmailNotificador) Log(nivel, mensaje string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] EMAIL [%s]: %s\n", timestamp, nivel, mensaje)
}

func (e *EmailNotificador) LogError(err error) { 
	e.Log("ERROR", err.Error()) 
} 
func (e *EmailNotificador) LogInfo(mensaje string) { 
	e.Log("INFO", mensaje) 
}


// ==========================================

// SMSNotificador - Otra implementacion
type SMSNotificador struct {
	apiKey string
	proveedor string
	registros map[string]*RegistroNotificacion
}

func NuevoSMSNotificador(apiKey, proveedor string) *SMSNotificador{
	return &SMSNotificador{
		apiKey: apiKey,
		proveedor: proveedor,
		registros: make(map[string]*RegistroNotificacion),
	}
}

// Implementa Notificador
func (s *SMSNotificador) EnviarNotificacion(destinatario, mensaje string) error {
	if err := s.ValidarDestinantario(destinatario); err != nil {
		return err
	}
	if err := s.ValidarMensaje(mensaje); err!= nil {
		return err
	}

	id := fmt.Sprintf("sms_%d", time.Now().UnixNano())
	registro := &RegistroNotificacion{
		ID:           id,
		Tipo:         SMS,
		Destinatario: destinatario,
		Mensaje:      mensaje,
		Estado:       Pendiente,
		Timestamp:    time.Now(),
		Intentos:     1,
	}
	s.registros[id] = registro
	s.LogInfo(fmt.Sprintf("Enviando SMS a %s via %s", destinatario, s.proveedor))
	time.Sleep(50 * time.Millisecond)

	// SMS mas confiable (95% de exito)
	if time.Now().UnixNano()%20 == 0{
		registro.Estado = Fallida
		registro.Error = errors.New("Numero no valido")
		s.LogError(registro.Error)
		return errors.New("fallo al enviar SMS")
	}

	registro.Estado = Enviado
	s.LogInfo(fmt.Sprintf("SMS exitosamente enviado a %s", id))
	return nil
}

// Implementa ValidadorMensaje
func (s *SMSNotificador) ValidarMensaje(mensaje string) error {
	if len(mensaje) == 0 {
		return errors.New("mensaje no puede estar vacío")
	}
	if len(mensaje) > 160 {
		return  errors.New("mensaje demasiado largo")
	}
	return nil
}

func (s *SMSNotificador) ValidarDestinantario(destinatario string) error {
	if len(destinatario) < 10 {
		return errors.New("numero de telefono muy corto")
	}
	if !strings.HasPrefix(destinatario, "+") && !strings.HasPrefix(destinatario, "0") {
		return errors.New("numero debe empezar con + o 0")
	}
	return nil
}

// Implementa Rastreador
func (s *SMSNotificador) ObtenerEstado(id string) (string, error) {
	if registro, existe := s.registros[id]; existe {
		return string(registro.Estado), nil
	}
	return "", errors.New("SMS no encontrado")
}

func (s *SMSNotificador) ObtenerEstadisticas() map[string]int {
	stats := map[string]int {
		"total": 0,
		"enviados": 0,
		"fallidos": 0,
		"pendientes": 0,
	}
	for _, registro := range s.registros {
		stats["total"]++
		switch registro.Estado {
		case Enviado:
			stats["enviados"]++
		case Fallida:
			stats["fallidos"]++
		case Pendiente:
			stats["pendientes"]++
		}
	}
	return stats
}

// Implementa Logger 
func (s *SMSNotificador) Log(nivel, mensaje string) { 
	timestamp := time.Now().Format("2006-01-02 15:04:05") 
	fmt.Printf("[%s] SMS [%s]: %s\n", timestamp, nivel, mensaje) 
} 
func (s *SMSNotificador) LogError(err error) { 
	s.Log("ERROR", err.Error()) 
}

func (s *SMSNotificador) LogInfo(mensaje string) { 
	s.Log("INFO", mensaje) 
}

// ========================================== 
// SlackNotificador - Implementación más simple 

func main() {

}
