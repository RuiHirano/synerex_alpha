package main

// multiple service provider

// combine complex demands/supply into ride-share demand.

// user demand
// taxi supply
// ad demand
//  -> select all of them then provide supply

// multi suppl to user

import (
	"context"
	"flag"
	"fmt"
	"github.com/synerex/synerex_alpha/api/fleet"
	"log"
	"strconv"
	"sync"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	idlist     []uint64
	dmMap      map[uint64]*sxutil.DemandOpts
	myPrice    int
)

func init() {
	idlist = make([]uint64, 0)
	dmMap = make(map[uint64]*sxutil.DemandOpts)
}

// callback for each Supply
func supplyCallback(clt *sxutil.SMServiceClient, sp *pb.Supply) {
	// check if supply is match with my demand.
	log.Println("Got supply callback")
	// choice is supply for me? or not.
	if clt.IsSupplyTarget(sp, idlist) { //
		// always select Supply
		clt.SelectSupply(sp)
	}
}

// callback for each Ad Demand
func adDemandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {
	// check if supply is match with my demand.
	log.Println("Got ad demand callback")

	if dm.GetDemandName() == "" { // this is Select!
		log.Println("getSelect!")
		clt.Confirm(sxutil.IDType(dm.GetId()))
	} else {
		// select any ride share demand!
		sp := &sxutil.SupplyOpts{Target: dm.GetId()} // ターゲットにDemand.Idを設定 (利用者側のSupplyチェックで使用)

		clt.ProposeSupply(sp)
	}
}

// callback for each Demand
func rideDemandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand, id sxutil.IDType) {
	// check if supply is match with my demand.
	log.Println("Got rideshare demand callback")

	// we need to ask other provider for lowest ...

	if dm.TargetId != 0 { // this is Select!
		log.Println("getSelect!")
		clt.Confirm(sxutil.IDType(dm.GetId()))
	} else {
		// select any ride share demand!

		fleet := fleet.Fleet{
			VehicleId: int32(10),
			Angle:     float32(100),
			Speed:     int32(20),
			Status:    int32(0),
			Coord: &fleet.Fleet_Coord{
				Lat: float32(34.874364),
				Lon: float32(137.1474168),
			},
		}
		sp := &sxutil.SupplyOpts{
			Target: dm.GetId(),
			Name:   "RideShare by Taxi",
			JSON:   `{"Price":` + strconv.Itoa(myPrice) + `,"Distance": 5200, "Arrival": 300, "Destination": 500, "Position":{"Latitude":36.6, "Longitude":135}}`,
			Fleet: &fleet,
		} // set TargetID as Demand.Id (User will check by them)

		clt.ProposeSupply(sp)
	}
}

func subscribeAdSupply(client *sxutil.SMServiceClient) {
	// this function should be run under goroutine!
	ctx := context.Background() // required?
	client.SubscribeSupply(ctx, supplyCallback)
	// comes here if channel closed
}

func subscribeAdDemand(client *sxutil.SMServiceClient) {
	// this function should be run under goroutine!
	ctx := context.Background() // required?
	client.SubscribeDemand(ctx, adDemandCallback)
	// comes here if channel closed
}

func subscribeRideDemand(client *sxutil.SMServiceClient) {
	// this function should be run under goroutine!
	ctx := context.Background() // required？
	client.SubscribeDemand(ctx, rideDemandCallback)
	// comes here if channel closed
}

func sendDemand(sclient *sxutil.SMServiceClient, nm string, js string) {
	opts := &sxutil.DemandOpts{Name: nm, JSON: js}
	id := sclient.RegisterDemand(opts)
	idlist = append(idlist, id)
	dmMap[id] = opts
}

//
//
func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "MultiProvider", false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	wg := sync.WaitGroup{}

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	sxutil.RegisterDeferFunction(func() { conn.Close() })

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:Multi_AD}")
	argJson2 := fmt.Sprintf("{Client:Multi_RIDE_SHARE}")

	sclient := sxutil.NewSMServiceClient(client, pb.ChannelType_AD_SERVICE, argJson) // connection for AD

	ride_sclient := sxutil.NewSMServiceClient(client, pb.ChannelType_RIDE_SHARE, argJson2) // create sclient for each type

	wg.Add(1)
	go subscribeAdDemand(sclient) // from ad company
	//	go subscribeAdSupply(sclient) // from taxi or other media
	wg.Add(1)
	go subscribeRideDemand(ride_sclient)

	//	sendDemand(sclient, "Ad to any Girls!", "{Target: \"female\", AgeFrom: 30, AgeTo:49, Spec: 0}")
	wg.Wait()

	sxutil.CallDeferFunctions() // cleanup!

}
