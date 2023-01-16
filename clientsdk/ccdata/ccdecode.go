package ccdata

// GetSecValueFinalIndexByValue GetSecValueFinalIndexByValue
func GetSecValueFinalIndexByValue(data []byte, value byte) byte {
	tmpLen := len(data)
	for index := 0; index < tmpLen; index++ {
		if data[index] == value {
			return byte(index)
		}
	}
	return byte(0)
}

// DecodeSecValueFinal DecodeSecValueFinal
func DecodeSecValueFinal(encryptRecordData []byte, saeTableFinal []byte) []byte {
	for index := 0; index < 4; index++ {
		for secIndex := 0; secIndex < 4; secIndex++ {
			recordIndex := index*4 + secIndex
			tmpOffset := index*0x400 + secIndex*0x100
			tmpIndex := GetSecValueFinalIndexByValue(saeTableFinal[tmpOffset:tmpOffset+256], encryptRecordData[recordIndex])
			encryptRecordData[recordIndex] = tmpIndex
		}
	}

	return encryptRecordData
}

// DecodeCircleShift DecodeCircleShift
func DecodeCircleShift(data []byte, offset uint32, pos uint32) []byte {
	retData := []byte{}
	retData = append(retData, data[0:]...)
	if pos == 1 {
		retData[offset+1] = data[offset+0]
		retData[offset+2] = data[offset+1]
		retData[offset+3] = data[offset+2]
		retData[offset+0] = data[offset+3]
	}

	if pos == 2 {
		retData[offset+2] = data[offset+0]
		retData[offset+0] = data[offset+2]
		retData[offset+2] = data[offset+1]
		retData[offset+1] = data[offset+3]
	}

	if pos == 3 {
		retData[offset+3] = data[offset+0]
		retData[offset+0] = data[offset+1]
		retData[offset+1] = data[offset+2]
		retData[offset+2] = data[offset+3]
	}

	return retData
}

// DecodeShiftRows DecodeShiftRows
func DecodeShiftRows(data []byte) []byte {
	retData := DecodeCircleShift(data, 4, 1)
	retData = DecodeCircleShift(retData, 8, 2)
	retData = DecodeCircleShift(retData, 12, 3)
	return retData
}
