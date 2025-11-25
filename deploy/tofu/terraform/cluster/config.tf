terraform {
  required_providers {
    kafka = {
      source  = "Mongey/kafka"
    }
  }
}

variable "kafka_hosts" {
  type = string
}

provider "kafka" {
  bootstrap_servers = [var.kafka_hosts]

  tls_enabled = false
}

resource "kafka_topic" "notification-sender_sent-event" {
  name = "notification-sender.sent-event"
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
