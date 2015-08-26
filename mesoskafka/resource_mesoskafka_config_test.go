package mesoskafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
	logfourj_options = "file:log4j.properties"
	options = "log.dirs=/tmp/kafka/"
	constraints = "hostname=unique"
	failover_delay = "17s"
	failover_max_delay = "14m"
	failover_max_tries = 42
}
`

const mesosKafkaClusterResource_optionals_basic_update = `
resource "mesoskafka_cluster" "broker-example" {
  broker_count = 1
	cpus = 1
	memory = 8096
	heap = 256
	jvm_options = ""
	logfourj_options = ""
	options = ""
	constraints = ""
	failover_delay = "99s"
	failover_max_delay = "5m"
	failover_max_tries = 5
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
					testAccCheckBrokerCount(1),
					testAccCheckBrokerAttributes_basic(),
				),
			},
		},
	})
}

func TestAccMesosKafkaCluster_addMoreBrokers(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDeleteCluster(),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: mesosKafkaClusterResource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccReadCluster("mesoskafka_cluster.broker-example"),
					testAccCheckBrokerCount(1),
				),
			},
			resource.TestStep{
				Config: mesosKafkaClusterResource_add_brokers,
				Check: resource.ComposeTestCheckFunc(
					testAccReadCluster("mesoskafka_cluster.broker-example"),
					testAccCheckBrokerCount(2),
				),
			},
		},
	})
}

func TestAccMesosKafkaCluster_removeBrokers(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDeleteCluster(),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: mesosKafkaClusterResource_add_brokers,
				Check: resource.ComposeTestCheckFunc(
					testAccReadCluster("mesoskafka_cluster.broker-example"),
					testAccCheckBrokerCount(2),
				),
			},
			resource.TestStep{
				Config: mesosKafkaClusterResource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccReadCluster("mesoskafka_cluster.broker-example"),
					testAccCheckBrokerCount(1),
				),
			},
		},
	})
}

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
					testAccCheckBrokerCount(1),
					testAccCheckBrokerAttributes_optionals(),
				),
			},
		},
	})
}

func TestAccMesosKafkaCluster_optionals_basic_update(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDeleteCluster(),
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: mesosKafkaClusterResource_optionals_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccReadCluster("mesoskafka_cluster.broker-example"),
					testAccCheckBrokerCount(1),
					testAccCheckBrokerAttributes_optionals(),
				),
			},
			resource.TestStep{
				Config: mesosKafkaClusterResource_optionals_basic_update,
				Check: resource.ComposeTestCheckFunc(
					testAccReadCluster("mesoskafka_cluster.broker-example"),
					testAccCheckBrokerCount(1),
					testAccCheckBrokerAttributes_optionals_update(),
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

func testAccCheckBrokerCount(broker_count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		client := testAccProvider.Meta().(Client)

		status, err := client.ApiBrokersStatus()
		if err != nil {
			panic(fmt.Errorf("Error during backends read: %v", err))
		}

		if len(status.Brokers) != broker_count {
			return fmt.Errorf("Create Cluster Failed: wrong number of brokers %#v", status.Brokers)
		}
		return nil
	}
}

func testAccCheckBrokerAttributes_basic() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// TODO: figure out how to get the current state instead of hardcoding things like cpu amounts
		// s.RootModule().Resources[""].

		client := testAccProvider.Meta().(Client)

		status, err := client.ApiBrokersStatus()
		if err != nil {
			panic(fmt.Errorf("Error during backends read: %v", err))
		}

		for _, broker := range status.Brokers {
			if broker.Cpus != float64(1) {
				return fmt.Errorf("Create Cluster Failed: wrong number of cpus %#v", status.Brokers)
			}

			if broker.Memory != int(2048) {
				return fmt.Errorf("Create Cluster Failed: wrong amount of memory %#v", status.Brokers)
			}

			if broker.Heap != int(1024) {
				return fmt.Errorf("Create Cluster Failed: wrong amount of heap #%v", status.Brokers)
			}

		}

		return nil
	}
}

func testAccCheckBrokerAttributes_optionals() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// TODO: figure out how to get the current state instead of hardcoding things like cpu amounts
		// s.RootModule().Resources[""].

		client := testAccProvider.Meta().(Client)

		status, err := client.ApiBrokersStatus()
		if err != nil {
			panic(fmt.Errorf("Error during backends read: %#v", err))
		}

		for _, broker := range status.Brokers {
			if broker.Cpus != float64(0.1) {
				return fmt.Errorf("Create Cluster Failed: wrong number of cpus %#v", status.Brokers)
			}

			if broker.Memory != int(256) {
				return fmt.Errorf("Create Cluster Failed: wrong amount of memory %#v", status.Brokers)
			}

			if broker.Heap != int(128) {
				return fmt.Errorf("Create Cluster Failed: wrong amount of heap %#v", status.Brokers)
			}

			if broker.JVMOptions != "-Xms128m" {
				return fmt.Errorf("Create Cluster Failed: wrong jvm-options %#v", status.Brokers)
			}

			if broker.Log4jOptions != "file:log4j.properties" {
				return fmt.Errorf("Create Cluster Failed: wrong logfourj_options %#v", status.Brokers)
			}

			if broker.Options != "log.dirs=/tmp/kafka/" {
				return fmt.Errorf("Create Cluster Failed: wrong options %#v", status.Brokers)
			}

			if broker.Constraints != "hostname=unique" {
				return fmt.Errorf("Create Cluster Failed: wrong constraints %#v", status.Brokers)
			}

			if broker.Failover.Delay != "17s" {
				return fmt.Errorf("Create Cluster Failed: wrong failover-delay %#v", status.Brokers)
			}

			if broker.Failover.MaxDelay != "14m" {
				return fmt.Errorf("Create Cluster Failed: wrong failover-max-delay %#v", status.Brokers)
			}

			if broker.Failover.MaxTries != int(42) {
				return fmt.Errorf("Create Cluster Failed: wrong failover-max-tries %#v", status.Brokers)
			}

		}

		return nil
	}
}

func testAccCheckBrokerAttributes_optionals_update() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// TODO: figure out how to get the current state instead of hardcoding things like cpu amounts
		// s.RootModule().Resources[""].

		client := testAccProvider.Meta().(Client)

		status, err := client.ApiBrokersStatus()
		if err != nil {
			panic(fmt.Errorf("Error during backends read: %#v", err))
		}

		for _, broker := range status.Brokers {
			if broker.Cpus != float64(1) {
				return fmt.Errorf("Create Cluster Failed: wrong number of cpus %#v", status.Brokers)
			}

			if broker.Memory != int(8096) {
				return fmt.Errorf("Create Cluster Failed: wrong amount of memory %#v", status.Brokers)
			}

			if broker.Heap != int(256) {
				return fmt.Errorf("Create Cluster Failed: wrong amount of heap %#v", status.Brokers)
			}

			if broker.JVMOptions != "" {
				return fmt.Errorf("Create Cluster Failed: wrong jvm-options %#v", status.Brokers)
			}

			if broker.Log4jOptions != "" {
				return fmt.Errorf("Create Cluster Failed: wrong logfourj_options %#v", status.Brokers)
			}

			if broker.Options != "" {
				return fmt.Errorf("Create Cluster Failed: wrong options %#v", status.Brokers)
			}

			if broker.Constraints != "" {
				return fmt.Errorf("Create Cluster Failed: wrong constraints %#v", status.Brokers)
			}

			if broker.Failover.Delay != "99s" {
				return fmt.Errorf("Create Cluster Failed: wrong failover-delay %#v", status.Brokers)
			}

			if broker.Failover.MaxDelay != "5m" {
				return fmt.Errorf("Create Cluster Failed: wrong failover-max-delay %#v", status.Brokers)
			}

			if broker.Failover.MaxTries != int(5) {
				return fmt.Errorf("Create Cluster Failed: wrong failover-max-tries %#v", status.Brokers)
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
