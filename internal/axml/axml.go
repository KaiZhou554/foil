// Package axml provides manipulation of Android Binary XML (AXML) files.
package axml

import (
	"encoding/binary"
	"fmt"
)

// Document represents a parsed binary AXML document.
type Document struct {
	raw   []byte
	pool  *stringPool
}

type stringPool struct {
	chunkOffset int
	count       uint32
	flags       uint32
	dataOffset  int // byte offset from file start to string data area
	headerSize  uint16
	entries     []poolEntry
}

type poolEntry struct {
	index  uint32 // index in the pool
	charCount uint16
	offset int    // byte offset from dataOffset to the character data (after charCount prefix)
}

// Parse parses a raw binary AXML file.
func Parse(data []byte) (*Document, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("axml: data too small")
	}
	chunkType := binary.LittleEndian.Uint16(data[0:2])
	if chunkType != 0x0003 {
		return nil, fmt.Errorf("axml: not an AXML file (type=0x%04x)", chunkType)
	}

	doc := &Document{raw: make([]byte, len(data))}
	copy(doc.raw, data)

	pos := uint32(8)
	for pos+8 <= uint32(len(data)) {
		ct := binary.LittleEndian.Uint16(data[pos:])
		cs := binary.LittleEndian.Uint32(data[pos+4:])
		if cs == 0 {
			break
		}
		if ct == 0x0001 {
			pool, err := parseStringPool(data, int(pos))
			if err != nil {
				return nil, fmt.Errorf("axml: string pool: %w", err)
			}
			doc.pool = pool
			return doc, nil
		}
		pos += cs
	}
	return nil, fmt.Errorf("axml: no string pool found")
}

func parseStringPool(fileData []byte, chunkOff int) (*stringPool, error) {
	chunk := fileData[chunkOff:]
	if len(chunk) < 28 {
		return nil, fmt.Errorf("chunk too small")
	}

	sp := &stringPool{
		chunkOffset: chunkOff,
		headerSize:  binary.LittleEndian.Uint16(chunk[2:4]),
		count:       binary.LittleEndian.Uint32(chunk[8:12]),
		flags:       binary.LittleEndian.Uint32(chunk[16:20]),
	}

	stringStart := binary.LittleEndian.Uint32(chunk[20:24])
	sp.dataOffset = chunkOff + int(stringStart)

	sp.entries = make([]poolEntry, sp.count)
	for i := uint32(0); i < sp.count; i++ {
		// offset from stringStart to the start of this string entry
		strOff := binary.LittleEndian.Uint32(chunk[28+i*4:])
		entryStart := sp.dataOffset + int(strOff)

		if entryStart+2 > len(fileData) {
			continue
		}

		// First 2 bytes = character count
		charCount := binary.LittleEndian.Uint16(fileData[entryStart:])

		sp.entries[i] = poolEntry{
			index:  i,
			charCount: charCount,
			offset: entryStart + 2, // point to actual character data, skip charCount
		}
	}

	return sp, nil
}

// StringAt returns the decoded string at the given pool index.
func (d *Document) StringAt(idx int) string {
	if idx < 0 || idx >= len(d.pool.entries) {
		return ""
	}
	e := d.pool.entries[idx]
	return decodeUTF16Slice(d.raw[e.offset : e.offset+int(e.charCount)*2])
}

// FindString returns the pool index of a string, or -1 if not found.
func (d *Document) FindString(s string) int {
	for i, e := range d.pool.entries {
		decoded := decodeUTF16Slice(d.raw[e.offset : e.offset+int(e.charCount)*2])
		if decoded == s {
			return i
		}
	}
	return -1
}

// SetString replaces a string in the pool in-place.
// The new string must be exactly the same length (in characters) as the old one.
// Returns error if not found or length mismatch.
func (d *Document) SetString(oldValue, newValue string) error {
	idx := d.FindString(oldValue)
	if idx < 0 {
		return fmt.Errorf("string %q not found in pool", oldValue)
	}

	e := d.pool.entries[idx]
	oldChars := int(e.charCount)
	newChars := len([]rune(newValue))

	if newChars != oldChars {
		return fmt.Errorf("new value %q (%d chars) != old %q (%d chars)",
			newValue, newChars, oldValue, oldChars)
	}

	// Write UTF-16LE character data starting at e.offset
	// (which already skips the charCount prefix)
	writeUTF16At(d.raw[e.offset:], newValue)
	return nil
}

// SetIntAttribute finds an attribute by name in the first start tag and updates
// its typed integer value. Uses pos + attributeStart (matching AXML spec where
// attributeStart is relative to chunk start).
func (d *Document) SetIntAttribute(attrName string, value int32) error {
	attrIdx := d.FindString(attrName)
	if attrIdx < 0 {
		return fmt.Errorf("attribute %q not found in pool", attrName)
	}

	// Scan through chunks after the AXML header
	pos := uint32(8) // skip AXML header (8 bytes)
	for pos+8 < uint32(len(d.raw)) {
		ct := binary.LittleEndian.Uint16(d.raw[pos:])
		cs := binary.LittleEndian.Uint32(d.raw[pos+4:])
		if cs == 0 {
			break
		}

		if ct != 0x0102 {
			pos += cs
			continue
		}

		// Start tag (0x0102) layout (header includes NS+Name = 16 bytes):
		//   pos+0: type(2) + headerSize(2) + chunkSize(4) = 8
		//   pos+8: ns(4) + name(4) = 8
		//   pos+16: attributeStart(2) + attributeSize(2) + attributeCount(2)
		// Attributes start at pos + attributeStart

		attrStart := binary.LittleEndian.Uint16(d.raw[pos+16:])
		attrSize := binary.LittleEndian.Uint16(d.raw[pos+18:])
		attrCount := binary.LittleEndian.Uint16(d.raw[pos+20:])

		for a := uint32(0); a < uint32(attrCount); a++ {
			attrOff := pos + uint32(attrStart) + a*uint32(attrSize)
			if int(attrOff+20) > len(d.raw) {
				continue
			}
			nameIdx := binary.LittleEndian.Uint32(d.raw[attrOff+4:])
			if nameIdx == uint32(attrIdx) {
				// TypedValue at attrOff + 12:
				//   Size(2) + Res0(1) + DataType(1) + Data(4) = 8 bytes
				d.raw[attrOff+15] = 0x10 // TYPE_INT_DEC
				binary.LittleEndian.PutUint32(d.raw[attrOff+16:], uint32(value))
				return nil
			}
		}
		break // only first start tag
	}
	return fmt.Errorf("attribute %q not found", attrName)
}

// Raw returns the modified binary AXML data.
func (d *Document) Raw() []byte { return d.raw }

// ── Helpers ────────────────────────────────────────────────────────────────

func decodeUTF16Slice(b []byte) string {
	n := len(b) / 2
	chars := make([]rune, n)
	for i := 0; i < n; i++ {
		chars[i] = rune(binary.LittleEndian.Uint16(b[i*2:]))
	}
	return string(chars)
}

func writeUTF16At(dst []byte, s string) {
	for i, r := range s {
		if i*2+1 < len(dst) {
			dst[i*2] = byte(r)
			dst[i*2+1] = byte(r >> 8)
		}
	}
}
