FROM ubuntu:latest

ADD ./vinfast-project .

RUN chmod +x vinfast-project

ENTRYPOINT ["/vinfast-project"]