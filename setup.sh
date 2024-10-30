#!/bin/bash

INSERT_API_KEY=$ADMIN_API_KEY
echo INSERT_API_KEY: $INSERT_API_KEY

num_users=10
num_buyers=100
num_merchants=100


##################################################################
# set up a bunch of test users
for i in $(seq 1 $num_users); do
    json_payload="{
        \"name\": \"user$i\",
        \"username\": \"user$i\",
        \"password\": \"pass\",
        \"sendEmail\": false,
        \"rootRole\": 1
    }"

    curl --location --request POST 'http://localhost:4242/api/admin/user-admin' \
        --header "Authorization: $INSERT_API_KEY" \
        --header 'Content-Type: application/json' \
        --data-raw "$json_payload"
done



##################################################################
# set up merchant ID context
json_array=()
for i in $(seq 1 $num_merchants); do
  json_array+=("{
    \"value\": \"m$i\",
    \"description\": \"merchant $i desc\"
  }")
done

legal_values="[$(IFS=,; echo "${json_array[*]}")]"

json_payload="{
  \"name\": \"merchantID\",
  \"description\": \"Merchant IDs\",
  \"legalValues\": $legal_values,
  \"stickiness\": false
}"

curl --location --request POST 'http://localhost:4242/api/admin/context' \
    --header "Authorization: $INSERT_API_KEY" \
    --header 'Content-Type: application/json' \
    --data-raw "$json_payload"


##################################################################
# set up buyer ID context
json_array=()
for i in $(seq 1 $num_merchants); do
  json_array+=("{
    \"value\": \"b$i\",
    \"description\": \"buyer $i desc\"
  }")
done

# Join the array elements with commas to form a valid JSON array
legal_values="[$(IFS=,; echo "${json_array[*]}")]"

# Construct the full JSON payload for the curl request
json_payload="{
  \"name\": \"buyerID\",
  \"description\": \"Buyer IDs\",
  \"legalValues\": $legal_values,
  \"stickiness\": false
}"

curl --location --request POST 'http://localhost:4242/api/admin/context' \
    --header "Authorization: $INSERT_API_KEY" \
    --header 'Content-Type: application/json' \
    --data-raw "$json_payload"

##################################################################
# set up wave 1 merchants
curl --location --request POST 'http://localhost:4242/api/admin/segments' \
    --header "Authorization: $INSERT_API_KEY" \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "name": "Wave 1",
    "description": "wave 1 merchants",
    "project": "default",
    "constraints": [
        {
        "contextName": "merchantID",
        "operator": "IN",
        "values": [
            "m1",
            "m2",
            "m3"
        ],
        "caseInsensitive": false,
        "inverted": false
        }
    ]
}'

##################################################################
# set up wave 2 merchants
curl --location --request POST 'http://localhost:4242/api/admin/segments' \
    --header "Authorization: $INSERT_API_KEY" \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "name": "Wave 2",
    "description": "wave 2 merchants",
    "project": "default",
    "constraints": [
        {
        "values": [
            "m4",
            "m5",
            "m6"
        ],
        "inverted": false,
        "operator": "IN",
        "contextName": "merchantID",
        "caseInsensitive": false
        }
    ]
}'

##################################################################
# set up demo feature flag
curl --location --request POST 'http://localhost:4242/api/admin/projects/default/features' \
    --header "Authorization: $INSERT_API_KEY" \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "type": "release",
        "name": "demo-ff",
        "description": "Feature Flag for demo",
        "impressionData": false
    }'


##################################################################
# set up strategy for feature flag
curl --location --request POST 'http://localhost:4242/api/admin/projects/default/features/demo-ff/environments/development/strategies' \
    --header "Authorization: $INSERT_API_KEY" \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "name": "flexibleRollout",
        "constraints": [],
        "parameters": {
            "rollout": "100",
            "stickiness": "default",
            "groupId": "demo-ff"
        },
        "variants": [
            {
            "stickiness": "default",
            "name": "some_json_for_wave_1",
            "weight": 1000,
            "payload": {
                "type": "json",
                "value": "{\n  \"somejson\": \"for wave 1\"\n}"
            },
            "weightType": "variable"
            }
        ],
        "segments": [
            2
        ],
        "disabled": false
        }'