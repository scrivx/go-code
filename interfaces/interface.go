package interfaces

import (
	"encoding/json"
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

// Rastreador permite hacer seguimiento de notificaciones
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
func (e *EmailNotificador) EnviarNotificacion(destinatario, mensaje string) error {
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
func (e *EmailNotificador) ValidarMensaje(mensaje string) error {
	if len(mensaje) == 0 {
		return errors.New("mensaje no puede estar vacío")
	}
	if len(mensaje) > 1000 {
		return errors.New("mensaje demasiado largo")
	}
	return nil
}

func (e *EmailNotificador) ValidarDestinantario(destinatario string) error {
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
	stats := map[string]int{
		"total":      0,
		"enviados":   0,
		"fallidos":   0,
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
	apiKey    string
	proveedor string
	registros map[string]*RegistroNotificacion
}

func NuevoSMSNotificador(apiKey, proveedor string) *SMSNotificador {
	return &SMSNotificador{
		apiKey:    apiKey,
		proveedor: proveedor,
		registros: make(map[string]*RegistroNotificacion),
	}
}

// Implementa Notificador
func (s *SMSNotificador) EnviarNotificacion(destinatario, mensaje string) error {
	if err := s.ValidarDestinantario(destinatario); err != nil {
		return err
	}
	if err := s.ValidarMensaje(mensaje); err != nil {
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
	if time.Now().UnixNano()%20 == 0 {
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
		return errors.New("mensaje demasiado largo")
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
	stats := map[string]int{
		"total":      0,
		"enviados":   0,
		"fallidos":   0,
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

type SlackNotificador struct {
	webhook string
	canal   string
}

func NuevoSlackNotificador(webhook, canal string) *SlackNotificador {
	return &SlackNotificador{
		webhook: webhook,
		canal:   canal,
	}
}

// Solo implementa Notificador
func (sl *SlackNotificador) EnviarNotificacion(destinatario, mensaje string) error {
	fmt.Printf("🔔Slack -> Canal: %s | Usuario: %s | Mensaje: %s\n", sl.canal, destinatario, mensaje)

	//Simular envio instantaneo
	time.Sleep(10 * time.Millisecond)
	return nil
}

// ==========================================
// PASO 5: SERVICIO PRINCIPAL
// ==========================================

type ServicioNotificaciones struct {
	notificadores []Notificador
	logger        Logger
}

func NuevoServicioNotificaciones() *ServicioNotificaciones {
	return &ServicioNotificaciones{
		notificadores: make([]Notificador, 0),
	}
}

func (sn *ServicioNotificaciones) AgregarNotificador(notificador Notificador) {
	sn.notificadores = append(sn.notificadores, notificador)
	if sn.logger != nil {
		sn.logger.LogInfo(fmt.Sprintf("Notificador agregado: %T", notificador))
	}
}

func (sn *ServicioNotificaciones) EnviarATodos(destinatario, mensaje string) map[string]error {
	resultados := make(map[string]error)

	if sn.logger != nil {
		sn.logger.LogInfo(fmt.Sprintf("Enviando a %d notificadores", len(sn.notificadores)))
	}

	for _, notificador := range sn.notificadores {
		tipoNotificador := fmt.Sprintf("%T", notificador)
		err := notificador.EnviarNotificacion(destinatario, mensaje)
		resultados[tipoNotificador] = err

		if sn.logger != nil {
			if err != nil {
				sn.logger.LogError(fmt.Errorf("%s fallo: %v", tipoNotificador, err))
			} else {
				sn.logger.LogInfo(fmt.Sprintf("%s exitoso", tipoNotificador))
			}
		}

	}
	return resultados
}

// Enviar solo a notificadores que implementan ValidadorMensaje

func (sn *ServicioNotificaciones) EnviarConvalidaciones(destinatario, mensaje string) map[string]error {
	resultados := make(map[string]error)

	for _, notificador := range sn.notificadores {
		tipoNotificador := fmt.Sprintf("%T", notificador)

		if validador, implementa := notificador.(ValidadorMensaje); implementa {
			// validar antes de enviar
			if err := validador.ValidarMensaje(mensaje); err != nil {
				resultados[tipoNotificador] = fmt.Errorf("Validacion fallo : %v", err)
				continue
			}
			if err := validador.ValidarDestinantario(destinatario); err != nil {
				resultados[tipoNotificador] = fmt.Errorf("Destinatario invalido : %v", err)
				continue
			}
		}
		//Enviar Notificacion
		err := notificador.EnviarNotificacion(destinatario, mensaje)
		resultados[tipoNotificador] = err
	}
	return resultados
}

// ==========================================
// PASO 6: FUNCIONES DE UTILIDAD
// ==========================================

// Función que acepta cualquier Notificador
func ProbarNotificador(n Notificador, destinatario, mensaje string) {
	fmt.Printf("\n🧪Probando %T:\n", n)
	fmt.Println("   Enviando:", mensaje)

	err := n.EnviarNotificacion(destinatario, mensaje)

	if err != nil {
		fmt.Printf("❌ Error : %v\n", err)
	} else {
		fmt.Printf("✅Enviado correctamente\n")
	}
}

// Funcion que verifica capacidades usando type assertions
func AnalizarCapacidadesNotificador(n Notificador) {
	fmt.Printf("\n🔍Analizando capacidades de %T:\n", n)
	capacidades := []string{}

	// verificar cada interface
	if _, implementa := n.(Notificador); implementa {
		capacidades = append(capacidades, "✅Notificador (envio basico)")
	}

	if _, implementa := n.(ValidadorMensaje); implementa {
		capacidades = append(capacidades, "✅ValidadorMensaje (validacion)")
	}

	if _, implementa := n.(Rastreador); implementa {
		capacidades = append(capacidades, "✅Rastreador (seguimiento))")
	}

	if _, implementa := n.(Logger); implementa {
		capacidades = append(capacidades, "✅Logger (registro)")
	}

	if _, implementa := n.(NotificadorCompleto); implementa {
		capacidades = append(capacidades, "✅NotificadorCompleto")
	}

	for capacidad := range capacidades {
		fmt.Printf("   %s\n", capacidad)
	}
}

// Type Switch para manejar diferentes tipos
func ProcesarNotificacionPorTipo(n Notificador, destinatario, mensaje string) {
	switch notificador := n.(type) {
	case *EmailNotificador:
		fmt.Println("📧 Procesando como EmailNotificador...")
		fmt.Printf("   notificador.puerto) Servidor: %s:%d\n", notificador.servidor, notificador.puerto)
		notificador.EnviarNotificacion(destinatario, mensaje)

	case *SMSNotificador:
		fmt.Println("📱 Procesando como SMSNotificador...")
		fmt.Printf("Proveedor: %s\n", notificador.proveedor)
		notificador.EnviarNotificacion(destinatario, mensaje)

	case *SlackNotificador:
		fmt.Println("💬Procesando como SlackNotificador...")
		fmt.Printf("  Canal %s\n", notificador.canal)
		notificador.EnviarNotificacion(destinatario, mensaje)

	default:
		fmt.Printf("❓Tipo desconocido: %T\n", notificador)
		notificador.EnviarNotificacion(destinatario, mensaje)
	}
}

// ==========================================
// PASO 7: FUNCIÓN PRINCIPAL DEMOSTRATIVA
// ==========================================

func main() {
	fmt.Println("🔔 SISTEMA DE NOTIFICACIONES - INTERFACES EN ACCIÓN")
	fmt.Println("=" + strings.Repeat("=", 60))

	// Crear diferentes notificadores
	email := NuevoEmailNotificador("smtp.gmail.com", 587, "app@empresa.com", "password")
	sms := NuevoSMSNotificador("api-key-123", "Twilio")
	slack := NuevoSlackNotificador("https://hooks.slack.com/...", "#general")

	// Crear servicio principal
	servicio := NuevoServicioNotificaciones()
	servicio.EstablecerLogger(email) // email tambien funciona como logger

	// Agregar notificadores
	servicio.AgregarNotificador(email)
	servicio.AgregarNotificador(sms)
	servicio.AgregarNotificador(slack)

	fmt.Println("\n 📋 1. POLIMORFISMO BÁSICO:")
	fmt.Println(strings.Repeat("-", 40))

	// Todos son tratados como notificadores
	notificadores := []Notificador{email, sms, slack}

	for _, n := range notificadores {
		ProbarNotificador(n, "usuario@gmail.com", "Hola, desde GO!")
	}

	fmt.Println("\n 📋2. TYPE ASSERTIONS Y CAPACIDADES:")
	fmt.Println(strings.Repeat("-", 40))

	// Analizar capacidades de cada notificador
	for _, n := range notificadores {
		AnalizarCapacidadesNotificador(n)
	}

	fmt.Println("\n 📋 3. TYPE SWITCH EN ACCIÓN:")
	fmt.Println(strings.Repeat("-", 40))

	// Usar type switch para la logica específica por tipo
	for _, n := range notificadores {
		ProcesarNotificacionPorTipo(n, "+51914909703", "Mensaje tipo especifico")
		fmt.Println()
	}

	fmt.Println("\n 📋 4. INTERFACES COMPUESTAS:")
	fmt.Println(strings.Repeat("-", 40))

	verificarInterfaceCompuesta := func(n Notificador) {
		nombre := fmt.Sprintf("%T", n)

		if completo, esCompleto := n.(NotificadorCompleto); esCompleto {
			fmt.Println("✅ %s implementa NotificadorCompleto\n", nombre)
			// Puede usar todas las funciones de NotificadorCompleto
			completo.ValidarMensaje("test")
			completo.EnviarNotificacion("test@test.com", "test")
		} else {
			fmt.Println("❌ %s no implementa NotificadorCompleto\n", nombre)
		}

		if avanzado, esAvanzado := n.(NotificadorAvanzado); esAvanzado {
			fmt.Println("✅ %s implementa NotificadorAvanzado\n", nombre)
			stats := avanzado.ObtenerEstadisticas()
			fmt.Printf(" Estadísticas: %v\n", stats)

		} else {
			fmt.Println("❌ %s no implementa NotificadorAvanzado\n", nombre)
		}
		fmt.Println()
	}

	for _, n := range notificadores {
		verificarInterfaceCompuesta(n)
	}

	fmt.Println("\n 📋 5. SERVICIO CON MULTIPLES NOTIFICADORES:")
	fmt.Println(strings.Repeat("-", 40))

	// Enviar a todos
	fmt.Println("📤 Enviando a todos...")
	resultados := servicio.EnviarATodos("crivera@gmail.com", "Sistema iniciado correctamente")

	for tipo, err := range resultados {
		if err != nil {
			fmt.Printf("❌ %s: %v\n", tipo, err)
		} else {
			fmt.Printf("✅ %s: Exito\n", tipo)
		}
	}

	fmt.Println("\n 📋 6. ESTADISTICAS Y RASTREABILIDAD:")
	fmt.Println(strings.Repeat("-", 40))

	// Mostrar estadisticas solo notificadores que implementan Rastreador
	for _, n := range notificadores {
		if rastreador, implementa := n.(Rastreador); implementa {
			nombre := fmt.Sprintf("%T", n)
			stats := rastreador.ObtenerEstadisticas()
			fmt.Printf("✅ Estadísticas de %s: %v\n", nombre)

			stastJSON, _ := json.MarshalIndent(stats, "", "  ")
			fmt.Println("   %s\n\n", string(stastJSON))
		}
	}

	fmt.Println(" 🎯 CONCEPTOS DEMOSTRADOS:")
	fmt.Println(strings.Repeat("-", 40))

	conceptos := []string{
		" ✅ Definición de interfaces simples y compuestas",
		" ✅ Implementación implícita de interfaces",
		" ✅ Polimorfismo con múltiples implementaciones",
		" ✅ Type assertions para verificar capacidades",
		" ✅ Type switches para lógica específica por tipo",
		" ✅Composición de interfaces",
		" ✅ Interfaces como contratos flexibles",
		" ✅ Uso práctico en arquitectura de servicios",
	}

	for _, concepto := range conceptos {
		fmt.Printf("   %s\n", concepto)
	}

	fmt.Printf("\n 🎊 ¡Ejemplo completado!")

}
