#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char const *argv[])
{
    FILE * fp;
    char c;
    int i = 0;
    int j = 0;
    int hasQuote = 0;
    fp = fopen ("test.c", "r+");
    //fp = fopen ("lab01.c", "r+");
    
    //Get the number of chars in file
    while(1)
    {
        c = fgetc(fp);
        if(feof(fp))
        {
            break;
        }
        i++;
    }
    fseek(fp, 0, SEEK_SET);
    char *buffer = (char *)malloc(sizeof(char)*i); 
    //Fill buffer of size i with i chars
    while(1)
    {
        c = fgetc(fp);
        if(feof(fp))
        {
            break;
        }
        buffer[j] = c;
        j++;
    }
    //Manipulate file now that we have it in array
    for(int k = 0; k < i; k++)
    {
        if(buffer[k] == 34)
        {
            hasQuote = !hasQuote;
        }
        else
        {
            if(((buffer[k] == '/') && (buffer[k+1] == '/')) && !hasQuote)
            {
                while(1)
                {
                    fseek(fp, k, SEEK_SET);
                    fputc(' ', fp);
                    if(buffer[k] == 10)
                    {
                        fseek(fp, k, SEEK_SET);
                        fputc(' ', fp);
                        break;
                    }
                    k++;
                }
            }
            if(((buffer[k] == '/') && (buffer[k+1] == '*')) && !hasQuote)
            {
                // Continue reading until */ is found
                while(1)
                {
                    fseek(fp, k, SEEK_SET);
                    fputc(' ', fp);
                    if((buffer[k-1] == '*') && (buffer[k] == '/'))
                    {
                        break;
                    }
                    k++;
                }
            }
        }
    }
    fclose(fp);
    return(0);
}