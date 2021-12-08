import json

import pulumi
import pulumi_aws as aws
from pulumi_aws import dynamodb, iam, appsync

table = dynamodb.Table(
    "outputs",
    hash_key="id",
    attributes=[dynamodb.TableAttributeArgs(name="id", type="S")],
    read_capacity=1,
    write_capacity=1,
)

role = iam.Role(
    f"{pulumi.get_stack()}-iam-role",
    assume_role_policy=json.dumps(
        {
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Action": "sts:AssumeRole",
                    "Principal": {"Service": "appsync.amazonaws.com"},
                    "Effect": "Allow",
                }
            ],
        }
    ),
)

policy = iam.Policy(
    f"{pulumi.get_stack()}-iam-policy",
    policy=table.arn.apply(
        lambda arn: json.dumps(
            {
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Action": [
                            "dynamodb:PutItem",
                            "dynamodb:GetItem",
                            "dynamodb:Scan",
                        ],
                        "Effect": "Allow",
                        "Resource": [arn],
                    }
                ],
            }
        )
    ),
)

attachment = iam.RolePolicyAttachment(
    f"{pulumi.get_stack()}-iam-policy-attachment", role=role.name, policy_arn=policy.arn
)

schema = """
type Query {
        result(id: ID!): Result
        results: [Result]
    }
    type Result {
        id: ID!
        status: String
    }
    schema {
        query: Query
    }
"""

api = appsync.GraphQLApi("api", authentication_type="API_KEY", schema=schema)
key = appsync.ApiKey(f"{pulumi.get_stack()}-api-key", api_id=api.id)


data_source = appsync.DataSource(
    f"{pulumi.get_stack()}-api-ds",
    name=f"{pulumi.get_stack()}_api_datasource",
    api_id=api.id,
    type="AMAZON_DYNAMODB",
    dynamodb_config=appsync.DataSourceDynamodbConfigArgs(
        table_name=table.name,
    ),
    service_role_arn=role.arn,
)

resolver = appsync.Resolver(
    f"{pulumi.get_stack()}-get-resolver",
    api_id=api.id,
    data_source=data_source.name,
    type="Query",
    field="results",
    request_template="""{
        "version": "2017-02-28",
        "operation": "Scan",
    }
    """,
    response_template="$util.toJson($ctx.result.items)",
)

pulumi.export("api_endpoint", api.uris["GRAPHQL"])
pulumi.export("api_key", key.key)
