#!/bin/bash
echo INSERT_API_KEY: $INSERT_API_KEY

curl -L -X GET 'http://localhost:4242/api/admin/addons' \
-H 'Accept: application/json' \
-H 'Authorization: bearer $INSERT_API_KEY'

# curl --location --request POST 'http://localhost:4242/api/admin/login' \
# --header 'Authorization: $INSERT_API_KEY' \
# --header 'Content-Type: application/json' \

# curl --location --request POST 'http://localhost:4242/api/admin/segments' \
# --header 'Authorization: $INSERT_API_KEY' \
# --header 'Content-Type: application/json' \
# --data-raw '{
#   "name": "buyers",
#   "description": "buyers group",
#   "project": "default",
#   "constraints": [
#     {
#       "values": [
#         "b1",
#         "b2",
#         "b3",
#         "b4"
#       ],
#       "inverted": false,
#       "operator": "IN",
#       "contextName": "buyerID",
#       "caseInsensitive": false
#     },
#     {
#       "values": [
#         "m1",
#         "m2",
#         "m3",
#         "m4"
#       ],
#       "inverted": false,
#       "operator": "IN",
#       "contextName": "merchantID",
#       "caseInsensitive": false
#     }
#   ]
# }'