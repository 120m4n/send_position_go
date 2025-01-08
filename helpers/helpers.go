package helpers

import (
	"bytes"
	// "context"
	"encoding/json"
	"fmt"

	// "io"
	"io/ioutil"
	"log"
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
				"latitude":  lat,
				"longitude": lon,
			},
			"fleet":      s.Fleet,
			"user_id":    s.Userid,
			"unique_id":  s.Uniqueid,
			"fleet_type": "camioneta",
			"avatar_ico": s.AvatarIcon,
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
					"latitude":  lat,
					"longitude": lon,
				},
				"fleet":      s.Fleet,
				"user_id":    s.Userid,
				"unique_id":  s.Uniqueid,
				"fleet_type": "camioneta",
				"avatar_ico": s.AvatarIcon,
			}
			output = append(output, obj)
		}
	}

	return output
}

func EnviarPOST(url string, obj map[string]interface{}, verbose bool) error {
	// Codificar el objeto en formato JSON
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Crear una nueva solicitud POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado de la respuesta
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("respuesta del servidor: %s", resp.Status)
	}

	// Imprimir el resultado del servidor
	if verbose {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf("Respuesta del servidor: %s\n", bodyString)
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

// create a function that take a orb.linestring and return the numbers of segments
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
