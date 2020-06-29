FROM mongo:3
ADD ./scripts /tmp/scripts

# Modify child mongo to use /data/db-accounts as dbpath (because /data/db wont persist the build because it is already a VOLUME)
RUN mkdir -p /data/db-posts \
    && echo "dbpath = /data/db-posts" > /etc/mongodb.conf \
    && chown -R mongodb:mongodb /data/db-posts

RUN chmod +x /tmp/scripts/mongo_create_insert.sh
RUN su - mongodb && mongod --fork --logpath /var/log/mongodb.log --dbpath /data/db-posts  \
    && /tmp/scripts/mongo_create_insert.sh \
    && mongod --dbpath /data/db-posts --shutdown \
    && chown -R mongodb /data/db-posts

# Make the new dir a VOLUME to persist it
VOLUME /data/db-posts

CMD ["mongod", "--config", "/etc/mongodb.conf", "--smallfiles"]
