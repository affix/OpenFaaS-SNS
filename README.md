# OpenFaaS SNS Subscriber

This project is an [OpenFaaS](http://openfaas.com) function to handle [AWS SNS](https://aws.amazon.com/sns/)
notifications and subscriptions.

## How do I set this up

- Deploy the function with `faas-cli deploy -f openfaas-sns.yml`
- Setup an AWS Topic
- Add an HTTP(s) Subscriber pointing to your function

## Caveats

- We do not calculate the Signatures of the messages
- The function return is not visible to SNS, this is intended to be used to trigger other
functions or do background processing. (A simple usecase would be to subscribe to
 S3 Bucket changes, Then process the files uploaded to s3).
