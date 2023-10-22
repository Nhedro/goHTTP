package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var delay time.Duration = 60

var times []int
var mux sync.Mutex

var fils *bufio.Writer
var filsMux sync.Mutex

func writeToFile(line string) error {

	filsMux.Lock()
	defer filsMux.Unlock()
	_, err := fils.WriteString(line)
	if err != nil {
		return err
	}

	err = fils.Flush()
	return err
}

func filterSlice(filter int) {
	// get Mutex
	mux.Lock()
	// Free mutex
	defer mux.Unlock()

	for i := 0; i < len(times); i++ {
		if times[i] >= filter {
			times = append(make([]int, 0), times[i:]...)
			break
		} else if i == len(times)-1 {
			times = make([]int, 0)
		}
	}

}

func request(w http.ResponseWriter, req *http.Request) {

	t := time.Now()

	tnum, err := strconv.Atoi(t.Format("20060102150405"))
	check(err)

	filter, err := strconv.Atoi(t.Add(-(time.Second * delay)).Format("20060102150405"))
	check(err)

	filterSlice(filter)

	times = append(times, tnum)

	err = writeToFile(strconv.FormatInt(int64(tnum), 10) + "\n")
	check(err)

	if err != nil {
		_, err = io.WriteString(w, "-1")
		check(err)
	} else {
		_, err = io.WriteString(w, strconv.FormatInt(int64(len(times)), 10))
		check(err)
	}

	fmt.Println("Print Request  " + strconv.FormatInt(int64(tnum), 10))
	for i := 0; i < len(times); i++ {
		fmt.Println(times[i])
	}
}

// receives 2 arguments: port to listen from and filename
func main() {

	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("Too few arguments, please include ip:port number of requests and delay between  requests")
		os.Exit(-1)
	}
	var err error

	filename := args[1]
	var file *os.File

	times = make([]int, 0)

	if fileExists(filename) {
		// process file into times
		filter, err := strconv.Atoi(time.Now().Add(-(time.Second * delay)).Format("20060102150405"))
		check(err)
		/*
			err = decrypt("../key", filename)
			defer encrypt("../key", filename)
		*/
		if err == nil {
			file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
			check(err)
			fileScanner := bufio.NewScanner(file)

			for fileScanner.Scan() {
				tnum, err := strconv.Atoi(fileScanner.Text())
				check(err)
				if tnum >= filter {
					times = append(times, tnum)
				}
			}
		}

		file, err = os.Create(filename)
		defer file.Close()
		check(err)

		filsMux.Lock()

		fils = bufio.NewWriter(file)

		for i := 0; i < len(times); i++ {
			_, err = fils.WriteString(strconv.FormatInt(int64(times[i]), 10) + "\n")
			check(err)
		}
		filsMux.Unlock()

	} else {
		file, err = os.Create(filename)
		defer file.Close()
		check(err)

		filsMux.Lock()
		fils = bufio.NewWriter(file)
		filsMux.Unlock()
	}

	http.HandleFunc("/", request)

	http.HandleFunc("/exit", request)

	err = http.ListenAndServe(":"+args[0], nil)
	check(err)

	fmt.Println("\n\nPrint current status")
	for i := 0; i < len(times); i++ {
		fmt.Println(times[i])
	}
	fmt.Println("\nDONE")

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func check(e error) {
	if e != nil {
		fmt.Println("There was an error:")
		fmt.Println(e.Error())
		//panic(e)
	}
}

/*
func encrypt(keyPath string, file string) error {
	var plainText []byte

	key, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	os.Create(file)

	err = os.WriteFile(file, cipherText, 0777)
	if err != nil {
		return err
	}

	return nil
}

func decrypt(keyPath string, file string) error {

	// Reading ciphertext file
	cipherText, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// Reading key
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}

	// Creating block of algorithm
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Creating GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Deattached nonce and decrypt
	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return err
	}

	os.Create(file)

	// Writing decryption content
	err = os.WriteFile(file, plainText, 0777)
	if err != nil {
		return err
	}

	return nil
}
*/
