FROM cassandra:3.0
RUN mkdir /data
COPY skilldirectoryschema.cql /data/
COPY setup-users.cql /data/

# Require user & pass for accessing Cassandra instance within container
RUN echo "authenticator: PasswordAuthenticator" >> /etc/cassandra/cassandra.yaml

# Create new superuser, and
# Use starwars superuser to "disable" default user w/ ridiculous password
#ENTRYPOINT ["/bin/bash", "-c", "\"cqlsh -u cassandra -p cassandra --execute=\"CREATE USER star WITH PASSWORD 'wars' SUPERUSER;\" && cqlsh -u star -p wars --execute=\"ALTER USER cassandra WITH PASSWORD 'k/xT<SEUAUZ/kV4k>6pA[VPvGk}&X`)+' NOSUPERUSER;\"\""]
