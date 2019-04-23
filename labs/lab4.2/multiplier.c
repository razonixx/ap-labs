#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>

#include "logger.h"

struct dotProductArgs {
    int  index;
    int   lock;
    int      j;
    int      k;
    long *vec1;
    long *vec2;
};

#define NUM_THREADS 2000
int NUM_BUFFERS;
long **buffers;
pthread_mutex_t *mutexes;
pthread_attr_t attr;
long result[4000000+1];
long *matA, *matB;

long *readMatrix(char *filename);
long *getColumn(int col, long *matrix);
long *getRow(int row, long *matrix);
int getLock(int threadID);
int releaseLock(int lock);
void* dotProduct(void *args);
void* dotProductSerial(void *args);
long *multiply(long *matA, long *matB);
int saveResultMatrix(long *result);

int main (int argc, char *argv[])
{
    ///////////////////////////////////////////////////////
    if(argc < 3)
    {
        errorf("Usage: %s -n <buffers>\n", argv[0]);
        return -1;
    }
    NUM_BUFFERS = atoi(argv[2]);
    ///////////////////////////////////////////////////////

    ///////////////////////////////////////////////////////
    buffers = calloc(NUM_BUFFERS, sizeof(long*));
    for(int i = 0; i < NUM_BUFFERS; i++)
    {
        buffers[i] = calloc(1, sizeof(long));
    }
    ///////////////////////////////////////////////////////

    ///////////////////////////////////////////////////////
    mutexes = calloc(NUM_BUFFERS, sizeof(pthread_mutex_t));
    for(int i = 0; i < NUM_BUFFERS; i++)
    {
        pthread_mutex_init(&mutexes[i], NULL);
    }
    ///////////////////////////////////////////////////////

    matA = readMatrix("matA.dat");
    matB = readMatrix("matB.dat");

    multiply(matA, matB);
    infof("Done\n");

    saveResultMatrix(result);
    free(buffers);
    free(matA);
    free(matB);
    for(int i = 0; i < NUM_BUFFERS; i++)
    {
        pthread_mutex_destroy(&mutexes[i]);
    }
    pthread_exit(NULL);
}

long *readMatrix(char *filename)
{
    int i = 1;
    size_t len = 0;
    ssize_t nread;
    char *line = NULL;
    long *res = calloc(4000000+1, sizeof(long));

    FILE *fp = fopen(filename, "rb");
    if(!fp)
    {
        panicf("Error opening file\n");
    } 

    while ((nread = getline(&line, &len, fp)) != -1) 
    {
        //fwrite(line, nread, 1, stdout);
        res[i++] = strtol(line, NULL, 10);
    }
    free(line);
    return res;
}

long *getRow(int row, long *matrix)
{
    long *res = calloc(2000+1, sizeof(long));

    for(int i = 1; i < 2000+1; i++)
    {
        res[i] = matrix[row * 2000 - 2000 + i];
    }
    return res;
}

long *getColumn(int col, long *matrix)
{
    long *res = calloc(2000+1, sizeof(long));

    for(int i = 1; i < 2000+1; i++)
    {
        res[i] = matrix[i * 2000 - 2000 + col];
    }
    return res;
}

int getLock(int threadID)
{
    int res = -1;
    int i = 0;
    while(res != 0)
    {
        res = pthread_mutex_trylock(&mutexes[i%NUM_BUFFERS]);
        if(res == 0)
        {
            //warnf("Thread %d has locked number %d\n", threadID, i%NUM_BUFFERS);  //Uncomment to show the mutexes being shared
            return i%NUM_BUFFERS;
        }
        i++;
    }
    return -1;
}

int releaseLock(int lock)
{
    int res = pthread_mutex_unlock(&mutexes[lock]);
    if(res == 0)
    {
        return 0;
    }
    return -1;
}

void* dotProduct(void *args)
{
    struct dotProductArgs *dpArgs = (struct dotProductArgs *)args;
    long temp = 0;
    int lock = getLock(dpArgs->index);    
    dpArgs->vec1 = getRow(dpArgs->j, matA);
    dpArgs->vec2 = getColumn(dpArgs->k, matB);
    for(int i = 1; i < 2000+1; i++)
    {
        temp += dpArgs->vec1[i]*dpArgs->vec2[i];
    }

    free(dpArgs->vec1);
    free(dpArgs->vec2);
    buffers[dpArgs->lock][0] = temp;
    result[dpArgs->index] = buffers[dpArgs->lock][0];
    releaseLock(lock);
    return NULL;
}

long *multiply(long *matA, long *matB)
{
    pthread_attr_init(&attr);
    pthread_t threads[NUM_THREADS];
    struct dotProductArgs *dpArgsArr = calloc(4000000+1, sizeof(struct dotProductArgs));

    int i = 1;
    for(int j = 1; j < 2000+1; j++)
    {
        for(int k = 1; k < 2000+1; k++)
        {
            dpArgsArr[i].index = i;
            dpArgsArr[i].j = j;
            dpArgsArr[i].k = k;
            int err = pthread_create(&threads[i%NUM_THREADS], &attr, dotProduct, (void *)&dpArgsArr[i]); 
            if(err != 0)
            {
                panicf("Error %d when creating thread %d", err, i);
                exit(-1);
            }
            pthread_detach(threads[i%NUM_THREADS]);
            i++;
        }
    }
}

int saveResultMatrix(long *result)
{
    FILE *fp = fopen("result.dat", "w");
    if(!fp)
    {
        panicf("Error opening file\n");
    }

    for(int i = 1; i < 4000000+1; i++)
    {
        char* temp = calloc(8, sizeof(char));
        sprintf(temp, "%ld", result[i]);
        fputs(temp, fp);
        fputc('\n', fp);
        free(temp);
    }
    return 0;
}

