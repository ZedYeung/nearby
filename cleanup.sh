#!/bin/bash
terraform destroy ./terraform
rm -rf ./terraform/vars.tf
rm -rf docker-compose.yml
rm -rf ingress.yml
echo "Done"