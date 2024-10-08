// Copyright 2024 Redpanda Data, Inc.
//
// Licensed as a Redpanda Enterprise file under the Redpanda Community
// License (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
// https://github.com/redpanda-data/connect/blob/main/licenses/rcl.md

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/charmbracelet/log"
	"github.com/usedatabrew/pglogicalstream"
	"gopkg.in/yaml.v3"
)

func main() {
	var config pglogicalstream.Config
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	pgStream, err := pglogicalstream.NewPgStream(config)
	if err != nil {
		panic(err)
	}

	pgStream.OnMessage(func(message pglogicalstream.Wal2JsonChanges) {
		fmt.Println(message.Changes)
		if message.Lsn != nil {
			// Snapshots dont have LSN
			pgStream.AckLSN(*message.Lsn)
		}
	})
}
