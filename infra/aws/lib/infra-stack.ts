import * as cdk from '@aws-cdk/core';
import * as s3 from '@aws-cdk/aws-s3';
import * as lambda from '@aws-cdk/aws-lambda';
import events = require('@aws-cdk/aws-events');
import targets = require('@aws-cdk/aws-events-targets');
import assets = require("@aws-cdk/aws-s3-assets");
import iam = require("@aws-cdk/aws-iam");
import path = require("path");

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
        "CLOUD_PROVIDER": "aws",
        "IAM_USER": "GithubIAMUser",
        "SECRETS_STORE": "github",
        "SECRET_NAME": "TESTING",
        "REPO_OWNER": "dorneanu",
        "REPO_NAME": "test",
        "TOKEN_CONFIG_STORE_PATH": "github-token",
        "GITHUB_APP_ID": "114149",
        "GITHUB_INST_ID": "16758104",
    }

    // Create IAM role  to be assumed by the lambda
    const lambdaIAMRole = new iam.Role(this, 'AcccessKeyRotatorIAMRole', {
        assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
        description: 'IAM Role to be assumed by the lambda',
    });
      
        
    // be able to get SSM parameters 
    lambdaIAMRole.addToPolicy(new iam.PolicyStatement({
        actions: ['ssm:GetParameter'],
        resources: ['arn:aws:ssm:eu-central-1:451556475769:parameter/github-token'],
    }));

    // be able to list, create, delete IAM access keys
    lambdaIAMRole.addToPolicy(new iam.PolicyStatement({
        actions: ['IAM:ListAccessKeys', 'IAM:CreateAccessKey', 'IAM:DeleteAccessKey'],
        resources: ['arn:aws:ssm:eu-central-1:451556475769:parameter/github-token'],
    }));

    const handler = new lambda.Function(this, "AccessKeyRotatorLambda", {
        runtime: lambda.Runtime.GO_1_X,
        handler: "access-key-rotator.lambda",
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
        code: lambda.Code.fromAsset('../../build/'),
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
