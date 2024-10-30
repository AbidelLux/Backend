package Comandos

import (
	"Backend/Structs"
	"bytes"
	"encoding/binary"
	"os"
	"strings"

	//"time"
	"unsafe"
)

func GetFree(spr Structs.SuperBloque, pth string, t string) int64 {
	ch := '2'
	file, err := os.Open(strings.ReplaceAll(pth, "\"", ""))

	if err != nil {
		Error("MKDIR", "No se ha encontrado el disco.")
		return -1
	}
	if t == "BI" {
		file.Seek(spr.S_bm_inode_start, 0)
		for i := 0; i < int(spr.S_inodes_count); i++ {
			data := leerBytes(file, int(unsafe.Sizeof(ch)))
			buffer := bytes.NewBuffer(data)
			err_ := binary.Read(buffer, binary.BigEndian, &ch)
			if err_ != nil {
				Error("MKDIR", "Error al leer el archivo")
				return -1
			}
			if ch == '0' {
				file.Close()
				return int64(i)
			}
		}
	} else {
		file.Seek(spr.S_bm_block_start, 0)
		for i := 0; i < int(spr.S_blocks_count); i++ {
			data := leerBytes(file, int(unsafe.Sizeof(ch)))
			buffer := bytes.NewBuffer(data)
			err_ := binary.Read(buffer, binary.BigEndian, &ch)
			if err_ != nil {
				Error("MKDIR", "Error al leer el archivo")
				return -1
			}
			if ch == '0' {
				file.Close()
				return int64(i)
			}
		}
	}
	file.Close()
	return -1
}
func GetPath(path string) []string {
	var result []string
	if path == "" {
		return result
	}
	aux := strings.Split(path, "/")
	for i := 1; i < len(aux); i++ {
		result = append(result, aux[i])
	}

	return result
}
