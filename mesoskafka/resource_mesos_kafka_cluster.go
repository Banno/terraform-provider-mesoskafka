package mesoskafka

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
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

	broker_ids := []int{}
	for i := 0; i < broker_count; i++ {
		broker_ids = append(broker_ids, i)
	}

	err := c.ApiBrokersCreate(broker_ids)

	if err != nil {
		return err
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
	c := meta.(Client)
	status, err := c.ApiBrokersStatus()
	if err != nil {
		return err
	}

	broker_count := d.Get("broker_count").(int)
	current_count := len(status.Brokers)

	if current_count > broker_count {
		// remove some brokers
		how_many := current_count - broker_count
		fmt.Println(how_many)
	} else if broker_count > current_count {
		// add some brokers
		how_many := broker_count - current_count

		max_broker_id := 0

		for _, broker := range status.Brokers {
			_id, _ := strconv.ParseInt(broker.Id, 10, 0)
			id := int(_id)
			if id > max_broker_id {
				max_broker_id = id
			}
		}

		broker_ids := []int{}
		for j := max_broker_id + 1; j <= max_broker_id+how_many; j++ {
			broker_ids = append(broker_ids, j)
		}

		err := c.ApiBrokersCreate(broker_ids)
		if err != nil {
			return err
		}

	} else {
		fmt.Println("Broker counts are the same")
	}

	return nil
}

func resourceMesosKafkaClusterDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(Client)

	status, _ := c.ApiBrokersStatus()

	broker_ids := []int{}
	for i := 0; i < len(status.Brokers); i++ {
		broker_ids = append(broker_ids, i)
	}

	err := c.ApiBrokersDelete(broker_ids)
	if err != nil {
		return err
	}

	return nil
}
