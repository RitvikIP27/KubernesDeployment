output "vpc_id" {
  value = aws_vpc.main.id
}

output "subnet_id" {
  value = aws_subnet.public.id
}

output "security_group_id" {
  value = aws_security_group.skillpulse_web_sg.id
}

output "instance_id" {
  value = aws_instance.app.id
}

output "instance_public_ip" {
  value = aws_instance.app.public_ip
}