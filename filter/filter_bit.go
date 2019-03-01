package filter

/*
func bit_cal(allJson AllJson, score TrustInfo , allThreshold AllThreshold, isUser bool) string{

	node_GT := score.GroupScore
	node_PT := score.PrivateScore
	node_TS := score.TrustScore
	if isUser{
		threshold := allThreshold.UT
		argJson := allJson.UJ
		//log.Printf("------------------------------------- \n")
		//log.Printf("json is:  %v\n", argJson)
		//log.Printf("score is:  %v\n", score)
		//log.Printf("threshold is:  %v\n", threshold)
		//log.Printf("------------------------------------- \n")
		t := reflect.TypeOf(threshold)
		v := reflect.ValueOf(threshold) //Value
		//vJson := reflect.ValueOf(argJson)
		for i := 0; i < t.NumField(); i++ {
			//fmt.Printf("Name=%s , tag=%d, score=%d\n", t.Field(i).Name, v.Field(i), score)
			threshold_GT := v.Field(i).Field(2).Uint()
			threshold_PT := v.Field(i).Field(1).Uint()
			threshold_TS := v.Field(i).Field(0).Uint()

			// location
			if  node_GT & threshold_GT == threshold_GT  || (node_TS > threshold_TS && node_PT & threshold_PT == threshold_PT){
				//log.Printf("Pass!: \n")
			}else{
				//log.Printf("threshold is:  %v\n", vJson.Field(i))
				//log.Printf("Not Pass!: \n")
			}
		}
		//json化
		argJson2, _ := json.Marshal(&argJson)
		//log.Printf("json2 is:  %v\n", &argJson)
		//log.Printf("json2 is:  %v\n", string(argJson2))
		return string(argJson2)

	}else {
		threshold := allThreshold.TT
		argJson := allJson.TJ
		//log.Printf("------------------------------------- \n")
		//log.Printf("json is:  %v\n", argJson)
		//log.Printf("score is:  %v\n", score)
		//log.Printf("threshold is:  %v\n", threshold)
		//log.Printf("------------------------------------- \n")
		t := reflect.TypeOf(threshold)
		v := reflect.ValueOf(threshold) //Value
		//vJson := reflect.ValueOf(argJson)
		for i := 0; i < t.NumField(); i++ {
			//fmt.Printf("Name=%s , tag=%d, score=%d\n", t.Field(i).Name, v.Field(i), score)
			threshold_GT := v.Field(i).Field(2).Uint()
			threshold_PT := v.Field(i).Field(1).Uint()
			threshold_TS := v.Field(i).Field(0).Uint()

			// location
			if  node_GT & threshold_GT == threshold_GT  || (node_TS > threshold_TS && node_PT & threshold_PT == threshold_PT){
				//log.Printf("Pass!: \n")
			}else{
				//vJson.Field(i).SetBool(false)
				//log.Printf("threshold is:  %v\n", vJson.Field(i))
				//log.Printf("Not Pass!: \n")
			}
		}
		//json化
		argJson2, _ := json.Marshal(&argJson)
		//log.Printf("json2 is:  %v\n", &argJson)
		//log.Printf("json2 is:  %v\n", string(argJson2))
		return string(argJson2)
	}
}

}*/