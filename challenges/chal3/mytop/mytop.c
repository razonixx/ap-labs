#include <stdio.h>
#include <string.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <dirent.h>
#include <stdlib.h> 
#include <unistd.h>
#include <time.h>
#include <signal.h>



struct entry
{
    char* pid;
    char* ppid;
    char* name;
    char* state;
    long memory;
    int threads;
    int files;
};


void clear();
void printTable(struct entry* ps, int c);
void ssTable(FILE* file, struct entry* ps, int c);

int k = 0;
struct entry processes[1000];
char ssPath[50];
FILE* ssFp;


static void handler(int sig, siginfo_t *si, void *unused)
{
    /* Note: calling printf() from a signal handler is not safe
        (and should not be done in production programs), since
        printf() is not async-signal-safe; see signal-safety(7).
        Nevertheless, we use printf() here as a simple way of
        showing that the handler was called. */
    ssFp = fopen(ssPath, "w");
    ssTable(ssFp, processes, k);
    fclose(ssFp);
}


int main()
{
    struct sigaction sa;

    sa.sa_flags = SA_SIGINFO;
    sa.sa_sigaction = handler;
    if (sigaction(SIGINT, &sa, NULL) == -1)
    {
        printf("sigaction");
        return 1;
    }

    time_t t = time(NULL);
    struct tm tm = *localtime(&t);

    char path[50];
    strcat(path, "/proc/");
    struct dirent *dir;
    FILE* fp;
    int i = 0;
    sprintf(ssPath, "mytop_status_%d_%d_%d.txt", tm.tm_year + 1900, tm.tm_mon + 1, tm.tm_mday);

    long rss;
    void* d; //dummy

    
    while(1) {
        DIR* proc = opendir(path);
        i = 0;
        rewinddir(proc);
        char* tempPath;
        char* statPath;

        if (proc != NULL)
        {
            while ((dir = readdir(proc)))
            {
                if(atoi(dir->d_name) > 0) // get pid as dirName
                {
                    char c = 0;

                    tempPath = calloc(50, sizeof(char));
                    statPath = calloc(50, sizeof(char));
                    strcat(tempPath, path);
                    strcat(tempPath, dir->d_name);
                    strcat(statPath, tempPath);
                    strcat(tempPath, "/status");
                    strcat(statPath, "/stat");
                    fp = fopen(tempPath, "r"); //open /proc/x/status
                    FILE* fpStat = fopen(statPath, "r");
                    fscanf(fpStat, "%d %s %c %d %d %d %d %d %u %lu %lu %lu %lu %lu %ld %ld %ld %ld %ld %ld %llu %lu %ld %ld", d, d, d, d, d, d, d, d, d, d, d, d, d, d, d, d, d, d, d, d, d, d, d, &rss);
                    //               1  2  3  4  5  6  7  8  9   10 11  12 13  14  15   16  17  18  19  20 21    22  23  24

                    fclose(fpStat);
                    free(tempPath);                    
                    free(statPath);
                    long memSize = (rss*4096)/1024; // # of pages * pagesize = size in bytes, /1024 = size in kbytes
                    //char* queEstaPasando = calloc(1, sizeof(char));
                    //free(queEstaPasando);
                    processes[i].memory = memSize; 
                    char* tempString1 = calloc(50, sizeof(char));
                    char* tempString2 = calloc(50, sizeof(char));

                    // Get pid and parent pid
                    int z = 0;
                    while((c = getc(fp)) != 'P'){}
                    sprintf(tempString1, "%c", 'P');
                    while(z < 2) 
                    {
                        tempString1 = calloc(50, sizeof(char));
                        tempString2 = calloc(50, sizeof(char));
                        fgetc(fp);
                        fgetc(fp);
                        fgetc(fp);
                        fgetc(fp);
                        if(z == 1)
                        {
                            fgetc(fp);
                            fgetc(fp);
                        }
                        while((c = getc(fp)) != '\n') {
                            sprintf(tempString2, "%c", c);
                            strcat(tempString1, tempString2);                    
                        }
                        if(z == 0)
                        {
                            processes[i].pid = tempString1;
                        }
                        if(z == 1)
                        {
                            if(strlen(tempString1) > 1)
                            {
                                processes[i].ppid = tempString1;
                            }
                            else
                            {
                                processes[i].ppid = "0";
                            }
                        }
                        z++;
                    }
                    rewind(fp);
                    while((c = getc(fp)) != EOF)
                    {
                        if(c == 'N')
                        {
                            if(getc(fp) == 'a')
                            {
                                if(getc(fp) == 'm')
                                {
                                    if(getc(fp) == 'e')
                                    {
                                        fgetc(fp); //skip :
                                        fgetc(fp); //skip \t
                                        tempString1 = calloc(50, sizeof(char));
                                        tempString2 = calloc(50, sizeof(char));
                                        while((c = getc(fp)) != '\n')
                                        {
                                            sprintf(tempString2, "%c", c);
                                            strcat(tempString1, tempString2);
                                        }
                                        processes[i].name = tempString1;
                                        fgetc(fp); //skip \n
                                    }
                                }
                            }
                        }
                    }
                    rewind(fp);
                    while((c = getc(fp)) != EOF)
                    {
                        if(c == 'S')
                        {
                            if(getc(fp) == 't')
                            {
                                if(getc(fp) == 'a')
                                {
                                    if(getc(fp) == 't')
                                    {
                                        if(getc(fp) == 'e')
                                        {
                                            fgetc(fp); //skip :
                                            fgetc(fp); //skip \t
                                            fgetc(fp); //skip (
                                            fgetc(fp); //?
                                            fgetc(fp); //?
                                            tempString1 = calloc(50, sizeof(char));
                                            tempString2 = calloc(50, sizeof(char));
                                            while((c = getc(fp)) != ')')
                                            {
                                                sprintf(tempString2, "%c", c);
                                                strcat(tempString1, tempString2);
                                            }
                                            processes[i].state = tempString1;
                                        }
                                    }
                                }
                            }
                        }
                    }
                    rewind(fp);
                    while((c = getc(fp)) != EOF)
                    {
                        if(c == 'T')
                        {
                            if(getc(fp) == 'h')
                            {
                                if(getc(fp) == 'r')
                                {
                                    if(getc(fp) == 'e')
                                    {
                                        if(getc(fp) == 'a')
                                        {
                                            fgetc(fp); //skip d
                                            fgetc(fp); //skip s
                                            fgetc(fp); //skip :
                                            fgetc(fp); // \t
                                            tempString1 = calloc(50, sizeof(char));
                                            tempString2 = calloc(50, sizeof(char));
                                            while((c = getc(fp)) != '\n')
                                            {
                                                sprintf(tempString2, "%c", c);
                                                strcat(tempString1, tempString2);
                                            }
                                            processes[i].threads = atoi(tempString1);
                                        }
                                    }
                                }
                            }
                        }
                    }
                    int fileCount = 0;
                    struct dirent * entry;
                    char* dirTempPath = calloc(50, sizeof(char));
                    strcat(dirTempPath, path);
                    strcat(dirTempPath, dir->d_name);
                    strcat(dirTempPath, "/fd/");
                    DIR* tempDir = opendir(dirTempPath);
                    if(readdir(tempDir) != NULL)
                    {
                        while ((entry = readdir(tempDir)) != NULL) 
                        {
                                fileCount++;                    
                        }
                        fileCount--; //Substract 1 for .
                        processes[i].files = fileCount;
                    }
                    else
                    {
                        processes[i].files = 0;
                    }
                    //printf("Proc: %s Size of proc: %ld K\n",processes[i].name, memSize); 
                    free(dirTempPath);
                    free(tempString1);
                    free(tempString2);
                    fclose(fp);
                    closedir(tempDir);
                    i++;
                    k++;
                    c = 0;
                }
            }
        }
        clear();
        printTable(processes, k);
        sleep(2);
        closedir(proc);
        k = 0;
    }
    return 0;
}

void printTable(struct entry* ps, int c)
{
    printf("+-------+--------+------------------------------------------+------------+----------+----------+------------+\n");
    printf("|  PID  | Parent |                   Name                   |   State    |  Memory  | #Threads | Open Files |\n");
    printf("+-------+--------+------------------------------------------+------------+----------+----------+------------+\n");
    for (int j = 0; j < c; j++) 
    {
        printf("| %-5s | %-6s | %-40s | %-10s | %-7ld K | %-8d | %-10d |\n", ps[j].pid, ps[j].ppid, ps[j].name, ps[j].state, ps[j].memory, ps[j].threads, ps[j].files);
    }
    printf("+-------+--------+------------------------------------------+------------+-----------+----------+------------+\n");
}

void ssTable(FILE* file, struct entry* ps, int c)
{
    rewind(file);
    char header[115] = "+-------+--------+------------------------------------------+------------+----------+----------+------------+\n";
    char header2[115]= "|  PID  | Parent |                   Name                   |   State    |  Memory  | #Threads | Open Files |\n";
    fwrite(header, strlen(header), 1, file);
    fwrite(header2, strlen(header2),1, file);
    fwrite(header, strlen(header), 1, file);
    for (int j = 0; j < c; j++) 
    {
        char* temp = calloc(200, sizeof(char));
        sprintf(temp, "| %-5s | %-6s | %-40s | %-10s | %-6ld K | %-8d | %-10d |\n", ps[j].pid, ps[j].ppid, ps[j].name, ps[j].state, ps[j].memory, ps[j].threads, ps[j].files);
        fwrite(temp, strlen(temp), 1, file);
        free(temp);
    }
    fwrite(header, strlen(header), 1, file);
}

void clear() 
{
    printf("\e[1;1H\e[2J"); 
}