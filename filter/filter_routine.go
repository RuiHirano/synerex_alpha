package filter

func routine_cal(json map[string]interface{},score map[string]uint64 , threshold map[string]map[string]uint64) map[string]interface{}{
	//is_path_location := false
	//is_path_company := false
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
	go func(){
		if node_TS > location_TS {
			//is_path_location = true
		}else if node_GT & location_GT == location_GT {
			//is_path_location = true
		}else if node_PT & location_PT == location_PT{
			//is_path_location = true
		}else {
			//is_path_location = false
		}
	}()

	func(){
		// company
		if node_TS > company_TS {
			//is_path_company = true
		}else if node_GT & company_GT == company_GT {
			//is_path_company = true
		}else if node_PT & company_PT == company_PT{
			//is_path_company = true
		}else {
			//is_path_company = false
		}
	}()*/

	return json
	//fmt.Printf("location: %t company: %t\n",is_path_location, is_path_company)
}


