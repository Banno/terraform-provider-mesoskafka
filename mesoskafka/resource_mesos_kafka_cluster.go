package mesoskafka

import "github.com/hashicorp/terraform/helper/schema"

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
	return nil
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
