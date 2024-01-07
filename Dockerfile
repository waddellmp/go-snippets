FROM mysql:latest
ENV MYSQL_ROOT_PASSWORD=test
ENV MYSQL_DATABASE=snippetbox
VOLUME "/var/lib/mysql"
EXPOSE 3306
