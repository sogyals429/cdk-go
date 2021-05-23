package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/awskms"
	"github.com/aws/aws-cdk-go/awscdk/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/awss3assets"
	"github.com/aws/aws-cdk-go/awscdk/awssns"
	"github.com/aws/aws-cdk-go/awscdk/awssnssubscriptions"
	"github.com/aws/aws-cdk-go/awscdk/awssqs"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type CdkGoStackProps struct {
	awscdk.StackProps
}

type IQueue interface {
}

func NewCdkGoStack(scope constructs.Construct, id string, props *CdkGoStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	common_key := awskms.Key_FromKeyArn(stack, jsii.String("common_key"), jsii.String("arn:aws:kms:ap-southeast-2:186680617253:key/09b80389-053d-4e3c-9202-aff737614463"))

	topic := awssns.NewTopic(stack, jsii.String("MyTopic"), &awssns.TopicProps{
		DisplayName: jsii.String("MyCoolTopic"),
		TopicName:   jsii.String("mytopic"),
	})

	queue := awssqs.NewQueue(stack, jsii.String("NewQueue"), &awssqs.QueueProps{
		QueueName:           jsii.String("cdk-queue"),
		DataKeyReuse:        awscdk.Duration_Seconds(jsii.Number(600)),
		Encryption:          awssqs.QueueEncryption_KMS_MANAGED,
		EncryptionMasterKey: common_key,
	})

	subscription := awssnssubscriptions.NewSqsSubscription(queue, &awssnssubscriptions.SqsSubscriptionProps{})

	topic.AddSubscription(subscription)

	myFunc := awslambda.NewFunction(stack, jsii.String("MyFunction"), &awslambda.FunctionProps{
		FunctionName: jsii.String("mygofunc"),
		Description:  jsii.String("Function made using go CDK"),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(10)),
		Handler:      jsii.String("main"),
		Runtime:      awslambda.Runtime_GO_1_X(),
		Code:         awslambda.Code_FromAsset(jsii.String("./funcs/function.zip"), &awss3assets.AssetOptions{}),
	})

	awslambda.NewEventSourceMapping(stack, jsii.String("new source"), &awslambda.EventSourceMappingProps{
		BatchSize:      jsii.Number(10),
		EventSourceArn: queue.QueueArn(),
		Target:         myFunc,
	})

	sqsPolicy := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"sqs:DeleteMessage",
			"sqs:GetQueueUrl",
			"sqs:ListDeadLetterSourceQueues",
			"sqs:DeleteMessageBatch",
			"sqs:ReceiveMessage",
			"sqs:GetQueueAttributes",
			"sqs:ListQueueTags",
		),
		Resources: jsii.Strings(*queue.QueueArn()),
	})

	myFunc.AddToRolePolicy(sqsPolicy)

	kmsPolicy := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect: awsiam.Effect_ALLOW,
		Actions: jsii.Strings(
			"kms:Decrypt",
			"kms:GenerateDataKey",
		),
		Resources: jsii.Strings(*common_key.KeyArn()),
	})

	myFunc.AddToRolePolicy(kmsPolicy)

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewCdkGoStack(app, "CdkGoStack", &CdkGoStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
