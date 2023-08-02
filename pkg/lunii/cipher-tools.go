package lunii

import (
	"encoding/binary"
	"os"
)

var key = []byte{0x91, 0xbd, 0x7a, 0x0a, 0xa7, 0x54, 0x40, 0xa9, 0xbb, 0xd4, 0x9d, 0x6c, 0xe0, 0xdc, 0xc0, 0xe3}

func min[T ~int](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// needed for asset an index file signature
func cipherFirstBlockCommonKey(data []byte) []byte {
	cipheredData := data

	// extract first block (512 bytes)
	firstBlockLength := min(512, len(data))
	firstBlock := data[0:firstBlockLength]

	// encrypt this block
	dataInt := bytesToInt32(firstBlock, binary.LittleEndian)
	keyInt := bytesToInt32(key, binary.BigEndian)
	encryptedIntData := Btea(dataInt, min(128, len(data)/4), keyInt)
	encryptedBlock := int32sToBytes(encryptedIntData, binary.LittleEndian)

	// replace first block and return
	copy(cipheredData, encryptedBlock)
	return cipheredData
}

func cipherFileFirstBlockCommonKey(file *os.File) error {
	data := make([]byte, 512)

	_, err := file.Read(data)
	if err != nil {
		return err
	}

	block := cipherFirstBlockCommonKey(data)

	_, err = file.WriteAt(block, 0)
	if err != nil {
		return err
	}

	return nil
}

// needed for boot file generation
func decipherFirstBlockCommonKey(data []byte) []byte {
	decipheredData := data

	// extract first block (512 bytes)
	firstBlockLength := min(512, len(data))
	firstBlock := data[:firstBlockLength]

	// decrypt this block
	dataInt := bytesToInt32(firstBlock, binary.LittleEndian)
	keyInt := bytesToInt32(key, binary.BigEndian)
	encryptedIntData := Btea(dataInt, -min(128, len(data)/4), keyInt)
	decryptedBlock := int32sToBytes(encryptedIntData, binary.LittleEndian)

	// replace first block and return
	copy(decipheredData, decryptedBlock)

	return decipheredData
}

func computeSpecificKeyFromUUID(uuid []byte) []byte {
	specificKey := decipherFirstBlockCommonKey(uuid)
	reorderedSpecificKey := []byte{
		specificKey[11], specificKey[10], specificKey[9], specificKey[8],
		specificKey[15], specificKey[14], specificKey[13], specificKey[12],
		specificKey[3], specificKey[2], specificKey[1], specificKey[0],
		specificKey[7], specificKey[6], specificKey[5], specificKey[4],
	}
	return reorderedSpecificKey
}
func cipherBlockSpecificKey(data []byte) []byte {
	device, _ := GetDevice()
	dataInt := bytesToInt32(data, binary.LittleEndian)
	keyInt := bytesToInt32(device.SpecificKey, binary.BigEndian)
	encryptedIntData := Btea(dataInt, min(128, len(data)/4), keyInt)
	return int32sToBytes(encryptedIntData, binary.LittleEndian)
}
