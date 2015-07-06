package mesoskafka

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MESOS_KAFKA_URL", nil),
				Description: "Mesos Kafka Scheduler Api Url",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"mesoskafka_cluster": resourceMesosKafkaCluster(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Url: d.Get("url").(string),
	}

	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}

	return config.client, nil
}
