
provider "aws" {
  region     = "${var.region}"
}


resource "aws_eip" "ip" {
  instance = "${aws_instance.cb1.id}"
  vpc = true
}

resource "aws_security_group" "default" {
  name        = "eip_cb1"
  description = "Used in the terraform"

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP access from anywhere
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP access to Couchbase
  ingress {
    from_port   = 8091
    to_port     = 8091
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8092
    to_port     = 8092
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 11207
    to_port     = 11207
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 11210
    to_port     = 11210
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 11211
    to_port     = 11211
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

   

  # outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "cb1" {
  ami           = "ami-076e276d85f524150"
  #instance_type = "t2.large"
  instance_type = "t3.xlarge"
  key_name = "${var.key_name}"
  security_groups = ["${aws_security_group.default.name}"]

  user_data = "${file("scripts/cb_install.sh")}"

  #Instance tags
  tags {
    Name = "dev-cb1"
  }
}


output "cb1_ip" {
  value = "${aws_eip.ip.public_ip}"
}


