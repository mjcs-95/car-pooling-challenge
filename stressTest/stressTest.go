package main

import (
	"fmt"
	"io/ioutil"
	"main/v2/server"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var groupAmount = 150000
var carsAmount = 100000

func stressTest() {
	fmt.Printf("generating car request --- ")
	start := time.Now()
	carsBody := "[\n    "
	for i := 1; i < carsAmount; i++ {
		carsBody += fmt.Sprintf("{ \"id\": %d, \"seats\": %d },\n    ", i, rand.Intn(3)+4)
	}
	carsBody += fmt.Sprintf("{ \"id\": %d, \"seats\": %d }\n", carsAmount, rand.Intn(3)+4)
	carsBody += "]"
	carsReqBody := string(carsBody)
	payload := strings.NewReader(carsReqBody)
	carsReq, _ := http.NewRequest("PUT", "http://localhost:9091/cars", payload)
	carsReq.Header.Add("Content-Type", server.ContentTypeJSON)

	elapsed := time.Since(start)
	fmt.Printf(" %f seconds\n", elapsed.Seconds())

	fmt.Printf("generating %d journey requests ---", groupAmount)
	start = time.Now()

	journeysREQ := []*http.Request{}
	for i := 0; i < groupAmount; i++ {
		payload = strings.NewReader(fmt.Sprintf("{ \"id\": %d, \"people\": %d }", i, rand.Intn(5)+1))
		Req, _ := http.NewRequest("POST", "http://localhost:9091/journey", payload)
		Req.Header.Add("Content-Type", server.ContentTypeJSON)
		journeysREQ = append(journeysREQ, Req)
	}

	elapsed = time.Since(start)
	fmt.Printf(" %f seconds\n", elapsed.Seconds())

	fmt.Printf("generating %d locate requests for random groups ---", groupAmount/2)
	start = time.Now()

	locateREQ := []*http.Request{}
	for i := 0; i < groupAmount; i++ {
		payload = strings.NewReader(fmt.Sprintf("ID=%d", rand.Intn(groupAmount)))
		Req, _ := http.NewRequest("POST", "http://localhost:9091/locate", payload)
		Req.Header.Add("Content-Type", server.ContentTypeURLENCODED)
		locateREQ = append(locateREQ, Req)
	}

	elapsed = time.Since(start)
	fmt.Printf(" %f seconds\n", elapsed.Seconds())

	fmt.Printf("generating %d dropoff requests ---", groupAmount)
	start = time.Now()
	dropoffREQ := []*http.Request{}
	for i := 0; i < groupAmount; i++ {
		payload = strings.NewReader(fmt.Sprintf("ID=%d", i))
		Req, _ := http.NewRequest("POST", "http://localhost:9091/dropoff", payload)
		Req.Header.Add("Content-Type", server.ContentTypeURLENCODED)
		dropoffREQ = append(dropoffREQ, Req)
	}
	elapsed = time.Since(start)
	fmt.Printf(" %f seconds\n", elapsed.Seconds())

	client := &http.Client{}
	/**/
	fmt.Printf("executing %d car request --- ", carsAmount)
	start = time.Now()
	_, err := client.Do(carsReq)
	if err != nil {
		fmt.Println(err)
		return
	}
	elapsed = time.Since(start)
	fmt.Printf(" %f seconds\n", elapsed.Seconds())
	time.Sleep(2 * time.Second)
	fmt.Printf("executing journey requests --- ")
	start = time.Now()
	for i := 0; i < len(journeysREQ); i++ {
		res, err := client.Do(journeysREQ[i])
		if err != nil {
			fmt.Println(err)
			break
		}
		defer res.Body.Close()
		ioutil.ReadAll(res.Body)
	}
	elapsed = time.Since(start)
	fmt.Printf(" %f seconds\n", elapsed.Seconds())
	/**/
	time.Sleep(2 * time.Second)
	fmt.Printf("executing locate requests  --- ")
	start = time.Now()
	for i := 0; i < len(locateREQ); i++ {
		res, err := client.Do(locateREQ[i])
		if err != nil {
			fmt.Println(err)
			break
		}
		defer res.Body.Close()
		ioutil.ReadAll(res.Body)
	}
	elapsed = time.Since(start)
	fmt.Printf(" %f seconds\n", elapsed.Seconds())
	time.Sleep(2 * time.Second)
	fmt.Printf("executing dropoff requests --- ")
	start = time.Now()
	for i := 0; i < len(dropoffREQ); i++ {
		res, err := client.Do(dropoffREQ[i])
		if err != nil {
			fmt.Println(err)
			break
		}
		defer res.Body.Close()
		ioutil.ReadAll(res.Body)
	}
	elapsed = time.Since(start)
	fmt.Printf(" %f seconds\n", elapsed.Seconds())
}

func main() {
	stressTest()
}
