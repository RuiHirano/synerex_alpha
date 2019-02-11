package filter

import (
	//"github.com/synerex/synerex_alpha/nodeapi"
	"log"
)
func bit_cal(json map[string]interface{},score map[string]uint64 , threshold map[string]map[string]uint64) map[string]interface{}{
	//is_path_location := false
	//is_path_company := false
	log.Printf("------------------------------------- \n")
	log.Printf("json is:  %v\n", json)
	log.Printf("score is:  %v\n", score)
	log.Printf("threshold is:  %v\n", threshold)
	log.Printf("------------------------------------- \n")

	/*node_GT := node_id.GroupTrust
	node_PT := node_id.PrivateTrust
	node_TS := node_id.TrustScore
	location_GT := exchange.Location.GroupTrust
	location_PT := exchange.Location.PrivateTrust
	location_TS := exchange.Location.TrustScore
	company_GT := exchange.Company.GroupTrust
	company_PT := exchange.Company.PrivateTrust
	company_TS := exchange.Company.TrustScore


	// location
	if node_TS > location_TS {
		//is_path_location = true
	}else if node_GT & location_GT == location_GT {
		//is_path_location = true
	}else if node_PT & location_PT == location_PT{
		//is_path_location = true
	}else {
		//is_path_location = false
	}

	// company
	if node_TS > company_TS {
		//is_path_company = true
	}else if node_GT & company_GT == company_GT {
		//is_path_company = true
	}else if node_PT & company_PT == company_PT{
		//is_path_company = true
	}else {
		//is_path_company = false
	}*/

	//fmt.Printf("location: %t company: %t\n",is_path_location, is_path_company)
	return json
}