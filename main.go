package main

/*
#cgo pkg-config: codec2
#include <stdlib.h>
#include <codec2.h>
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("usage: %s InputRawSpeechFile OutputRawSpeechFile\n", os.Args[0])
		os.Exit(1)
	}

	fin, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input speech file: %s\n", os.Args[1])
		os.Exit(1)
	}
	defer fin.Close()

	fout, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening output speech file: %s\n", os.Args[2])
		os.Exit(1)
	}
	defer fout.Close()

	codec2 := C.codec2_create(C.CODEC2_MODE_3200)
	defer C.codec2_destroy(codec2)

	nsam := int(C.codec2_samples_per_frame(codec2))
	speechSamples := make([]C.short, nsam)
	compressedBytes := make([]byte, C.codec2_bytes_per_frame(codec2))

	for {
		n, err := fin.Read((*(*[1 << 30]byte)(unsafe.Pointer(&speechSamples[0])))[:nsam*2])
		if err != nil || n != nsam*2 {
			break
		}

		C.codec2_encode(codec2, (*C.uchar)(&compressedBytes[0]), &speechSamples[0])
		C.codec2_decode(codec2, &speechSamples[0], (*C.uchar)(&compressedBytes[0]))

		fout.Write((*(*[1 << 30]byte)(unsafe.Pointer(&speechSamples[0])))[:nsam*2])
	}
}
