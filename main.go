package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-spatial/proj"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"

	"send_position/geometry"
)

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

func readFileIntoByteSlice(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func enviarPOST(url string, obj map[string]interface{}) error {
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

func main() {
	// Verificar si se proporciona un argumento (ruta al archivo geojson)
	if len(os.Args) <= 2 {
		fmt.Println("Uso: miaplicacion <ruta_al_archivo_geojson>")
		os.Exit(1)
	}

	// Obtener la ruta del archivo geojson del argumento de la línea de comandos
	filePath := os.Args[1]

	unique_id := "918422f729de0567"

	// URL del servidor donde realizar la solicitud POST
	url := os.Getenv("url_endpoint") + unique_id // Reemplaza esto con tu URL real

	// Configurar el tiempo de espera en segundos
	tiempoEspera := 2 // Cambia esto al tiempo de espera deseado en segundos


	// Leer el archivo geojson
	rawJSON, err := readFileIntoByteSlice(filePath)
	if err != nil {
		log.Fatalf("Error al leer el archivo geojson: %v", err)
	}

	// Parsear el archivo geojson
	feature, err := geojson.UnmarshalFeatureCollection(rawJSON)
	if err != nil {
		log.Fatalf("Error al parsear el archivo geojson: %v", err)
	}

	// Verificar si la feature es del tipo LineString
	if feature.Features[0].Geometry.GeoJSONType() != geojson.TypeLineString {
		fmt.Println("El archivo geojson no contiene un elemento LineString")
		os.Exit(1)
	}

	// temp := feature.Features[0].Geometry.GeoJSONType()
	// fmt.Println(string(temp))

	// Obtener las coordenadas del LineString
	coordinates := feature.Features[0].Geometry.(orb.LineString)
    inputVector := make([][]float64, len(coordinates))
	for j, coor := range coordinates {
		var lonlat = []float64{coor.Lon(), coor.Lat()}
		xy, err := proj.Convert(proj.EPSG3395, lonlat)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%.2f, %.2f\n", xy[0], xy[1])
		inputVector[j] = xy
	}

	output := SubsampleVector(inputVector)
	newpoints := geometry.GeneratePoints(output, 5)

	outpuLonLat := geometry.ConverToLonLat(newpoints)
	PrintMatrix(outpuLonLat)

	//geometry.ExampleUsage()

	// Crear una función para mapear las coordenadas a un objeto
	result := make([]map[string]interface{}, len(coordinates))
	for i, coord := range coordinates {
		lat := coord[1]
		lon := coord[0]

		obj := map[string]interface{}{
			"position": map[string]float64{
				"lat": lat,
				"lon": lon,
			},
			"fleet":  "camioneta",
			"userid": "G2012/roman",
		}

		result[i] = obj
	}

	// Convertir el resultado a JSON y mostrarlo
	// jsonData, err := json.Marshal(result)
	// if err != nil {
	//     log.Fatalf("Error al convertir el resultado a JSON: %v", err)
	// }

	// fmt.Println(string(jsonData))

	// Iterar sobre result y enviar cada elemento a la función enviarPOST
	for _, elem := range result {
		if err := enviarPOST(url, elem); err != nil {
			fmt.Printf("Error al enviar el elemento: %v\n", err)
		}
		
		// Esperar el tiempo configurado antes de la siguiente solicitud
		time.Sleep(time.Duration(tiempoEspera) * time.Second)
	}
	// Realizar la solicitud POST
	// if err := enviarPOST(url, obj); err != nil {
	// 	fmt.Printf("Error al enviar la solicitud POST: %v\n", err)
	// }
}
