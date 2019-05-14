FROM rust

WORKDIR /usr/local/bin/sonic
RUN apt-get update && apt-get install -y clang
RUN curl -sL https://github.com/valeriansaliou/sonic/archive/v1.2.0.tar.gz | tar xz \
	&& mv ./sonic-1.2.0/* .
RUN cargo build --release

CMD "./target/release/sonic" "-c" "/etc/sonic.cfg"
