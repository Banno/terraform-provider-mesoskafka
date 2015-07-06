package mesoskafka

import (
	"fmt"
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

	return resourceMesosKafkaBrokerRead(d, meta)
}

func resourceMesosKafkaClusterRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMesosKafkaClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMesosKafkaClusterDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMesosKafkaBrokerRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(Client)
	status, _ := c.ApiBrokersStatus()

	fmt.Println(status)

	return nil
}
