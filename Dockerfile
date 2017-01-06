FROM ubuntu:14.04
ARG DEBUG_FLAG
ADD skilldirectory /bin/skilldirectory

# ENV/ARG vars do not persist into runtime, so can't expand them in the ENTRYPOINT command.
# So, instead, we'll store its value in a file, and pull it out of the file when we startup.
RUN mkdir /env-var-storage/
RUN echo "${DEBUG_FLAG}" > /env-var-storage/debug.flag
ENTRYPOINT ["/bin/bash", "-c", "/bin/skilldirectory -debug=$(cat /env-var-storage/debug.flag)"]
