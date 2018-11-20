package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/nkawa/gtfsparser"
	"github.com/nkawa/gtfsparser/gtfs"
	"github.com/synerex/synerex_alpha/api"
	"github.com/synerex/synerex_alpha/api/common"
	"github.com/synerex/synerex_alpha/api/ptransit"
	"google.golang.org/grpc"
	"math"
	"strconv"

	//	"github.com/synerex/synerex_alpha/api/common"
	//	"github.com/synerex/synerex_alpha/api/ptransit"
	"github.com/synerex/synerex_alpha/sxutil"
	//	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:10000", "The synerex server address in the format of host:port")
	nodesrv    = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	feedName      = flag.String("feed", "", "GTFS Feed Filename")
)

func demandPTCallback(clt *sxutil.SMServiceClient, sp *api.Demand) {
 sp.GetArg_PTService()
}


func subscribePTDemand(client *sxutil.SMServiceClient) {
	ctx := context.Background() //
	err := client.SubscribeDemand(ctx, demandPTCallback)
	// not finish.
	log.Printf("Error:Demand %s\n",err.Error())
}

// find closest index from shapepoints
func pointDistance(p1 *common.Point, lat2 float32, lon2 float32) float32{

	p2 := &common.Point{
		Latitude: float64(lat2),
		Longitude: float64(lon2),
	}
	d, _ := p1.Distance(p2)
	return float32(d)
}

func getClosestPoints(pts gtfs.ShapePoints, lat float32, lon float32 ) (int, float32 ){
	dist := float32(math.MaxFloat32)
	di := -1
	p := &common.Point{
		Latitude:  float64(lat),
		Longitude: float64(lon),
	}
	dists := []float32{}
	idx  :=[]int{}
	for i, pt := range (pts) {
		d := pointDistance(p, pt.Lat, pt.Lon)
		if d < dist {
			di = i
			dist = d
		}
		if d < 15 { //less than 5m
			dists = append(dists,d)
			idx = append(idx,i)
		}
	}
	if len(idx) ==0{
		dists = append(dists,dist)
		idx = append(idx,di)
		fmt.Printf("Distance stop!")
	}
	fmt.Printf	("Got nearest %d %v %v\n",len(idx),idx,dists)
	return di, dist
}

func getLatLonFromRatio(shp *gtfs.Shape, from_stop *gtfs.Stop,to_stop *gtfs.Stop, ratio float32) (float32, float32){
	// try to findout nearest point from shp.
	if ratio < 0.001 {
//		fmt.Printf("small ratio %f\n",ratio)
		return from_stop.Lat, from_stop.Lon
	}
	if ratio > 0.999 {
//		fmt.Printf("large ratio %f\n",ratio)
		return to_stop.Lat, to_stop.Lon
	}

	c1,d1 := getClosestPoints(shp.Points, from_stop.Lat, from_stop.Lon )
	c2,d2 := getClosestPoints(shp.Points, to_stop.Lat, to_stop.Lon )

	if c1 == -1 || c2 ==- 1{
		fmt.Printf("Umm. failuer")
		fmt.Printf("Find closest index len:%d, %d, %.1f  -> %d, %.1f\n",len(shp.Points), c1,d1, c2, d2)
	}
	if  d2 > 15 {
		fmt.Printf("Distance %f, %s\n", d2, to_stop.Name)
	}



	var totalDist float32
	totalDist = 0
	if c1 < c2 {// ok order
		for i := c1; i <c2; i++ {
			pt := shp.Points[i]
			pt2 := shp.Points[i+1]
			totalDist += pointDistance(&common.Point{Latitude:float64(pt.Lat),Longitude:float64(pt.Lon)},
							pt2.Lat, pt2.Lon)
		}
//		fmt.Printf("Dist:%f, Ratio %f \n",distance, ratio)
		ratioDist := totalDist * ratio // we need to step this.
		var distance  float32
		distance = 0
		for i := c1; i <c2; i++ {
			pt := shp.Points[i]
			pt2 := shp.Points[i+1]
			diff := pointDistance(&common.Point{Latitude:float64(pt.Lat),Longitude:float64(pt.Lon)},
				pt2.Lat, pt2.Lon)
			if distance + diff > ratioDist { // now We came here!
				restDist := ratioDist - distance
				ptRatio := restDist / ratioDist // get ratio of pt -pt2
				fmt.Printf("lat,lon dist %.1f , %.1f, total %.1f ratio %.3f\n",restDist, distance, totalDist, ptRatio)
				return (pt.Lat*(1-ptRatio)+pt2.Lat*ptRatio), (pt.Lon *(1-ptRatio)+ pt2.Lon*ptRatio)
			}
			distance += diff
		}
		return	to_stop.Lat, to_stop.Lon
	}else{
		fmt.Printf("Reverse order ")
		fmt.Printf("Find closest index len:%d, %d, %.1f  -> %d, %.1f\n",len(shp.Points), c1,d1, c2, d2)
	}

	return 0,0
}



//   st[ix-1].Departure_time __ bus(t) __   st[ix].Arrival_time
// return bus location with index and time.
func locWithTime(feed *gtfsparser.Feed, trip_id string,tp *gtfs.Trip, ix int,  t gtfs.Time) (rix int, lat float32, lon float32){
	if ix < 0 {return ix, 0, 0} // just for error check
	st := tp.StopTimes
	if ix == 0 && t.Minus(st[0].Arrival_time) <= 0 { // before start
		return 0, st[0].Stop.Lat, st[0].Stop.Lon
	}

	if ix >= len(st) {
		fmt.Printf("Index Error %d at %s\n",ix, trip_id)
		return -1, 0, 0
	}
	for ; t.Minus(st[ix].Arrival_time) >= 0;  { // arrived current dest.
		ix ++
		if len(st)==ix	{ // it was final station.
			return -1, st[ix-1].Stop.Lat, st[ix-1].Stop.Lon
		}
	}
	//
	duration := st[ix].Arrival_time.Minus(st[ix-1].Departure_time)
	dt := t.Minus(st[ix-1].Departure_time)

	shape_id := tp.Route.Id
	shapes , ok := feed.Shapes[shape_id]
	if !ok {
		fmt.Printf("Can't find shape from shape_id %s\n", shape_id)
	}
	ratio := float32(dt)/float32(duration)
//	if shape_id == "111001" {
//		fmt.Printf("Get Lonlat %d  time: %d, duration:%d  ratio:%f\n", ix, dt, duration, ratio)
//	}
	lat, lon = getLatLonFromRatio(shapes,st[ix-1].Stop, st[ix].Stop, ratio)

	// we should use shape index. but just use interpolated..

	return ix, lat, lon

}
func deg2rad(deg float32) float64 {
	return float64(math.Pi *deg / 180.0)
}
func calcAngle(lat1 float32, lon1 float32, lat2 float32, lon2 float32) float32{
	a := 6378137.000
	b := 6356752.314
	e := math.Sqrt((math.Pow(a, 2) - math.Pow(b, 2)) / math.Pow(a, 2))

	lat0 :=(deg2rad(lat1)+deg2rad(lat2))/2.0
	dlat := deg2rad(lat1)-deg2rad(lat2)
	dlon := deg2rad(lon1)-deg2rad(lon2)

	W := math.Sqrt(1 - math.Pow(e, 2)*math.Pow(math.Sin(lat0), 2))
	M := a * (1 - math.Pow(e, 2)) / math.Pow(W, 3)
	N := a / W

	ddi := dlat * M
	ddk := dlon * N * math.Cos(lat0)

//	ret := float32(180.0*math.Atan2(ddi,ddk)/math.Pi)
//	fmt.Printf("Ret %f\n",ret)
	ret := float32(180 + 180.0*math.Atan2(ddk, ddi)/math.Pi)
	if ret >= 360 {
		ret -= 360
	}

	return ret

}

func supplyPTransitFeed(clt *sxutil.SMServiceClient, feed *gtfsparser.Feed){

//	go subscribePTDemand(clt) // wait for demand to give "TimeTable Info"

	tripStatus := make(map[string]int)

//	for { // run for every 2 secs.
//		time.Sleep(time.Second * 2)

//  we start think from current time
//	now := time.Now()

	now := time.Date(2018,11,10,8,40,0,0,time.Local)
	ct := 0
//		year, month, date := now.Date()
	lastLat := make(map[string]float32)
	lastLon := make(map[string]float32)
	lastAng := make(map[string]float32)
	for {
		t := gtfs.Time{
			Hour: int8(now.Hour()),
			Minute: int8(now.Minute()),
			Second: int8(now.Second()),
		}
		for k,v := range(feed.Trips){
			if tripStatus[k] < 0 {
				continue
			}
			rid , _ := strconv.ParseInt(feed.Trips[k].Route.Id,10,32)
			if rid != 111001 {
				continue
			}
			st ,lat, lon:=	locWithTime(feed, k, v, tripStatus[k], t)
			tripStatus[k] = st
			fmt.Printf("%d:st %d: %s \n",rid, st, feed.Trips[k].StopTimes[st].Stop.Name)
			if st > 0 {
				var angle float32
				if lastLat[k] == lat && lastLon[k] == lat {
					angle = lastAng[k]
				} else {
					angle = calcAngle(lastLat[k], lastLon[k], lat, lon)
				}
				fmt.Printf("%s: %s, %d: %d, (%.4f,%.4f)-(%.4f,%.4f) angle:%.2f\n",now.Format("15:04:05"), k,rid, st, lastLat[k], lastLon[k],lat, lon, angle)

				place := common.NewPlace().WithPoint(&common.Point{
					Latitude: float64(lat),
					Longitude: float64(lon),
				})
				pts := &ptransit.PTService{
					VehicleId: int32(rid),
					Angle: float32(angle),
					Speed: int32(0.0),
					CurrentLocation: place,
				}
				lastLat[k]=lat
				lastLon[k]=lon
				smo := sxutil.SupplyOpts{
					Name:  "GTFS Supply",
					PTService: pts,
					JSON: "",
				}
				clt.RegisterSupply(&smo)

			}else{
				if lat != 0.0 {
					lastLat[k] = lat
					lastLon[k] = lon
//					fmt.Printf("%s: %s, %d: %d, (%.4f,%.4f)-\n",now.Format("15:04:05"), k,rid, st, lastLat[k], lastLon[k])
				}
			}
		}
		time.Sleep(time.Millisecond * 500)
		now = now.Add(time.Second * 10)
		ct++
		if ct > 5000 {
			break
		}
	}

}


// start reading gtfs
func main(){
	flag.Parse()


	if len(*feedName) ==0 {
		fmt.Printf("Please speficy GTFS Feed name.")
		os.Exit(0)
	}

	feed := gtfsparser.NewFeed()

	err := feed.Parse(*feedName)
	if err != nil{
		fmt.Printf("Error %s\n",err.Error())
	}
	fmt.Printf("Done, parsed %d agencies, %d stops, %d routes, %d trips, %d fare attributes\n",
		len(feed.Agencies),len(feed.Stops), len(feed.Routes), len(feed.Trips), len(feed.FareAttributes))

//	for k, v := range feed.Stops{
//		fmt.Printf("[%s] %s (@ %f, %f)\n", k, v.Name, v.Lat, v.Lon)
//	}



// here for usual provider
	sxutil.RegisterNodeName(*nodesrv, "PTransit:"+*feedName, false)

	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	var opts []grpc.DialOption
	//	wg := sync.WaitGroup{} // for syncing other goroutines

	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to connect synerex server: %v", err)
	}

	client := api.NewSynerexClient(conn)
	sclient := sxutil.NewSMServiceClient(client, api.ChannelType_PT_SERVICE,"")

	supplyPTransitFeed(sclient, feed)

	sxutil.CallDeferFunctions() // cleanup!

}
