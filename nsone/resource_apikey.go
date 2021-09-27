package nsone

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	nsone "gopkg.in/sarguru/ns1-go.v18"
)

func apikeyResource() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"key": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"teams": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
	s = addPermsSchema(s)
	return &schema.Resource{
		Schema: s,
		Create: ApikeyCreate,
		Read:   ApikeyRead,
		Update: ApikeyUpdate,
		Delete: ApikeyDelete,
	}
}

func apikeyToResourceData(d *schema.ResourceData, u *nsone.Apikey) error {
	d.SetId(u.Id)
	d.Set("name", u.Name)
	d.Set("key", u.Key)
	d.Set("teams", u.Teams)
	permissionsToResourceData(d, u.Permissions)
	return nil
}

func resourceDataToPermissions(d *schema.ResourceData) nsone.PermissionsMap {
	var p nsone.PermissionsMap
	if v, ok := d.GetOk("dns_view_zones"); ok {
		p.Dns.ViewZones = v.(bool)
	}
	if v, ok := d.GetOk("dns_manage_zones"); ok {
		p.Dns.ManageZones = v.(bool)
	}
	if v, ok := d.GetOk("dns_zones_allow_by_default"); ok {
		p.Dns.ZonesAllowByDefault = v.(bool)
	}
	if v, ok := d.GetOk("dns_zones_deny"); ok {
		deny_raw := v.([]interface{})
		p.Dns.ZonesDeny = make([]string, len(deny_raw))
		for i, deny := range deny_raw {
			p.Dns.ZonesDeny[i] = deny.(string)
		}
	} else {
		p.Dns.ZonesDeny = make([]string, 0)
	}
	if v, ok := d.GetOk("dns_zones_allow"); ok {
		allow_raw := v.([]interface{})
		p.Dns.ZonesAllow = make([]string, len(allow_raw))
		for i, allow := range allow_raw {
			p.Dns.ZonesAllow[i] = allow.(string)
		}
	} else {
		p.Dns.ZonesAllow = make([]string, 0)
	}
	if v, ok := d.GetOk("data_push_to_datafeeds"); ok {
		p.Data.PushToDatafeeds = v.(bool)
	}
	if v, ok := d.GetOk("data_manage_datasources"); ok {
		p.Data.ManageDatasources = v.(bool)
	}
	if v, ok := d.GetOk("data_manage_datafeeds"); ok {
		p.Data.ManageDatafeeds = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_users"); ok {
		p.Account.ManageUsers = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_payment_methods"); ok {
		p.Account.ManagePaymentMethods = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_plan"); ok {
		p.Account.ManagePlan = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_teams"); ok {
		p.Account.ManageTeams = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_apikeys"); ok {
		p.Account.ManageApikeys = v.(bool)
	}
	if v, ok := d.GetOk("account_manage_account_settings"); ok {
		p.Account.ManageAccountSettings = v.(bool)
	}
	if v, ok := d.GetOk("account_view_activity_log"); ok {
		p.Account.ViewActivityLog = v.(bool)
	}
	if v, ok := d.GetOk("account_view_invoices"); ok {
		p.Account.ViewInvoices = v.(bool)
	}
	if v, ok := d.GetOk("monitoring_manage_lists"); ok {
		p.Monitoring.ManageLists = v.(bool)
	}
	if v, ok := d.GetOk("monitoring_manage_jobs"); ok {
		p.Monitoring.ManageJobs = v.(bool)
	}
	if v, ok := d.GetOk("monitoring_view_jobs"); ok {
		p.Monitoring.ViewJobs = v.(bool)
	}
	return p
}

func resourceDataToApikey(u *nsone.Apikey, d *schema.ResourceData) error {
	u.Id = d.Id()
	u.Name = d.Get("name").(string)
	if v, ok := d.GetOk("teams"); ok {
		teams_raw := v.([]interface{})
		u.Teams = make([]string, len(teams_raw))
		for i, team := range teams_raw {
			u.Teams[i] = team.(string)
		}
	} else {
		u.Teams = make([]string, 0)
	}
	u.Permissions = resourceDataToPermissions(d)
	return nil
}

func ApikeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	mj := nsone.Apikey{}
	if err := resourceDataToApikey(&mj, d); err != nil {
		return err
	}
	if err := client.CreateApikey(&mj); err != nil {
		return err
	}
	return apikeyToResourceData(d, &mj)
}

func ApikeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	mj, err := client.GetApikey(d.Id())
	if err != nil {
		return err
	}
	apikeyToResourceData(d, &mj)
	return nil
}

func ApikeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	err := client.DeleteApikey(d.Id())
	d.SetId("")
	return err
}

func ApikeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	mj := nsone.Apikey{
		Id: d.Id(),
	}
	if err := resourceDataToApikey(&mj, d); err != nil {
		return err
	}
	if err := client.UpdateApikey(&mj); err != nil {
		return err
	}
	apikeyToResourceData(d, &mj)
	return nil
}
