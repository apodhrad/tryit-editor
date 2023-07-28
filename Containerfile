#  ____        _ _     _           
# | __ ) _   _(_) | __| | ___ _ __ 
# |  _ \| | | | | |/ _` |/ _ \ '__|
# | |_) | |_| | | | (_| |  __/ |   
# |____/ \__,_|_|_|\__,_|\___|_|   
#                                  
FROM registry.access.redhat.com/ubi9/toolbox:9.2 as builder

RUN dnf install -y go

RUN mkdir /tmp/tryit-editor

COPY ./ /tmp/tryit-editor/

RUN cd /tmp/tryit-editor && make build

#  __  __       _       
# |  \/  | __ _(_)_ __  
# | |\/| |/ _` | | '_ \ 
# | |  | | (_| | | | | |
# |_|  |_|\__,_|_|_| |_|
#
FROM registry.access.redhat.com/ubi9/ubi-minimal:9.2

RUN microdnf install --nodocs -y shadow-utils &&\
    microdnf clean all

COPY --from=builder /tmp/tryit-editor/target/tryit-editor /usr/local/bin/

RUN useradd -M tryit-editor

USER tryit-editor

EXPOSE 8080

ENTRYPOINT ["tryit-editor"]

CMD ["start"]
