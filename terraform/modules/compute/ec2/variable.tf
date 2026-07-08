variable "vpc_id" {
  type = string
}

variable "subnet_id" {
  type = string
}

variable "ami" {
  type = string
}

variable "instance_type" {
  type = string
}

variable "key_name" {
  type = string
}

variable "instance_name" {
  type = string
}

variable "security_group_name" {
  type = string
}

variable "elastic_ip_allocation_id" {
  description = "Existing Elastic IP Allocation ID"
  type        = string
  default     = null
}