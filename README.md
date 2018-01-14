# Enable basic controls for roku via alexa

# How to run

## install
```
go get github.com/randomtask1155/rokuremote
go get github.com/randomtask1155/alexaroku
```

## set environment variables via bash `export`

|  Variable | Description  |
|---|---|
| ALEXAAPPID  | App ID found in amazon developer dashboard for the given alexa skill  |
| PORT  | Listening port  |
| ROKUIP | IP address of roku device that you want to control  |

```
export ALEXAAPPID="app123456"
export PORT="3030"
export ROKUIP="10.1.1.1"
```

## start

```
nohup alexaroku > /tmp/alexaroku.log &
```

# Recommend NGINX frontend

## create file `/etc/nginx/sites-available/domain.net`

```

upstream backends {
    server 10.1.1.1:3000;
}

server {
    listen              443 ssl; # 'ssl' parameter tells NGINX to decrypt the traffic
    server_name         domain.net;
    ssl_certificate     /etc/nginx/ssl/cert.pem; # The certificate file
    ssl_certificate_key /etc/nginx/ssl/private.pem; # The private key file
    ssl on;


    location / {
        proxy_pass http://backends;
	proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## create soft link

```
ln -s /etc/nginx/sites-available/domain.net /etc/nginx/sites-enabled/domain.net
```

## restart nginx

```
sudo nginx restart
```
