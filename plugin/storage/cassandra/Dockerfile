FROM cassandra:4.1

RUN apt-get update && apt-get install -y dos2unix 

COPY schema/* /cassandra-schema/
# ENV CQLSH_HOST=cassandra1


RUN dos2unix /cassandra-schema/create.sh /cassandra-schema/docker.sh && apt-get --purge remove -y dos2unix && rm -rf /var/lib/apt/lists/*

RUN groupadd -g 65532 nonroot && \
    useradd -u 65532 -g nonroot nonroot --create-home
# /opt/cassandra/bin/cqlsh  cassandra 9042 -u cassandra -p cassandra
# /opt/cassandra/bin/cqlsh 172.20.20.111 9042 -u cassandra -p cassandra

USER 65532:65532
ENTRYPOINT ["/cassandra-schema/docker.sh"]
