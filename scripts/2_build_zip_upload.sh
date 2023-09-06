#aws s3 rm s3://backend-lambda-golang/v0.0.1/ --recursive --profile personaljulian
bucketname=`cat bucketname`
cd ../lambda-go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
zip lambda-go.zip main
echo $(aws s3 ls s3://backend-lambda-golang --profile personaljulian)
echo 'Please type version :'
read version
aws s3 cp lambda-go.zip s3://$bucketname/v$version/lambda-go.zip --profile personaljulian
