FROM golang:1.16-buster
RUN true \
  && wget -qO /tmp/cf.tgz "https://packages.cloudfoundry.org/stable?release=linux64-binary&version=7.2.0&source=docker" \
  && echo "e7e6284ebc3c9002681a924c199ea4a79ff22f8a5820afe24a1bae8142143332  /tmp/cf.tgz" | shasum -a 256 -c \
  && cd /usr/local/bin \
  && tar -xzf /tmp/cf.tgz cf cf7 \
  && rm /tmp/cf.tgz
RUN true \
  && curl -sLo /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 \
  && echo 'af986793a515d500ab2d35f8d2aecd656e764504b789b66d7e1a0b727a124c44  /usr/local/bin/jq' | sha256sum -c \
  && chmod +x /usr/local/bin/jq
