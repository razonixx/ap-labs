#define _XOPEN_SOURCE 500
#include <stdio.h>
#include <sys/inotify.h>
#include <ftw.h>
#include <limits.h>

#include "logger.h"

void displayInotifyEvent(struct inotify_event *i);
int nftwFunc(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf);

int inotifyFd;
char* path;
int flags = 0;
int cookie1 = 0, cookie2 = 0;

int main(int argc, char* argv[])
{
    if (argc < 3)
    {
        errorf("Usage: %s <log> <path>\nlog: Type of log desired, 'syslog' for syslog and 'stdout' for stdout\npath: path of directory(s) to watch\n ", argv[0]);
        return -1;
    }

    initLogger(argv[1]);

    path = argv[2];

    ssize_t numRead;
    char *chPointer;
    int BUF_LEN = (10 * (sizeof(struct inotify_event) + NAME_MAX + 1));
    char buffer[BUF_LEN];
    struct inotify_event *event;
    
    if (strchr(argv[2], 'd') != NULL)
        flags |= FTW_DEPTH;
    if (strchr(argv[2], 'p') != NULL)
        flags |= FTW_PHYS;


    inotifyFd = inotify_init();                 /* Create inotify instance */
    if (inotifyFd == -1)
    {
        panicf("inotify_init");
    }

    if (nftw(argv[2], nftwFunc, 20, flags) == -1)
    {
        panicf("nftw");
    }
    printf("\n");

    while(1)
    {
        numRead = read(inotifyFd, buffer, BUF_LEN);
        if (numRead == 0)
            panicf("read() from inotify file descriptor returned 0!");

        if (numRead == -1)
            panicf("read");


        /* Process all of the events in buffer returned by read() */

        for (chPointer = buffer; chPointer < buffer + numRead; ) {
            event = (struct inotify_event *) chPointer;
            displayInotifyEvent(event);
            printf("\n");
            chPointer += sizeof(struct inotify_event) + event->len;
        }
    }

    return 0;
}

void displayInotifyEvent(struct inotify_event *i)
{
    if (i->len > 0)
    {
        infof("Name: %s, Watch Descriptor: %d;\n", i->name, i->wd);
    }
    else
    {
        infof("Watch Descriptor: %d;\n", i->wd);
    }

    if (i->mask & IN_CREATE)        infof("IN_CREATE\n");
    if (i->mask & IN_DELETE)        infof("IN_DELETE\n"); 
    if (i->mask & IN_DELETE_SELF)   infof("IN_DELETE_SELF\n");
    if (i->mask & IN_MOVED_FROM)    infof("IN_MOVED_FROM\n");
    if (i->mask & IN_MOVED_TO)      infof("IN_MOVED_TO\n");

    if ((i->mask & IN_CREATE) && (i->mask & IN_ISDIR))
    {
        warnf("Adding %s to inotify FD.\n\n", i->name);
        nftw(path, nftwFunc, 20, flags); 
    }
    if (i->mask & IN_DELETE_SELF)
    {
        warnf("Directory with watch descriptor %d was removed from inotify FD.\n", i->wd);
        inotify_rm_watch(inotifyFd, i->wd);
    }
    if (i->cookie > 0)
    {
        if(cookie1 == 0)
        {
            cookie1 = i->cookie;
        }
        else if(cookie2 == 0)
        {
            cookie2 = i->cookie;
        }
        if(cookie1 != i->cookie && cookie2 != i->cookie)
        {
            cookie1=0;
            cookie2=0;
        }
        if(cookie1 == cookie2)
        {
            warnf("File was renamed\n");
        }
    }
}

int nftwFunc(const char *fpath, const struct stat *sb, int tflag, struct FTW *ftwbuf)
{
    if(tflag == FTW_D || tflag == FTW_DP)
    {
        int wd = inotify_add_watch(inotifyFd, fpath, IN_CREATE|IN_DELETE|IN_DELETE_SELF|IN_MOVED_FROM|IN_MOVED_TO);
        if (wd == -1)
        {
            panicf("inotify_add_watch");
        }
        infof("Watching %s using Watch Descriptor %d\n", fpath, wd);
    }
    return 0;
}