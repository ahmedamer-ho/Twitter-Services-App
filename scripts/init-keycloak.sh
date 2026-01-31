#!/bin/bash

# Configuration
KEYCLOAK_URL="http://localhost:8085"
ADMIN_USER="admin"
ADMIN_PASS="admin"
REALM="myrealm"
CLIENT_ID="myclient"
CLIENT_SECRET="mysecret"

# 1. Get Admin Token
TOKEN=$(curl -s -d "client_id=admin-cli" -d "username=$ADMIN_USER" -d "password=$ADMIN_PASS" -d "grant_type=password" \
    "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" | jq -r '.access_token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "Failed to get admin token"
    exit 1
fi

# 2. Create Realm
curl -s -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
    -d "{\"realm\": \"$REALM\", \"enabled\": true}" \
    "$KEYCLOAK_URL/admin/realms"

echo "Realm $REALM ensured"

# 3. Create Client (with Service Accounts Enabled)
curl -s -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
    -d "{
        \"clientId\": \"$CLIENT_ID\",
        \"secret\": \"$CLIENT_SECRET\",
        \"enabled\": true,
        \"publicClient\": false,
        \"directAccessGrantsEnabled\": true,
        \"serviceAccountsEnabled\": true,
        \"authorizationServicesEnabled\": true,
        \"redirectUris\": [\"*\"],
        \"webOrigins\": [\"*\"]
    }" \
    "$KEYCLOAK_URL/admin/realms/$REALM/clients"

echo "Client $CLIENT_ID ensured"

# 4. Assign Roles to Service Account
# Get Internal Client ID (UUID)
CLIENT_UUID=$(curl -s -H "Authorization: Bearer $TOKEN" "$KEYCLOAK_URL/admin/realms/$REALM/clients?clientId=$CLIENT_ID" | jq -r '.[0].id')
# Get Service Account User ID
USER_ID=$(curl -s -H "Authorization: Bearer $TOKEN" "$KEYCLOAK_URL/admin/realms/$REALM/clients/$CLIENT_UUID/service-account-user" | jq -r '.id')
# Get realm-management client UUID
RM_CLIENT_UUID=$(curl -s -H "Authorization: Bearer $TOKEN" "$KEYCLOAK_URL/admin/realms/$REALM/clients?clientId=realm-management" | jq -r '.[0].id')
# Get manage-users role
ROLE_JSON=$(curl -s -H "Authorization: Bearer $TOKEN" "$KEYCLOAK_URL/admin/realms/$REALM/clients/$RM_CLIENT_UUID/roles" | jq -c '.[] | select(.name=="manage-users")')

# Assign role
curl -s -X POST -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
    -d "[$ROLE_JSON]" \
    "$KEYCLOAK_URL/admin/realms/$REALM/users/$USER_ID/role-mappings/clients/$RM_CLIENT_UUID"

echo "Permissions assigned to $CLIENT_ID service account"
