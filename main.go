package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"

	"send_position/geometry"
	"send_position/helpers"
	"send_position/singleton"
)

func readFileIntoByteSlice(filename string) ([]byte, error) {
    return ioutil.ReadFile(filename)
}



func main() {

	ciclicoPtr := flag.Bool("ciclico", false, "procesar el geojson de forma ciclica")
	fleetPtr := flag.String("fleet", "", "avatar")
	geojsonPtr := flag.String("json", "", "archivo geojson para procesar")
    puntosPtr := flag.Int("puntos", 3, "numero de puntos por segmento de linea")
	parametroPtr := flag.String("parametro", "", "parametro de la ruta http")
	timePtr := flag.Int("pausa", 2000, "tiempo en mili-segundos entre request sucesivos")
    urlPtr  := flag.String("url", "", "url endpoint para envio de datos")
	useridPtr := flag.String("userid", "", "Ggoland/roman")
	verbosePtr := flag.Bool("verbose", false, "verbose")
	flag.Parse()
	
	if (*geojsonPtr == "") {
		fmt.Println("Uso: send_position -json=<ruta_al_archivo_geojson>")
		os.Exit(1)
	} 

	// URL del servidor donde realizar la solicitud POST
    if *urlPtr == "" {
        log.Fatal("Uso: send_position -json=<ruta_al_archivo_geojson> -url=http://localhost:3000")
    } else {
        u, err := url.Parse(*urlPtr)
        if err != nil {
            log.Fatalf("Error parsing URL: %v", err)
        }
        if !strings.HasSuffix(u.Path, "/") {
            u.Path += "/"
        }
        *urlPtr = u.String()
    }
	
	if *parametroPtr != "" {  
       *urlPtr = *urlPtr + *parametroPtr
	} 

	if *useridPtr == "" {
		*useridPtr = "Ggoland"
	}

	if *fleetPtr == "" {
		*fleetPtr = "avatar"
	}

	s := singleton.GetInstance()
	s.SetFleet(*fleetPtr)
	s.SetUserid(*useridPtr)

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
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			for _, elem := range result {
				err := helpers.EnviarPOST(ctx, *urlPtr, elem, *verbosePtr)
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						log.Println("Request timed out, retrying with a longer deadline...")
						// Retry with a longer deadline
						ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
						defer cancel()
						err = helpers.EnviarPOST(ctx, *urlPtr, elem, *verbosePtr)
					}
					if err != nil {
						log.Printf("Error al enviar el elemento: %v\n", err)
					}
				}

				// Esperar el tiempo configurado antes de la siguiente solicitud
				time.Sleep(time.Duration(*timePtr) * time.Millisecond)
			}
		}
	} else {
		fmt.Println("Enviando coordenadas de forma secuencial...")
		result := helpers.CreateListOfPoints(outpuLonLat)

		// Iterar sobre result y enviar cada elemento a la función enviarPOST
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		for _, elem := range result {
			err := helpers.EnviarPOST(ctx, *urlPtr, elem, *verbosePtr)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					log.Println("Request timed out, retrying with a longer deadline...")
					// Retry with a longer deadline
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()
					err = helpers.EnviarPOST(ctx, *urlPtr, elem, *verbosePtr)
				}
				if err != nil {
					log.Printf("Error al enviar el elemento: %v\n", err)
				}
			}

			// Esperar el tiempo configurado antes de la siguiente solicitud
			time.Sleep(time.Duration(*timePtr) * time.Millisecond)
		}
	}

}

