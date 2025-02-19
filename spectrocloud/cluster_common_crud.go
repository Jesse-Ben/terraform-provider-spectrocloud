package spectrocloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spectrocloud/palette-sdk-go/client"
)

var resourceClusterReadyPendingStates = []string{
	"NotReady",
}

var resourceClusterDeletePendingStates = []string{
	"Pending",
	"Provisioning",
	"Running",
	"Deleting",
	"Importing",
}
var resourceClusterCreatePendingStates = []string{
	"Unknown",
	"Pending",
	"Provisioning",
	"Importing",
}

var virtualClusterLifecycleStates = []string{
	"Resuming",
	"Pausing",
	"Paused",
}

func waitForClusterReady(ctx context.Context, d *schema.ResourceData, uid string, diags diag.Diagnostics, c *client.V1Client) (diag.Diagnostics, bool) {
	d.SetId(uid)

	stateConf := &retry.StateChangeConf{
		Pending:    resourceClusterReadyPendingStates,
		Target:     []string{"Ready"},
		Refresh:    resourceClusterReadyRefreshFunc(c, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate) - 1*time.Minute,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err), true
	}
	return nil, false
}

func waitForVirtualClusterLifecyclePause(ctx context.Context, d *schema.ResourceData, uid string, diags diag.Diagnostics, c *client.V1Client) (diag.Diagnostics, bool) {
	d.SetId(uid)
	stateConf := &retry.StateChangeConf{
		Pending:    virtualClusterLifecycleStates,
		Target:     []string{"Paused"},
		Refresh:    resourceVirtualClusterLifecycleStateRefreshFunc(c, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate) - 1*time.Minute,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err), true
	}
	return nil, false
}
func waitForVirtualClusterLifecycleResume(ctx context.Context, d *schema.ResourceData, uid string, diags diag.Diagnostics, c *client.V1Client) (diag.Diagnostics, bool) {
	d.SetId(uid)
	stateConf := &retry.StateChangeConf{
		Pending:    virtualClusterLifecycleStates,
		Target:     []string{"Running"},
		Refresh:    resourceVirtualClusterLifecycleStateRefreshFunc(c, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate) - 1*time.Minute,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err), true
	}
	return nil, false
}

func resourceClusterReadyRefreshFunc(c *client.V1Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cluster, err := c.GetClusterWithoutStatus(id)
		if err != nil {
			return nil, "", err
		} else if cluster == nil || cluster.Status == nil {
			return nil, "NotReady", nil
		}
		return cluster, "Ready", nil
	}
}

func waitForClusterCreation(ctx context.Context, d *schema.ResourceData, uid string, diags diag.Diagnostics, c *client.V1Client, initial bool) (diag.Diagnostics, bool) {
	d.SetId(uid)

	if initial { // only skip_completion when initally creating a cluster, do not skip when attach addon profile
		if d.Get("skip_completion") != nil && d.Get("skip_completion").(bool) {
			return diags, true
		}

		if _, found := toTags(d)["skip_completion"]; found {
			return diags, true
		}
	}

	diagnostics, isError := waitForClusterReady(ctx, d, uid, diags, c)
	if isError {
		return diagnostics, true
	}

	stateConf := &retry.StateChangeConf{
		Pending:    resourceClusterCreatePendingStates,
		Target:     []string{"Running"},
		Refresh:    resourceClusterStateRefreshFunc(c, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate) - 1*time.Minute,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	// Wait, catching any errors
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err), true
	}
	return nil, false
}

//	var resourceClusterUpdatePendingStates = []string{
//		"backing-up",
//		"modifying",
//		"resetting-master-credentials",
//		"upgrading",
//	}
func waitForClusterDeletion(ctx context.Context, c *client.V1Client, id string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:    resourceClusterDeletePendingStates,
		Target:     nil, // wait for deleted
		Refresh:    resourceClusterStateRefreshFunc(c, id),
		Timeout:    timeout,
		MinTimeout: 10 * time.Second,
		Delay:      30 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func resourceClusterStateRefreshFunc(c *client.V1Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cluster, err := c.GetCluster(id)
		if err != nil {
			return nil, "", err
		} else if cluster == nil {
			return nil, "Deleted", nil
		}

		state := cluster.Status.State
		log.Printf("Cluster state (%s): %s", id, state)

		return cluster, state, nil
	}
}

func resourceVirtualClusterLifecycleStateRefreshFunc(c *client.V1Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cluster, err := c.GetCluster(id)
		if err != nil {
			return nil, "", err
		} else if cluster == nil {
			return nil, "Deleted", nil
		}

		state := cluster.Status.Virtual.LifecycleStatus.Status
		log.Printf("Cluster state (%s): %s", id, state)

		return cluster, state, nil
	}
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.V1Client)

	var diags diag.Diagnostics

	err := c.DeleteCluster(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := waitForClusterDeletion(ctx, c, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
