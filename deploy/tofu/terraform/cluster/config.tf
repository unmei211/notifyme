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

locals {
  kafka_hosts = split(",", var.kafka_hosts)
}

provider "kafka" {
  bootstrap_servers = local.kafka_hosts

  tls_enabled = false
}

resource "kafka_topic" "notification-sender_sent-event" {
  name = "notification-sender.sent-event"
  partitions         = 1
  replication_factor = 3
  config = {
    "cleanup.policy" = "compact"
  }
}