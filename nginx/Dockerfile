FROM nginx

# COPY ./nginx-entrypoint.sh /usr/local/bin/nginx-entrypoint.sh

COPY ./default.conf /etc/nginx/conf.d/default.conf
COPY ./nginx.conf /etc/nginx/nginx.conf

# RUN ["chmod", "+x", "/usr/local/bin/nginx-entrypoint.sh"]
# ENTRYPOINT ["/usr/local/bin/nginx-entrypoint.sh"]