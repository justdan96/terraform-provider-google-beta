package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cloudidentity "google.golang.org/api/cloudidentity/v1beta1"
)

func dataSourceGoogleCloudIdentityGroupMemberships() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceCloudIdentityGroupMembership().Schema)

	return &schema.Resource{
		Read: dataSourceGoogleCloudIdentityGroupMembershipsRead,

		Schema: map[string]*schema.Schema{
			"memberships": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `List of Cloud Identity group memberships.`,
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
			"group": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The name of the Group to get memberships from.`,
			},
		},
	}
}

func dataSourceGoogleCloudIdentityGroupMembershipsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	result := []map[string]interface{}{}
	err := config.clientCloudIdentity.Groups.Memberships.List(d.Get("group").(string)).View("FULL").Pages(config.context, func(resp *cloudidentity.ListMembershipsResponse) error {

		for _, member := range resp.Memberships {
			result = append(result, map[string]interface{}{
				"name":                 member.Name,
				"roles":                flattenCloudIdentityGroupMembershipsRoles(member.Roles),
				"member_key":           flattenCloudIdentityGroupsEntityKey(member.MemberKey),
				"preferred_member_key": flattenCloudIdentityGroupsEntityKey(member.PreferredMemberKey),
			})
		}

		return nil
	})
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("CloudIdentityGroupMemberships %q", d.Id()))
	}

	d.Set("memberships", result)
	d.SetId(time.Now().UTC().String())
	return nil
}

func flattenCloudIdentityGroupMembershipsRoles(roles []*cloudidentity.MembershipRole) []interface{} {
	transformed := []interface{}{}

	for _, role := range roles {
		transformed = append(transformed, map[string]interface{}{
			"name": role.Name,
		})
	}
	return transformed
}
