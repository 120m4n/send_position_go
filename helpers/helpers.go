package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"send_position/singleton"

	"net/http"
	//"time"

	"github.com/go-spatial/proj"
	"github.com/paulmach/orb"
)

func CreatePoints(input [][2]float64) []map[string]interface{} {
	output := make([]map[string]interface{}, 0)

	s := singleton.GetInstance()

	for _, pair := range input {
		lat := pair[1]
		lon := pair[0]

		obj := map[string]interface{}{
			"coordinates": map[string]float64{
				"latitude": lat,
				"longitude": lon,
			},
			"fleet":  s.Fleet,
			"user_id": s.Userid,
			"fleet_type": "camioneta",
		}
		output = append(output, obj)
	}
	
	return output	
}

// CreateListOfPoints Crea una lista que mapea las coordenadas a un objeto
func CreateListOfPoints(input [][][2]float64) []map[string]interface{} {
	output := make([]map[string]interface{}, 0)
	s := singleton.GetInstance()
	
	for _, fila := range input {
		for _, pair := range fila {
			lat := pair[1]
			lon := pair[0]

			obj := map[string]interface{}{
				"coordinates": map[string]float64{
					"latitude": lat,
					"longitude": lon,
				},
				"fleet":  s.Fleet,
				"user_id": s.Userid,
				"fleet_type": "camioneta",
			}
			output = append(output, obj)
		}
	}

	return output
}

func EnviarPOST(ctx context.Context, url string, obj map[string]interface{}, verbose bool) error {
	// Codificar el objeto en formato JSON
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Crear una nueva solicitud POST
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Realizar la solicitud POST
	client := &http.Client{}
	resp, err := client.Do(req)
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
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("respuesta del servidor: %s , %s", resp.Status, respData)
	}

	// Imprimir el resultado del servidor
	if verbose {
		fmt.Printf("Respuesta del servidor: %s\n", respData)
	}

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

//create a function that take a orb.linestring and return the numbers of segments
func GetSegmentCount(coordinates orb.LineString) int {
	return len(coordinates) - 1
}

// ExtraeSegmentos devuelve los segmentos de linea de cada tramo de un LineString
func ExtraeSegmentos(coordinates orb.LineString) [][]float64 {
	outputVector := make([][]float64, len(coordinates))
	for j, coor := range coordinates {
		var lonlat = []float64{coor.Lon(), coor.Lat()}
		xy, err := proj.Convert(proj.EPSG3395, lonlat)
		if err != nil {
			panic(err)
		}
		outputVector[j] = xy
	}
	return outputVector
}

// SubsampleVector toma una matriz de puntos de entrada y devuelve pares de puntos consecutivos
func SubsampleVector(input [][]float64) [][][2]float64 {
	// Verificar que haya al menos dos puntos en la entrada
	if len(input) < 2 {
		return nil
	}

	// Inicializar la matriz de salida
	output := make([][][2]float64, len(input)-1)

	// Crear pares de puntos consecutivos
	for i := 0; i < len(input)-1; i++ {
		pair := [][2]float64{
			{input[i][0], input[i][1]},
			{input[i+1][0], input[i+1][1]},
		}
		output[i] = pair
	}

	return output
}