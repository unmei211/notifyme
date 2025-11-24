terraform {
  required_providers {
    kafka = {
      source  = "Mongey/kafka"
    }
  }
}

provider "kafka" {
  bootstrap_servers = ["infra-kafka:9092"]

  tls_enabled = false
}

resource "kafka_topic" "payment_create-command" {
  name = "payment.create-command"
  partitions         = 1
  replication_factor = 1
  config = {
    "cleanup.policy" = "compact"
  }
}

resource "kafka_topic" "payment_created-event" {
  name = "payment.created-event"
  partitions         = 1
  replication_factor = 1
  config = {
    "cleanup.policy" = "compact"
  }
}


resource "kafka_topic" "subscription_charge-requested" {
  name               = "subscription.charge-requested"
  partitions         = 1
  replication_factor = 1
  config = {
    "cleanup.policy" = "compact"
  }
}

resource "kafka_topic" "subscription_created" {
  name               = "subscription.created"
  partitions         = 1
  replication_factor = 1
  config = {
    "cleanup.policy" = "compact"
  }
}

resource "kafka_topic" "subscription_disabled" {
  name               = "subscription.disabled"
  partitions         = 1
  replication_factor = 1
  config = {
    "cleanup.policy" = "compact"
  }
}
