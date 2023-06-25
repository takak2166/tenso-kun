#!/bin/bash

cd $HOME/tenso-kun
docker-compose run --rm certbot renew
docker-compose exec -T nginx nginx -s reload
