// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"encoding/json"

	pb "github.com/datacommonsorg/mixer/proto"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "127.0.0.1:12345", "Address of grpc server.")
)

func main() {
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMixerClient(conn)
	ctx := context.Background()

	// Get arcs for City.
	if getProperty(ctx, c, []string{"City", "County", "State"}) != nil {
		log.Fatalf("Error: %v", err)
	}

	// Get arcs for Mountain View.
	if getProperty(ctx, c, []string{"geoId/0649670"}) != nil {
		log.Fatalf("Error: %v", err)
	}

	// Get arcs for an EncodeExperiment.
	if getProperty(ctx, c, []string{"dc/bsmvthtq89217"}) != nil {
		log.Fatalf("Error: %v", err)
	}
}

func getProperty(ctx context.Context, c pb.MixerClient, dcids []string) error {
	// Create the request
	req := &pb.GetPropertyLabelsRequest{Dcids: dcids}
	fmt.Printf("Requesting { %s}\n", req)

	// Invoke GetPropertyLabels
	start := time.Now()
	res, err := c.GetPropertyLabels(ctx, req)
	elapsed := time.Since(start)
	if err != nil {
		return err
	}

	/// Format the payload
	jsonByte := []byte(res.GetPayload())
	var jsonFmt bytes.Buffer
	err = json.Indent(&jsonFmt, jsonByte, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonFmt.Bytes()))
	fmt.Printf("Request took: %s\n\n", elapsed)
	return nil
}
