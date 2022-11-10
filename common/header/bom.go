package header

import (
	"io"
)

func ReadBom(r io.Reader) []byte {
	var bom = make([]byte, 4)
	_, err := r.Read(bom)
	if err != nil && err != io.EOF {
		return nil
	}

	if len(bom) >= 1 && bom[0] == 0x00 &&
		len(bom) >= 2 && bom[1] == 0x00 &&
		len(bom) >= 3 && bom[2] == 0xFE &&
		len(bom) >= 4 && bom[3] == 0xFF {
		return bom
	}

	if len(bom) >= 1 && bom[0] == 0xFF &&
		len(bom) >= 2 && bom[1] == 0xFE &&
		len(bom) >= 3 && bom[2] == 0x00 &&
		len(bom) >= 4 && bom[3] == 0x00 {
		return bom
	}

	if len(bom) >= 1 && bom[0] == 0xEF &&
		len(bom) >= 2 && bom[1] == 0xBB &&
		len(bom) >= 3 && bom[2] == 0xBF {
		return bom[:3]
	}

	if len(bom) >= 1 && bom[0] == 0xFE &&
		len(bom) >= 2 && bom[1] == 0xFF {
		return bom[:2]
	}

	if len(bom) >= 1 && bom[0] == 0xFF &&
		len(bom) >= 2 && bom[1] == 0xFE {
		return bom[:2]
	}

	return nil
}
