package filter

import (
	"encoding/json"
	"fmt"
	"log"
	"math/bits"
	"reflect"
)

func if_cal(allJson AllJson, score TrustInfo , allThreshold AllThreshold, isUser bool) string{

	node_GT := score.GroupScore
	node_PT := score.PrivateScore
	node_TS := score.TrustScore
	if isUser{
		threshold := allThreshold.UT
		argJson := allJson.UJ
		log.Printf("------------------------------------- \n")
		log.Printf("json is:  %v\n", argJson)
		log.Printf("score is:  %v\n", score)
		log.Printf("threshold is:  %v\n", threshold)
		log.Printf("------------------------------------- \n")
		t := reflect.TypeOf(threshold)
		v := reflect.ValueOf(threshold) //Value
		vJson := reflect.ValueOf(argJson)
		for i := 0; i < t.NumField(); i++ {
			fmt.Printf("Name=%s , tag=%d, score=%d\n", t.Field(i).Name, v.Field(i), score)
			threshold_GT := v.Field(i).Field(2).Uint()
			threshold_PT := v.Field(i).Field(1).Uint()
			threshold_TS := v.Field(i).Field(0).Uint()
			// nodeのトラストスコア、nodeのbitの長さがlocationのbitの長さよりも短い
			if node_TS < threshold_TS || bits.Len64(node_PT) < bits.Len64(threshold_PT) || bits.Len64(node_GT) < bits.Len64(threshold_GT) {
				vJson.Field(i)
				log.Printf("Not Pass!: \n")
			}else {
				for k := 0; k < bits.Len64(threshold_PT); k++ {
					// 最右部の1までの距離を比較
					if bits.TrailingZeros64(node_PT) > bits.TrailingZeros64(threshold_PT) {
						vJson.Field(i)
						log.Printf("Not Pass!: \n")
						break
					}
					//右へシフトする
					threshold_PT = bits.RotateLeft64(threshold_PT, -1)
					node_PT = bits.RotateLeft64(node_PT, -1)
				}
				for k := 0; k < bits.Len64(threshold_GT); k++ {
					// 最右部の1までの距離を比較
					if bits.TrailingZeros64(node_GT) > bits.TrailingZeros64(threshold_GT) {
						vJson.Field(i)
						log.Printf("Not Pass!: \n")
						break
					}
					//右へシフトする
					threshold_GT = bits.RotateLeft64(threshold_GT, -1)
					node_GT = bits.RotateLeft64(node_GT, -1)
				}
			}
		}

		//json化
		argJson2, _ := json.Marshal(&argJson)
		log.Printf("json2 is:  %v\n", argJson)
		log.Printf("json2 is:  %v\n", string(argJson2))
		return string(argJson2)

	}else {
		threshold := allThreshold.TT
		argJson := allJson.TJ
		log.Printf("------------------------------------- \n")
		log.Printf("json is:  %v\n", argJson)
		log.Printf("score is:  %v\n", score)
		log.Printf("threshold is:  %v\n", threshold)
		log.Printf("------------------------------------- \n")
		t := reflect.TypeOf(threshold)
		v := reflect.ValueOf(threshold) //Value
		vJson := reflect.ValueOf(argJson)
		for i := 0; i < t.NumField(); i++ {
			fmt.Printf("Name=%s , tag=%d, score=%d\n", t.Field(i).Name, v.Field(i), score)
			threshold_GT := v.Field(i).Field(2).Uint()
			threshold_PT := v.Field(i).Field(1).Uint()
			threshold_TS := v.Field(i).Field(0).Uint()

			// nodeのトラストスコア、nodeのbitの長さがlocationのbitの長さよりも短い
			if node_TS < threshold_TS || bits.Len64(node_PT) < bits.Len64(threshold_PT) || bits.Len64(node_GT) < bits.Len64(threshold_GT) {
				vJson.Field(i)
				log.Printf("Not Pass!: \n")
			}else {
				for i := 0; i < bits.Len64(threshold_PT); i++ {
					// 最右部の1までの距離を比較
					if bits.TrailingZeros64(node_PT) > bits.TrailingZeros64(threshold_PT) {
						vJson.Field(i)
						log.Printf("Not Pass!: \n")
						break
					}
					//右へシフトする
					threshold_PT = bits.RotateLeft64(threshold_PT, -1)
					node_PT = bits.RotateLeft64(node_PT, -1)
				}
				for i := 0; i < bits.Len64(threshold_GT); i++ {
					// 最右部の1までの距離を比較
					if bits.TrailingZeros64(node_GT) > bits.TrailingZeros64(threshold_GT) {
						vJson.Field(i)
						log.Printf("Not Pass!: \n")
						break
					}
					//右へシフトする
					threshold_GT = bits.RotateLeft64(threshold_GT, -1)
					node_GT = bits.RotateLeft64(node_GT, -1)
				}
			}
		}

		//json化
		argJson2, _ := json.Marshal(&argJson)
		log.Printf("json2 is:  %v\n", &argJson)
		log.Printf("json2 is:  %v\n", string(argJson2))
		return string(argJson2)
	}

}

