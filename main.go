package main

/*
#cgo pkg-config: codec2
#include <stdlib.h>
#include <codec2.h>
*/
import "C"
import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
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

	fout, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening output speech file: %s\n", os.Args[2])
		os.Exit(1)
	}

	codec2 := C.codec2_create(C.CODEC2_MODE_3200)

	nsam := int(C.codec2_samples_per_frame(codec2))
	speechSamples := make([]int16, nsam)
	compressedBytes := make([]byte, C.codec2_bytes_per_frame(codec2))

	for {
		err := binary.Read(fin, binary.LittleEndian, &speechSamples)
		if err == io.EOF {
			os.Exit(0)
		}

		C.codec2_encode(codec2, (*C.uchar)(&compressedBytes[0]), (*C.short)(&speechSamples[0]))
		C.codec2_decode(codec2, (*C.short)(&speechSamples[0]), (*C.uchar)(&compressedBytes[0]))

		err = binary.Write(fout, binary.LittleEndian, speechSamples)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
	}
}
