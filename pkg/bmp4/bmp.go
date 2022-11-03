package bmp4

import (
	"bytes"
	"image"

	"encoding/binary"
)

func GetBitmap(source *image.Gray) []byte {

	width := source.Rect.Size().X
	height := source.Rect.Size().Y

	//////////////////////////////////////

	type encodePart struct {
		length uint32
		color  uint8
	}

	imgBuf := new(bytes.Buffer)

	// Start of Bitmap Data
	for i := 0; i < height; i++ {
		inverse := height - i
		var lineFeed encodePart

		for j := 0; j < width; j++ {
			pix := source.GrayAt(j, inverse).Y / 16
			if j == 0 {
				lineFeed = encodePart{length: 1, color: pix}
				continue
			}

			if lineFeed.color == pix && lineFeed.length < 255 {
				// If this is the same color as the last pixel, just auglent the length
				lineFeed.length++
			} else {
				// otherwise, write the code to the file an start a new count
				// Double color
				color8 := lineFeed.color*16 + lineFeed.color
				binary.Write(imgBuf, binary.LittleEndian, int8(lineFeed.length))
				binary.Write(imgBuf, binary.LittleEndian, int8(color8))

				// add a new color
				lineFeed = encodePart{length: 1, color: pix}
			}
		}

		// commit data
		// Double color
		color8 := lineFeed.color*16 + lineFeed.color
		binary.Write(imgBuf, binary.LittleEndian, int8(lineFeed.length))
		binary.Write(imgBuf, binary.LittleEndian, int8(color8))

		// end of line, but not for the last line
		if i < height-1 {
			imgBuf.Write([]byte{0x00, 0x00})
		}
	}
	// end of file
	imgBuf.Write([]byte{0x00, 0x01})

	imgBytes := imgBuf.Bytes()

	///////////////////////////////////////

	// global settings
	headerSize := 54
	dataSize := len(imgBytes)
	paletteSize := 16 * 4
	fullSize := headerSize + paletteSize + dataSize

	outBuf := new(bytes.Buffer)
	// 	0h	2	42 4D	"BM"	Magic Number (unsigned integer 66, 77)
	outBuf.Write([]byte{'B', 'M'})
	// 2h	4	46 00 00 00	70 Bytes	Size of the BMP file
	binary.Write(outBuf, binary.LittleEndian, uint32(fullSize))
	// 6h	2	00 00	Unused	Application Specific
	// 8h	2	00 00	Unused	Application Specific
	outBuf.Write(make([]byte, 4))
	// Ah	4	36 00 00 00	54 bytes	The offset where the bitmap data (pixels) can be found.
	binary.Write(outBuf, binary.LittleEndian, uint32(headerSize+paletteSize))

	// Eh	4	28 00 00 00	40 bytes	The number of bytes in the header (from this point).
	binary.Write(outBuf, binary.LittleEndian, uint32(40))
	// 12h	4	02 00 00 00	2 pixels	The width of the bitmap in pixels
	binary.Write(outBuf, binary.LittleEndian, uint32(width))
	// 16h	4	02 00 00 00	2 pixels	The height of the bitmap in pixels
	binary.Write(outBuf, binary.LittleEndian, uint32(height))
	// 1Ah	2	01 00	1 plane	Number of color planes being used.
	binary.Write(outBuf, binary.LittleEndian, uint16(1))
	// 1Ch	2	18 00	24 bits	The number of bits/pixel.
	binary.Write(outBuf, binary.LittleEndian, uint16(4))
	// 1Eh	4	00 00 00 00	2	RLE_4, Use RLE_4 compression
	binary.Write(outBuf, binary.LittleEndian, uint32(2))

	// 22h	4	10 00 00 00	16 bytes	The size of the raw BMP data (after this header)
	binary.Write(outBuf, binary.LittleEndian, uint32(dataSize))
	// 26h	4	13 0B 00 00	2,835 pixels/meter	The horizontal resolution of the image
	binary.Write(outBuf, binary.LittleEndian, uint32(0))
	// 2Ah	4	13 0B 00 00	2,835 pixels/meter	The vertical resolution of the image
	binary.Write(outBuf, binary.LittleEndian, uint32(0))
	// 2Eh	4	00 00 00 00	0 colors	Number of colors in the palette
	binary.Write(outBuf, binary.LittleEndian, uint32(0))
	// 32h	4	00 00 00 00	0 important colors	Means all colors are important
	binary.Write(outBuf, binary.LittleEndian, uint32(0))

	// Palette
	// Grayscale, 16 levels
	for i := 0; i < 16; i++ {
		step := (255 / 16) * i
		// red
		binary.Write(outBuf, binary.LittleEndian, uint8(step))
		// green
		binary.Write(outBuf, binary.LittleEndian, uint8(step))
		// blue
		binary.Write(outBuf, binary.LittleEndian, uint8(step))
		//reserved
		binary.Write(outBuf, binary.LittleEndian, uint8(0))
	}

	// get bytes
	return append(outBuf.Bytes(), imgBytes...)

}
