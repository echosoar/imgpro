package processor

import (
	"encoding/binary"
	"hash/crc32"
)

// apngFrame holds the parsed metadata and compressed image data for a single APNG frame.
type apngFrame struct {
	width   int
	height  int
	xOffset int
	yOffset int
	// idatData is the raw compressed pixel data (equivalent to PNG IDAT content).
	idatData []byte
}

// parseAPNGFrames parses an APNG file and returns each animation frame together
// with the raw IHDR chunk bytes and any ancillary chunks (PLTE, tRNS, etc.)
// required to reconstruct valid per-frame PNG images.
func parseAPNGFrames(fileBytes []byte) (frames []apngFrame, ihdrChunk []byte, ancillaryChunks []byte) {
	offset := 8 // skip PNG signature
	var current *apngFrame
	firstFctlSeen := false

	for offset+12 <= len(fileBytes) {
		chunkLen := int(binary.BigEndian.Uint32(fileBytes[offset : offset+4]))
		if offset+12+chunkLen > len(fileBytes) {
			break
		}
		chunkType := string(fileBytes[offset+4 : offset+8])
		chunkData := fileBytes[offset+8 : offset+8+chunkLen]
		rawChunk := fileBytes[offset : offset+12+chunkLen]

		switch chunkType {
		case "IHDR":
			ihdrChunk = rawChunk
		case "PLTE", "tRNS", "gAMA", "cHRM", "sRGB", "iCCP", "bKGD", "sBIT":
			ancillaryChunks = append(ancillaryChunks, rawChunk...)
		case "fcTL":
			if current != nil {
				frames = append(frames, *current)
			}
			// fcTL layout: seq(4) + width(4) + height(4) + x_offset(4) + y_offset(4) + ...
			f := apngFrame{
				width:   int(binary.BigEndian.Uint32(chunkData[4:8])),
				height:  int(binary.BigEndian.Uint32(chunkData[8:12])),
				xOffset: int(binary.BigEndian.Uint32(chunkData[12:16])),
				yOffset: int(binary.BigEndian.Uint32(chunkData[16:20])),
			}
			current = &f
			firstFctlSeen = true
		case "IDAT":
			// When the first fcTL precedes the first IDAT, those IDAT chunks
			// belong to the first animation frame.
			if firstFctlSeen && current != nil {
				current.idatData = append(current.idatData, chunkData...)
			}
		case "fdAT":
			// fdAT layout: seq(4) + compressed_image_data
			if current != nil && len(chunkData) > 4 {
				current.idatData = append(current.idatData, chunkData[4:]...)
			}
		case "IEND":
			if current != nil {
				frames = append(frames, *current)
				current = nil
			}
		}

		offset += 12 + chunkLen
	}
	return
}

// buildAPNGFramePNG reconstructs a minimal, self-contained PNG file for the
// given frame so that it can be decoded with image.Decode.
func buildAPNGFramePNG(frame apngFrame, ihdrChunk []byte, ancillaryChunks []byte) []byte {
	pngSig := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	// Build a new IHDR using the frame's own width/height but the original
	// bit-depth, colour-type, compression, filter and interlace values.
	originalIHDRData := ihdrChunk[8 : len(ihdrChunk)-4] // strip length, type and CRC
	newIHDRData := make([]byte, 13)
	binary.BigEndian.PutUint32(newIHDRData[0:4], uint32(frame.width))
	binary.BigEndian.PutUint32(newIHDRData[4:8], uint32(frame.height))
	copy(newIHDRData[8:13], originalIHDRData[8:13])

	var buf []byte
	buf = append(buf, pngSig...)
	buf = append(buf, buildPNGChunk("IHDR", newIHDRData)...)
	buf = append(buf, ancillaryChunks...)
	buf = append(buf, buildPNGChunk("IDAT", frame.idatData)...)
	buf = append(buf, buildPNGChunk("IEND", nil)...)
	return buf
}

// buildPNGChunk creates a complete PNG chunk (length + type + data + CRC).
func buildPNGChunk(chunkType string, data []byte) []byte {
	buf := make([]byte, 4+4+len(data)+4)
	binary.BigEndian.PutUint32(buf[0:4], uint32(len(data)))
	copy(buf[4:8], chunkType)
	copy(buf[8:], data)
	c := crc32.NewIEEE()
	c.Write(buf[4 : 8+len(data)])
	binary.BigEndian.PutUint32(buf[8+len(data):], c.Sum32())
	return buf
}
