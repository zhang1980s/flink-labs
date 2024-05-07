# flink-labs
## flink official
https://github.com/apache/flink-training/tree/master

## AWS blog
https://aws.amazon.com/blogs/big-data/get-started-with-flink-sql-apis-in-amazon-kinesis-data-analytics-studio/ (The overall lab cannot be processed due to errors in the table schema, but the provided SQL code can be used as a reference)
https://aws.amazon.com/blogs/big-data/top-10-flink-sql-queries-to-try-in-amazon-kinesis-data-analytics-studio/   (SQL reference)

https://aws.amazon.com/blogs/big-data/automate-deployment-and-version-updates-for-amazon-kinesis-data-analytics-applications-with-aws-codepipeline/ (CICD codepipeline)

## AWS workshop studio - Starters Guide to Local Development with Apache Flink
https://catalog.us-east-1.prod.workshops.aws/workshops/429cec9e-3222-4943-82f7-1f45c45ed99a/en-US

## Product page of Ali Cloud
https://cn.aliyun.com/product/bigdata/sc?from_alibabacloud=

## My lab

### path
cd mylab/generator/taxi-fare

```
docker build -t taxi-fare-generator:latest .
```

docker (TBD)


### architecture
generator (ec2/fargate) -> kinesis data stream -> amazon managed apache flink 

### SQL in flink notebook
0. create table
```
%flink.ssql
CREATE TABLE `hive`.`flinklab`.`taxifareflink` (
  `id` VARCHAR(2147483647),
  `starttime` TIMESTAMP(3),
  WATERMARK FOR starttime AS starttime - INTERVAL '1' SECOND,
  `endtime` TIMESTAMP(3),
  `startlocation` VARCHAR(2147483647),
  `endlocation` VARCHAR(2147483647),
  `distance` FLOAT,
  `passengercount` INT,
  `fareamount` FLOAT
) WITH (
  'format' = 'json',
  'connector' = 'kinesis',
  'aws.region' = 'ap-southeast-1',
  'stream' = 'FlinkLabTaxiFareStream'
)
```

1. link data in stream
```
%flink.ssql(type=update)
select * from taxifareflink;
```

2. aggrigate data with tumble

```
%flink.ssql(type=update)
SELECT 
   TUMBLE_START(starttime, INTERVAL '1' MINUTES) AS window_start,
   AVG(fareamount) AS avg_fare, 
   COUNT(*) AS num_rides
FROM taxifareflink
GROUP BY TUMBLE(starttime, INTERVAL '1' MINUTES); 
```

### Tips
1. the stream schema should not include uppercase 