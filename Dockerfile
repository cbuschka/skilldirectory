FROM ubuntu:14.04
ADD skilldirectory /bin/skilldirectory
ENTRYPOINT ["/bin/bash", "-c", "/bin/skilldirectory"]
