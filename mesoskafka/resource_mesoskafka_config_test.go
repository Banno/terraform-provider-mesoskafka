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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckBackendDelete(backendName, &backendsResult),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: basicMesosKafkaClusterResource,
				Check: resource.ComposeTestCheckFunc(
					testAccReadCluster("mesoskafka_cluster.broker-example"),
					// testCheckCreate(backendName, &backendsResult),
					// testIfBackendIsPublic(backendName, &backendsResult, false),
				),
			},
		},
	})
}

func testAccReadCluster(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("mesoskafka_cluster resource not found: %s", name)
		}

		client := testAccProvider.Meta().(Client)

		status, err := client.ApiBrokersStatus()

		if err != nil {
			return fmt.Errorf("Error during backends read: %v", err)
		}

		time.Sleep(5 * time.Second)

		if len(status.Brokers) != 1 {
			return fmt.Errorf("Add Brokers Failed: wrong number of brokers %v", status.Brokers)
		}

		return nil
	}
}
