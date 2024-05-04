import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as kinesis from 'aws-cdk-lib/aws-kinesis';
import * as glue from 'aws-cdk-lib/aws-glue';
import { StreamEventSource } from 'aws-cdk-lib/aws-lambda-event-sources';

export class MylabStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    // 1. VPC Creation
    const vpc = new ec2.Vpc(this, 'FlinkLabVpc', {
      cidr: '10.6.0.0/16',
      maxAzs: 3, 
      enableDnsHostnames: true,
      enableDnsSupport: true,
      subnetConfiguration: [
        { 
          name: 'private-1',
          cidrMask: 24, 
          subnetType: ec2.SubnetType.PRIVATE_ISOLATED
        }
      ],
    });

    const taxiFarestream = new kinesis.Stream(this, 'FlinkLabTaxiFareStream', {
      streamName: 'FlinkLabTaxiFareStream',
      streamMode: kinesis.StreamMode.ON_DEMAND,   
    })

    const flinkLabDatabase = new glue.CfnDatabase(this, 'flinklab', {
      catalogId: this.account,
      databaseInput: {
        name: 'flinklab'
      } 
    })

    // const taxiFareTable = new glue.CfnTable(this, 'taxiFareTable', {
    //   databaseName: 'flinklab', //FIXME
    //   catalogId: this.account,
    //   tableInput: {
    //     name: 'taxifare',
    //     description: 'Taxifare table',
    //     storageDescriptor: {
    //       location: 'FlinkLabTaxiFareStream', //FIXME
    //       inputFormat: 'org.apache.hadoop.mapred.TextInputFormat',
    //       outputFormat: 'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat',
    //       compressed: false,
    //       serdeInfo: {
    //         serializationLibrary: 'org.openx.data.jsonserde.JsonSerDe',
    //       },
    //       columns: [
    //         {
    //           name: 'id',
    //           type: 'string',
    //         },
    //         {
    //           name: 'start_time',
    //           type: 'timestamp',
    //         },
    //         {
    //           name: 'end_time',
    //           type: 'timestamp',
    //         },
    //         {
    //           name: 'start_location',
    //           type: 'string',
    //         },
    //         {
    //           name: 'end_location',
    //           type: 'string',
    //         },
    //         {
    //           name: 'distance',
    //           type: 'float',
    //         },
    //         {
    //           name: 'passenger_count',
    //           type: 'int',
    //         },s
    //         {
    //           name: 'fare_amount',
    //           type: 'float',
    //         }
    //       ],
    //       parameters: {
    //         streamname: 'FlinkLabTaxiFareStream',
    //         typeofdata: 'kinesis'
    //       }
    //     }
    //   }
    // })
  }
}