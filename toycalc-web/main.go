// web_server/main.go
package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	toycalc_core "github.com/vladimirck/toycalc/toycalc-core" // <-- Importa tu librería ToyCalc
)

// Estructura para pasar datos a las plantillas HTML
type CalculatorData struct {
	Expression string
	Result     string
	Error      string
}

var tmpl *template.Template // Para cargar las plantillas una sola vez

func init() {
	// Carga las plantillas HTML una sola vez al inicio
	// Ajusta la ruta a tus archivos de plantilla
	var err error
	// Asegúrate de que estas rutas sean relativas al lugar donde ejecutarás 'go run' o el binario
	tmpl, err = template.ParseFiles(
		filepath.Join("templates", "index.html"),
		filepath.Join("templates", "display.html"),
	)
	if err != nil {
		log.Fatalf("Error cargando plantillas: %v", err)
	}
}

func main() {
	// Ruta principal para servir el HTML y manejar las peticiones POST de la calculadora
	http.HandleFunc("/", calculatorHandler)

	// Sirve archivos estáticos (como el CSS de Tailwind si lo compilas)
	// Ajusta la ruta si tu 'static' está en otro lugar
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join("web_server", "static")))))

	log.Println("Servidor Go escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) // Puerto donde Caddy hará reverse_proxy
}

// calculatorHandler maneja la lógica de la calculadora web
func calculatorHandler(w http.ResponseWriter, r *http.Request) {
	data := CalculatorData{Expression: "0", Result: "", Error: ""}

	if r.Method == "POST" {
		input := r.FormValue("input")
		currentExpression := r.FormValue("current_expression")

		var newExpression string
		if input == "C" {
			newExpression = "0"
		} else if input == "=" {
			// Cuando se presiona '=', usa la expresión actual para calcular
			resultStr, err := toycalc_core.CalculateExpression(currentExpression) // <--- ¡Usa tu librería!
			if err != nil {
				data.Error = err.Error()
				data.Expression = currentExpression // Mantén la expresión visible en caso de error
			} else {
				data.Result = resultStr
				data.Expression = resultStr // Muestra el resultado como la nueva expresión
			}
		} else {
			// Construye la expresión para mostrarla en pantalla
			if currentExpression == "0" && input != "." {
				newExpression = input
			} else if currentExpression == "" { // Para el primer input después de un reset o carga
				newExpression = input
			} else {
				newExpression = currentExpression + input
			}
			data.Expression = newExpression
		}

		// HTMX espera HTML. Renderiza solo la sección del display o todo el body.
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.ExecuteTemplate(w, "display.html", data) // Renderiza solo la plantilla del display
		if err != nil {
			log.Printf("Error ejecutando plantilla display.html: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}

		// Importante: Actualizar el valor del campo oculto current_expression para la siguiente petición
		// HTMX by default sends all named inputs from the parent form/element.
		// If you only update the display, the hidden input's value won't change unless you target it too.
		// A common pattern is to re-render a small form or container that includes the hidden input.
		// For now, the `hx-swap="outerHTML"` on the #calculator-display div means its *parent container*
		// needs to include the hidden input if it's outside. Let's adjust the template to include the hidden input.
		// (See updated template below)

	} else { // GET request (carga inicial de la página)
		data.Expression = "0"
		err := tmpl.ExecuteTemplate(w, "index.html", data) // Renderiza la plantilla principal
		if err != nil {
			log.Printf("Error ejecutando plantilla index.html: %v", err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			return
		}
	}
}