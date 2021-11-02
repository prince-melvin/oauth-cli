# oauth-cli
A simple oauth-complaint cli to fetch, introspect access tokens 

# Usage
## user-token
 - oauth.exe usertoken -u <username> -p <password>     [get all the claims]
 - oauth.exe usertoken -u <username> -p <password> -a  [get only access token]

## service-token
 - oauth.exe servicetoken -s <serviceID> -p "<private-key>"      [get all the claims, `enclose private-key in double quotes`]
 - oauth.exe servicetoken -s <serviceID> -p "<private-key>" -a   [get only access token, ``enclose private-key in double quotes`]
 - oauth.exe servicetoken -s <serviceID> --private-key-file <path-to-private-key> -j   [get only JWT Bearer Token, read privte key from file]
