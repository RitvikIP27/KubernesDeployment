variable "aws_region" {
  description = "AWS Region"
  type        = string
  default     = "ap-south-1"
}

variable "vpc_cidr" {
  type = string
}

variable "public_subnet_cidr" {
  type = string
}

variable "availability_zone" {
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

variable "associate_public_ip_address" {
  description = "Whether to assign a public IP to the EC2 instance at launch"
  type        = bool
  default     = true
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

variable "instance_name" {
  type = string
}

variable "security_group_name" {
  type = string
}



variable "elastic_ip_allocation_id" {
  description = "Optional allocation ID of a manually created Elastic IP"
  type        = string
  default     = null
}

variable "environment" {
  type = string
}

variable "project_name" {
  type = string
}

variable "owner" {
  type = string
}
