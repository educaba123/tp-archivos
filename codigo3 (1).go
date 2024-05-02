package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

const (
    numPuntos       = 1000000
    numGoroutines   = 4
    numEjecuciones  = 1000
    mostrarResultados = 50
)

type Point struct {
    X float64
    Y float64
}

// Generar datos aleatorios 
func generarDatos(n int) []Point {
    rand.Seed(time.Now().UnixNano())
    puntos := make([]Point, n)
    for i := range puntos {
        x := rand.Float64() * 100  // X se genera aleatoriamente.
        y := 2*x + 10 + rand.NormFloat64()*10 // Y = 2X + 10
        puntos[i] = Point{
            X: x,
            Y: y,
        }
    }
    return puntos
}


// Versión secuencial de la regresion lineal
func linearRegression(points []Point) (m float64, b float64) {
    var sumX, sumY, sumXY, sumX2, n float64
    n = float64(len(points))
    for _, p := range points {
        sumX += p.X
        sumY += p.Y
        sumXY += p.X * p.Y
        sumX2 += p.X * p.X
    }
    // Cálculo de la pendiente (m) y el intercepto (b) de la regresion lineal
    if denom := (n*sumX2 - sumX*sumX); denom != 0 { // Asegurar que el denominador no sea cero
        m = (n*sumXY - sumX*sumY) / denom
        b = (sumY - m*sumX) / n
    } else {
        fmt.Println(" revisar datos de entrada")
    }
    return
}

// Versión concurrente de la regresion lineal
func linearRegressionConcurrente(points []Point) (m float64, b float64) {
    var sumX, sumY, sumXY, sumX2 float64
    var wg sync.WaitGroup
    var mutex sync.Mutex// Este proceso asegura que cuando multiples goroutines actualizan las variables compartidas (sumX, sumY, sumXY, sumX2)
    n := float64(len(points))
    segmentSize := len(points) / numGoroutines

    for i := 0; i < numGoroutines; i++ {
        start := i * segmentSize
        end := start + segmentSize
        if i == numGoroutines-1 {
            end = len(points) // Asegura que cubrimos todos los puntos
        }

        wg.Add(1)
        go func(pts []Point) {
            defer wg.Done()
            var sumXLocal, sumYLocal, sumXYLocal, sumX2Local float64
            for _, p := range pts {
                sumXLocal += p.X
                sumYLocal += p.Y
                sumXYLocal += p.X * p.Y
                sumX2Local += p.X * p.X
            }
			//// Mutex lock y unlock para proteger la actualizacion de las sumas globales
            mutex.Lock()
            sumX += sumXLocal
            sumY += sumYLocal
            sumXY += sumXYLocal
            sumX2 += sumX2Local
            mutex.Unlock()
        }(points[start:end])
    }

    wg.Wait() // Esperar a que todas las goroutines terminen

    if denom := (n*sumX2 - sumX*sumX); denom != 0 { // Asegurar que el denominador no sea cero
        m = (n*sumXY - sumX*sumY) / denom
        b = (sumY - m*sumX) / n
    } else {
        fmt.Println("Denominador cero detectado, revisar datos de entrada")
    }
    return
}

func main() {
    var totalMSec, totalBSec, totalMCon, totalBCon float64
    var totalTiempoSec, totalTiempoCon time.Duration

    for i := 0; i < numEjecuciones; i++ {
        puntos := generarDatos(numPuntos)
		// En la función main(), justo después de generar los datos:
		if i == 0 { // Solo para la primera ejecución
    		fmt.Println("Primeros 100 puntos de la primera ejecución:")
    		for j := 0; j < 100; j++ {
        	fmt.Printf("Punto %d: X = %.2f, Y = %.2f\n", j+1, puntos[j].X, puntos[j].Y)
    		}
		}

        // Calculo secuencial
        start := time.Now()
        mSec, bSec := linearRegression(puntos)
        durationSec := time.Since(start)
        totalTiempoSec += durationSec

        // Calculo concurrente
        start = time.Now()
        mCon, bCon := linearRegressionConcurrente(puntos)
        durationCon := time.Since(start)
        totalTiempoCon += durationCon

        // Acumulación de resultados para promedios
        totalMSec += mSec
        totalBSec += bSec
        totalMCon += mCon
        totalBCon += bCon

        // Muestra resultados para las primeras 50 ejecuciones
        if i < mostrarResultados {
            fmt.Printf("Ejecución %d - Secuencial: Pendiente = %.2f, Intercepto = %.2f, Tiempo = %v\n", i+1, mSec, bSec, durationSec)
            fmt.Printf("Ejecución %d - Concurrente: Pendiente = %.2f, Intercepto = %.2f, Tiempo = %v\n", i+1, mCon, bCon, durationCon)
        }
    }

    // Calcular promedios
    promedioMSec := totalMSec / float64(numEjecuciones)
    promedioBSec := totalBSec / float64(numEjecuciones)
    promedioMCon := totalMCon / float64(numEjecuciones)
    promedioBCon := totalBCon / float64(numEjecuciones)

    // Mostrar el promedio de los resultados y el tiempo total
    fmt.Printf("Promedio Secuencial - Pendiente: %.2f, Intercepto: %.2f, Tiempo total : %v\n", promedioMSec, promedioBSec, totalTiempoSec)
    fmt.Printf("Promedio Concurrente - Pendiente: %.2f, Intercepto: %.2f, Tiempo total : %v\n", promedioMCon, promedioBCon, totalTiempoCon)
}
