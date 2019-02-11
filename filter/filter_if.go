package filter

func if_cal(json map[string]interface{},score map[string]uint64 , threshold map[string]map[string]uint64) map[string]interface{} {
	//is_path_location := true
	//is_path_company := true

	/*
	node_GT := node_id.GroupTrust
	node_PT := node_id.PrivateTrust
	node_TS := node_id.TrustScore
	location_GT := exchange.Location.GroupTrust
	location_PT := exchange.Location.PrivateTrust
	location_TS := exchange.Location.TrustScore
	company_GT := exchange.Company.GroupTrust
	company_PT := exchange.Company.PrivateTrust
	company_TS := exchange.Company.TrustScore

	// location
	// nodeのトラストスコア、nodeのbitの長さがlocationのbitの長さよりも短い
	if node_TS < location_TS || bits.Len64(node_PT) < bits.Len64(location_PT) || bits.Len64(node_GT) < bits.Len64(location_GT) {
		//is_path_location = false
	}else {
		for i := 0; i < bits.Len64(location_PT); i++ {
			// 最右部の1までの距離を比較
			if bits.TrailingZeros64(node_PT) > bits.TrailingZeros64(location_PT) {
				//is_path_location = false
				break
			}
			//右へシフトする
			location_PT = bits.RotateLeft64(location_PT, -1)
			node_PT = bits.RotateLeft64(node_PT, -1)
		}
		for i := 0; i < bits.Len64(location_GT); i++ {
			// 最右部の1までの距離を比較
			if bits.TrailingZeros64(node_GT) > bits.TrailingZeros64(location_GT) {
				//is_path_location = false
				break
			}
			//右へシフトする
			location_GT = bits.RotateLeft64(location_GT, -1)
			node_GT = bits.RotateLeft64(node_GT, -1)
		}
	}

	// nodeのトラストスコア、nodeのbitの長さがlocationのbitの長さよりも短い
	if node_TS < company_TS || bits.Len64(node_PT) < bits.Len64(company_PT) || bits.Len64(node_GT) < bits.Len64(company_GT) {
		//is_path_company = false
	}else {
		for i := 0; i < bits.Len64(company_PT); i++ {
			// 最右部の1までの距離を比較
			if bits.TrailingZeros64(node_PT) > bits.TrailingZeros64(company_PT) {
				//is_path_company = false
				break
			}
			//右へシフトする
			company_PT = bits.RotateLeft64(company_PT, -1)
			node_PT = bits.RotateLeft64(node_PT, -1)
		}
		for i := 0; i < bits.Len64(company_GT); i++ {
			// 最右部の1までの距離を比較
			if bits.TrailingZeros64(node_GT) > bits.TrailingZeros64(company_GT) {
				//is_path_company = false
				break
			}
			//右へシフトする
			company_GT = bits.RotateLeft64(company_GT, -1)
			node_GT = bits.RotateLeft64(node_GT, -1)
		}
	}*/


	//fmt.Printf("location: %t company: %t\n",is_path_location, is_path_company)
	return json
}
