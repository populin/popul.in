FROM php:7.1-fpm

RUN apt-get update && apt-get install -y \
    openssl \
    git \
    unzip \
    libicu-dev \
    libssh2-1-dev

# Set timezone
RUN ln -snf /usr/share/zoneinfo/Europe/Paris /etc/localtime && echo Europe/Paris > /etc/timezone
RUN printf '[PHP]\ndate.timezone = "%s"\n', Europe/Paris > /usr/local/etc/php/conf.d/tzone.ini
RUN "date"

RUN pecl install ssh2-1.0 apcu-5.1.7 && docker-php-ext-enable ssh2
RUN docker-php-ext-configure intl && docker-php-ext-install intl

RUN echo 'alias sf="php bin/console"' >> ~/.bashrc
