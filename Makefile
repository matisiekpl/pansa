package:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bootstrap cmd/main.go
	zip func.zip bootstrap
	aws lambda update-function-code --function-name pansa --zip-file fileb://func.zip --region eu-central-2
