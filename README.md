<u><strong>Go Application Deployment to AWS using Docker, ECR, ECS, Fargate and GitHub Actions</strong></u>

This project demonstrates how to containerize a Go application using Docker and deploy it to AWS with ECR and ECS, while automating the deployment process using GitHub Actions. The goal is to provide a seamless CI/CD pipeline for deploying containerized applications to the cloud

**The Go Application – The Foundation**

First things first, I had to make sure the Go application itself was solid. This involved writing the actual code, of course, handling dependencies, and doing a lot of local testing. I won't go into the specifics of the Go code itself here, as the focus was on the deployment side of things. Let's just say it was a standard Go project, nothing too fancy, but it did what it was supposed to do.

**Docker – Building the Container**

The next big step was creating a Docker image. This was where I really started to get into the details of containerization. I created a Dockerfile, which is like a recipe for building the image. I started with a base image – golang:1.20.5-alpine3.18 for the build stage. Alpine Linux is tiny, which makes the resulting image smaller and faster. Inside the Dockerfile, I set the working directory, copied my Go code, and then ran go mod download to grab all the necessary dependencies.   

Then came the build process: go build -o main .. This compiled my Go application into a binary called main. For security, I also created a non-root user and group within the image using addgroup and adduser. This is a best practice to avoid running the application as root inside the container. Finally, I switched to a smaller image alpine:latest for the runtime stage, copied the binary, and set the correct permissions and the user. I also exposed port 8080, which my application uses.

Building the image with docker build -t golang-app . was pretty straightforward, but I spent some time making sure everything was set up correctly within the Dockerfile. I ran the image locally with docker run to confirm that the application was running inside the container as expected. This was really helpful for catching early issues.

The Dockerfile for the project below...

![dockerfile](https://github.com/user-attachments/assets/135a9384-ba34-4734-bc3c-be2573e47dc9)

Local build of Docker...

![docker build process](https://github.com/user-attachments/assets/279d5d92-a623-4b92-b260-3de8088a51c9)


**GitHub Actions – Automating the Process**

To automate my Go app's deployment, I built a GitHub Actions workflow that triggers on pushes to the main branch. It first checks out the code and securely configures AWS credentials using an IAM role and stored secrets. Then, it logs into ECR, builds the Docker image, tags it with the correct ECR URI, and finally pushes it to my ECR repository, ready for deployment to ECS. This setup automates the entire build-and-push process, streamlining my workflowECR – Storing the Image in the Cloud:

With the Docker image built, I needed a place to store it in the cloud. That's where Amazon ECR comes in. I went into the AWS console and created a private ECR repository named golang-app.

Then, I had to authenticate my local Docker client with ECR using the aws ecr get-login-password command. This was a bit tricky at first, getting the AWS CLI configured correctly, but once I had that sorted, it was smooth sailing. I tagged my local Docker image with the full ECR URI – something like 827950560876.dkr.ecr.us-east-1.amazonaws.com/golang-app:latest. Finally, I pushed the image to ECR with docker push. It felt pretty cool to see my image up there in the cloud!

creating the public repo for the project...
![create github public repo](https://github.com/user-attachments/assets/f2831559-be36-4b1b-918b-a611fb5bbfe0)

Storing my AWS credentials in GitHub Secrets is a best practice for security reasons...

![github action secret setup](https://github.com/user-attachments/assets/7122f931-911d-45f2-83d4-f5e4f47532d1)

Cloning my local git repo...

![cloning my github repo](https://github.com/user-attachments/assets/424ecc02-97f3-486e-a0a5-5e0e9f512bbd)




GitHub Actions workflow file that automates the deployment process


![workflow yml file](https://github.com/user-attachments/assets/eeb8d5b4-650c-4629-a6cc-88446f15d672)

And commiting the repo..

![git commit2](https://github.com/user-attachments/assets/f4a19b1d-8011-44dd-bf9c-a5dfa2a7a01e)


Screeshot of Github workflow action successfully completed...


![github action completed](https://github.com/user-attachments/assets/1384d719-e22a-436c-9cb5-4a142ff03b58)



**ECS – Orchestrating the Deployment**

The final piece of the puzzle was deploying the application to ECS. This was the most complex part, but also the most rewarding. I started by creating an ECS cluster, which could be  Fargate or EC2 servers  that run the containers. I chose Fargate, which is serverless, so I didn't have to manage any EC2 instances myself.   

Then came IAM roles. I created two roles: a task execution role and a task role. The task execution role allows ECS to pull images from ECR and manage logs. The task role gives the application running inside the container the permissions it needs. Next, I created an ECS task definition. This is where you define the container configuration: which image to use (the one I pushed to ECR), port mappings, resource limits, logging configurations and the IAM roles.

Finally, I created an ECS service. This service manages the desired number of running tasks. I configured the service to use the task definition I created, specified the cluster, and set the desired number of tasks to 1 (my project is a small one so 1 is fine). 

Screehshot of the Cluster creation...
![cluster creation](https://github.com/user-attachments/assets/4fd42717-5f65-4285-8536-dbbedaa08db0)

Task Definition for the cluster...
![task definition 2](https://github.com/user-attachments/assets/32b3e792-0e65-4b21-8702-7e7e6433be6b)

Peeping the ongoing process at Cloudformation...
![cloud formation complete](https://github.com/user-attachments/assets/042c5100-65b1-4749-b622-8439cfd0f893)

![configuration ip address](https://github.com/user-attachments/assets/f1511a42-115e-4fc9-9be2-9447703335c7)

![final](https://github.com/user-attachments/assets/b9a69ee5-61e6-4f6f-adf1-5cf5526cf155)



**Challenges and Learnings**

There were a few bumps along the way. Getting the IAM roles and policies correct was probably the trickiest part. I had some issues with permissions at first, but with a bit of debugging and careful checking of the ARNs, I managed to get it working. I also learned a lot about Docker best practices, particularly in multi-stage builds and non-root users.

I'm planning to add a database to the application to make it more useful. I'm thinking of using Amazon Aurora PostgreSQL. This means I'll need to figure out how to connect the app to the database, set up the tables and data structure, and make sure everything is secure. Adding a database will let the app actually save and retrieve information, which will open up a lot more possibilities for what it can do.

**Conclusion**

This project was a great learning experience. I now have a much better understanding of how to containerize applications with Docker and deploy them to AWS using ECR and ECS. The automated CI/CD pipeline with GitHub Actions makes the deployment process much smoother and more efficient. I'm excited to use these skills in future projects.   

