#include <stdio.h>
#include <stdlib.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/stat.h>
#include <string.h>
#include <ctype.h>

struct nlist { /* table entry: */
    struct nlist *next; /* next entry in chain */
    char *name; /* defined name */
    char *defn; /* replacement text */
};

#define HASHSIZE 1000000
#define BUFF_SIZE 1000000
static struct nlist *hashtab[HASHSIZE]; /* pointer table */

unsigned hash(char *s);
struct nlist *lookup(char *s);
char *myStrdup(char *);
struct nlist *install(char *name, char *defn);
void push(struct nlist *head, char *name, char *defn);
void printNode(struct nlist *node);
void printList(struct nlist *head);
void fillBuffer(char *buffer, char *file);
int myStrcmp(char *origin, char *substr);
char* myStradd(char *origin, char *addition);
void remCharToLower(char *s, char chr);

int main(int argc, char **argv) 
{
    if(argc < 2){
        printf("You need to pass 1 file to analyze");
        return -1;
    }

    char buffer[100000];
    char *lines [510];
    fillBuffer(buffer, argv[1]);
    //printf(buffer);

    int lineCounter = 0;
    char *lineCounterBuffer = calloc(4, sizeof(char));
    sprintf(lineCounterBuffer, "%d", lineCounter);

    char *line = strtok(buffer, "\n");
	while(line != NULL)
	{
        lineCounter++;
        sprintf(lineCounterBuffer, "%d", lineCounter);
        lines[lineCounter] = line;
		line = strtok(NULL, "\n");
	}
    struct nlist *head = install("HEAD", "0");
    for(int i = 1; i <= lineCounter; i++)
    {
        char *lineSpace = strtok(lines[i], " ");
        while(lineSpace != NULL)
        {
            remCharToLower(lineSpace, ',');
            remCharToLower(lineSpace, '.');
            remCharToLower(lineSpace, '(');
            remCharToLower(lineSpace, ')');
            remCharToLower(lineSpace, ';');
            remCharToLower(lineSpace, '-');
            remCharToLower(lineSpace, '!');
            remCharToLower(lineSpace, '"');
            remCharToLower(lineSpace, ':');
            sprintf(lineCounterBuffer, "%d", i);
            push(head, lineSpace, lineCounterBuffer);
            //printf("Word: %s\nLine: %s\n", lineSpace, lineCounterBuffer);
            lineSpace = strtok(NULL, " ");
        }
    }

    printList(head);
}

/* hash: form hash value for string s */
unsigned hash(char *s)
{
    unsigned hashval;
    for (hashval = 0; *s != '\0'; s++)
        hashval = *s + 1300711 * hashval;
    return hashval % HASHSIZE;
}

/* lookup: look for s in hashtab */
struct nlist *lookup(char *s)
{
    struct nlist *np;
    for (np = hashtab[hash(s)]; np != NULL; np = np->next)
        if (strcmp(s, np->name) == 0)
          return np; /* found */
    
    return NULL; /* not found */
}

/* install: put (name, defn) in hashtab */
struct nlist *install(char *name, char *defn)
{
    struct nlist *np;
    unsigned hashval;
    if ((np = lookup(name)) == NULL) { /* not found */
        np = (struct nlist *) malloc(sizeof(*np));
        if (np == NULL || (np->name = myStrdup(name)) == NULL)
        {
            printf("SE QUEDO SIN MEMORIA\n");
            exit(-1);
            return NULL;
        }
        hashval = hash(name);
        np->next = hashtab[hashval];
        hashtab[hashval] = np;
    } 
    else /* already there */
    {
        char *tmp = np->defn;
        //strcat(np->defn, ", ");
        myStradd(np->defn, ",");
        //free((void *) np->defn); /*free previous defn */        
    }
    if ((np->defn = myStrdup(defn)) == NULL)
    {
        return NULL;
    }
    return np;
}

char *myStrdup(char *s) /* make a duplicate of s */
{
    char *p;
    p = (char *) malloc(strlen(s)+10000000); /* +1 for ’\0’ */
    if (p != NULL)
       strcpy(p, s);
    return p;
}

void push(struct nlist *head, char *name, char *defn) 
{
    struct nlist *tmp;
    //Old word, new Line
    if((tmp = lookup(name)) != NULL)
    {
        char *lineNumber;
        lineNumber = strcat(tmp->defn, ", ");
        tmp->defn = strcat(lineNumber, defn);
    }
    //New word
    else
    {
        struct nlist *current = head;
        while (current->next != NULL) {
            current = current->next;
        }
        /* now we can add a new variable */
        current->next = install(name, defn);
        current->next->next = NULL;
    }
}

void printNode(struct nlist *node)
{
    printf("Word: %s\nLine: %s\n", node->name, node->defn);  
}

void printList(struct nlist *head)
{
    struct nlist *current = head;
    //printNode(current);
    while (current->next != NULL) {
        current = current->next;
        if(strcmp(current->name, "the") == 0 || strcmp(current->name, "") == 0
        || strcmp(current->name, "is")  == 0 || strcmp(current->name, "a") == 0
        || strcmp(current->name, "of") == 0  || strcmp(current->name, "in") == 0
        || strcmp(current->name, "and") == 0 || strcmp(current->name, "with") == 0) //Remove wordds
        {
            
        }
        else
        {
            printNode(current);
        }
    }
}

void fillBuffer(char *buffer, char *name)
{
    int ofd = open(name, O_RDONLY); 
    if (ofd < 0) { 
        perror("Error opening file"); 
        exit(-1);
    } 
    int rf = read(ofd, buffer, BUFF_SIZE); 
    buffer[rf] = '\0'; 
    close(ofd);
}

int myStrcmp(char *origin, char *substr)
{
    int j = 0;
    for(int i = 0; i < strlen(origin); i++)
    {
        if(j != 0 && origin[i] != substr[j])
        {
            j = 0;
        }
        if(origin[i] == substr[j])    
        {
            j++;
        }
        if(j == strlen(substr))
        {
            return 1;
        }
    }
    return 0;
}

char* myStradd(char *origin, char *addition)
{
    int lenOrigin = strlen(origin);
    int lenAddition = strlen(addition);
    char *res = (char *)malloc(lenOrigin + lenAddition);
    
    int i;
    for(i = 0; i < lenOrigin; i++)
    {
        res[i] = origin[i];
    }

    for(int j = 0; j < lenAddition; j++)
    {
        res[i] = addition[j];
        i++;
    }
    return res;
}

void remCharToLower(char *s, char chr)
{
   int i, j = 0;
   for ( i = 0; s[i] != '\0'; i++ ) /* 'i' moves through all of original 's' */
   {
      s[j] = tolower(s[j]);
      if ( s[i] != chr )
      {
         s[j] = s[i]; /* 'j' only moves after we write a non-'chr' */
         j++;
      }
   }
   s[j] = '\0'; /* re-null-terminate */
}