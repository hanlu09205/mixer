// Copyright 2020 Google LLC
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

package e2etest

import (
	"context"
	"io/ioutil"
	"path"
	"runtime"
	"testing"

	pb "github.com/datacommonsorg/mixer/proto"
	"github.com/datacommonsorg/mixer/server"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestGetStatAll(t *testing.T) {
	ctx := context.Background()

	memcacheData, err := loadMemcache()
	if err != nil {
		t.Fatalf("Failed to load memcache %v", err)
	}

	client, err := setup(server.NewMemcache(memcacheData))
	if err != nil {
		t.Fatalf("Failed to set up mixer and client")
	}
	_, filename, _, _ := runtime.Caller(0)
	goldenPath := path.Join(
		path.Dir(filename), "../golden_response/staging/get_stat_all")

	for _, c := range []struct {
		statVars   []string
		places     []string
		goldenFile string
	}{
		{
			[]string{"Count_Person", "Count_CriminalActivities_CombinedCrime", "Amount_EconomicActivity_GrossNationalIncome_PurchasingPowerParity_PerCapita"},
			[]string{"country/USA", "geoId/06", "geoId/0649670"},
			"result.json",
		},
	} {
		resp, err := client.GetStatAll(ctx, &pb.GetStatAllRequest{
			StatVars: c.statVars,
			Places:   c.places,
		})
		if err != nil {
			t.Errorf("could not GetStatAll: %s", err)
			continue
		}
		goldenFile := path.Join(goldenPath, c.goldenFile)
		if generateGolden {
			updateProtoGolden(resp, goldenFile)
			continue
		}
		var expected pb.GetStatAllResponse
		file, _ := ioutil.ReadFile(goldenFile)
		err = protojson.Unmarshal(file, &expected)
		if err != nil {
			t.Errorf("Can not Unmarshal golden file")
			continue
		}

		if diff := cmp.Diff(resp, &expected, protocmp.Transform()); diff != "" {
			t.Errorf("payload got diff: %v", diff)
			continue
		}
	}
}
