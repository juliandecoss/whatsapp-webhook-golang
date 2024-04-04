#aws s3 rm s3://backend-lambda-golang/v0.0.1/ --recursive --profile personaljulian
bucketname=`cat bucketname`
cd ../lambda-go
zip_name=`go list -m`
make package
echo $(aws s3 ls s3://backend-lambda-golang --profile personaljulian)
echo 'Please type version :'
read version
aws s3 cp bin/bootstrap.zip s3://$bucketname/v$version/bootstrap.zip --profile personaljulian
