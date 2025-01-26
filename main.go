package main

import (
	"context"
	"errors"
	"flag"
	// "fmt"
	// "io/ioutil"
	"log"
	"net/url"
	"os"
    "os/signal"
    "syscall"
	"strings"
	"time"

	"github.com/google/uuid"
	// "github.com/paulmach/orb"
	// "github.com/paulmach/orb/geojson"

	// "send_position/geometry"
	"send_position/helpers"
	"send_position/singleton"
)

// func readFileIntoByteSlice(filename string) ([]byte, error) {
//     return ioutil.ReadFile(filename)
// }



func main() {

	latPtr := flag.Float64("lat", 7.1299, "latitude")
    lonPtr := flag.Float64("lon", -73.111929, "longitude")
	fleetPtr := flag.String("fleet", "", "avatar")
	// geojsonPtr := flag.String("json", "", "archivo geojson para procesar")
    puntosPtr := flag.Int("puntos", 113, "numero de puntos por segmento de linea")
	uniqueidPtr := flag.String("uniqueid", "", "identificacion unica del usuario")
	timePtr := flag.Int("pausa", 2000, "tiempo en mili-segundos entre request sucesivos")
    urlPtr  := flag.String("url", "", "url endpoint para envio de datos")
	useridPtr := flag.String("userid", "", "Ggoland/roman")
	verbosePtr := flag.Bool("verbose", false, "verbose")
	flag.Parse()
	

	// URL del servidor donde realizar la solicitud POST
    if *urlPtr == "" {
        log.Fatal("Uso: send_position -json=<ruta_al_archivo_geojson> -url=http://localhost:3000")
    } else {
        u, err := url.Parse(*urlPtr)
        if err != nil {
            log.Fatalf("Error parsing URL: %v", err)
        }
        // if !strings.HasSuffix(u.Path, "/") {
        //     u.Path += "/"
        // }
        *urlPtr = u.String()
    }
	
	if *uniqueidPtr == "" {
		uuidWithHyphen := uuid.New()
		uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
		*uniqueidPtr = uuid
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
	s.SetUniqueid(*uniqueidPtr)

	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive SIGINT and SIGTERM signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a channel to signal when the program should exit
	done := make(chan bool, 1)

	// Start a goroutine that will set done to true when a signal is received
	go func() {
		sig := <-sigs
		log.Printf("Received signal: %v", sig)
		done <- true
	}()


    // Iterate from 0 to *puntosPtr, generate a CreateRandomPoint(-73,12) 
    // and send each element to the EnviarPOST function
    for i := 0; i < *puntosPtr; i++ {
        select {
        case <-done:
            log.Println("Exiting due to interrupt signal")
            return
        default:
            elem := helpers.CreateRandomPoint(*lonPtr, *latPtr)

            err := helpers.EnviarPOST(*urlPtr, elem, *verbosePtr)
            if err != nil {
                if errors.Is(err, context.DeadlineExceeded) {
                    log.Println("Request timed out, retrying with a longer deadline...")
                } else {
                    log.Printf("Error al enviar el elemento: %v\n", err)
                }
            }

            // Wait the configured time before the next request
            time.Sleep(time.Duration(*timePtr) * time.Millisecond)
        }
    }
}

