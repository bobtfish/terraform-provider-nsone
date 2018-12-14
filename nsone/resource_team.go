package nsone

import (
	"github.com/hashicorp/terraform/helper/schema"
	nsone "gopkg.in/sarguru/ns1-go.v18"
)

func teamResource() *schema.Resource {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}
	s = addPermsSchema(s)
	return &schema.Resource{
		Schema: s,
		Create: TeamCreate,
		Read:   TeamRead,
		Update: TeamUpdate,
		Delete: TeamDelete,
	}
}

func teamToResourceData(d *schema.ResourceData, t *nsone.Team) error {
	d.SetId(t.Id)
	d.Set("name", t.Name)
	permissionsToResourceData(d, t.Permissions)
	return nil
}

func resourceDataToTeam(u *nsone.Team, d *schema.ResourceData) error {
	u.Id = d.Id()
	u.Name = d.Get("name").(string)
	u.Permissions = resourceDataToPermissions(d)
	return nil
}

func TeamCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	mj := nsone.Team{}
	if err := resourceDataToTeam(&mj, d); err != nil {
		return err
	}
	if err := client.CreateTeam(&mj); err != nil {
		return err
	}
	return teamToResourceData(d, &mj)
}

func TeamRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	mj, err := client.GetTeam(d.Id())
	if err != nil {
		return err
	}
	teamToResourceData(d, &mj)
	return nil
}

func TeamDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	err := client.DeleteTeam(d.Id())
	d.SetId("")
	return err
}

func TeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	mj := nsone.Team{
		Id: d.Id(),
	}
	if err := resourceDataToTeam(&mj, d); err != nil {
		return err
	}
	if err := client.UpdateTeam(&mj); err != nil {
		return err
	}
	teamToResourceData(d, &mj)
	return nil
}
