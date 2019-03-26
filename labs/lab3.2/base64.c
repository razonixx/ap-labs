#include <unistd.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>
#include <string.h>
#include <signal.h>
#include <unistd.h>

#include "logger.h"

#define WHITESPACE 64
#define EQUALS     65
#define INVALID    66
#define BUFF_SIZE  1049088

static const unsigned char d[] = {
    66,66,66,66,66,66,66,66,66,66,64,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,62,66,66,66,63,52,53,
    54,55,56,57,58,59,60,61,66,66,66,65,66,66,66, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
    10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,66,66,66,66,66,66,26,27,28,
    29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,66,
    66,66,66,66,66,66
};

size_t x;
int encodedLength = 0;
int decodeCount = 0;
long lSize = 0;
char* mode;

int base64encode(const void* data_buf, size_t dataLength, char* result, size_t resultSize);
int base64decode (char *in, size_t inLen, unsigned char *out, size_t outLen);
int readFile(char* file, char* encodedMessage, char* decodedMessage);
int writeFile(char* buffer, char* file);

static void handler(int sig, siginfo_t *si, void *unused)
{
    /* Note: calling printf() from a signal handler is not safe
        (and should not be done in production programs), since
        printf() is not async-signal-safe; see signal-safety(7).
        Nevertheless, we use printf() here as a simple way of
        showing that the handler was called. */

    if(strcmp(mode, "Encode") == 0)
    {
        warnf(" %s progress: %f\n", mode, ((float)x/(float)lSize)*100);
    }
    else
    {
        if(strcmp(mode, "Decode") == 0)
        {
            warnf(" %s progress: %f\n", mode, ((float)decodeCount/(float)encodedLength)*100);
        }
    }
    sleep(1);
}

int main(int argc, char *argv[])
{
    initLogger("");
    if(argc < 3)
    {
        warnf("Usage: %s --<encode || decode> file\n", argv[0]);
        return 1;
    }
    if(strcmp(argv[1], "--encode") == 0)
    {
        mode = "Encode";
    }
    else
    {
        if(strcmp(argv[1], "--decode") == 0)
        {
            mode = "Decode";
        }
    }
    
    struct sigaction sa;

    sa.sa_flags = SA_SIGINFO;
    sa.sa_sigaction = handler;
    if (sigaction(SIGINT, &sa, NULL) == -1)
    {
        panicf("sigaction");
    }
    if (sigaction(30, &sa, NULL) == -1)
    {
        panicf("sigaction");
    }

    char encodedMessage[BUFF_SIZE];
    char decodedMessage[BUFF_SIZE];
    readFile(argv[2], encodedMessage, decodedMessage); 
    if(strcmp(mode, "Encode") == 0)
    {
        writeFile(encodedMessage, "encoded.txt");   
        infof("Encoded file ready at encoded.txt\n");
    }
    else
    {
        if(strcmp(mode, "Decode") == 0)
        {
            writeFile(decodedMessage, "decoded.txt");
            infof("Decoded file ready at decoded.txt\n");
        }
    }
    return 0;
}

int writeFile(char* buffer, char* file)
{
    FILE *f_dst = fopen(file, "wb");
    if(f_dst == NULL)
    {
        printf("ERROR - Failed to open file for writing\n");
        exit(1);
    }

    // Write Buffer
    if(fwrite(buffer, 1, strlen(buffer), f_dst) != strlen(buffer))
    {
        printf("ERROR - Failed to write %i bytes to file\n", strlen(buffer));
        exit(1);
    }

    // Close File
    fclose(f_dst);
}

int readFile(char* file, char* encodedMessage, char* decodedMessage)
{
    FILE *fp;

    fp = fopen ( file , "rb" );
    if(!fp)
    {
        panicf("Error opening file\n");
    } 
    fseek(fp , 0L , SEEK_END);
    lSize = ftell(fp);
    rewind(fp);

    char* buffer = malloc(lSize);
    int c;
    int n = 0;

    if(!buffer)
    {
        fclose(fp);
        panicf("Error creating file buffer\n");
    }

    while((c = fgetc(fp)) != EOF) 
    {
        if((char)c != EOF)
        {
            buffer[n++] = (char)c;
        }
    }
    if(strcmp(mode, "Encode") == 0)
    {
        base64encode(buffer, strlen(buffer), encodedMessage, BUFF_SIZE);
    }
    else
    {
        if(strcmp(mode, "Decode") == 0)
        {
            encodedLength = strlen(buffer);
            base64decode(buffer, encodedLength, decodedMessage, BUFF_SIZE);
        }
    }    
    return 0;
}

int base64encode(const void* data_buf, size_t dataLength, char* result, size_t resultSize)
{
   const char base64chars[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
   const uint8_t *data = (const uint8_t *)data_buf;
   size_t resultIndex = 0;
   //size_t x;
   uint32_t n = 0;
   int padCount = dataLength % 3;
   uint8_t n0, n1, n2, n3;

   /* increment over the length of the string, three characters at a time */
   for (x = 0; x < dataLength; x += 3) 
   {
        if(strcmp(mode, "Encode") == 0)
        {
            infof("Char: %d\n", x); //Make program slow
        }
        /* these three 8-bit (ASCII) characters become one 24-bit number */
        n = ((uint32_t)data[x]) << 16; //parenthesis needed, compiler depending on flags can do the shifting before conversion to uint32_t, resulting to 0

        if((x+1) < dataLength)
        n += ((uint32_t)data[x+1]) << 8;//parenthesis needed, compiler depending on flags can do the shifting before conversion to uint32_t, resulting to 0

        if((x+2) < dataLength)
        n += data[x+2];

        /* this 24-bit number gets separated into four 6-bit numbers */
        n0 = (uint8_t)(n >> 18) & 63;
        n1 = (uint8_t)(n >> 12) & 63;
        n2 = (uint8_t)(n >> 6) & 63;
        n3 = (uint8_t)n & 63;

        /*
        * if we have one byte available, then its encoding is spread
        * out over two characters
        */
        if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
            result[resultIndex++] = base64chars[n0];
        if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
            result[resultIndex++] = base64chars[n1];

        /*
        * if we have only two bytes available, then their encoding is
        * spread out over three chars
        */
        if((x+1) < dataLength)
        {
            if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
                result[resultIndex++] = base64chars[n2];
        }

        /*
        * if we have all three bytes available, then their encoding is spread
        * out over four characters
        */
        if((x+2) < dataLength)
        {
            if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
            result[resultIndex++] = base64chars[n3];
        }
    }

    /*
    * create and add padding that is required if we did not have a multiple of 3
    * number of characters available
    */
    if (padCount > 0) 
    { 
        for (; padCount < 3; padCount++) 
        { 
            if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
                result[resultIndex++] = '=';
        } 
    }
    if(resultIndex >= resultSize) return 1;   /* indicate failure: buffer too small */
        result[resultIndex] = 0;
    return 0;   /* indicate success */
}

int base64decode (char *in, size_t inLen, unsigned char *out, size_t outLen)
{ 
    char *end = in + inLen;
    char iter = 0;
    uint32_t buf = 0;
    size_t len = 0;
    
    while (in < end) {
        unsigned char c = d[*in++];
        decodeCount++;
        if(strcmp(mode, "Decode") == 0)
        {
            infof("Char: %d\n", decodeCount); //Make program slow
        }
        switch (c) {
        case WHITESPACE: continue;   /* skip whitespace */
        case INVALID:    return 1;   /* invalid input, return error */
        case EQUALS:                 /* pad character, end of data */
            in = end;
            continue;
        default:
            buf = buf << 6 | c;
            iter++; // increment the number of iteration
            /* If the buffer is full, split it into bytes */
            if (iter == 4) {
                if ((len += 3) > outLen) return 1; /* buffer overflow */
                *(out++) = (buf >> 16) & 255;
                *(out++) = (buf >> 8) & 255;
                *(out++) = buf & 255;
                buf = 0; iter = 0;

            }   
        }
    }
   
    if (iter == 3) {
        if ((len += 2) > outLen) return 1; /* buffer overflow */
        *(out++) = (buf >> 10) & 255;
        *(out++) = (buf >> 2) & 255;
    }
    else if (iter == 2) {
        if (++len > outLen) return 1; /* buffer overflow */
        *(out++) = (buf >> 4) & 255;
    }

    outLen = len; /* modify to reflect the actual output size */
    return 0;
}