FROM alpine:3.14

COPY db_script_generator /db_script_generator

ENTRYPOINT ["/db_script_generator"]