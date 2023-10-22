package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// 2 arguments, number of requests and delay between them
func main() {
	fmt.Println("aa")

	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Println("Too few arguments, please include ip:port number of requests and delay between  requests")
		os.Exit(-1)
	}

	var err error
	var its, delay int

	//ipport := args[0]
	its, err = strconv.Atoi(args[1])
	check(err)
	delay, err = strconv.Atoi(args[2])
	check(err)

	servURL, err := url.Parse(args[0])
	check(err)

	fmt.Println("Starting requests to server")
	for i := 0; i < its; i++ {
		resp, err := http.Get(servURL.String())
		if err != nil {
			check(err)
		}
		defer resp.Body.Close()

		fmt.Println("Response status:", resp.Status)

		scanner := bufio.NewScanner(resp.Body)
		for i := 0; scanner.Scan() && i < 5; i++ {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			check(err)
		}

		time.Sleep(time.Second * time.Duration(rand.Int31n(int32(delay))))
	}
}

func check(e error) {
	if e != nil {
		fmt.Println("There was an error:")
		//fmt.Println(e.Error())
		panic(e)
	}
}
