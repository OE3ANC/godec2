# Godec2

Testing codec2 with cgo (c2demo from codec 2 examples in go)

Build and install https://github.com/drowe67/codec2

Run main go and use the included "test.raw" (Signed 16 bit Little Endian, Rate 8000 Hz, Mono):
```bash
go run ./main.go test.raw test.raw.out
aplay -f S16_LE test.raw.out
```
