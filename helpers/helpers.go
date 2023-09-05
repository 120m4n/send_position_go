package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//
// CreateListOfPoints Crea una lista que mapea las coordenadas a un objeto
func CreateListOfPoints(input [][][2]float64) []map[string]interface{} {
	output := make([]map[string]interface{}, 0)

	for _, fila := range input {
		for _, pair := range fila {
			lat := pair[1]
			lon := pair[0]

			obj := map[string]interface{}{
				"position": map[string]float64{
					"lat": lat,
					"lon": lon,
				},
				"fleet":  "camioneta",
				"userid": "G2012/roman",
			}
			output = append(output, obj)
		}
	}

	return output
}

// enviarPOST toma elementos de un vector y los envia a un endpoint
func EnviarPOST(url string, obj map[string]interface{}) error {
	// Codificar el objeto en formato JSON
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Realizar la solicitud POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Leer la respuesta del servidor
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Verificar el código de estado de la respuesta
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("respuesta del servidor: %s", resp.Status)
	}

	// Imprimir el resultado del servidor
	fmt.Printf("Respuesta del servidor: %s\n", respData)

	return nil
}

func PrintMatrix(matrix [][][2]float64) {
	// Imprimir la cabecera de la matriz
	fmt.Println("[")
	// Imprimir los elementos de la matriz
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ {
			// Formatear el número con dos decimales
			fmt.Printf("[ %.6f, %.6f ]", matrix[i][j][0], matrix[i][j][1])
		}
		fmt.Println()
	}

	// Imprimir el cierre de la matriz
	fmt.Println("]")
}