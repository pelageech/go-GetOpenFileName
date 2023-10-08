package main

import (
	"fmt"
	"log"
	"ofn"
	"syscall"
	"unicode/utf16"
)

func main() {
	if err := ofn.Init(); err != nil {
		log.Fatalln(err)
	}
	defer ofn.Release()
	s := "PNG file (*.png)\u0000*.PNG\u0000\u0000"
	lpcstr := utf16.Encode([]rune(s))

	filePath := make([]uint16, 256)

	if b := ofn.ChooseFileSimple(&lpcstr[0], 0, filePath); !b {
		log.Fatalln("GetOpenFileName operation unsuccessful")
	}

	fmt.Println("Opened file: ", syscall.UTF16ToString(filePath))
}
