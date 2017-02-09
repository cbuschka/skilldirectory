FROM node:alpine
RUN apk update
RUN apk add bash curl
RUN npm install -g newman
ADD * /postman/
RUN ls /postman
ENTRYPOINT ["/bin/bash", "-c", "cd /postman && /postman/test-reporter"]
