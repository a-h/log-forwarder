package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogsdestinations"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	logGroup := awslogs.NewLogGroup(stack, jsii.String("LogStashGroup"), &awslogs.LogGroupProps{
		LogGroupName:  jsii.String("LogStashGroup"),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		Retention:     awslogs.RetentionDays_ONE_DAY,
	})

	f := awslambda.NewDockerImageFunction(stack, jsii.String("LogStash"), &awslambda.DockerImageFunctionProps{
		MemorySize: jsii.Number(1024),
		Timeout:    awscdk.Duration_Minutes(jsii.Number(15.0)),
		Code:       awslambda.DockerImageCode_FromImageAsset(jsii.String("../function/"), &awslambda.AssetImageCodeProps{}),
	})
	awslogs.NewSubscriptionFilter(stack, jsii.String("SubscriptionFilter"), &awslogs.SubscriptionFilterProps{
		Destination:   awslogsdestinations.NewLambdaDestination(f, &awslogsdestinations.LambdaDestinationOptions{}),
		FilterPattern: awslogs.FilterPattern_AllEvents(),
		LogGroup:      logGroup,
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "LogStashPipe", &CdkStackProps{
		awscdk.StackProps{},
	})

	app.Synth(nil)
}
