module "dynamodb_schemas_table" {
  source  = "terraform-aws-modules/dynamodb-table/aws"
  version = "4.0.0"

  name         = "survey-invitations"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "phoneNumber"
  #range_key    = "CreationDate"

  attributes = [
    {
      name = "phoneNumber"
      type = "S"
    }
  ]
}

# module "dynamodb_topics_table" {
#   source  = "terraform-aws-modules/dynamodb-table/aws"
#   version = "3.1.2"

#   tags         = local.tagging
#   name         = local.naming.dynamo_topics
#   billing_mode = "PAY_PER_REQUEST"
#   hash_key     = "Name"
#   range_key    = "CreationDate"

#   attributes = [
#     {
#       name = "Name"
#       type = "S"
#     },
#     {
#       name = "CreationDate"
#       type = "S"
#     }
#   ]
# }
