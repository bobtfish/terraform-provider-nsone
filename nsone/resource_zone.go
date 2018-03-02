package nsone

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	nsone "gopkg.in/sarguru/ns1-go.v15"
)

func zoneResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"link": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"nx_ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"refresh": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"retry": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"expiry": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"hostmaster": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dns_servers": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"networks": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Default:  "0",
			},
			"primary": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
		Create: ZoneCreate,
		Read:   ZoneRead,
		Update: ZoneUpdate,
		Delete: ZoneDelete,
	}
}

func zoneToResourceData(d *schema.ResourceData, z *nsone.Zone) {
	d.SetId(z.Id)
	d.Set("hostmaster", z.Hostmaster)
	d.Set("ttl", z.Ttl)
	d.Set("nx_ttl", z.Nx_ttl)
	d.Set("refresh", z.Refresh)
	d.Set("retry", z.Retry)
	d.Set("expiry", z.Expiry)
	d.Set("dns_servers", strings.Join(z.Dns_servers[:], ","))
	d.Set("networks", strings.Join(int2StringSlice(z.Networks)[:], ","))
	if z.Secondary != nil && z.Secondary.Enabled {
		d.Set("primary", z.Secondary.Primary_ip)
	}
	if z.Link != "" {
		d.Set("link", z.Link)
	}
}

func resourceToZoneData(z *nsone.Zone, d *schema.ResourceData) {
	z.Id = d.Id()
	if v, ok := d.GetOk("hostmaster"); ok {
		z.Hostmaster = v.(string)
	}
	if v, ok := d.GetOk("ttl"); ok {
		z.Ttl = v.(int)
	}
	if v, ok := d.GetOk("nx_ttl"); ok {
		z.Nx_ttl = v.(int)
	}
	if v, ok := d.GetOk("refresh"); ok {
		z.Refresh = v.(int)
	}
	if v, ok := d.GetOk("retry"); ok {
		z.Retry = v.(int)
	}
	if v, ok := d.GetOk("expiry"); ok {
		z.Expiry = v.(int)
	}
	if v, ok := d.GetOk("primary"); ok {
		z.MakeSecondary(v.(string))
	}
	if v, ok := d.GetOk("link"); ok {
		z.LinkTo(v.(string))
	}
	if v, ok := d.GetOk("networks"); ok {
		networkSlice := strings.Split(v.(string), ",")
		networkSliceInt := string2IntSlice(networkSlice)
		z.Networks = networkSliceInt
	}
}

func ZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	z := nsone.NewZone(d.Get("zone").(string))
	resourceToZoneData(z, d)
	if err := client.CreateZone(z); err != nil {
		return err
	}
	zoneToResourceData(d, z)
	return nil
}

func ZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	z, err := client.GetZone(d.Get("zone").(string))
	if err != nil {
		return err
	}
	zoneToResourceData(d, z)
	return nil
}

func ZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	err := client.DeleteZone(d.Get("zone").(string))
	d.SetId("")
	return err
}

func ZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	z := nsone.NewZone(d.Get("zone").(string))
	resourceToZoneData(z, d)
	if err := client.UpdateZone(z); err != nil {
		return err
	}
	zoneToResourceData(d, z)
	return nil
}

func int2StringSlice(intSl []int) []string {
	var newStringSlice []string

	for _, v := range intSl {
		newStringSlice = append(newStringSlice, strconv.Itoa(v))
	}
	return newStringSlice
}

func string2IntSlice(stringSl []string) []int {
	var newIntSlice []int

	for _, v := range stringSl {
		intV, _ := strconv.Atoi(v)
		newIntSlice = append(newIntSlice, intV)
	}
	return newIntSlice
}
