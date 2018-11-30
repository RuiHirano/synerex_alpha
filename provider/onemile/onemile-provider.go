package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/mtfelian/golang-socketio"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/sxutil"
	"google.golang.org/grpc"
)

var (
	version = "0.01"

	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	client     api.SynerexClient

	port   = flag.Int("port", 7777, "Onemile Provider Listening Port")
	ioserv *gosocketio.Server

	n      = flag.Int("n", 1, "Number of taxi (or display)")
	dispWg sync.WaitGroup
)

// vehicle
type vehicle struct {
	VehicleId   string     `json:"vehicle_id"`   // unique id
	VehicleType string     `json:"vehicle_type"` // [onemile | bus | train | ...]
	Status      string     `json:"status"`       // [pickup | free | ride]
	Coord       [2]float64 `json:"coord"`        // current position (lat/lon)
}

// managed vehicles by onemile-provider
var vehicleMap = make(map[string]*vehicle)

// display
type display struct {
	dispId  string              // display id
	channel *gosocketio.Channel // Socket.IO channel
	wg      sync.WaitGroup      // for synchronization to display ad and enquate
}

// taxi/display mapping
var dispMap = make(map[string]*display)

// register OnemileProvider to NodeServer
func registerOnemileProvider() {
	sxutil.RegisterNodeName(*nodesrv, "OnemileProvider", false)
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)
	go sxutil.HandleSigInt()
}

// create SMServiceClient for a given ChannelType
func createSMServiceClient(ch api.ChannelType, arg string) *sxutil.SMServiceClient {
	// create grpc client (at onece)
	if client == nil {
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithInsecure())

		conn, err := grpc.Dial(*serverAddr, opts...)
		if err != nil {
			log.Fatalf("Fail to Connect Synerex Server: %v", err)
		}

		client = api.NewSynerexClient(conn)
	}

	// create SMServiceClient
	return sxutil.NewSMServiceClient(client, ch, arg)
}

// TODO: 乗車シーケンス
// subscribe rideshare channel
func subscribeRideShare(rdClient, rtClient *sxutil.SMServiceClient) {
	ctx := context.Background()
	rdClient.SubscribeDemand(ctx, func(clt *sxutil.SMServiceClient, dm *api.Demand) {
		if dm.GetDemandName() == "" {
			// Confirm
			// TODO: 迎車処理 (routing-providerからのSelectSupply受信〜乗車まで)
		} else {
			// ProposeSupply
			// TODO: 経路取得 (routing-providerからのRegisterDemand受信〜ProposeSupplyまで)
		}
	})
}

// subscribe marketing channel
func subscribeMarketing(mktClient *sxutil.SMServiceClient) {
	// wait until completing display registration
	dispWg.Wait()

	ctx := context.Background()
	seen := make(map[string]struct{})

	mktClient.SubscribeDemand(ctx, func(clt *sxutil.SMServiceClient, dm *api.Demand) {
		if dm.GetDemandName() == "" {
			// Confirm
			log.Printf("Receive SelectSupply [id: %d, name: %s]\n", dm.GetId(), dm.GetDemandName())
			clt.Confirm(sxutil.IDType(dm.GetId()))

			// SubscribeMbus
			clt.SubscribeMbus(context.Background(), func(clt *sxutil.SMServiceClient, msg *api.MbusMsg) {
				// emit display event for each display
				for taxi := range dispMap {
					dispMap[taxi].wg.Add(1)
					go func(taxi, name string, payload interface{}) {
						// wait unti a taxi will depart
						dispMap[taxi].wg.Wait()
						// emit event
						dispMap[taxi].channel.Emit(name, payload)
						log.Printf("Emit [taxi: %s, name: %s, payload: %s]\n", taxi, name, payload)
					}(taxi, "disp_start", msg.ArgJson)
				}
			})
		} else {
			// ProposeSupply
			if _, ok := seen[dm.GetDemandName()]; !ok {
				seen[dm.GetDemandName()] = struct{}{}
				log.Printf("Receive RegisterDemand [id: %d, name: %s]\n", dm.GetId(), dm.GetDemandName())
				sp := &sxutil.SupplyOpts{
					Target: dm.GetId(),
					Name:   "a display for advertising and enqueting",
				}
				clt.ProposeSupply(sp)
			}
		}
	})
}

// run Socket.IO server for Onemile-Client and Onemile-Display-Client
func runSocketIOServer(rdClient, mktClient *sxutil.SMServiceClient) {
	ioserv := gosocketio.NewServer()

	ioserv.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Connected from %s as %s\n", c.IP(), c.Id())
	})

	ioserv.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Disconnected from %s as %s\n", c.IP(), c.Id())
	})

	// TODO: ログイン
	ioserv.On("clt_login", func(c *gosocketio.Channel, data interface{}) {
		log.Printf("Receive clt_login from %s [%v]\n", c.Id(), data)

		taxi := data.(map[string]interface{})["device_id"].(string)

		if v, ok := vehicleMap[taxi]; ok {
			ret := map[string]interface{}{
				"act":  "clt_login",
				"code": 0,
				"results": map[string]interface{}{
					"provider_id": "onemile-provider",
					"vehicle_id":  v.VehicleId,
					"token":       "1234567890",
				},
			}

			// TODO: use Socket.IO Acknowledgement
			c.Emit("clt_login_res", ret)
			log.Printf("Emit [taxi: %s, name: %s, payload: %v]\n", taxi, "clt_login_res", ret)
		}
	})

	// TODO: 位置情報報告 (定期的に)
	// [Rideshare]
	ioserv.On("xxxxx", func(c *gosocketio.Channel, data interface{}) {
	})

	// TODO: 移動処理 (乗車〜降車まで)
	ioserv.On("xxxxx", func(c *gosocketio.Channel, data interface{}) {
	})
	ioserv.On("xxxxx", func(c *gosocketio.Channel, data interface{}) {
	})

	// [Marketing] register taxi and display mapping
	ioserv.On("disp_register", func(c *gosocketio.Channel, data interface{}) {
		log.Printf("Receive disp_register from %s [%v]\n", c.Id(), data)

		taxi := data.(map[string]interface{})["taxi"].(string)
		disp := data.(map[string]interface{})["disp"].(string)

		if _, ok := dispMap[taxi]; !ok {
			dispMap[taxi] = &display{dispId: disp, channel: c, wg: sync.WaitGroup{}}
			log.Printf("Register display [taxi: %s => display: %v]\n", taxi, dispMap[taxi])
			dispWg.Done()
		}
	})

	// [Marketing] complete ad and enquate
	ioserv.On("disp_complete", func(c *gosocketio.Channel, data interface{}) {
		log.Printf("Receive disp_complete from %s [%v]\n", c.Id(), data)

		// marshal json
		bytes, err := json.Marshal(data)
		if err != nil {
			log.Printf("Marshal error: %s\n", err)
		}

		// send results via Mbus
		mktClient.SendMsg(context.Background(), &api.MbusMsg{ArgJson: string(bytes)})
	})

	// [DEBUG] (simulate departure or arrive of taxi in disp-test.html)
	ioserv.On("depart", func(c *gosocketio.Channel, data interface{}) {
		log.Printf("Receive depart from %s [%v]\n", c.Id(), data)

		taxi := data.(map[string]interface{})["taxi"].(string)

		dispMap[taxi].wg.Done()
	})
	ioserv.On("arrive", func(c *gosocketio.Channel, data interface{}) {
		log.Printf("Receive arrive from %s [%v]\n", c.Id(), data)
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", ioserv)
	serveMux.Handle("/", http.FileServer(http.Dir("./display-client")))

	log.Printf("Starting Socket.IO Server %s on port %d", version, *port)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), serveMux)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	// init vehicles
	for i := 0; i < *n; i++ {
		var id = fmt.Sprintf("%02d", i+1)
		vehicleMap[id] = &vehicle{"vehicle" + id, "onemile", "free", [2]float64{0.0, 0.0}}
	}

	// set number of display
	dispWg.Add(*n)

	// register onemile-provider
	registerOnemileProvider()

	var wg sync.WaitGroup

	wg.Add(1)
	// subscribe rideshare channel
	rdClient := createSMServiceClient(api.ChannelType_RIDE_SHARE, "")
	rtClient := createSMServiceClient(api.ChannelType_ROUTING_SERVICE, "")
	go subscribeRideShare(rdClient, rtClient)

	wg.Add(1)
	// subscribe marketing channel
	mktClient := createSMServiceClient(api.ChannelType_MARKETING_SERVICE, "")
	go subscribeMarketing(mktClient)

	wg.Add(1)
	// start Websocket Server
	go runSocketIOServer(rdClient, mktClient)

	wg.Wait()
}
