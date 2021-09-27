package nsone

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	nsone "gopkg.in/sarguru/ns1-go.v18"
)

func dataFeedResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
		Create: DataFeedCreate,
		Read:   DataFeedRead,
		Update: DataFeedUpdate,
		Delete: DataFeedDelete,
	}
}

func dataFeedToResourceData(d *schema.ResourceData, df *nsone.DataFeed) {
	d.SetId(df.Id)
	d.Set("name", df.Name)
	d.Set("config", df.Config)
}

func resourceDataToDataFeed(d *schema.ResourceData) *nsone.DataFeed {
	df := nsone.NewDataFeed(d.Get("source_id").(string))
	df.Name = d.Get("name").(string)
	config := make(map[string]string)
	for k, v := range d.Get("config").(map[string]interface{}) {
		config[k] = v.(string)
	}
	df.Config = config
	return df
}

func DataFeedCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	df := resourceDataToDataFeed(d)
	if err := client.CreateDataFeed(df); err != nil {
		return err
	}
	dataFeedToResourceData(d, df)
	return nil
}

func DataFeedRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	df, err := client.GetDataFeed(d.Get("source_id").(string), d.Id())
	if err != nil {
		return err
	}
	dataFeedToResourceData(d, df)
	return nil
}

func DataFeedDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	err := client.DeleteDataFeed(d.Get("source_id").(string), d.Id())
	d.SetId("")
	return err
}

func DataFeedUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	df := resourceDataToDataFeed(d)
	df.Id = d.Id()
	if err := client.UpdateDataFeed(df); err != nil {
		return err
	}
	dataFeedToResourceData(d, df)
	return nil
}
