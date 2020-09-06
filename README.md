# microService configuration validator
Verify docker micro-service's status

Verify configuration between docker docker-Compose and nginx 
   
Easy debugging

# Static verification 
   
	Checks for docker-dockerCompose-nginx configuration mismatch

# Runtime verification    
   
	Detects mismatch between code's port-configuration and docker's port-configuration
   
# Refer services.txt for information

static:
   
	./runParser.sh
runtime:
   
	./runParser.sh 1
