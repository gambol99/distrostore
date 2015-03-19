/*
Copyright 2014 Rohith All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"math/rand"
	"log"
	"syscall"
	"time"

	ds "github.com/gambol99/distrostore"

	"github.com/alecthomas/kingpin"
)

var (
	bootstrap = kingpin.Flag("bootstrap", "whether we are the bootstrap").Bool()
	members = kingpin.Flag("member", "add a member to the list").Strings()
	offset = kingpin.Flag("offset", "add the offset to the ports").Required().Int()
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randonString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func tempDirectory() string {
	dir, err := ioutil.TempDir("", "consul")
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	return dir
}

func main() {
	kingpin.Parse()
	config := ds.DefaultContext()
	config.LogOutput = os.Stdout
	config.DataDir = tempDirectory()
	config.Bootstrap = *bootstrap
	config.Members = *members
	config.PortsConfig.ApplyIndex(*offset)
	log.Println("Ports: %s", config.PortsConfig)

	store, err := ds.New(config)
	if err != nil {
		log.Fatalf("Failed to create the distributed data store, error: %s", err)
	}

	time.Sleep(time.Duration(3) * time.Second)

	// step: check if a key exists
	key_name := "my_store_key"
	if value, found, err := store.Get(key_name); err != nil {
		log.Fatalf("Failed to check the status of the key, error: %s", err)
	} else if !found {
		log.Println("Setting a random key in the cluster")
		store.Set(key_name, randonString(32))
	} else {
		log.Printf("found the store key: %s", value)
	}

	defer func() {
		// delete the directory
		log.Printf("Removing the data directory: %s\n", config.DataDir)
		if err := os.RemoveAll(config.DataDir); err != nil {
			log.Fatalf("Failed to remove the temporary directory: %s, error: %s", config.DataDir, err)
		}
	}()

	// wait for the signal to terminate
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChannel
}
