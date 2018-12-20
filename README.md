# Postgrest-auth

This project is inspired of [postgrest-auth](https://www.npmjs.com/package/postgrest-auth). But it's writting in golang, it's actively maintained, and email are using the [hermes](https://github.com/matcornic/hermes) library to be prettier.

The goal of this project is to provide the whole authentication features for a postgrest-prowered API. It must be deployed alongside your API and share the same jwt secret with your postgrest instance.

## Installation

Using docker:

```bash
docker run -p 3001:3001 \
    -e POSTGREST_AUTH_DB_CONNECTIONSTRING=postgres://user:pass@localhost/db \
    -e POSTGREST_AUTH_EMAIL_AUTH_PASS=pass \
    [...]
    alexandrevilain/postgrest-auth
```

## API

#### Sign in

POST /signin

```bash
curl -X POST http://localhost:3001/signin \
  -H 'Content-Type: application/json' \
  -d '{ "email": "myemail@me.com", "password": "password" }'
```

#### Sign up

POST /signup

```bash
curl -X POST http://localhost:3001/signup \
  -H 'Content-Type: application/json' \
  -d '{ "email": "myemail@me.com", "password": "password" }'
```

#### Confirm email address

GET /confirm/{id}?token={token}

#### Ask for password reset

POST /reset

```bash
curl -X POST http://localhost:3001/reset \
  -H 'Content-Type: application/json' \
  -d '{ "email": "myemail@me.com" }'
```

#### Reset password

POST /reset/:token

```bash
curl -X POST http://localhost:3001/reset/{token} \
  -H 'Content-Type: application/json' \
  -d '{ "password": "mynewpassword" }'
```

## Configuration

Many environment variables are availables to custom your postgrest-auth instance:

| Name                               | Description                                                                                                                                      | Default                              |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------ |
| POSTGREST_AUTH_API_PORT            | The listening port of the service                                                                                                                | 3001                                 |
| POSTGREST_AUTH_API_TOKEN           | The secret used to create the reset password token                                                                                               | supersecret                          |
| POSTGREST_AUTH_LINKS_RESET         | The reset password link sent by email ("%v" will be replaced with the token)                                                                     | http://localhost/reset/%v            |
| POSTGREST_AUTH_LINKS_CONFIRM       | The confirm account link sent by email (The first %v will be replaced by the user's id and the second %v will be replaced by the confirm token ) | http://localhost/confirm/%v?token=%v |
| POSTGREST_AUTH_JWT_EXP             | The token expiration (in hours)                                                                                                                  | X                                    |
| POSTGREST_AUTH_JWT_SECRET          | The shared secret with postgrest                                                                                                                 | X                                    |
| POSTGREST_AUTH_DB_CONNECTIONSTRING | Your dd connection string                                                                                                                        | X                                    |
| POSTGREST_AUTH_DB_ROLES_ANONYMOUS  | The role for anonymous users                                                                                                                     | X                                    |
| POSTGREST_AUTH_DB_ROLES_USER       | The role when users are authenticated                                                                                                            | X                                    |
| POSTGREST_AUTH_APP_NAME            | The application's name where postgrest-auth is installed (your band name)                                                                        | X                                    |
| POSTGREST_AUTH_APP_LINK            | Your appplication's website                                                                                                                      | X                                    |
| POSTGREST_AUTH_APP_LOGO            | Your application's logo                                                                                                                          | X                                    |
| POSTGREST_AUTH_EMAIL_FROM          |                                                                                                                                                  | X                                    |
| POSTGREST_AUTH_EMAIL_HOST          |                                                                                                                                                  | X                                    |
| POSTGREST_AUTH_EMAIL_PORT          |                                                                                                                                                  | X                                    |
| POSTGREST_AUTH_EMAIL_AUTH_USER     |                                                                                                                                                  | X                                    |
| POSTGREST_AUTH_EMAIL_AUTH_PASS     |                                                                                                                                                  | X                                    |
| POSTGREST_AUTH_API_ALLOWEDDOMAINS  | The list of allowed email domains for signup (comma-separated)                                                                                   | X                                    |

## Integration with postgreSQL

This service automatically creates a schema named "auth" and roles defined used environment variables.
It provides you an helper fonction `auth.current_user_id()` that you can for instance use in your POLICES:

```sql
CREATE POLICY questions_update ON questions FOR UPDATE
    USING (user_id = auth.current_user_id())
    WITH CHECK (user_id = auth.current_user_id());
```

## TODO

- Unit tests
- Support oAuth2

## Contributing

Feel free to send PRs!
