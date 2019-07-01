FROM debian:stretch-slim

ADD ./statser /statser

CMD ["/statser"]
