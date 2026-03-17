FROM dunglas/frankenphp:1.12.1-builder-php8.3-trixie AS builder

# Copy xcaddy in the builder image
COPY --from=caddy:2.10.2-builder /usr/bin/xcaddy /usr/bin/xcaddy

# CGO must be enabled to build FrankenPHP
# MODIFICATION: we added build module caching (--mount=...) here, + GOMODCACHE + GOCACHE declarations
RUN --mount=type=bind,source=./,target=./exciplex \
    --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    CGO_ENABLED=1 \
    XCADDY_SETCAP=1 \
    XCADDY_GO_BUILD_FLAGS="-ldflags='-w -s' -tags=nobadger,nomysql,nopgx" \
    CGO_CFLAGS="-D_GNU_SOURCE $(php-config --includes)" \
    CGO_LDFLAGS="$(php-config --ldflags) $(php-config --libs)" \
    GOMODCACHE=/go/pkg/mod \
    GOCACHE=/root/.cache/go-build \
    xcaddy build v2.10.2 \
        --output /usr/local/bin/frankenphp \
        --with github.com/dunglas/frankenphp=./ \
        --with github.com/dunglas/frankenphp/caddy=./caddy/ \
        --with github.com/dunglas/caddy-cbrotli \
        --with github.com/sandstorm/exciplex=./exciplex

FROM dunglas/frankenphp:1.12.1-php8.3-trixie AS php-base

# reference: https://github.com/mlocati/docker-php-extension-installer
RUN install-php-extensions \
    excimer

# Replace the official binary by the one contained your custom modules
COPY --from=builder /usr/local/bin/frankenphp /usr/local/bin/frankenphp

# run as webserver
ARG USER=www-data

RUN \
	useradd ${USER}; \
	# Give write access to /config/caddy and /data/caddy \
	chown -R ${USER}:${USER} /config/caddy /data/caddy && \
    touch /var/run/caretakerd.key && chown ${USER}:${USER} /var/run/caretakerd.key

ADD root/ /

# HOTFIX for ARM64 Architectures and VIPS (see https://github.com/opencv/opencv/issues/14884#issuecomment-706725583 for details)
# only needed for development on Apple Silicon Macs
RUN echo '. /etc/bash.vips-arm64-hotfix.sh' >>  /etc/bash.bashrc

COPY --from=composer /usr/bin/composer /usr/local/bin/composer

# add colored shell env to distinguish environments properly
ENV SHELL_ENV_DISPLAY=dev
RUN echo '. /etc/bash.colorprompt.sh' >>  /etc/bash.bashrc

# Add Caretaker (Startup manager)
ARG TARGETARCH
RUN curl -SL "https://caretakerd.echocat.org/latest/download/caretakerd-linux-${TARGETARCH}.tar.gz" \
    | tar -xz --exclude caretakerd.html -C /usr/bin

ADD ./test /app/Web

# performance profiling
RUN mkdir -p /app/tracing/_traces/ \
    && chown -R ${USER}:${USER} /app/tracing

# cleanup & chown -> for DEV, the full /app dir is writable
RUN chown -R ${USER} /app /var/www

WORKDIR /app
USER ${USER}
ENTRYPOINT [ "/usr/bin/caretakerd", "run" ]