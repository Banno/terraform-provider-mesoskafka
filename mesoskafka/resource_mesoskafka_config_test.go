package mesoskafka

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"testing"
)

const mesosKafkaClusterResource_basic = `
resource "mesoskafka_cluster" "broker-example" {
  broker_count = 1
}
`

const mesosKafkaClusterResource_add_brokers = `
resource "mesoskafka_cluster" "broker-example" {
  broker_count = 2
}
`
const mesosKafkaClusterResource_optionals_basic = `
resource "mesoskafka_cluster" "broker-example" {
  broker_count = 1
	cpus = 0.1
	memory = 256
	heap = 128
	jvm_options = "-Xms128m"
}
`

func TestAccMesosKafkaCluster_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDeleteCluster(),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: mesosKafkaClusterResource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccReadCluster("mesoskafka_cluster.broker-example"),
					testAccCheckBroker(1, 1, 2048, 1024, ""),
				),
			},
		},
	})
}

// func TestAccMesosKafkaCluster_addMoreBrokers(t *testing.T) {
//
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccDeleteCluster(),
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: mesosKafkaClusterResource_basic,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccReadCluster("mesoskafka_cluster.broker-example"),
// 					testAccCheckBrokerCount(1),
// 				),
// 			},
// 			resource.TestStep{
// 				Config: mesosKafkaClusterResource_add_brokers,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccReadCluster("mesoskafka_cluster.broker-example"),
// 					testAccCheckBrokerCount(2),
// 				),
// 			},
// 		},
// 	})
// }
//
// func TestAccMesosKafkaCluster_removeBrokers(t *testing.T) {
//
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccDeleteCluster(),
// 		Steps: []resource.TestStep{
// 			resource.TestStep{
// 				Config: mesosKafkaClusterResource_add_brokers,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccReadCluster("mesoskafka_cluster.broker-example"),
// 					testAccCheckBrokerCount(2),
// 				),
// 			},
// 			resource.TestStep{
// 				Config: mesosKafkaClusterResource_basic,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccReadCluster("mesoskafka_cluster.broker-example"),
// 					testAccCheckBrokerCount(1),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccMesosKafkaCluster_optionals_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDeleteCluster(),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: mesosKafkaClusterResource_optionals_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccReadCluster("mesoskafka_cluster.broker-example"),
					testAccCheckBroker(1, 0.1, 256, 128, "-Xms128m"),
				),
			},
		},
	})
}

//Helpers
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

func testAccCheckBroker(broker_count int, cpus float64, memory int, heap int, jvmOptions string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(Client)

		status, err := client.ApiBrokersStatus()

		if err != nil {
			return fmt.Errorf("Error during backends read: %v", err)
		}

		time.Sleep(5 * time.Second)

		if len(status.Brokers) != broker_count {
			return fmt.Errorf("Create Cluster Failed: wrong number of brokers %v", status.Brokers)
		}

		// TODO: figure out how to get the current state instewad of hardcoding things like cpu amounts
		// s.RootModule().Resources[""].

		for _, broker := range status.Brokers {
			if broker.Cpus != cpus {
				return fmt.Errorf("Create Cluster Failed: wrong number of cpus %v", status.Brokers)
			}

			if broker.Memory != memory {
				return fmt.Errorf("Create Cluster Failed: wrong amount of memory %v", status.Brokers)
			}

			if broker.Heap != heap {
				return fmt.Errorf("Create Cluster Failed: wrong amount of heap %v", status.Brokers)
			}

			if broker.JVMOptions != jvmOptions {
				return fmt.Errorf("Create Cluster Failed: wrong jvm-options %v", status.Brokers)
			}

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
