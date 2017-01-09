FROM cassandra:3.0
RUN mkdir /data
COPY skilldirectoryschema.cql /data/

# Require user & pass for accessing Cassandra instance within container
RUN echo "authenticator: PasswordAuthenticator" >> /etc/cassandra/cassandra.yaml
