# Flink Application Handling Logic

## Application path

```
s3://<BucketName>/flinkapps/<appname>
```

## Application role

### Permissions
TBD

#### CW Log
TBD

#### CW Metrics
TBD

#### Code Source （S3）
**粗粒度**
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "s3:DeleteObject",
                "s3:GetObject",
                "s3:ListBucket",
                "s3:PutObject"
            ],
            "Resource": "arn:aws:s3:::<BucketName>/*",
            "Effect": "Allow"
        }
    ]
}
```
**细粒度**

TBD

#### Application Source and Destination
TBD


### Trust relationships

**粗粒度**
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "kinesisanalytics.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
```

**细粒度**
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "kinesisanalytics.amazonaws.com"
            },
            "Action": "sts:AssumeRole",
            "Condition": {
                "StringEquals": {
                    "aws:SourceAccount": "<AccountID>"
                },
                "ArnEquals": {
                    "aws:SourceArn": "arn:aws:kinesisanalytics:<AccountID>:application/<ApplicationName>"
                }
            }
        }
    ]
}
```


## Application action 

```
flinkApplication --action create
https://docs.aws.amazon.com/managed-flink/latest/apiv2/API_CreateApplication.html

flinkApplication --action start (TBD start from snapshot)
https://docs.aws.amazon.com/managed-flink/latest/apiv2/API_StartApplication.html

flinkApplication --action status
https://docs.aws.amazon.com/managed-flink/latest/apiv2/API_DescribeApplication.html

flinkApplication --action stop
https://docs.aws.amazon.com/managed-flink/latest/apiv2/API_StopApplication.html

flinkApplication --action update (TBD)
https://docs.aws.amazon.com/managed-flink/latest/apiv2/API_UpdateApplication.html

flinkApplication --action snapshot
https://docs.aws.amazon.com/managed-flink/latest/apiv2/API_CreateApplicationSnapshot.html

```


# Appendix
---

## Resources
Role: arn:aws:iam::123456789012:role/temp-flink-application
s3 bucket: arn:aws:s3:::mylabstack-flinklabapp

## CMD

### Create application
```
./flinkApplication --action create \
--applicationName sdkapp1 \
--codeS3Bucket arn:aws:s3:::mylabstack-flinklabapp \
--codeS3FileKey flinkapps/sdkapp1/kds-to-s3-datastream-java-1.0.1.jar \
--applicationRole arn:aws:iam::123456789012:role/temp-flink-application  
```

### Start application
```
./flinkApplication --action start \
--applicationName  sdkapp1
```

### Stop application
```
./flinkApplication --action stop \
--applicationName  sdkapp1
```

### Describe application
```
./flinkApplication --action describe \
--applicationName  sdkapp1
```

### Create application snapshot
```
./flinkApplication --action create-snapshot \
--applicationName  sdkapp1 \
--snapshotName snapshot1

```