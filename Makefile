.PHONY: deps clean build_and_deploy build package deploy local build_and_run

deps:
	go get -u ./...

clean: 
	rm -rf ./go-path-router/go-path-router

build_and_deploy: build package deploy

build:
	GOOS=linux GOARCH=amd64 go build -o go-path-router/go-path-router ./go-path-router

package:
	sam package --output-template-file packaged.yaml --s3-bucket ${S3_BUCKET}

deploy:
	sam deploy --template-file packaged.yaml --stack-name ${STACK_NAME} --capabilities CAPABILITY_IAM --region ${AWS_REGION} --parameter-overrides CertificateArn=${ACM_CERT} CodePath=${CODE_PATH} Domain=${DomainName} HostedZoneName=${HOSTED_ZONE}

build_and_run: build local

local:
	sam local start-api