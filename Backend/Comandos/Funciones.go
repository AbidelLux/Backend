package Comandos

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"
)

func Comparar_Cadenas(a string, b string) bool {
	if strings.ToUpper(a) == strings.ToUpper(b) { //compara las dos cadenas, se pone en mayuscula las dos
		return true

	}
	return false
}

func Error(op string, mensaje string) string {
	fmt.Println("\tERROR: " + op + "\n\tTIPO: " + mensaje)
	return "\tERROR: " + op + "\n\tTIPO: " + mensaje + "\n"
}

func Mensaje(op string, mensaje string) string {
	fmt.Println("\tCOMANDO: " + op + "\n\tTIPO: " + mensaje)
	return "\tCOMANDO: " + op + "\n\tTIPO: " + mensaje + "\n"
}

func Confirmar(mensaje string) bool {
	fmt.Println(mensaje + " (y/n)")
	var respuesta string
	fmt.Scanln(&respuesta)
	return Comparar_Cadenas(respuesta, "y")
}

func ArchivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}

func EscribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func leerDisco(path string) *Structs.MBR {
	m := Structs.MBR{}

	if _, err := os.Stat(strings.ReplaceAll(path, "\"", "")); os.IsNotExist(err) {
		// Si el archivo no existe, mostrar un mensaje de error y devolver nil
		Error("FDISK", "El archivo especificado no existe")
		return nil
	}

	file, err := os.Open(strings.ReplaceAll(path, "\"", ""))
	defer file.Close()
	if err != nil {
		Error("FDISK", "Error al abrir el archivo")
		return nil
	}
	file.Seek(0, 0)
	data := leerBytes(file, int(unsafe.Sizeof(Structs.MBR{})))
	buffer := bytes.NewBuffer(data)
	err_ := binary.Read(buffer, binary.BigEndian, &m)
	if err_ != nil {
		Error("FDSIK", "Error al leer el archivo")
		return nil
	}
	var mDir *Structs.MBR = &m
	return mDir
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes

	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
func ruta() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Error al obtener la ruta del archivo actual")
		return ""
	}

	// Obtenemos la ruta del directorio padre del archivo actual
	dirpre := filepath.Dir(filename)
	dir := filepath.Join(dirpre, "..")
	// Concatenamos la ruta del directorio padre con el nombre de la carpeta "Discos"
	discosDir := filepath.Join(dir, "Discos")

	// Verificamos si la carpeta "Discos" existe
	if _, err := os.Stat(discosDir); os.IsNotExist(err) {
		fmt.Println("La carpeta 'Discos' no existe en la ruta especificada.")
		return ""
	}
	//efmt.Println(discosDir)
	return discosDir

}

func verificarArchivos(discosDir string) string {
	archivos := []string{"A.dsk", "B.dsk", "C.dsk", "D.dsk", "E.dsk", "F.dsk", "G.dsk", "H.dsk", "I.dsk", "J.dsk", "K.dsk", "L.dsk", "M.dsk", "N.dsk", "O.dsk", "P.dsk", "Q.dsk", "R.dsk", "S.dsk", "T.dsk", "U.dsk", "V.dsk", "W.dsk", "X.dsk", "Y.dsk", "Z.dsk"}

	// Itera sobre cada archivo en la lista de archivos y verifica si existe dentro de discosDir.
	for _, archivo := range archivos {
		rutaArchivo := filepath.Join(discosDir, archivo)
		if _, err := os.Stat(rutaArchivo); err == nil {
			fmt.Print("")

		} else if os.IsNotExist(err) {
			return archivo
		} else {
			fmt.Printf("Error al verificar el archivo %s: %v\n", archivo, err)
		}
	}
	return "Z.dk"
}
