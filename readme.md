## Skinatar

Experimental minecraft avatar service for Simple Discord Link, Simple RPC and other mods.

As of now, this system is in testing, and not quite recommended for public use yet.

You are however, free to self-host this, for your own use.

### What is this?

This is a simple Minecraft avatar api, that allows you to retrieve Minecraft Avatars, Heads or Full Body renders with either a UUID, Username or Texture Hash. This service should support Geyser Bedrock Accounts, as well as mods like Fabric Tailor.

### Caching

The server itself, makes use of Caching to prevent spamming Mojang servers. Additionally, our production instance also uses Cloudflare caching.

Skins and Usernames are cached for 2 hours (both on the server, and via Cloudflare), before they are pulled again.

### Ratelimits

Requests are limited to 5 requests per second. This is just to prevent frequent calls for skins, when trying to bypass the cloudflare cache, which would put additional load on our servers.

---

### How do I host this?

Simple. Make a copy of the `docker-compose.yml` file found in this repo. Then, just run `docker compose up -d`. 

You will need to handle your own URL management for this, either via Cloudflare Tunnels, Apache, NGINX or Nginx Proxy Manager. By default, the server uses port 8080. Caching is also up to you to configure. If you need help with this, we might be able to assist you via our discord server.

In the future, we might allow public usage of our hosted instance, once it's proven battle ready

### Endpoints

The server includes the following endpoints for skins:

- Head -> `SERVERURL/head/{UUID|USERNAME|TEXTUREHASH}` or `SERVERURL/isometric/{UUID|USERNAME|TEXTUREHASH}`
- Full Body -> `SERVERURL/body/{UUID|USERNAME|TEXTUREHASH}`
- Avatar -> `SERVERURL/avarar/{UUID|USERNAME|TEXTUREHASH}`

---

### License

This application and code is licensed under MIT.