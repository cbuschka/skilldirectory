FROM ubuntu:14.04
RUN apt-get update
RUN apt-get install -y ca-certificates
ARG DEBUG_FLAG
ADD skilldirectory /bin/skilldirectory
ENTRYPOINT ["/bin/bash", "-c", "/bin/skilldirectory -debug=true"]
