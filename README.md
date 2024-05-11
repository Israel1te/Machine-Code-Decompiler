# Machine Code Disassembler

## Table of Contents

- [About](#about)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## About

This is a machine code decompiler that converts machine language code into Assembly Language
It is programmed using the Go Language and it is a single file console application that is run through the terminal to provide arguments
This project is part of the `Computer Architecture` course required projects

---------
Author : Israel Ibinayin - israelibinayin69@gmail.com
---------


## Usage

There are no additional steps or dependencies to install before using however the program must be run from the terminal

- To use, make sure there is an available input file, preferable with the .txt extention
- Open a new commandline or terminal and type
    `go run decompiler.go -i input.txt -o output`
- Replace input.txt with the file containing the input
- Replace output with what you would like to name the output file.

## Output

# There will be TWO separate output files and the names are depending on the `-o` argument
-   Using the default settings, the input files will be called 
    'output_dis.txt' and 'output_sim.txt'

- The first file will contain the disassembled output and will now show the input file but with the instructions. 
- The second file will contain a simulation of what running the assembly language code would look like theoretically. 

