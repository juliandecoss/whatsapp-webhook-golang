echo "Please type your bucket name"
read bucketname
cat > bucketname <<EOF
$bucketname
EOF
bucketname=`cat bucketname`
aws s3api create-bucket --bucket $bucketname --region us-west-2 --create-bucket-configuration LocationConstraint=us-west-2 --profile personaljulian
cd ../terraform
cat > variables.tf <<EOF
variable "app_version" {}

variable "s3_bucket" {
  default = "$bucketname"
}

variable "file_zip_name" {
  default = ""
}
EOF

