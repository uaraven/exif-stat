package exif

func nikonV3Detector(data []byte) bool {
	header := []byte{'N', 'i', 'k', 'o', 'n', 0x00, 0x02, 0x10, 0x00, 0x00}
	for i, v := range header {
		if v != data[i] {
			return false
		}
	}
	return true
}

func nikonV3VariantDetector(data []byte) bool {
	header := []byte{'N', 'i', 'k', 'o', 'n', 0x00, 0x02, 0x00, 0x00, 0x00}
	for i, v := range header {
		if v != data[i] {
			return false
		}
	}
	return true
}

func nikonV3Reader(file File, entry ifdEntry) (*ifd, error) {
	mainTiffHeaderOffset := file.GetTiffHeaderOffset()
	mainOrder := file.GetOrder()
	defer func() {
		file.SetTiffHeaderOffset(mainTiffHeaderOffset)
		file.SetOrder(mainOrder)
	}()

	offset := file.GetTiffHeaderOffset() + int64(entry.Data) + 10 // 10 bytes of nikon signature
	_, err := file.seek(offset)
	if err != nil {
		return nil, err
	}
	file.SetTiffHeaderOffset(offset)
	file.SetOrder(BigEndian)
	err = readTiffHeader(file)
	if err != nil {
		return nil, err
	}
	return readIfd(file, -1, entry.IfdIndex)
}

type makerNoteReader struct {
	CanRead func([]byte) bool
	Reader  func(File, ifdEntry) (*ifd, error)
}

var makerNoteReaders = []makerNoteReader{
	{nikonV3Detector, nikonV3Reader},
	{nikonV3VariantDetector, nikonV3Reader},
}

func readMakerNotes(file File, entry ifdEntry) (*ifd, error) {
	for _, reader := range makerNoteReaders {
		if reader.CanRead(entry.ValueBytes) {
			return reader.Reader(file, entry)
		}
	}
	return nil, nil
}
