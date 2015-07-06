package mesoskafka

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"
)

const basicMesosKafkaClusterResource = `
resource "mesoskafka_cluster" "broker-example" {
   broker_count = 1
}
`

func TestAccMesosKafkaCluster_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDeleteCluster(),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: basicMesosKafkaClusterResource,
				Check: resource.ComposeTestCheckFunc(
					testAccReadCluster("mesoskafka_cluster.broker-example"),
					testAccCreateCluster(),
				),
			},
		},
	})
}

func testAccReadCluster(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("mesoskafka_cluster resource not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("mesoskafka_cluster resource id not set correctly: %s", name)
		}

		return nil
	}
}

func testAccCreateCluster() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(Client)

		status, err := client.ApiBrokersStatus()

		if err != nil {
			return fmt.Errorf("Error during backends read: %v", err)
		}

		time.Sleep(5 * time.Second)

		if len(status.Brokers) != 1 {
			return fmt.Errorf("Create Cluster Failed: wrong number of brokers %v", status.Brokers)
		}

		return nil
	}
}

func testAccDeleteCluster() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(Client)

		status, err := client.ApiBrokersStatus()

		if err != nil {
			return fmt.Errorf("Error during backends read: %v", err)
		}

		if len(status.Brokers) != 0 {
			return fmt.Errorf("Cluster not deleted!")
		}

		return nil
	}
}
