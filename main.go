package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-spatial/proj"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"

	"send_position/geometry"
	"send_position/helpers"
)


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



func main() {
	// Verificar si se proporciona un argumento (ruta al archivo geojson)
	if len(os.Args) != 2 {
		fmt.Println("Uso: send_position <ruta_al_archivo_geojson>")
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
    // fmt.Printf("%.2f, %.2f\n", xy[0], xy[1])
    inputVector := ExtraeSegmentos(coordinates)


	output := SubsampleVector(inputVector)
	newpoints := geometry.GeneratePoints(output, 15)

	outpuLonLat := geometry.ConverToLonLat(newpoints)
	helpers.PrintMatrix(outpuLonLat)

	result := helpers.CreateListOfPoints(outpuLonLat)

	// Convertir el resultado a JSON y mostrarlo
	// jsonData, err := json.Marshal(result)
	// if err != nil {
	//     log.Fatalf("Error al convertir el resultado a JSON: %v", err)
	// }

	// fmt.Println(string(jsonData))

	// Iterar sobre result y enviar cada elemento a la función enviarPOST
	for _, elem := range result {
		if err := helpers.EnviarPOST(url, elem); err != nil {
			fmt.Printf("Error al enviar el elemento: %v\n", err)
		}
		
		// Esperar el tiempo configurado antes de la siguiente solicitud
		time.Sleep(time.Duration(tiempoEspera) * time.Second)
	}

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
