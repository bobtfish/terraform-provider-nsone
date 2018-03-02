package nsone

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	nsone "gopkg.in/sarguru/ns1-go.v15"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func recordResource() *schema.Resource {
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
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !regexp.MustCompile(`^(A|AAAA|ALIAS|AFSDB|CNAME|DNAME|HINFO|MX|NAPTR|NS|PTR|RP|SPF|SRV|TXT)$`).MatchString(value) {
						es = append(es, fmt.Errorf(
							"only A, AAAA, ALIAS, AFSDB, CNAME, DNAME, HINFO, MX, NAPTR, NS, PTR, RP, SPF, SRV, TXT allowed in %q", k))
					}
					return
				},
			},
			"meta": metaSchema(),
			"link": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"use_client_subnet": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"answers": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"answer": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"meta": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"feed": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										//ConflictsWith: []string{"value"},
									},
									"value": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										//ConflictsWith: []string{"feed"},
									},
								},
							},
							Set: metaToHash,
						},
					},
				},
				Set: answersToHash,
			},
			"regions": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"georegion": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
								value := v.(string)
								if !regexp.MustCompile(`^(US-WEST|US-EAST|US-CENTRAL|EUROPE|AFRICA|ASIAPAC|SOUTH-AMERICA)$`).MatchString(value) {
									es = append(es, fmt.Errorf(
										"only US-WEST, US-EAST, US-CENTRAL, EUROPE, AFRICA, ASIAPAC, SOUTH-AMERICA allowed in %q", k))
								}
								return
							},
						},
						"country": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"us_state": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"longitude": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"latitude": &schema.Schema{
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"up": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
				Set: regionsToHash,
			},
			"filters": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"disabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"config": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},
		},
		Create: RecordCreate,
		Read:   RecordRead,
		Update: RecordUpdate,
		Delete: RecordDelete,
	}
}

func regionsToHash(v interface{}) int {
	var buf bytes.Buffer
	r := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", r["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", r["georegion"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", r["country"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", r["us_state"].(string)))
	buf.WriteString(fmt.Sprintf("%f-", r["latitude"].(float64)))
	buf.WriteString(fmt.Sprintf("%f-", r["longitude"].(float64)))
	buf.WriteString(fmt.Sprintf("%t-", r["up"].(bool)))
	return hashcode.String(buf.String())
}

func answersToHash(v interface{}) int {
	var buf bytes.Buffer
	a := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", a["answer"].(string)))
	if a["region"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", a["region"].(string)))
	}
	metas := make([]int, 0)
	switch t := a["meta"].(type) {
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	case *schema.Set:
		for _, meta := range t.List() {
			metas = append(metas, metaToHash(meta))
		}
	case []map[string]interface{}:
		for _, meta := range t {
			metas = append(metas, metaToHash(meta))
		}
	}
	sort.Ints(metas)
	for _, metahash := range metas {
		buf.WriteString(fmt.Sprintf("%d-", metahash))
	}
	hash := hashcode.String(buf.String())
	log.Println("Generated answersToHash %d from %+v", hash, a)
	return hash
}

func metaToHash(v interface{}) int {
	var buf bytes.Buffer
	s := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", s["field"].(string)))
	if v, ok := s["feed"]; ok && v.(string) != "" {
		buf.WriteString(fmt.Sprintf("feed%s-", v.(string)))
	}
	if v, ok := s["value"]; ok && v.(string) != "" {
		buf.WriteString(fmt.Sprintf("value%s-", v.(string)))
	}

	hash := hashcode.String(buf.String())
	log.Println("Generated metaToHash %d from %+v", hash, s)
	return hash
}

func recordToResourceData(d *schema.ResourceData, r *nsone.Record) error {
	d.SetId(r.Id)
	d.Set("domain", r.Domain)
	d.Set("zone", r.Zone)
	d.Set("type", r.Type)
	d.Set("ttl", r.Ttl)
	if r.Link != "" {
		d.Set("link", r.Link)
	}
	if len(r.Filters) > 0 {
		filters := make([]map[string]interface{}, len(r.Filters))
		for i, f := range r.Filters {
			m := make(map[string]interface{})
			m["filter"] = f.Filter
			if f.Disabled {
				m["disabled"] = true
			}
			if f.Config != nil {
				m["config"] = f.Config
			}
			filters[i] = m
		}
		d.Set("filters", filters)
	}
	if len(r.Answers) > 0 {
		ans := &schema.Set{
			F: answersToHash,
		}
		log.Printf("Got back from nsone answers: %+v", r.Answers)
		for _, answer := range r.Answers {
			ans.Add(answerToMap(answer))
		}
		log.Printf("Setting answers %+v", ans)
		err := d.Set("answers", ans)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting answers for: %s, error: %#v", r.Domain, err)
		}
	}
	if len(r.Regions) > 0 {
		regions := make([]map[string]interface{}, 0, len(r.Regions))
		for region_name, region := range r.Regions {
			new_region := make(map[string]interface{})
			new_region["name"] = region_name
			if len(region.Meta.GeoRegion) > 0 {
				new_region["georegion"] = region.Meta.GeoRegion[0]
			}
			if len(region.Meta.Country) > 0 {
				new_region["country"] = region.Meta.Country[0]
			}
			if len(region.Meta.USState) > 0 {
				new_region["us_state"] = region.Meta.USState[0]
			}
			if region.Meta.Latitude > 0 {
				new_region["latitude"] = strconv.FormatFloat(region.Meta.Latitude, 'f', 1, 64)
			}
			if region.Meta.Longitude > 0 {
				new_region["longitude"] = strconv.FormatFloat(region.Meta.Longitude, 'f', 1, 64)
			}
			if region.Meta.Up {
				new_region["up"] = region.Meta.Up
			} else {
				new_region["up"] = false
			}
			regions = append(regions, new_region)
		}
		log.Printf("Setting regions %+v", regions)
		err := d.Set("regions", regions)
		if err != nil {
			return fmt.Errorf("[DEBUG] Error setting regions for: %s, error: %#v", r.Domain, err)
		}
	}
	return nil
}

func answerToMap(a nsone.Answer) map[string]interface{} {
	m := make(map[string]interface{})
	m["meta"] = make([]map[string]interface{}, 0)
	m["answer"] = strings.Join(a.Answer, " ")
	if a.Region != "" {
		m["region"] = a.Region
	}
	if a.Meta != nil {
		metas := &schema.Set{
			F: metaToHash,
		}
		for k, v := range a.Meta {
			meta := make(map[string]interface{})
			meta["field"] = k
			switch t := v.(type) {
			case map[string]interface{}:
				meta["feed"] = t["feed"].(string)
			case string:
				meta["value"] = t
			case []interface{}:
				var val_array []string
				for _, pref := range t {
					val_array = append(val_array, pref.(string))
				}
				sort.Strings(val_array)
				string_val := strings.Join(val_array, ",")
				meta["value"] = string_val
			case bool:
				int_val := btoi(t)
				meta["value"] = strconv.Itoa(int_val)
			case float64:
				int_val := int(t)
				meta["value"] = strconv.Itoa(int_val)
			}
			metas.Add(meta)
		}
		m["meta"] = metas
	}
	return m
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func resourceDataToRecord(r *nsone.Record, d *schema.ResourceData) error {
	r.Id = d.Id()
	if answers := d.Get("answers").(*schema.Set); answers.Len() > 0 {
		al := make([]nsone.Answer, answers.Len())
		for i, answer_raw := range answers.List() {
			answer := answer_raw.(map[string]interface{})
			a := nsone.NewAnswer()
			v := answer["answer"].(string)
			if d.Get("type") != "TXT" {
				a.Answer = strings.Split(v, " ")
			} else {
				a.Answer = []string{v}
			}
			if v, ok := answer["region"]; ok {
				a.Region = v.(string)
			}
			if metas := answer["meta"].(*schema.Set); metas.Len() > 0 {
				for _, meta_raw := range metas.List() {
					meta := meta_raw.(map[string]interface{})
					key := meta["field"].(string)
					if value, ok := meta["feed"]; ok && value.(string) != "" {
						a.Meta[key] = nsone.NewMetaFeed(value.(string))
					}
					if value, ok := meta["value"]; ok && value.(string) != "" {
						meta_array := strings.Split(value.(string), ",")
						if len(meta_array) > 1 {
							sort.Strings(meta_array)
							a.Meta[key] = meta_array
						} else {
							a.Meta[key] = value.(string)
						}
					}
				}
			}
			al[i] = a
		}
		r.Answers = al
		if _, ok := d.GetOk("link"); ok {
			return errors.New("Cannot have both link and answers in a record")
		}
	}
	if v, ok := d.GetOk("ttl"); ok {
		r.Ttl = v.(int)
	}
	if v, ok := d.GetOk("link"); ok {
		r.LinkTo(v.(string))
	}
	useClientSubnetVal := d.Get("use_client_subnet").(bool)
	if v := strconv.FormatBool(useClientSubnetVal); v != "" {
		r.UseClientSubnet = useClientSubnetVal
	}

	if rawFilters := d.Get("filters").([]interface{}); len(rawFilters) > 0 {
		f := make([]nsone.Filter, len(rawFilters))
		for i, filter_raw := range rawFilters {
			fi := filter_raw.(map[string]interface{})
			config := make(map[string]interface{})
			filter := nsone.Filter{
				Filter: fi["filter"].(string),
				Config: config,
			}
			if disabled, ok := fi["disabled"]; ok {
				filter.Disabled = disabled.(bool)
			}
			if raw_config, ok := fi["config"]; ok {
				for k, v := range raw_config.(map[string]interface{}) {
					if v.(string) == "true" || v.(string) == "false" {
						filter.Config[k], _ = strconv.ParseBool(v.(string))
					} else {
						if i, err := strconv.Atoi(v.(string)); err == nil {
							filter.Config[k] = i
						} else {
							filter.Config[k] = v
						}
					}
				}
			}
			f[i] = filter
		}
		r.Filters = f
	}
	if regions := d.Get("regions").(*schema.Set); regions.Len() > 0 {
		rm := make(map[string]nsone.Region)
		for _, region_raw := range regions.List() {
			region := region_raw.(map[string]interface{})
			nsone_r := nsone.Region{
				Meta: nsone.RegionMeta{},
			}
			if g := region["georegion"].(string); g != "" {
				nsone_r.Meta.GeoRegion = []string{g}
			}
			if g := region["country"].(string); g != "" {
				nsone_r.Meta.Country = []string{g}
			}
			if g := region["us_state"].(string); g != "" {
				nsone_r.Meta.USState = []string{g}
			}
			if g := region["latitude"].(float64); g != 0 {
				nsone_r.Meta.Latitude = g
			}
			if g := region["longitude"].(float64); g != 0 {
				nsone_r.Meta.Longitude = g
			}
			if g := region["up"].(bool); g {
				nsone_r.Meta.Up = g
			}

			rm[region["name"].(string)] = nsone_r
		}
		r.Regions = rm
	}
	return nil
}

func setToMapByKey(s *schema.Set, key string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, rawData := range s.List() {
		data := rawData.(map[string]interface{})
		result[data[key].(string)] = data
	}

	return result
}

func RecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	r := nsone.NewRecord(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err := resourceDataToRecord(r, d); err != nil {
		return err
	}
	if err := client.CreateRecord(r); err != nil {
		return err
	}
	return recordToResourceData(d, r)
}

func RecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	r, err := client.GetRecord(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err != nil {
		return err
	}
	recordToResourceData(d, r)
	return nil
}

func RecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	err := client.DeleteRecord(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	d.SetId("")
	return err
}

func RecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*nsone.APIClient)
	r := nsone.NewRecord(d.Get("zone").(string), d.Get("domain").(string), d.Get("type").(string))
	if err := resourceDataToRecord(r, d); err != nil {
		return err
	}
	if err := client.UpdateRecord(r); err != nil {
		return err
	}
	recordToResourceData(d, r)
	return nil
}
