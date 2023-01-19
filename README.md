Developing Mini-s3 server
New developers coding task: developing Mini-s3 server
•	Things you should know by now
•	Your task
•	Commands details
•	The auth.json file
•	Tasks high-level architecture
•	Task breakdown 
o	Test you server
o	The client side
o	Example for the expected flow for ./lc mb command
o	Full flow example
•	Unit tests
•	Deploy your application
•	Live session with explanation
Things you should know by now
1.	MinIo client (mc) basic commands.
2.	What is a bucket in terms of cloud storage.
3.	What is our product.
________________________________________
Your task
Create your own Simple storage service.
You will create a command-line program that will communicate with a server through HTTP requests. Your program will get cli arguments, and those arguments will be one of those 4 commands.
1.	./lc alias set - create alias to the following data: alias name, endpoint, access key, secret key. Store this data in a json file.
2.	./lc mb - create a new bucket (in our case - an empty folder).
3.	./lc cp - copy data from your local storage to the server. The files must be in the projects folder (see examples below).
4.	./lc rm - remove an object from a bucket (objects only - not buckets).
•	You are free to choose the task data model, as long as it satisfies the commands above.
•	You need to handle failure cases (such as - copy object to unknown bucket? deleting unknown object, and so on..)
•	Don’t “overkill” this task. You can assume that the arguments are being sent in the correct order, don’t add features like in the actual MinIO client, only support functionality you are being ask to support.
Download the basic project structure files:
 
________________________________________
Commands details
Bold text denotes text that will be entered as-is, italics denote arguments that will be replaced by a value.
> ./lc alias set aliasName http://localhost:port accessKey secretKey
alias aliasName was created successfully.
•	This command adds a new alias to a user profile.
•	alias name, endpoint, access key, secret key are all strings.
•	You should save a file named “aliasName.json” that holds the data in your client side.
•	Some of the alias data will be stored in a file called auth.json in the server side, in a folder named “buckets”. We will talk about auth.json later in detail.
> ./lc mb aliasName/bucket-name
bucket bucket-name was created successfully.
•	This command will create a new bucket (folder) in our server, inside the “buckets” folder.
•	We can’t have two buckets with the same name.
•	Bucket names are strings that contain only: [a-zA-Z0-9 and “-“] chars.
> ./lc cp [file.extension] aliasName/bucket-name
file [file.extension] copied successfully to aliasName/bucket-name
•	This command will copy data from your local machine (from the projects directory only) to your storage server.
•	We cannot copy data to unknown buckets.
•	For the sake of simplicity, the data will be files only, not folders.
•	Copying a file with the same name to a bucket will override the older file, no need to inform the user for such action.
> ./lc rm aliasName/bucket-name/[file.extension]
file [file.extension] deleted successfully from aliasName/bucket-name
•	This command will remove files from your buckets (files only - we do not remove buckets).
•	We cannot remove non-exist files from buckets.
•	We cannot remove files from non-exist bucket.
________________________________________
The auth.json file
We would like that ONLY the bucket owner (i.e aliasName) would be able to access the bucket and modify it. In order to simulate such authentication before accessing buckets, we will store the authentication data inside the auth.json file.
For example, after typing the following commands:
➜  ~ ./lc alias set testAlias http://localhost:8080 accKey secKey
alias testAlias was created successfully.

➜  ~ ./lc mb testAlias/test-bucket-1
bucket test-bucket-1 was created successfully.

➜  ~ ./lc mb testAlias/test-bucket-2
bucket test-bucket-2 was created successfully.
The auth.json file content should be:
{
  "testAlias": {
    "AccessKey": "accKey",
    "SecretKey": "secKey",
    "Buckets": [
      "test-bucket-1",
      "test-bucket-2"
    ]
  }
}
Thus, when the user testAlias will try to upload a file to the test-bucket-1, we would be able to check if the secret key and access key are correct, and if the user is the owner of test-bucket-1.
Note that the commands rm and cp should not modify the auth.json file.
________________________________________
Tasks high-level architecture
 
________________________________________
Task breakdown
Read the instructions below to complete the task.
You can develop the server/client first, or develop the client and server side in parallel. For your server, use either http/net package or Gin(I used Gin). In Your server, as we already said, you should have a folder named “buckets”, that will store all the users buckets.
Notice: The following json is sent from the client to the server with the below commands.
{
  “aliasName”:userProfileName”,
  “accessKey":"access",
  “secretKey":"shhh",
  “bucketName:my_bucket”
}
Test you server
If you decide to develop the server first, you can use the following curl commands to test it (or create your own):
mb:
➜  ~ curl -X POST http://localhost:8080/mb -H 'Content-Type: application/json' -d ‘{"aliasName":"userProfileName","secretKey":"shhh","accessKey":"access","bucketName":"my-bucket"}'
bucket my-bucket was created successfully.
cp:
➜  ~ curl -X PUT -F upload=@testfile.test -H 'Content-Type: multipart/form-data' -F credentials='{"aliasName":"userProfileName","secretKey":"shhh","accessKey":"access","bucketName":"my-bucket"}'  http://localhost:8080/cp/testfile.test
file testfile.test copied successfully to userProfileName/my-bucket
rm:
➜  ~ curl -X DELETE  -H 'Content-Type: application/json' -d '{"aliasName":"userProfileName","secretKey":"shhh","accessKey":"access","bucketName":"my-bucket"}' http://localhost:8080/rm/testfile.test
file testfile.test deleted successfully from userProfileName/my-bucket
After your server Is capable of handling the above commands, develop the client side.
________________________________________
The client side
You client side should support the commands that were listed in page 1.
Use the files in the basic project structure that was supplied to you and complete the functions implementation.
Tip: before implementing the cp command, visit this helpful link : https://gist.github.com/schollz/f25d77afc9130b72390748bdbce0d9a3 
________________________________________
Example for the expected flow for ./lc mb command
 
________________________________________
Full flow example
First, we will create some dummy files that we will use later:

➜  My-lc ls    
lc      main.go

➜  My-lc touch YOU_ARE_AWSOME.WOW YOU_MADE_IT.WOW

➜  My-lc ls
YOU_ARE_AWSOME.WOW.  YOU_MADE_IT.WOW lc main.go

➜  My-lc ./lc alias sef onboarding http://localhost:8080 acc sec
Unknown command: alias sef

➜  My-lc ./lc alias set onboarding http://localhost:8080 acc sec
alias onboarding was created successfully.

➜  My-lc ls
YOU_ARE_AWSOME.WOW  YOU_MADE_IT.WOW  lc   main.go onboarding.json

➜  My-lc ./lc mb onboarding/the-end
bucket the-end was created successfully.

➜  My-lc ./lc mb onboarding/the-end
Unable to create bucket! please Try another bucket name

➜  My-lc ./lc cp YOU_ARE_AWSOME.WOW onboarding/the-end
file YOU_ARE_AWSOME.WOW copied successfully to onboarding/the-end

➜  My-lc ./lc cp YOU_MADE_IT.WOW onboarding/the-end
file YOU_MADE_IT.WOW copied successfully to onboarding/the-end

➜  My-lc touch file_that_will_be_removed.txt

➜  My-lc ./lc cp file_that_will_be_removed.txt onboarding/the-end
file file_that_will_be_removed.txt copied successfully to onboarding/the-end

➜  My-lc ./lc rm onboarding/the-end/file_that_will_be_removed.txt
file file_that_will_be_removed.txt deleted successfully from onboarding/the-end
At the end of these command sequence, the content of auth.json was:
{
  "onboarding": {
    "AccessKey": "acc",
    "SecretKey": "sec",
    "Buckets": [
      "the-end"
    ]
  }
}
If your server and clients yield these same outputs, and the files and buckets are being created correctly, move on to the next step.
________________________________________
Unit tests
You should write unit-tests for your application (both the cli tool and the server) . Please follow Go’s unit tests conventions. 
________________________________________
Deploy your application
We would like to be able to view the servers state with Graphic user interface. In order to do that, you will need to deploy your server into k8s cluster.
Create the frontend-deployment/service.yaml files. Use the following image alfital2/task-front-deployment and port should be 3000. 
Create the backend-deployment and service.yaml files (for the backend image you should create dockerfile and push the image to docker hub). 
The front end app would send an http get request to:
http://backend:8080/ls (notice the backend service name and port in service.yaml file). 
 
Live session with explanation
Link
Call with Tal Alfi-20221108_161832-Meeting Recording (1).mp4 
