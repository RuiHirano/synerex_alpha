package main

// Simple Taxi Provider demo

import (
	"context"
	"flag"
	"reflect"

	//"github.com/synerex/synerex_alpha/api/fleet"
	"log"
	"sync"
	"time"

	pb "github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
	"fmt"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	price    = flag.Int("price", 100, "Taxi price")
	idlist     []uint64
	spMap      map[uint64]*sxutil.SupplyOpts
	mu		sync.Mutex
	sclient *sxutil.SMServiceClient
	stTaxiCommunication uint64
	ftTaxiCommunication uint64
	stUserCommunication uint64
	timesTaxiCommunication float64
	timesAll float64
	count int
)

func init(){
	idlist = make([]uint64, 0)
	spMap = make(map[uint64]*sxutil.SupplyOpts)
}

// callback for each Demand
//ユーザーなどほかのプロバイダが走ると呼ばれる関数
//ユーザーの情報を取得する
func demandCallback(clt *sxutil.SMServiceClient, dm *pb.Demand) {

	//ここまでが通信時間
	stUserCommunication = dm.GetSt()
	stTaxiCommunication = dm.GetStTaxi()
	ftTaxiCommunication = uint64(time.Now().UnixNano())
	timeTaxiCommunication := float64(ftTaxiCommunication - stTaxiCommunication)/1000000
	timeAll := float64(ftTaxiCommunication - stUserCommunication)/1000000
	//処理時間
	fmt.Println("-------------------------------------------------\n")
	log.Printf("time taxi communication is:  %f ms\n", timeTaxiCommunication)
	log.Printf("time all is:  %f ms\n", timeAll)
	fmt.Println("---------------------------------------------------\n")
	//全体平均
	timesTaxiCommunication += timeTaxiCommunication
	timesAll += timeAll
	count += 1
	iterNum := 100
	if count == iterNum{
		fmt.Println("平均\n")
		log.Printf("avarage taxi communication is:  %f\n", timesTaxiCommunication/float64(iterNum))
		log.Printf("avarage all is:  %f\n", timesAll/float64(iterNum))
		timesAll = 0
		timesTaxiCommunication = 0
		count = 0
	}

	//go subscribeDemand(sclient)
	/*
	// check if demand is match with my supply.
	log.Println("Got ride share demand callback")
	if dm.TargetId != 0 { // this is Select!
		log.Println("getSelect!")

		clt.Confirm(sxutil.IDType(dm.GetId()))

	}else { // not select
		// select any ride share demand!
		// should check the type of ride..

		log.Printf("Provider dm %v\n", dm.GetId())
		// create dummy fleet
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
		//id := clt.getChannel().ClientId
		sp := &sxutil.SupplyOpts{
			ID: uint64(clt.ClientID),
			Target: dm.GetId(),
			Name: "RideShare by Taxi",
			JSON: `{"Price":`+strconv.Itoa(*price)+`,"Distance": 5200, "Arrival": 300, "Destination": 500, "Position":{"Latitude":36.6, "Longitude":135}}`,
			Fleet: &fleet,
		} // set TargetID as Demand.Id (User will check by them)

		mu.Lock()
		//log.Printf("Taxi SPaa ID %v\n\n", sp.ID)
		pid := clt.ProposeSupply(sp)
		idlist = append(idlist,pid)
		spMap[pid] = sp
		mu.Unlock()
	}*/
}

func subscribeDemand(client *sxutil.SMServiceClient) {
	// goroutine!
	ctx := context.Background() //
	client.SubscribeDemand(ctx, demandCallback)
	// comes here if channel closed
	log.Printf("Server closed... on taxi provider")
}

func oldproposeSupply(client pb.SynerexClient, targetNum uint64) {
	dm := pb.Supply{Id: 200, SenderId: 555, TargetId: targetNum, Type: pb.ChannelType_RIDE_SHARE, SupplyName: "Taxi"}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.ProposeSupply(ctx, &dm)
	if err != nil {
		log.Fatalf("%v.Propose Supply err %v", client, err)
	}
	log.Println(resp)

}

func threshold() string{

	return `{"Price": {"TrustScore": 53, "PrivateScore": 23, "GroupScore": 38}, "Distance": {"TrustScore": 23, "PrivateScore": 24, "GroupScore": 42}, "Arrival": {"TrustScore": 43, "PrivateScore": 14, "GroupScore": 23}, "Destination": {"TrustScore": 62, "PrivateScore": 23, "GroupScore": 33}, "Position": {"TrustScore": 34, "PrivateScore": 43, "GroupScore": 25}}`

}

func main() {
	flag.Parse()
	sxutil.RegisterNodeName(*nodesrv, "TaxiProvider", true)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	client := pb.NewSynerexClient(conn)
	argJson := fmt.Sprintf("{Client:TaxiProVider, Price: %d}",*price)
	sclient = sxutil.NewSMServiceClient(client, pb.ChannelType_RIDE_SHARE,argJson)
	log.Print(reflect.TypeOf(sclient))

	wg.Add(1)
	go subscribeDemand(sclient)
	wg.Wait()
	sxutil.CallDeferFunctions() // cleanup!

}
