package mesoskafka

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMesosKafkaCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceMesosKafkaClusterCreate,
		Read:   resourceMesosKafkaClusterRead,
		Update: resourceMesosKafkaClusterUpdate,
		Delete: resourceMesosKafkaClusterDelete,

		Schema: map[string]*schema.Schema{
			"broker_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceMesosKafkaClusterCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(Client)

	broker_count := d.Get("broker_count").(int)

	for i := 0; i < broker_count; i++ {
		_, err := c.ApiBrokersAdd(i)

		if err != nil {
			panic(err)
		}

		_, err = c.ApiBrokersStart(i)

		if err != nil {
			panic(err)
		}
	}

	status, _ := c.ApiBrokersStatus()
	// TODO: this should probably not be this. just a placeholder
	d.SetId(status.FrameworkID)

	return resourceMesosKafkaClusterRead(d, meta)
}

func resourceMesosKafkaClusterRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(Client)
	status, err := c.ApiBrokersStatus()
	if err != nil {
		return err
	}

	d.Set("broker_count", len(status.Brokers))

	return nil
}

func resourceMesosKafkaClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMesosKafkaClusterDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(Client)

	status, _ := c.ApiBrokersStatus()

	for i := 0; i < len(status.Brokers); i++ {
		_, err := c.ApiBrokersStop(i)

		if err != nil {
			return err
		}

		_, err = c.ApiBrokersRemove(i)
		if err != nil {
			return err
		}
	}

	return nil
}
