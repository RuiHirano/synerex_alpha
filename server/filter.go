package main

import (
	"math/bits"
)

func bit_cal() bool {
	infoNum := 1000000
	isPass := false
	groupBit := uint64(63)
	for i:=0; i<infoNum; i++{
		isPass = false
		infoBit := uint64(72)
		if infoBit & groupBit == infoBit{
			isPass = true
		}
	}

	return isPass
}

func if_cal() bool{

	infoNum := 1000000
	isPass := false
	groupBit := uint64(64)
	for i:=0; i<infoNum; i++{
		isPass = false
		infoBit := uint64(75)
		for k := 0; k < int(infoBit); k++ {
			// 最右部の1までの距離を比較
			if bits.TrailingZeros64(uint64(infoBit))==0{
				if bits.TrailingZeros64(uint64(groupBit))==0{
					isPass = true
				}
			}
			//右へシフトする
			infoBit = bits.RotateLeft64(uint64(infoBit), -1)
			groupBit = bits.RotateLeft64(uint64(groupBit), -1)
		}
	}

	return isPass
}
