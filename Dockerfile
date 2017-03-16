FROM ubuntu:14.04

ADD ./statser /statser

CMD ["/statser"]
