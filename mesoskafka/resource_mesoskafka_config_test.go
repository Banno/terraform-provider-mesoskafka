package mesoskafka

import (
	"fmt"

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
					testAccReadBrokers("mesoskafka_cluster.broker-example"),
					// testCheckCreate(backendName, &backendsResult),
					// testIfBackendIsPublic(backendName, &backendsResult, false),
				),
			},
		},
	})
}

func testAccReadBrokers(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("mesoskafka_cluster resource not found: %s", name)
		}

		client := testAccProvider.Meta().(*Client)
		fmt.Println(client)

		// backendRead, err := client.BackendsRead()
		// if err != nil {
		// 	return fmt.Errorf("Error during backends read: %v", err)
		// }

		// time.Sleep(5 * time.Second)

		// *backendsResult = *backendRead

		return nil
	}
}
