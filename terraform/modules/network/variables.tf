variable "cidr_block" {
  description = "VPC CIDR range"
  type        = string
}

variable "subnet_cidr" {
  description = "Public subnet CIDR range"
  type        = string
}

variable "availability_zone" {
  type = string
}

variable "vpc_name" {
  type = string
}

variable "subnet_name" {
  type = string
}

variable "igw_name" {
  type = string
}

variable "route_table_name" {
  type = string
}
