resource "aws_instance" "selefra_instance" {
  ami                    = "ami-09208e69ff3feb1db"
  instance_type          = "t2.micro"
  subnet_id              = aws_subnet.selefra_subnet.id
  vpc_security_group_ids = [aws_security_group.selefra_security_group.id]
  depends_on = [
    aws_security_group.selefra_security_group
  ]
}

resource "aws_vpc" "selefra_vpc" {
  cidr_block = "10.7.0.0/16"
}

resource "aws_subnet" "selefra_subnet" {
  vpc_id     = aws_vpc.selefra_vpc.id
  cidr_block = "10.7.7.0/24"
}

resource "aws_security_group" "selefra_security_group" {
  name = "selefra_security_group"
}

resource "aws_security_group_rule" "open_22_port" {
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  ipv6_cidr_blocks  = []
  security_group_id = aws_security_group.selefra_security_group.id
}
