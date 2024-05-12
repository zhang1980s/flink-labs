import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as kinesis from 'aws-cdk-lib/aws-kinesis';
import * as glue from 'aws-cdk-lib/aws-glue';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as kinesisanalyticsv2 from 'aws-cdk-lib/aws-kinesisanalyticsv2';
import * as iam from 'aws-cdk-lib/aws-iam';

// import * as rds from 'aws-cdk-lib/aws-rds'


export class MylabStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    // 1. VPC Creation
    const flinkLabVpc = new ec2.Vpc(this, 'FlinkLabVpc', {
      ipAddresses: ec2.IpAddresses.cidr('10.6.0.0/16'),
      maxAzs: 3, 
      enableDnsHostnames: true,
      enableDnsSupport: true,
      subnetConfiguration: [
        { 
          name: 'private1',
          cidrMask: 24, 
          subnetType: ec2.SubnetType.PRIVATE_ISOLATED
        },
        {
          name: 'public1',
          cidrMask: 28,
          subnetType: ec2.SubnetType.PUBLIC
        }
      ],
    });

    // Create a bastionHost
    // const bastionHost = new ec2.BastionHostLinux(this,'BastionHost', {
    //   vpc: flinkLabVpc,
    //   instanceName: 'flinkLabBastion',
    //   subnetSelection: 
      
    // } )

    const taxiFareStream = new kinesis.Stream(this, 'FlinkLabTaxiFareStream', {
      streamName: 'FlinkLabTaxiFareStream',
      streamMode: kinesis.StreamMode.ON_DEMAND,   
    })

    const flinkLabDatabase = new glue.CfnDatabase(this, 'FlinkLabGlueDatabase', {
      catalogId: this.account,
      databaseInput: {
        name: 'flinklab'
      } 
    })

    const flinkApplicationBucket = new s3.Bucket(this,'FlinkLabApplicationBucket', {
      versioned: false,
      removalPolicy: cdk.RemovalPolicy.DESTROY
    })

    // Flink Applications

    // const flinkApp1Role = new iam.Role(this, 'FlinkApp1Role', {
    //   assumedBy: new iam.ServicePrincipal('kinesisanalytics.amazonaws.com'),
    //   managedPolicies: [
    //     iam.ManagedPolicy.fromAwsManagedPolicyName('AdministratorAccess')
    //   ]
    // })

    const app1 = new kinesisanalyticsv2.CfnApplication(this, 'app1', {
      runtimeEnvironment: 'FLINK-1_18',
      serviceExecutionRole: 'arn:aws:iam::894855526703:role/temp-flink-application',
      applicationName: 'app1',
    })



    // const flinkLabDBPrivateSubnetGroup = new rds.SubnetGroup(this, 'FlinkLabDBPrivateSubnetGroup',{
    //   vpc: flinkLabVpc,
    //   description: 'Subnet group for Flink Lab Aurora PostgreSQL (private)',
    //   subnetGroupName: 'FlinkLabPgsqlPrivateSubnetGroup',
    //   vpcSubnets: {
    //     subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
    //   },
    //   removalPolicy: cdk.RemovalPolicy.DESTROY
    // })

    // const flinkLabDBPublicSubnetGroup = new rds.SubnetGroup(this, 'FlinkLabDBPublicSubnetGroup',{
    //   vpc: flinkLabVpc,
    //   description: 'Subnet group for Flink Lab Aurora PostgreSQL (public)',
    //   subnetGroupName: 'FlinkLabPgsqlPublicSubnetGroup',
    //   vpcSubnets: {
    //     subnetType: ec2.SubnetType.PUBLIC,
    //   },
    //   removalPolicy: cdk.RemovalPolicy.DESTROY
    // })

    // const flinkLabDBSecurityGroup = new ec2.SecurityGroup(this, 'FlinkLabDBSecurityGroup',{
    //   vpc: flinkLabVpc,
    //   description: 'Security group for the Flink Lab Aurora PostgreSQL database',
    //   allowAllOutbound: true,
    // })

    // flinkLabDBSecurityGroup.addIngressRule(
    //   ec2.Peer.ipv4('10.6.0.0/16'),
    //   ec2.Port.POSTGRES,
    //   'Allow PostgreSQL from within the VPC'
    // )
    // const flinkLabManageFlinkStudioSecurityGroup = new ec2.SecurityGroup(this, 'FlinkLabManageFlinkStudioSecurityGroup',{
    //   vpc: flinkLabVpc,
    //   description: 'Security group for the Flink Lab Managed Flink Studio',
    //   allowAllOutbound: true,
    // })

    // const flinkLabBastionSecurityGroup = new ec2.SecurityGroup(this, 'FlinkLabBastionSecurityGroup',{
    //   vpc: flinkLabVpc,
    //   description: 'Security group for the Bastion EC2',
    //   allowAllOutbound: true,
    // })

    // flinkLabBastionSecurityGroup.addIngressRule(
    //   ec2.Peer.ipv4('0.0.0.0/0'),
    //   ec2.Port.SSH,
    //   'Allow SSH from public internet'

    // )
    /////////////////////////////////////////////////////////////
    // Provisioned Cluster
    // const flinkLabDBCluster = new rds.DatabaseCluster(this, 'FlinkLabDBCluster', {
    //   vpc: flinkLabVpc,
    //   engine: rds.DatabaseClusterEngine.auroraPostgres({
    //     version: rds.AuroraPostgresEngineVersion.VER_15_2
    //   }),
      
    //   writer: rds.ClusterInstance.provisioned('writer', {
    //     publiclyAccessible: false,
    //     instanceType: ec2.InstanceType.of(ec2.InstanceClass.T4G, ec2.InstanceSize.MEDIUM)
    //   }),
    //   subnetGroup: flinkLabDBSubnetGroup,
    //   securityGroups: [
    //     flinkLabDBSecurityGroup,
    //   ],
    // })

    // Serverless V1 Cluster
    // const flinkLabDBClusterServerless = new rds.ServerlessCluster(this, 'FlinkLabDBClusterServerless', {
    //   engine: rds.DatabaseClusterEngine.AURORA_POSTGRESQL,
    //   vpc: flinkLabVpc,
    //   subnetGroup: flinkLabDBSubnetGroup,
    //   securityGroups: [
    //     flinkLabDBSecurityGroup,
    //   ],
    //   deletionProtection: false,
    //   removalPolicy: cdk.RemovalPolicy.DESTROY,
    //   scaling: {
    //     minCapacity: rds.AuroraCapacityUnit.ACU_2,
    //     maxCapacity: rds.AuroraCapacityUnit.ACU_4,
    //   },
    //   parameterGroup: rds.ParameterGroup.fromParameterGroupName(this, 'ParameterGroup', 'default.aurora-postgresql15')
    // })

    //////////////////////////////////////////////////////////////////
    // Serverless V2 instance
    // const flinkLabDBCluster = new rds.DatabaseCluster(this, 'FlinkLabDBCluster', {
    //   vpc: flinkLabVpc,
    //   engine: rds.DatabaseClusterEngine.auroraPostgres({
    //     version: rds.AuroraPostgresEngineVersion.VER_15_2
    //   }),
    //   writer: rds.ClusterInstance.serverlessV2('writer'),
    //   subnetGroup: flinkLabDBPublicSubnetGroup,
    //   securityGroups: [
    //     flinkLabDBSecurityGroup,
    //   ],
    //   parameterGroup: rds.ParameterGroup.fromParameterGroupName(this,'ParameterGroup', 'default.aurora-postgresql15')
    // })










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