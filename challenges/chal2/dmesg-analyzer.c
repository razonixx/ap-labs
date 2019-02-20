#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/stat.h>
#include <string.h>
 
#define REPORT_FILE "report.txt"

const int MAP_SIZE = 0x00FFF;

struct Map {
    char* key;
    char* value;
};


void analizeLog(char *logFile, char *report);
void printMap(struct Map *map, struct Map *mapLogType, int fd, char* outB, int bufferSize);
void initMap(struct Map *map);
int checkForLog(struct Map *map, char *log);
void poblateMap(struct Map *map, char* buffer, int bufferSize);
void parseLine(char **line, struct Map *data, int lineNumber);
void sortMap(struct Map *map, int n);

int main(int argc, char **argv) {

    if (argc < 2) {
	printf("Usage:./dmesg-analizer.o logfile.txt\n");
	return 1;
    }

    analizeLog(argv[1], REPORT_FILE);

    return 0;
}

void analizeLog(char *logFile, char *report) {
    printf("Generating Report from: [%s] log file\n", logFile);
    // Implement your solution here.

    struct Map mapLogType[MAP_SIZE];
    initMap(mapLogType);

    //Report FD
    int reportfd = open(report, O_RDWR | O_CREAT, 0666); 
    //Log FD
    int logfd = open(logFile, O_RDONLY); 

    struct stat st;
    fstat(logfd, &st);
    int bufferSize = st.st_size;    

    char* readBuffer = (char *) calloc(bufferSize, sizeof(char));
    char writeBuffer[40000];

    read(logfd, readBuffer, bufferSize);

    struct Map map2[MAP_SIZE];
    initMap(map2);
    poblateMap(map2, readBuffer, bufferSize);
    printMap(map2, mapLogType, reportfd, writeBuffer, bufferSize);
    
    write(reportfd, writeBuffer, bufferSize);
    free(readBuffer);
    close(reportfd);
    close(logfd);
    printf("Report is generated at: [%s]\n", report);
}

void sortMap(struct Map *map, int n)
{
    char* tempKey = (char*)calloc(128, sizeof(char));
    char* tempValue = (char*)calloc(2048, sizeof(char));
    for (int j=0; j<n-1; j++) 
    { 
        for (int i=j+1; i<n; i++) 
        { 
            
            if( (strcmp(map[j].key, map[i].key) > 0) && ((strcmp(map[i].key, "NULL") != 0) && (strcmp(map[j].key, "NULL") != 0)))
            { 
                (tempKey = map[j].key); 
                (tempValue = map[j].value); 
                (map[j].key = map[i].key); 
                (map[j].value = map[i].value); 
                (map[i].key = tempKey);
                (map[i].value = tempValue); 
            } 
        } 
    } 
    free(tempKey);
    free(tempValue);
}

void poblateMap(struct Map *map, char* buffer, int bufferSize)
{
    char **line = (char**)calloc(1024, sizeof(char*));
    for(int i = 0; i < 1024; i++)
    {
        line[i] = "";        
    }
    
    int numLines = 0;
    for(int i = 0; i < bufferSize; i++)
    {
        if(buffer[i] == '\n')
        {
            numLines++;
        }
    }

    struct Map data[MAP_SIZE];
    initMap(data);
    
    int carry = 0;
    for(int i = 0; i < numLines; i++)
    {
        char *temp = (char*)calloc(1024, sizeof(char));
        int k = 0;
        while(buffer[carry] != '\n')
        {
            temp[k] = buffer[carry];
            k++;
            carry++;
        }
        carry++;
        line[i] = temp;
        parseLine(line, data, i); 
        map[i].key = data[i].key;
        map[i].value = data[i].value;
        free(temp);
    }
    sortMap(map, numLines);
    free(line);
}

void parseLine(char **line, struct Map *data, int lineNumber)
{
    int logFlag = 0;
    int logFlag2 = 0;
    char *log_type = (char*)calloc(1024, sizeof(char));
    char *nums = (char*)calloc(1024, sizeof(char));
    char *msg = (char*)calloc(1024, sizeof(char));
    msg[0] = '\0';
    int w = 0;
    int x = 0;
    int y = 0;
    int z = 0;
    while (line[lineNumber][z-1] != ']') {
        nums[x++] = line[lineNumber][z++];
    }
    if(line[lineNumber][z+1] == ' ' && line[lineNumber][z+2] == ' ' && (strcmp(data[lineNumber-1].key, "General")!=0))
    {
        logFlag = 1;
    }
    while (line[lineNumber][z] != '\0')
    {
        if (line[lineNumber][z] == ':' && line[lineNumber][z+1] == '\0') 
        {
            logFlag2 = 1;
        }
        if (line[lineNumber][z] == ':' && line[lineNumber][z+1] == ' ') 
        {
            break;
        }
        log_type[w++] = line[lineNumber][z];
        if (line[lineNumber][z] == ']')
        {
            break;
        }
        z++;
    }
    log_type[w] = '\0';
    while (line[lineNumber][z++] != '\0')
    {
        msg[y++] = line[lineNumber][z];
    }
    msg[y] = '\0';
    if(msg[0] == '\0')
    {
        strcpy(msg, log_type);
        strcpy(log_type, " General");
    }
    if(!logFlag)
    {
        data[lineNumber].key = log_type;
    }
    
    if(logFlag2)
    {
        data[lineNumber].key = msg;
    }
    else
    {
        if(strcmp(data[lineNumber].key, "") == 0 || strcmp(data[lineNumber].key, "NULL") == 0)
        {
            data[lineNumber].key = data[lineNumber-1].key;
        }
        data[lineNumber].value = strcat(nums, msg);
    }
}

void printMap(struct Map *map, struct Map *mapLogType, int fd, char* outB, int bufferSize)
{
    int j = 0; //MapLogType counter
    for(int i = 0; i < MAP_SIZE; i++)
    {
        if(map[i].key == "NULL")
        {
            break;
        }
        if(checkForLog(mapLogType, map[i].key) && strcmp(map[i].value, "NULL") != 0)
        {
            //printf("  %s\n",map[i].value);
            strcat(outB, map[i].value);
            strcat(outB, "\n");
        }
        else
        {
            mapLogType[j].key = map[i].key;
            if(strcmp(map[i].key, "\n") != 0)
            {
                //printf("%s\n",map[i].key);
                strcat(outB, map[i].key);
                strcat(outB, "\n");

            }
            if(strcmp(map[i].key, "NULL") != 0 && strcmp(map[i].value, "NULL") != 0)
            {
                //printf("  %s\n",map[i].value);
                strcat(outB, map[i].value);
                strcat(outB, "\n");
            }     
            j++;
        }
    }
}

int checkForLog(struct Map *map, char *log)
{
    for(int i = 0; i < MAP_SIZE; i++)
    {
        if(map[i].key != "NULL" && strcmp(map[i].key,log)==0)
        {
            return 1;
        }
        if(map[i].key == "NULL")
        {
            break;
        }
    }
    return 0;
}

void initMap(struct Map *map)
{
    for(int i = 0; i < MAP_SIZE; i++)
    {
        map[i].key = "NULL";
        map[i].value = "NULL";
    }
}
