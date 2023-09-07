package main

import (
	"flag"
	"fmt"
	"log"
	//"net/url"
	"os"
	//"path/filepath"
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
	geojsonPtr := flag.String("json", "", "archivo geojson para procesar")
	ciclicoPtr := flag.Bool("ciclico", false, "procesar el geojson de forma ciclica")
    puntosPtr := flag.Int("puntos", 3, "numero de puntos por segmento de linea")
	parametroPtr := flag.String("parametro", "", "parametro de la ruta http")
	timePtr := flag.Int("tiempo", 2, "tiempo en segundos entre request sucesivos")
    urlPtr  := flag.String("url", "", "url endpoint para envio de datos")
	flag.Parse()
	
	if (*geojsonPtr == "") {
		fmt.Println("Uso: send_position -json=<ruta_al_archivo_geojson>")
		os.Exit(1)
	} 

	// URL del servidor donde realizar la solicitud POST
	if *urlPtr == "" {
		fmt.Println("Uso: send_position -json=<ruta_al_archivo_geojson> -url=http://localhost:3000")
		os.Exit(1)
	}
	
	if *parametroPtr != "" {  
       *urlPtr = *urlPtr + *parametroPtr
	} 

	// Leer el archivo geojson
	rawJSON, err := readFileIntoByteSlice(*geojsonPtr)
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
	newpoints := geometry.GeneratePoints(output, *puntosPtr)

	outpuLonLat := geometry.ConverToLonLat(newpoints)
	//helpers.PrintMatrix(outpuLonLat)

	if *ciclicoPtr {
	
		iterator := helpers.NewCircularIterator(outpuLonLat)
		// Ejemplo de uso en un ciclo infinito
		for {
			currentData := iterator.Next()
			if currentData == nil {
				break // Salir del ciclo si no hay datos
			}

			// Procesar los datos actuales
			// for _, pair := range currentData {
			// 	fmt.Printf("(%f, %f) ", pair[0], pair[1])
			// }
			// fmt.Println()
			result := helpers.CreatePoints(currentData)
			// Iterar sobre result y enviar cada elemento a la función enviarPOST
			for _, elem := range result {
				if err := helpers.EnviarPOST(*urlPtr, elem); err != nil {
					fmt.Printf("Error al enviar el elemento: %v\n", err)
				}
				
				// Esperar el tiempo configurado antes de la siguiente solicitud
				time.Sleep(time.Duration(*timePtr) * time.Second)
			}
		}
	} else {
		result := helpers.CreateListOfPoints(outpuLonLat)

		// Iterar sobre result y enviar cada elemento a la función enviarPOST
		for _, elem := range result {
			if err := helpers.EnviarPOST(*urlPtr, elem); err != nil {
				fmt.Printf("Error al enviar el elemento: %v\n", err)
			}
			
			// Esperar el tiempo configurado antes de la siguiente solicitud
			time.Sleep(time.Duration(*timePtr) * time.Second)
		}
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
