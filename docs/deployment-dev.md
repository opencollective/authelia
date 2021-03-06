# Deployment for Dev

1. [Deploy With npm](#deploy-with-npm)
2. [Deploy With Docker](#deploy-with-docker)
3. [Deploy nginx](#deploy-nginx)
4. [Discard components](#discard-components)
    1. [Discard SQL Server](#discard-sql-server)
    2. [Discard Redis](#discard-redis)
    3. [Discard LDAP](#discard-ldap)
5. [FAQ](#faq)

**WARNING:** *the instructions given in this documentation are not meant
to be used for production environments since it will make **Authelia**
non resilient to failures.*

**NOTE:** If not done already, we highly recommend you first follow the
[Getting Started] documentation.

In some cases, like protecting personal projects/websites, it can be fine
to use **Authelia** in a non highly-available setup. This reduces the number
of components to only two: a reverse proxy such as nginx or Traefik and
Authelia as a companion of the proxy.

As for a regular deployment in production, you need to install **Authelia**
either by pulling the Docker image or building distributable version.

## Build and deploy the distributable version

    $ authelia-scripts build
    $ PUBLIC_DIR=./dist/public_html ./dist/authelia --config /path/to/your/configuration.yml

## Deploy with Docker

    $ docker pull authelia/authelia
    $ docker run -v /path/to/your/configuration.yml:/etc/authelia/configuration.yml authelia/authelia

## Deploy Nginx

You also need to install nginx and take
[example/compose/nginx/minimal/nginx.conf](./example/compose/nginx/minimal/nginx.conf)
as an example for your configuration.

## Deploy Traefik

TODO

## Discard components

### Discard SQL server

There is an option in the configuration file to avoid using an external
SQL server and use a local sqlite3 database instead. This option will
therefore prevent you from running multiple instances of **Authelia**
in parallel.
Consequently, this option is not meant to be used in production or at
least not one that should scale out.

Here is the configuration you should use:

    storage:
      # The file path of the sqlite3 file where data will be persisted
      local:
        path: /var/lib/authelia/db.sqlite3

### Discard Redis

There is an option in the configuration file to discard Redis and use the
memory of the server to store the KV data. This option will therefore
prevent you from running multiple instances of **Authelia** in parallel and
will make you lose user sessions if the application restarts. This
concretely means that all your users will need to authenticate again after
a restart of Authelia. Hence, this option is not meant to be used in production.

To use memory instead of a Redis backend, just comment out the Redis
connection details in the following block:

    session:
      ...      
      # # The redis connection details
      # redis:
      #  host: redis
      #  port: 6379
      #  password: authelia

### Discard LDAP

**Authelia** can use a file backend in order to store users instead of a
LDAP server or an Active Directory. This mode will therefore prevent you
from running multiple instances of **Authelia** in parallel and is therefore
discouraged for production environments.

To use a file backend instead of a LDAP server, you should first duplicate
the file [users_database.yml](../test/suites/basic/users_database.yml) and
edit it to add the users you want.

The content of this file is as follows:

    users:
      ...
      john:
        password: "$6$rounds=500000$jgiCMRyGXzoqpxS3$w2pJeZnnH8bwW3zzvoMWtTRfQYsHbWbD/hquuQ5vUeIyl9gdwBIt6RWk2S6afBA0DPakbeWgD/4SZPiS0hYtU/"
        email: john.doe@authelia.com
        groups:
          - admins
          - dev

The password is hashed and salted as it is in LDAP servers with salted SHA-512
(more hash algorithms such as Argon2 will be provided in the future).
Here is a one-liner to generate such hashed password:

    $ docker run authelia/authelia authelia hash-password yourpassword
    $6$rounds=50000$BpLnfgDsc2WD8F2q$PumMwig8O0uIe9SgneL8Cm1FvUniOzpqBrH.uQE3aZR4K1dHsQldu5gEjJZsXcO./v3itfz6CXTDTJgeh5e8t.

Copy this newly hashed password into your `users_database.yml` file.

Once the file is created, edit the configuration file with the following
block (as used in [configuration.yml](../test/suites/basic/configuration.yml)):

    authentication_backend:
      file:
        path: /path/to/the/users_database.yml

instead of (used in [config.template.yml](../config.template.yml)):

    authentication_backend:
      ldap:
        url: ldap://openldap
        base_dn: dc=example,dc=com

        additional_users_dn: ou=users
        users_filter: cn={0}

        additional_groups_dn: ou=groups
        groups_filter: (&(member={dn})(objectclass=groupOfNames))

        group_name_attribute: cn

        mail_attribute: mail

        user: cn=admin,dc=example,dc=com
        password: password

## FAQ

### Can you give more details on why this is not suitable for production environments?

This documentation gives instructions that will make **Authelia** non
highly-available and non scalable by preventing you from running multiple
instances of the application.
This means that **Authelia** won't be able to distribute the
load across multiple servers and it will prevent failover in case of a
crash or an hardware issue. Moreover, it will also prevent from reliably
persisting data and consequently fail access to your platform as the devices
registered by your users will be lost.

### Why aren't all those steps automated?

Well, as stated before those instructions are not meant to be applied for
a production environment. That being said, in some cases it is just fine and
writing an Ansible playbook to automate all this process is ok.
We would really be more than happy to review such a PR.
In the meantime, you can check the *Standalone* [suite](./suites.md) to see all this
in a real example.

[Getting Started]: ./getting-started.md
