package main

import (
	"bufio"
	"fmt"
	"os"
	"rerepolez/errores"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		elError := new(errores.ErrorParametros)
		fmt.Fprintf(os.Stderr, "%s\n", elError.Error())
	} else {
		archivos := args[1:]
		archivoLista := abrirArchivo(archivos[0])
		archivoPadron := abrirArchivo(archivos[1])
		s1 := bufio.NewScanner(archivoLista)
		s2 := bufio.NewScanner(archivoPadron)
		for s1.Scan() {
			fmt.Fprintf(os.Stdout, "%s\n", s1.Text())
		}
		for s2.Scan() {
			fmt.Fprintf(os.Stdout, "%s\n", s2.Text())
		}
	}
}

func abrirArchivo(ruta string) *os.File {
	archivo, err := os.Open(ruta)
	if err != nil {
		errorLeerArchivo := new(errores.ErrorLeerArchivo)
		fmt.Fprintf(os.Stderr, "%s", errorLeerArchivo.Error())
	}
	return archivo
}
