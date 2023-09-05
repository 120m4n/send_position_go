package geometry

import (
	"fmt"

	"github.com/go-spatial/proj"
)

// GeneratePoints toma una matriz de pares de puntos y genera puntos entre ellos
func GeneratePoints(input [][][2]float64, numSamples int) [][][2]float64 {
	output := make([][][2]float64, len(input))

	for i, pair := range input {
		points := GenerateLinePoints(pair[0], pair[1], numSamples)
		output[i] = points
	}

	return output
}

// GenerateLinePoints genera puntos a lo largo de una línea recta entre dos puntos
func GenerateLinePoints(start, end [2]float64, numSamples int) [][2]float64 {
	var points [][2]float64
	m := (end[1] - start[1]) / (end[0] - start[0])
	b := start[1] - m*start[0]

	step := (end[0] - start[0]) / float64(numSamples-1)

	for i := 0; i < numSamples; i++ {
		x := start[0] + float64(i)*step
		y := m*x + b
		points = append(points, [2]float64{x, y})
	}

	return points
}

// ConvertLinePoints convierte cada punto de la linea un punto proyectado lon,lat
func ConvertLinePoints(xy [2]float64 ) []float64 {
	// var points []float64

	var points = xy[:]
    
	lonlat, err := proj.Inverse(proj.EPSG3395, points)
		if err != nil {
			panic(err)
		}

	return lonlat
}

// ConvertToLonLat toma una matriz de coordenadas planas y devuelve su equivalente en coordenadas proyectadas lon,lat
func ConverToLonLat(input [][][2]float64) [][][2]float64 {
	output := make([][][2]float64, len(input))
	

	for i := 0; i < len(input); i++ {
		tmpLine := make([][2]float64, len(input[0]))
		for j := 0; j < len(input[i]); j++ {
			lonlat := ConvertLinePoints(input[i][j])
			var array [2]float64
			copy(array[:], lonlat)
			tmpLine[j] = array
		}
		output[i] = tmpLine
	}


	return output
}


// ExampleUsage es una función de ejemplo para usar las funciones del paquete
func ExampleUsage() {
	input := [][][2]float64{
		{{1.0, 2.0}, {4.0, 8.0}},
		{{2.0, 3.0}, {5.0, 12.0}},
	}

	numSamples := 5

	output := GeneratePoints(input, numSamples)

	fmt.Println(output)
}
