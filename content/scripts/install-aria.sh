#!/bin/bash

BLUE='\033[1;36m'
RED='\033[1;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}INSTALL ARIA${NC}"

# deploy proxy-service
echo ""
echo -e "${BLUE}Deploy PROXY${NC}"
scripts/deploy-proxy.sh false

# deploy nginx
echo ""
echo -e "${BLUE}Deploy static web-server${NC}"
scripts/deploy-nginx.sh false

# deploy publisher
echo ""
echo -e "${BLUE}Deploy publisher${NC}"
scripts/deploy-publisher.sh false

# sunc applications
echo ""
echo -e "${BLUE}Sync applications${NC}"
scripts/sync-applications

# deploy redis
kubectl create -f manifests/aria-redis.yaml

echo ""
echo -e "${RED}ARIA INSTALLED${NC}"