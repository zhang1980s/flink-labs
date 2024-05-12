import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as kinesis from 'aws-cdk-lib/aws-kinesis';
import * as s3 from 'aws-cdk-lib/aws-s3';

export class DatastreamApiLabStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const ExampleInputStream = new kinesis.Stream(this, 'ExampleInputStream', {
      streamName: 'ExampleInputStream',
      streamMode: kinesis.StreamMode.ON_DEMAND,   
    })

    const ExampleOutputStream = new kinesis.Stream(this, 'ExampleOutputStream', {
      streamName: 'ExampleOutputStream',
      streamMode: kinesis.StreamMode.ON_DEMAND,   
    })

    const flinkApplicationBucket = new s3.Bucket(this,'FlinkLabApplicationBucket', {
      versioned: false,
      removalPolicy: cdk.RemovalPolicy.DESTROY
    })

  }
}
