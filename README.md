# configurationValidator


# Golang utility 
   Verify docker micro-service's status
   
   Verify configuration between docker docker-Compose and nginx 
   
   Easy debugging

Static verification 
  -
   check for docker-dockerCompose-nginx configuration mismatch
   
   Parses services.txt to look for Dockerfile, docker-compose and nginx.conf path
   
   Verifies container details for each entry of service.txt

Runtime verification 
   -
   
   Verifies whether each container is running as per static verification
   
   Detects mismatch between code's port-configuration and docker's port-configuration
