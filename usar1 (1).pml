#define wait(s) atomic { s > 0 -> s-- }
#define signal(s) s++

#define NUM_GOROUTINES 4
#define TOTAL_POINTS 100

byte sumLock = 1;
int sumX = 0, sumY = 0;

active [NUM_GOROUTINES] proctype Worker() {
    int localSumX = 0, localSumY = 0;
    int i = 0;

    printf("Iniciando goroutines %d\n", _pid);
    //TOTAL_POINTS / NUM_GOROUTINES, asegurando una distribución equitativa del trabajo.
    do :: (i < TOTAL_POINTS / NUM_GOROUTINES) -> {
        localSumX = localSumX + i;
        localSumY = localSumY + i;
        // Imprime el estado actual de las sumas locales y el valor de i.
        printf("gorotines %d: localSumX = %d, localSumY = %d, i = %d\n", _pid, localSumX, localSumY, i);
        i = i + 1;
    }
    :: else -> {
        // Notifica la finalización del bucle.
        printf("gorotines %d terminando bucle con i = %d\n", _pid, i);
        break;
    }
    od;
    //múltiples hilos accedan simultáneamente a un recurso compartido
    atomic {
        wait(sumLock);
        // Notifica antes de actualizar las sumas globales.
        assert(sumLock == 0) //verifica q el semaforo esta tomando 0
        printf("goroutines %d esperando sumLock. Sumas globales actuales: sumX = %d, sumY = %d\n", _pid, sumX, sumY);
        sumX = sumX + localSumX;
        sumY = sumY + localSumY;
        // Notifica después de actualizar las sumas globales.
        printf("goroutines %d actualizó sumas globales: sumX = %d, sumY = %d\n", _pid, sumX, sumY);
        signal(sumLock);
    }
}
