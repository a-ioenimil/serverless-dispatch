from diagrams import Cluster, Diagram, Edge
from diagrams.aws.compute import Lambda
from diagrams.aws.database import Dynamodb
from diagrams.aws.engagement import SimpleEmailServiceSes
from diagrams.aws.network import APIGateway
from diagrams.aws.security import Cognito
from diagrams.aws.devtools import Codepipeline
from diagrams.aws.mobile import Amplify
from diagrams.onprem.client import User
from diagrams.onprem.vcs import Github
from diagrams.generic.blank import Blank


graph_attr = {
    "splines": "ortho",
    "nodesep": "0.6",
    "ranksep": "0.75",
    "fontsize": "18",
    "fontname": "Helvetica",
}

node_attr = {
    "fontsize": "12",
    "fontname": "Helvetica",
}

edge_attr = {
    "fontsize": "10",
    "fontname": "Helvetica",
}


def build_diagram() -> None:
    with Diagram(
        "Serverless Task Dispatch - AWS Architecture",
        show=False,
        direction="LR",
        graph_attr=graph_attr,
        node_attr=node_attr,
        edge_attr=edge_attr,
        filename="serverless-task-dispatch-arch",
        outformat="png",
    ):
        with Cluster("End User"):
            end_user = User("End User")
            admin_inbox = User("Admin Inbox\n(Admin: manage/assign/close)")
            member_inbox = User("Member Inbox\n(Member: view/update status)")

        with Cluster("Delivery/Automation"):
            github_actions = Github("GitHub Actions")
            terraform = Codepipeline("Terraform IaC")

        with Cluster("AWS Cloud"):
            with Cluster("Amplify Hosting"):
                amplify = Amplify("Amplify Hosting\nReact/Vite App")

            with Cluster("Identity (Cognito)"):
                cognito = Cognito("Cognito User Pool\nJWT Auth + Groups")
                post_confirm = Lambda("Post Confirmation\nPersist User Profile")

            with Cluster("API Layer (API Gateway)"):
                api_gw = APIGateway("API Gateway\nCognito Authorizer")

            with Cluster("Compute (Lambda)"):
                task_lambdas = Lambda("Task Lambdas\nCreate/List/Update/Close/Assign")
                notifier = Lambda("Async Notifier")

            with Cluster("Data & Events (DynamoDB + Streams)"):
                dynamo = Dynamodb("DynamoDB\nTasks + User Profiles")

            with Cluster("Messaging/Email (SES)"):
                ses = SimpleEmailServiceSes("Amazon SES")

            with Cluster("Legend"):
                legend_sync = Blank("Sync request/response")
                legend_async = Blank("Async event")
                legend_notify = Blank("Notification path")

                legend_sync >> Edge(color="black") >> Blank(" ")
                legend_async >> Edge(color="blue", style="dashed") >> Blank("  ")
                legend_notify >> Edge(color="green") >> Blank("   ")

        end_user >> Edge(label="1") >> amplify
        amplify >> Edge(label="2") >> cognito
        amplify >> Edge(label="3") >> api_gw
        api_gw >> Edge(label="4") >> task_lambdas
        task_lambdas >> Edge(label="5") >> dynamo
        dynamo >> Edge(label="6", color="blue", style="dashed") >> notifier
        notifier >> Edge(label="7", color="green") >> ses
        ses >> Edge(label="7", color="green") >> admin_inbox
        ses >> Edge(label="7", color="green") >> member_inbox
        cognito >> Edge(label="8") >> post_confirm
        post_confirm >> Edge(label="8") >> dynamo

        github_actions >> Edge(label="CI/CD") >> terraform
        terraform >> Edge(label="Provision/Deploy") >> amplify
        terraform >> Edge(label="Provision/Deploy") >> api_gw
        terraform >> Edge(label="Provision/Deploy") >> task_lambdas
        terraform >> Edge(label="Provision/Deploy") >> dynamo
        terraform >> Edge(label="Provision/Deploy") >> cognito
        terraform >> Edge(label="Provision/Deploy") >> ses


if __name__ == "__main__":
    build_diagram()
