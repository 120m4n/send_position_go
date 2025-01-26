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
	"math"
	"math/rand"

	"github.com/go-spatial/proj"
	"github.com/paulmach/orb"
)


type Point struct {
	X float64
	Y float64
}

func RandomPointInCircle(centerLon float64, centerLat float64, radiusInMeters float64) Point {
    // Convert radius from meters to degrees
    radiusInDegrees := radiusInMeters / 111139.0

    // Generate a random angle
    angle := rand.Float64() * 2 * math.Pi

    // Generate a random radius within the circle's radius
    r := radiusInDegrees * math.Sqrt(rand.Float64())

    // Calculate the coordinates of the random point
    deltaLon := r * math.Cos(angle)
    deltaLat := r * math.Sin(angle)

    // Add deltas to the center coordinates
    lon := centerLon + deltaLon
    lat := centerLat + deltaLat

    return Point{lon, lat}
}

//given a pair lon, lat coordinates , return a map[string]interface{} object
// in the circle around point at radius 100m

func CreateRandomPoint(longitude float64, latitude float64 ) map[string]interface{} {
	s := singleton.GetInstance()
	/// define temp like a cicle with center longitude, latitude pair and radious 100 meters
	temp := RandomPointInCircle(longitude, latitude, 100)


	obj := map[string]interface{}{
		"coordinates": map[string]float64{
			"latitude": temp.Y,
			"longitude": temp.X,
		},
		"fleet":  s.Fleet,
		"user_id": s.Userid,
		"unique_id": s.Uniqueid,
		"fleet_type": "car",
	}

	return obj

}

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
			"unique_id": s.Uniqueid,
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
				"unique_id": s.Uniqueid,
				"fleet_type": "camioneta",
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