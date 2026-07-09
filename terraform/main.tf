module "network" {
  source = "./modules/network"

  cidr_block        = var.vpc_cidr
  subnet_cidr       = var.public_subnet_cidr
  availability_zone = var.availability_zone

  vpc_name         = var.vpc_name
  subnet_name      = var.subnet_name
  igw_name         = var.igw_name
  route_table_name = var.route_table_name
}

module "compute" {
  source = "./modules/compute/ec2"

  vpc_id    = module.network.vpc_id
  subnet_id = module.network.public_subnet_id

  ami           = var.ami
  instance_type = var.instance_type
  key_name      = var.key_name


  instance_name            = var.instance_name
  security_group_name      = var.security_group_name
  elastic_ip_allocation_id = var.elastic_ip_allocation_id
}

