# DataStream  API lab
https://docs.aws.amazon.com/managed-flink/latest/java/getting-started.html

## lab steps
1. Create kinesis datastream named ExampleInputStream and ExampleOutputStream using CDK code
2. Create s3 bucket using CDK code
3. Create stock.py to generate the source events to ExampleInputStream
4. Copmile the application code 
```
git clone https://github.com/aws-samples/amazon-managed-service-for-apache-flink-examples.git

cp ~/aws-samples/amazon-managed-service-for-apache-flink-examples/java/GettingStarted/* ~/datastream-api-lab/code/datastream-api

```

5. Copy the JAR to s3 bucket

```
aws s3 cp target/amazon-msf-java-stream-app-1.0.jar s3://datastreamapilabstack-flinklabapplicationbucket/flinkapps/steamapiapp/amazon-msf-java-stream-app-1.0.jar
```

6. Create application
```
AWS_REGION='ap-southeast-1' ./flinkApplication --action create \
--applicationName steamapiapp \
--codeS3Bucket arn:aws:s3:::datastreamapilabstack-flinklabapplicationbucket \
--codeS3FileKey flinkapps/steamapiapp/amazon-msf-java-stream-app-1.0.jar \
--applicationRole arn:aws:iam::123456799012:role/temp-flink-application  
```

7. Start application
```
AWS_REGION='ap-southeast-1' ./flinkApplication --action start \
--applicationName  steamapiapp
```

8. Describe application
```
AWS_REGION='ap-southeast-1' ./flinkApplication --action describe \
--applicationName  steamapiapp
```

9. Stop application
```
AWS_REGION='ap-southeast-1' ./flinkApplication --action stop \
--applicationName  steamapiapp
```

10. Create snapshot
```
AWS_REGION='ap-southeast-1' ./flinkApplication --action create-snapshot \
--applicationName  steamapiapp \
--snapshotName snapshot1
```