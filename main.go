package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//go:embed VERSION
var VERSION string

func main() {
	println(VERSION)

	// Leer el contenido del conanfile.txt
	conanfile, err := os.Open("conanfile.txt")
	if err != nil {
		fmt.Println("Error al abrir el archivo conanfile.txt:", err)
		return
	}
	defer conanfile.Close()

	scanner := bufio.NewScanner(conanfile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "requires") {
			// Analizar la línea de requires para obtener las dependencias
			dependencies := strings.Fields(line)[1:]
			for _, dep := range dependencies {
				checkDependency(dep)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo conanfile.txt:", err)
	}
}

func checkDependency(dep string) {
	// Verificar la última versión de la dependencia en conan center
	cmd := exec.Command("conan", "search", dep, "-r", "conan-center")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error al buscar la versión de %s en conan center: %v\n", dep, err)
		return
	}

	// Analizar la salida para obtener la versión más reciente
	latestVersion := parseLatestVersion(string(output))

	// Imprimir los resultados
	fmt.Printf("Dependencia: %s, Última versión en conan center: %s\n", dep, latestVersion)
}

func parseLatestVersion(output string) string {
	lines := strings.Split(output, "\n")
	if len(lines) >= 2 {
		// La última versión se encuentra en la segunda línea
		return strings.TrimSpace(lines[1])
	}
	return "No se pudo determinar la versión"
}
