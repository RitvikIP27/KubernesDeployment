# I am building my CI WORKFLOW into .github.workflows.ci.yml
 ##Structure is as follows
 Name of workflow
            jobs (of workflow)
               build: job 1
                   uses-actions/checout@v4(means actions repo se features use jo code ko cehckout krke run kre)
                   (basically ek ubuntu runner me code clone)
                   uses: actions setupx (for dockerizing)
NOTE: this structure each action trigger is makinga runner to be run for  it like for for docker we used docker/setup buildexaction v4 is ALSO  an action 

step 2  Login docker
uses secrets.docker use
secrets ,socker pass extracted fro mthe github secrets we pushed earlier

Step 3 build the image 
for this use the push action 
 then with context: as ./backend and push:true will build the image
 adn tag: secrets.username:RitvikIP27  for the image to be tagged each time its ibuilt
 next i am going to execute the docker-ompose file so that both fronend backend images are built 

 AND PROCEED TO WRITE THE cd pipeline
 # CI/CD Pipeline

## CI Flow

Push
↓
Checkout
↓
Docker Build
↓
DockerHub Push

## CD Flow

Workflow Run
↓
SSH EC2
↓
Pull Images
↓
Restart Containers

## Secrets

DOCKERHUB_USER
DOCKERHUB_PASS
EC2_HOST
EC2_USER
EC2_SSH_KEY