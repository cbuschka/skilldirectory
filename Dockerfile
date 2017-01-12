FROM ubuntu:14.04
ARG DEBUG_FLAG
ADD skilldirectory /bin/skilldirectory
ENTRYPOINT ["/bin/bash", "-c", "/bin/skilldirectory -debug=true"]
