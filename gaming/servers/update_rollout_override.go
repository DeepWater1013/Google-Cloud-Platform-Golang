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

// [START cloud_game_servers_deployment_rollout_override]

import (
	"context"
	"fmt"
	"io"

	gaming "cloud.google.com/go/gaming/apiv1"
	gamingpb "google.golang.org/genproto/googleapis/cloud/gaming/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// updateRolloutOverrideConfig sets the override config for a realm.
func updateRolloutOverrideConfig(w io.Writer, projectID, realmLocation, realmID, deploymentID, configID string) error {
	// projectID := "my-project"
	// realmLocation := "global"
	// realmID := "myrealm"
	// deploymentID := "mydeployment"
	// configID := "myconfig"
	ctx := context.Background()
	client, err := gaming.NewGameServerDeploymentsClient(ctx)
	if err != nil {
		return fmt.Errorf("NewGameServerDeploymentsClient: %v", err)
	}
	defer client.Close()

	req := &gamingpb.UpdateGameServerDeploymentRolloutRequest{
		Rollout: &gamingpb.GameServerDeploymentRollout{
			Name: fmt.Sprintf("projects/%s/locations/global/gameServerDeployments/%s", projectID, deploymentID),

			GameServerConfigOverrides: []*gamingpb.GameServerConfigOverride{
				{
					Selector: &gamingpb.GameServerConfigOverride_RealmsSelector{
						RealmsSelector: &gamingpb.RealmSelector{
							Realms: []string{
								fmt.Sprintf("projects/%s/locations/%s/realms/%s", projectID, realmLocation, realmID),
							},
						},
					},
					Change: &gamingpb.GameServerConfigOverride_ConfigVersion{
						ConfigVersion: configID,
					},
				},
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"game_server_config_overrides",
			},
		},
	}

	op, err := client.UpdateGameServerDeploymentRollout(ctx, req)
	if err != nil {
		return fmt.Errorf("UpdateGameServerDeploymentRollout: %v", err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Wait: %v", err)
	}

	fmt.Fprintf(w, "Deployment rollout updated: %v", resp.Name)
	return nil
}

// [END cloud_game_servers_deployment_rollout_override]
