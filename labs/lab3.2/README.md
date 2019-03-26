Lab 3.2 - Progress Notifier with Signals
========================

Description
--------------------------
This is a progress notifier. It keeps track of the progress made when encoding and decoding a `.txt` file using the `base 64` algorithm.

Compilation
--------------------
To build the binary, use the included `Makefile` by running `make`.

The program recieves 2 arguments:
-  `--encode` or `--decode`, to designate the function to be performed on the file

- The file that will be worked on. If a normal textfile, it should be encoded, if a base64-encoded file, it should be decoded.

- If using the `--encode` parameter, a file `encoded.txt` will be created with the encoded text. This text can be decoded by passing the `encoded.txt` to the program along with the `--decode` parameter.

- 3 files have been included as test files. `aesop11.txt`, `sick-kid.txt`, and `vgilante.txt`.


Usage: 
- `./base64.out <--encode || --decode> file`

Examples:
- `./base64.out --encode aesop11.txt`
- `./base64.out --decode encoded.txt`
- `./base64.out --encode sick-kid.txt`
- `./base64.out --decode encoded.txt`
- `./base64.out --encode vgilante.txt`
- `./base64.out --decode encoded.txt`

