# Server configuration

Aqua reads YAML from the directory passed with `-config-dir` (default `./config`).
No configuration file is committed to this repository.

```bash
cp app.example.yaml app.yaml
cp admin.example.yaml admin.yaml
```

Then replace all placeholders. For local development, generate an RSA key pair for JWT:

```bash
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out jwt_private.pem
openssl rsa -pubout -in jwt_private.pem -out jwt_public.pem
```

Point `jwt.private_key_path` and `jwt.public_key_path` at those files, or mount
them at `/run/secrets/` in production. Never commit the generated keys.
