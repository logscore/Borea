# Borea

Philosophy:

-   Docs are everything
-   Focus on developers, not marketers
-   Open source to the core forever
-   Give back to the OSS community. We support those projects that we utilize

### How to install and run Borea:

#### Machine hardware requirements

#### Command to run the initialization

-   curl/webget command to run the install bash script
-   takes as user input:
    -   domain
    -   docker image version (default latest)
    -   admin username
    -   admin password
-   pulls the docker image
-   creates 1 db with 3 tables:
    -   admin users
        -   holds admin username and passwords (hashed in bcrypt)
    -   users
        -   holds unique user identifiers
    -   sessions
        -   holds all session info with foreign key linked to a user in user table.
-   generates and writes server key to .env
