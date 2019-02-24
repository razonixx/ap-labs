lab2.3

Description: 
My own implementation of a cross referencer that recieves a textfile and returns the words, along with the lines where the word is found. This program does not differentiate between letter case and excludes the following words/characters: this, is, of, and, a, in, with and the empty string; and punctuation characters. (., ,, !, (, ), -, ", :)

Compilation:
I've included a Makefile, run make and then run ./cross-ref.o, sending as parameter the .txt file you want to analyze. The included example files are: irving-little-573.txt and irving-london-598.txt
