package netbox

import (
	"strconv"

	"github.com/fbreckle/go-netbox/netbox/client"
	"github.com/fbreckle/go-netbox/netbox/client/ipam"
	"github.com/fbreckle/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetboxIpamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamRoleCreate,
		Read:   resourceNetboxIpamRoleRead,
		Update: resourceNetboxIpamRoleUpdate,
		Delete: resourceNetboxIpamRoleDelete,

		Description: `:meta:subcategory:IP Address Management (IPAM):From the [official documentation](https://docs.netbox.dev/en/stable/features/ipam/#prefixvlan-roles):

> A role indicates the function of a prefix or VLAN. For example, you might define Data, Voice, and Security roles. Generally, a prefix will be assigned the same functional role as the VLAN to which it is assigned (if any).`,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"weight": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 32767),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
func resourceNetboxIpamRoleCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	data := models.Role{}

	name := d.Get("name").(string)
	slugValue, slugOk := d.GetOk("slug")
	var slug string
	// Default slug to name attribute if not given
	if !slugOk {
		slug = name
	} else {
		slug = slugValue.(string)
	}
	weight := int64(d.Get("weight").(int))
	description := d.Get("description").(string)

	data.Name = &name
	data.Slug = &slug

	data.Weight = &weight
	data.Description = description
	data.Tags = []*models.NestedTag{}

	params := ipam.NewIpamRolesCreateParams().WithData(&data)
	res, err := api.Ipam.IpamRolesCreate(params, nil)
	if err != nil {
		return err
	}
	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxIpamRoleUpdate(d, m)
}

func resourceNetboxIpamRoleRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := ipam.NewIpamRolesReadParams().WithID(id)

	res, err := api.Ipam.IpamRolesRead(params, nil)
	if err != nil {
		errorcode := err.(*ipam.IpamRolesReadDefault).Code()
		if errorcode == 404 {
			// If the ID is updated to blank, this tells Terraform the resource no longer exists (maybe it was destroyed out of band). Just like the destroy callback, the Read function should gracefully handle this case. https://www.terraform.io/docs/extend/writing-custom-providers.html
			d.SetId("")
			return nil
		}
		return err
	}

	if res.GetPayload().Name != nil {
		d.Set("name", res.GetPayload().Name)
	}

	if res.GetPayload().Slug != nil {
		d.Set("slug", res.GetPayload().Slug)
	}

	if res.GetPayload().Weight != nil {
		d.Set("weight", res.GetPayload().Weight)
	}

	if res.GetPayload().Description != "" {
		d.Set("description", res.GetPayload().Description)
	}

	return nil
}

func resourceNetboxIpamRoleUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	data := models.Role{}

	name := d.Get("name").(string)
	slugValue, slugOk := d.GetOk("slug")
	var slug string
	// Default slug to name attribute if not given
	if !slugOk {
		slug = name
	} else {
		slug = slugValue.(string)
	}
	weight := int64(d.Get("weight").(int))
	description := d.Get("description").(string)

	data.Name = &name
	data.Slug = &slug

	data.Weight = &weight
	data.Description = description
	data.Tags = []*models.NestedTag{}

	params := ipam.NewIpamRolesUpdateParams().WithID(id).WithData(&data)
	_, err := api.Ipam.IpamRolesUpdate(params, nil)
	if err != nil {
		return err
	}
	return resourceNetboxIpamRoleRead(d, m)
}

func resourceNetboxIpamRoleDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := ipam.NewIpamRolesDeleteParams().WithID(id)
	_, err := api.Ipam.IpamRolesDelete(params, nil)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
