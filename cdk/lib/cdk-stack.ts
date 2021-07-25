import * as path from "path";
import * as cdk from "@aws-cdk/core";
import * as ec2 from "@aws-cdk/aws-ec2";
import * as cloudfront from "@aws-cdk/aws-cloudfront";
import * as elb from "@aws-cdk/aws-elasticloadbalancingv2";
import * as origin from "@aws-cdk/aws-cloudfront-origins";
import * as ecs from "@aws-cdk/aws-ecs";
import * as ecsp from "@aws-cdk/aws-ecs-patterns";
import * as ecra from "@aws-cdk/aws-ecr-assets";
import * as ecr from "@aws-cdk/aws-ecr";
import * as s3 from "@aws-cdk/aws-s3";
import * as iam from "@aws-cdk/aws-iam";
import * as autoscaling from "@aws-cdk/aws-applicationautoscaling";

export class CdkStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // The code that defines your stack goes here
    const vpcId = this.node.tryGetContext("vpc");
    console.log(vpcId);
    const vpc = ec2.Vpc.fromLookup(this, "vpc", {
      vpcId,
    });

    const storage = new s3.Bucket(this, "bucket", {});
    console.log(storage.bucketArn);

    // HACK: 一瞬作られて終わるので、メンテのときにどうするか決める必要がある
    const assets = new ecra.DockerImageAsset(this, "assets", {
      directory: path.resolve(__dirname, "..", "..", "server"),
    });

    const image = ecs.ContainerImage.fromDockerImageAsset(assets);

    const taskRole = new iam.Role(this, "taskRole", {
      assumedBy: new iam.ServicePrincipal("sns.amazonaws.com"),
    });

    const myBucketPolicy = new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      actions: ["s3:GetObject", "s3:PutObject", "s3:ListBucket"],
      principals: [taskRole],
      resources: [storage.bucketArn, storage.bucketArn + "/*"],
    });
    storage.addToResourcePolicy(myBucketPolicy);

    // const cluster = new ecs.Cluster(this, "cluster", {
    //   vpc,
    //   enableFargateCapacityProviders: true,
    // });

    // const appAutoScaling = new applicationautoscaling.ScalableTarget(
    //   this,
    //   "appAutoScalling",
    //   {
    //     minCapacity: 1,
    //     maxCapacity: 2,
    //     role:
    //   }
    // );

    // const policy = new applicationautoscaling.TargetTrackingScalingPolicy(this, "", {
    // })

    // cluster.addAsgCapacityProvider({});

    const fargate = new ecsp.ApplicationLoadBalancedFargateService(
      this,
      "imageServer",
      {
        vpc,
        taskImageOptions: {
          image,
          containerName: "image-server",
          containerPort: 8080,
          environment: {
            IMAGE_SERVER_BUCKET: storage.bucketName,
            // TODO: ...
          },
          taskRole,
        },
        desiredCount: 0,
        // cluster,
      }
    );

    const cdn = new cloudfront.Distribution(this, "distribution", {
      comment: "image server endpoint",
      defaultBehavior: {
        origin: new origin.LoadBalancerV2Origin(fargate.loadBalancer, {}),
      },
    });
  }
}
