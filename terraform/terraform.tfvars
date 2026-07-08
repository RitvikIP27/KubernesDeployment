
aws_region = "ap-south-1"

vpc_cidr           = "10.0.0.0/16"
public_subnet_cidr = "10.0.20.0/24"
availability_zone  = "ap-south-1a"



ami           = "ami-01a00762f46d584a1"
instance_type = "t3.small"
key_name      = "skillpulse-newkey"



vpc_name            = "helixacore-vpc"
subnet_name         = "helixacore-public-subnet"
igw_name            = "igw-1"
route_table_name    = "helixacore-route-table"
instance_name       = "helixacore-server-1"
security_group_name = "helixacore-sg"

environment  = "dev"
project_name = "helixacore"
owner        = "Ritvik Kant"
elastic_ip_allocation_id="eipalloc-01f344abe0bdb363a"