FROM golang
RUN true \
  && wget -qO /tmp/cf.tgz "https://packages.cloudfoundry.org/stable?release=linux64-binary&version=6.37.0&source=docker" \
  && echo "a9f34cb87572eca678b89f7303fe960059347373bb59f82411433a8ef5638445  /tmp/cf.tgz" | shasum -a 256 -c \
  && cd /usr/local/bin \
  && tar -xzf /tmp/cf.tgz cf \
  && rm /tmp/cf.tgz
RUN true \
  && curl -sLo /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 \
  && echo 'c6b3a7d7d3e7b70c6f51b706a3b90bd01833846c54d32ca32f0027f00226ff6d  /usr/local/bin/jq' | sha256sum -c \
  && chmod +x /usr/local/bin/jq
