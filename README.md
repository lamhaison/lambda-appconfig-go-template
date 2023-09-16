# README

# What is appconfig for featureflag?
AWS AppConfig is a service provided by Amazon Web Services (AWS) that helps you manage and deploy application configurations. While it's not specifically designed for feature flags, you can use AWS AppConfig to implement feature flags as part of your application's configuration management.

Here's how AWS AppConfig can be used for feature flags:

1. **Configuration Management**: AWS AppConfig allows you to define configurations for different parts of your application. These configurations can include feature flags or switches that control the behavior of your application.

2. **Environment Management**: You can create different environments within AWS AppConfig to manage configurations for different stages of your application's lifecycle (e.g., development, testing, production).

3. **Configuration Profiles**: AWS AppConfig lets you create configuration profiles within an environment. Each profile can represent a specific set of configurations, including feature flags.

4. **Deployment Strategies**: AWS AppConfig supports deployment strategies that enable you to roll out configuration changes gradually or to specific segments of your user base. This can be useful when toggling feature flags for A/B testing or canary deployments.

5. **Validation and Monitoring**: AWS AppConfig includes validation features, allowing you to enforce constraints on configuration values, including feature flag values. You can also monitor the health of your application's configurations.

6. **Integration**: AWS AppConfig integrates with AWS services like AWS Lambda and Amazon CloudWatch, allowing you to automate deployments and monitor configuration changes.

To implement feature flags using AWS AppConfig:

1. Define feature flags as part of your configurations within AWS AppConfig configuration profiles.

2. Use AWS AppConfig's deployment strategies to control when and how feature flags are rolled out to different environments and user segments.

3. Update your application's code to read and react to feature flag values from the configurations provided by AWS AppConfig.

4. Monitor the impact of feature flag changes using AWS AppConfig's monitoring and validation features.

AWS AppConfig provides a centralized and scalable way to manage and control feature flags, making it easier to experiment with different application behaviors and manage feature releases. It can be particularly useful in scenarios where you want to make runtime configuration changes without modifying your application's code or redeploying it.

# Overview about the solution
APIGW(Not yet in the demo) ---->Lambda ---->AppConfig(To get the setting)
* In this demo, I will not include APIGW however, You can integrate APIGW to Lambda easily.

# Build source code

1. Ensure GNU Make is installed
2. To build
    ```make```
3. Binary is build in the ```bin``` director

# Setting AWS Appconfig
To create an AWS AppConfig application with the name "demo," an environment named "demo," and a default configuration profile.

## For using AWS Console, following these steps:

Creating an AWS AppConfig configuration by using the AWS Management Console involves several steps. Here's a high-level walkthrough of how you can do it:

1. **Log in to the AWS Console**: Log in to your AWS Management Console using your credentials.

2. **Navigate to AWS AppConfig**:
   - Type "AppConfig" into the AWS Console search bar.
   - Click on "AWS AppConfig" under "Services."

3. **Create an Application**:
   - Click on "Create application."
   - Enter the name "demo" for the application.
   - Optionally, you can provide a description.
   - Click "Create application."

4. **Create an Environment**:
   - Inside the "demo" application, click on "Create environment."
   - Enter the name "demo" for the environment.
   - Optionally, provide a description.
   - Click "Create environment."

5. **Create a Configuration Profile**:
   - Inside the "demo" environment, click on "Create configuration profile."
   - Enter the name "default" for the profile.
   - Optionally, provide a description.
   - In the "Configuration source" section, choose "No configuration source."
   - Click "Next."

6. **Create a Feature Flag**:
   - Inside the "default" configuration profile, click "Add feature flag."
   - Enter the name "featureA" for the feature flag.
   - For the data type, choose "BOOLEAN."
   - For "Default value," set it to either "true" or "false" to enable or disable the feature flag by default.
   - Optionally, provide a description.
   - Click "Create feature flag."

7. **Deploy the Configuration**:
   - After creating the feature flag, you can deploy the configuration by selecting the "default" configuration profile and clicking "Start deployment."

8. **Review and Confirm**:
   - Review your configuration details to ensure they are correct.
   - Click "Create" or "Confirm" to create the configuration profile.

Please note that these are high-level steps and may vary slightly depending on the AWS Console's current user interface. Ensure that you have the necessary AWS permissions to create and manage AppConfig resources. Also, the Feature Flag's default value is the initial state of the feature when no explicit configuration has been set. You can later update the configuration to change the feature's state.


# Create lambda function to get config (The lambda can be trigger by APIGW)

## Create role for lambda function with Default Lambda Role and read appconfig permission

To create an AWS Identity and Access Management (IAM) role for a Lambda function with permissions to read AppConfig configurations and assume the default Lambda execution role, you can create a custom IAM role and attach the necessary policies. Here are the steps to create this IAM role:

1. **Create a Custom IAM Role**:

   You can create a custom IAM role using the AWS Management Console or the AWS CLI. Below is an example using the AWS CLI.

   ```bash
   aws iam create-role \
     --role-name LambdaAppConfigRole \
     --assume-role-policy-document '{
       "Version": "2012-10-17",
       "Statement": [
         {
           "Effect": "Allow",
           "Principal": {
             "Service": "lambda.amazonaws.com"
           },
           "Action": "sts:AssumeRole"
         }
       ]
     }'
   ```

   In this command:
   - `LambdaAppConfigRole` is the name of the custom IAM role.
   - `lambda.amazonaws.com` is the service that can assume this role.

2. **Attach AWS Managed Policies**:

   To grant the role the necessary permissions to read AppConfig configurations, you can attach the `AWSAppConfig_ReadOnlyAccess` AWS managed policy. Additionally, you can attach the default AWS managed policy for Lambda, which is `AWSLambda_FullAccess` to provide basic Lambda execution permissions. You can do this using the AWS CLI as follows:

   ```bash

   ACCOUNT_ID=813995029960

   aws iam create-policy \
     --policy-name AWSAppConfig_ReadOnlyAccess \
     --policy-document '{
       "Version": "2012-10-17",
       "Statement": [
         {
           "Effect": "Allow",
           "Action": [
             "appconfig:GetLatestConfiguration",
             "appconfig:StartConfigurationSession"
           ],
           "Resource": "*"
         }
       ]
     }'


   aws iam attach-role-policy \
     --role-name LambdaAppConfigRole \
     --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/AWSAppConfig_ReadOnlyAccess

   aws iam attach-role-policy \
     --role-name LambdaAppConfigRole \
     --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
   ```

   This attaches both policies to the custom role. The `AWSAppConfig_ReadOnlyAccess` policy grants read-only access to AppConfig, and `AWSLambdaBasicExecutionRole` provides basic Lambda permissions.

```bash
aws iam list-attached-role-policies --role-name LambdaAppConfigRole
```

3. **Use the Role in Your Lambda Function**:

   When you create or update your Lambda function, specify the `LambdaAppConfigRole` as the execution role.

By following these steps, you'll have a custom IAM role that allows a Lambda function to assume the role, read AppConfig configurations, and have basic Lambda execution permissions. Make sure to configure your Lambda function to use this role during its creation or update.

To create an AWS Lambda function with the name "lambda-appconfig-go" and upload a Lambda package from the "appconfig_demo.zip" file, you can use the AWS Management Console or the AWS Command Line Interface (CLI). Here's how to do it using the AWS CLI:

**Note:** Before proceeding, ensure that you have the AWS CLI installed and configured with the necessary AWS credentials.

## **Create the Lambda Function**:

   Use the AWS CLI to create the Lambda function. Replace `your-role-arn` with the ARN of the IAM role that your Lambda function should assume.

   ```bash
   aws lambda create-function \
     --function-name lambda-appconfig-go \
     --runtime go1.x \
     --role arn:aws:iam::${ACCOUNT_ID}:role/LambdaAppConfigRole \
     --handler appconfig_demo \
     --zip-file fileb://bin/appconfig_demo.zip


   aws lambda get-function --function-name lambda-appconfig-go
   ```

   - `--function-name` specifies the name of the Lambda function.
   - `--runtime` specifies the runtime environment. For Go, use "go1.x."
   - `--role` specifies the IAM role ARN that the function should assume.
   - `--handler` specifies the entry point for the Lambda function.
   - `--zip-file` specifies the location of the Lambda deployment package.

2. **Configure Lambda Environment Variables**:


   If your Lambda function requires environment variables, you can configure them using the AWS CLI. Replace `key1=value1 key2=value2` with your environment variables. In my case, i am using region ap-southeast-1

   ```bash
   aws lambda update-function-configuration \
     --function-name lambda-appconfig-go \
     --handler appconfig_demo \
     --layers "arn:aws:lambda:ap-southeast-1:421114256042:layer:AWS-AppConfig-Extension:91" \
     --environment Variables="{ENV=demo,project=demo}"

   aws lambda update-function-code \
      --function-name lambda-appconfig-go \
      --zip-file fileb://bin/appconfig_demo.zip
   ```

3. **Invoke the Lambda Function** (Optional):

   You can test your Lambda function using the AWS CLI as follows:

   ```bash
   aws lambda invoke \
     --function-name lambda-appconfig-go \
     --payload '{}' \
     output.txt

   aws logs tail --follow /aws/lambda/lambda-appconfig-go --since 120m
   ```

   This command will invoke the Lambda function with an empty JSON payload and store the result in "output.txt."

That's it! You have created an AWS Lambda function named "lambda-appconfig-go" and uploaded the Lambda package from the "appconfig_demo.zip" file. Remember to replace placeholders with your specific details, such as the IAM role ARN, handler function, and environment variables, as needed.


# Todo later
* Fully setup to integrate with CI, CD to make full workflow
* Using Terraform to build the whole resource.

# Reference 
* [Offical Document from AWS for lambda Integration with AppConfig](https://docs.aws.amazon.com/appconfig/latest/userguide/appconfig-integration-lambda-extensions.html)


