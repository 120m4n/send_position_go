package main

import "fmt"

type CircularIterator struct {
    data   [][][2]float64
    index  int
}

func NewCircularIterator(data [][][2]float64) *CircularIterator {
    return &CircularIterator{
        data:   data,
        index:  0,
    }
}

func (ci *CircularIterator) Next() [][2]float64 {
    if len(ci.data) == 0 {
        return nil
    }
    currentData := ci.data[ci.index]
    ci.index = (ci.index + 1) % len(ci.data)
    return currentData
}

func main() {
    // Ejemplo de datos: una colección de tres conjuntos de datos
    data := [][][2]float64{
		{{1.1, 2.2}, {3.3, 4.4}},
		{{5.5, 6.6}, {7.7, 8.8}},            
		{{9.9, 10.10}, {11.11, 12.12}},
		{{13.13, 14.14}, {15.15, 16.16}},
    }

    iterator := NewCircularIterator(data)

    // Ejemplo de uso en un ciclo infinito
    for {
        currentData := iterator.Next()
        if currentData == nil {
            break // Salir del ciclo si no hay datos
        }

        // Procesar los datos actuales
        //fmt.Println("Datos actuales:")
        for _, pair := range currentData {
            fmt.Printf("(%f, %f) ", pair[0], pair[1])
        }
        fmt.Println()
    }
}
