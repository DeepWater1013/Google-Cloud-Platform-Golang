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

package servers

// [START cloud_game_servers_cluster_list]

import (
	"context"
	"fmt"
	"io"

	gaming "cloud.google.com/go/gaming/apiv1"
	"google.golang.org/api/iterator"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1"
)

// listGameServerClusters lists the clusters registered with a realm.
func listGameServerClusters(w io.Writer, projectID, location, realmID string) error {
	// projectID := "my-project"
	// location := "global"
	// realmID := "myrealm"
	ctx := context.Background()
	client, err := gaming.NewGameServerClustersClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGameServerClustersClient: %v", err)
	}
	defer client.Close()

	req := &gamingpb.ListGameServerClustersRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/realms/%s", projectID, location, realmID),
		View:   gamingpb.GameServerClusterView_FULL,
	}

	it := client.ListGameServerClusters(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Next: %v", err)
		}

		fmt.Fprintf(w, "Cluster listed: %v\n", resp.Name)
		fmt.Fprintf(w, "Cluster installed Agones version: %v\n", resp.ClusterState.AgonesVersionInstalled)
		fmt.Fprintf(w, "Cluster installed Kubernetes version: %v\n", resp.ClusterState.KubernetesVersionInstalled)
		fmt.Fprintf(w, "Cluster installation state: %v\n", resp.ClusterState.InstallationState)
	}

	return nil
}

// [END cloud_game_servers_cluster_list]
