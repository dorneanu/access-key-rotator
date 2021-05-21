import * as cdk from '@aws-cdk/core';
import * as s3 from '@aws-cdk/aws-s3';
import * as lambda from '@aws-cdk/aws-lambda';
import events = require('@aws-cdk/aws-events');
import targets = require('@aws-cdk/aws-events-targets');
import assets = require("@aws-cdk/aws-s3-assets");
import iam = require("@aws-cdk/aws-iam");
import path = require("path");
import { SSL_OP_MICROSOFT_BIG_SSLV3_BUFFER } from 'constants';

export class InfraStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    // Golang binaries must have a place where they are uploaded to s3 as a .zip
    const asset = new assets.Asset(this, 'AccessKeyRotator', {
        path: path.join(__dirname, '../../../build/AccessKeyRotator.zip'),
    });

    // Define ENV variables
    var env = {
        "CLOUD_PROVIDER": this.node.tryGetContext("cloudProvider"),
        "IAM_USER": this.node.tryGetContext("iamUser"),
        "SECRETS_STORE": this.node.tryGetContext("secretsStore"),
        "SECRET_NAME": this.node.tryGetContext("secretName"),
        "REPO_OWNER": this.node.tryGetContext("repoOwner"),
        "REPO_NAME": this.node.tryGetContext("repoName"),
        "TOKEN_CONFIG_STORE_PATH": this.node.tryGetContext("ssmParam"),
        "GITHUB_APP_ID": this.node.tryGetContext("githubAppID"),
        "GITHUB_INST_ID": this.node.tryGetContext("githubInstID"),
    }

    // Create IAM role  to be assumed by the lambda
    const lambdaIAMRole = new iam.Role(this, 'AcccessKeyRotatorIAMRole', {
        assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
        description: 'IAM Role to be assumed by the lambda',
    });
      
        
    // be able to get SSM parameters 
    var ssmRessource = ['arn', 'aws', 'ssm', this.region, this.account, 'parameter/'+env.TOKEN_CONFIG_STORE_PATH].join(':'); 
    lambdaIAMRole.addToPolicy(new iam.PolicyStatement({
        actions: ['ssm:GetParameter'],
        resources: [ssmRessource],
    }));

    // be able to list, create, delete IAM access keys
    var iamRessource = ['arn', 'aws', 'iam', '', this.account, 'user/'+env.IAM_USER].join(':'); 
    lambdaIAMRole.addToPolicy(new iam.PolicyStatement({
        actions: ['iam:ListAccessKeys', 'iam:CreateAccessKey', 'iam:DeleteAccessKey'],
        resources: [iamRessource],
    }));

    const handler = new lambda.Function(this, "AccessKeyRotatorLambda", {
        runtime: lambda.Runtime.GO_1_X,
        handler: "build/access-key-rotator.lambda",
        code: lambda.Code.fromBucket(
            asset.bucket,
            asset.s3ObjectKey
        ),
        environment: env,
        role: lambdaIAMRole 
    })

    // This is used for debugging 
    const debugHandler = new lambda.Function(this, "DebugAccessKeyRotatorLambda", {
        runtime: lambda.Runtime.GO_1_X,
        handler: "access-key-rotator.lambda",
        code: lambda.Code.fromAsset('../../build/AccessKeyRotator.zip'),
        environment: env,
        role: lambdaIAMRole 
    });

    // Create s3 bucket
    new s3.Bucket(this, 'AccessKeyRotatorBucket', {
        versioned: true,
        removalPolicy: cdk.RemovalPolicy.DESTROY,
    });


    // Define cron job to run our Lambda
    // For expressions: https://docs.aws.amazon.com/lambda/latest/dg/services-cloudwatchevents-expressions.html
    // Run every day at 10:30
    const rule = new events.Rule(this, 'Rule', {
        schedule: events.Schedule.expression('cron(30 10 * * ? *)')
    });
    rule.addTarget(new targets.LambdaFunction(handler));
  }
}
