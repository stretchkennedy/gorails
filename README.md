# gorails

Decrypts Rails cookies and parses user id from devise/warden, then opens a SockJS echo server.

Doesn't have CORS, so you'll need to set up reverse proxy. Here's a sample nginx config:
```nginx
server {
    listen 80 default_server;
    listen [::]:80 default_server ipv6only=on;

    server_name localhost;

    location /ws/ {
        proxy_pass http://127.0.0.1:3001/;
        proxy_http_version 1.1;

        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    location / {
        proxy_pass         http://127.0.0.1:3000/;

        proxy_set_header   Host             $host;
        proxy_set_header   X-Real-IP        $remote_addr;
        proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
    }
}

```
