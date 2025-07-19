package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	toycalc_core "github.com/vladimirck/toycalc/toycalc-core"
)

// PageData contiene los datos que se pasarán a la plantilla HTML.
type PageData struct {
	Expression        string
	Result            string
	GoogleAnalyticsID string
}

func main() {
	// Define el manejador para la ruta principal.
	http.HandleFunc("/", handleCalculator)

	fmt.Println("Servidor escuchando en http://localhost:8080")
	// Inicia el servidor en el puerto 8080.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleCalculator se encarga de las peticiones a la página.
func handleCalculator(w http.ResponseWriter, r *http.Request) {
	// Parsea la plantilla HTML. Es importante manejar el error.
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error al cargar la página", http.StatusInternalServerError)
		log.Println("Error parsing template:", err)
		return
	}

	// Obtiene la expresión del formulario enviado (parámetro GET 'expression').
	expression := r.URL.Query().Get("expression")
	gaID := os.Getenv("GOOGLE_ANALYTICS_ID")
	
	data := PageData{
		Expression:        expression,
		GoogleAnalyticsID: gaID,
	}

	// Si hay una expresión, la calcula.
	if expression != "" {
		result, err := toycalc_core.CalculateExpression(expression)
		if err != nil {
			// Si hay un error en el cálculo, lo muestra como resultado.
			data.Result = "Error: " + err.Error()
		} else {
			// Si el cálculo es exitoso, muestra el resultado.
			data.Result = result
		}
	}

	// Ejecuta la plantilla, pasándole los datos (expresión y resultado).
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error al renderizar la página", http.StatusInternalServerError)
		log.Println("Error executing template:", err)
	}
}
