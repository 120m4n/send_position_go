package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"


	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"

	"send_position/geometry"
	"send_position/helpers"
)

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
	timePtr := flag.Int("pausa", 2000, "tiempo en mili-segundos entre request sucesivos")
    urlPtr  := flag.String("url", "", "url endpoint para envio de datos")
	verbosePtr := flag.Bool("verbose", false, "verbose")
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
	fmt.Println("Leer archivo geojson...")
	rawJSON, err := readFileIntoByteSlice(*geojsonPtr)
	if err != nil {
		log.Fatalf("Error al leer el archivo geojson: %v", err)
	}

	// Parsear el archivo geojson
	fmt.Println("Parsear archivo geojson...")
	feature, err := geojson.UnmarshalFeatureCollection(rawJSON)
	if err != nil {
		log.Fatalf("Error al parsear el archivo geojson: %v", err)
	}

	// Verificar si la feature es del tipo LineString
	fmt.Println("Verificar si la feature es del tipo LineString...")
	if feature.Features[0].Geometry.GeoJSONType() != geojson.TypeLineString {
		fmt.Println("El archivo geojson no contiene un elemento LineString")
		os.Exit(1)
	}

	//Contar el numero de segmentos en el LineString
	segmentCount := helpers.GetSegmentCount(feature.Features[0].Geometry.(orb.LineString))
	fmt.Printf("Numero de segmentos encontrados: %d\n", segmentCount)

	// Obtener las coordenadas del LineString
	fmt.Println("Obtener las coordenadas del LineString...")
	coordinates := feature.Features[0].Geometry.(orb.LineString)
    // procesando el archivo geojson
	fmt.Println("Procesando archivo geojson...")
	inputVector := helpers.ExtraeSegmentos(coordinates)
	output := helpers.SubsampleVector(inputVector)
	fmt.Printf("Subdividiendo segmentos...\n")
	newpoints := geometry.GeneratePoints(output, *puntosPtr, 3.0)
	fmt.Println("Total coordenadas a procesar: ", len(newpoints) * *puntosPtr)
	outpuLonLat := geometry.ConverToLonLat(newpoints)
	//helpers.PrintMatrix(outpuLonLat)

	if *ciclicoPtr {
	    fmt.Println("Enviando coordenadas de forma ciclica...")
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
				if err := helpers.EnviarPOST(*urlPtr, elem, *verbosePtr); err != nil {
					fmt.Printf("Error al enviar el elemento: %v\n", err)
				}
				
				// Sleep *timePtr milisegundos
				time.Sleep(time.Duration(*timePtr) * time.Millisecond)

			}
		}
	} else {
		fmt.Println("Enviando coordenadas de forma secuencial...")
		result := helpers.CreateListOfPoints(outpuLonLat)

		// Iterar sobre result y enviar cada elemento a la función enviarPOST
		for _, elem := range result {
			if err := helpers.EnviarPOST(*urlPtr, elem, *verbosePtr); err != nil {
				fmt.Printf("Error al enviar el elemento: %v\n", err)
			}
			
			// Esperar el tiempo configurado antes de la siguiente solicitud
			time.Sleep(time.Duration(*timePtr) * time.Millisecond)
		}
	}

}

