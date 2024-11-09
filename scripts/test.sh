#!/bin/bash

# Replace with your actual server URL
SERVER_URL="http://localhost:8080"

# Function to create a campaign
create_campaign() {
  echo "Creating Campaign..."
  RESPONSE=$(curl -s -X POST "$SERVER_URL/campaigns" \
    -H 'Content-Type: application/json' \
    -d '{
      "name": "First Login Promotion",
      "discount": 30.0,
      "max_users": 2,
      "start_date": "2024-11-01T00:00:00Z",
      "end_date": "2024-12-31T23:59:59Z",
      "description": "30% off for the first 100 users."
    }')
  echo "Response: $RESPONSE"
  CAMPAIGN_ID=$(echo $RESPONSE | jq -r '.id')
  echo "Campaign ID: $CAMPAIGN_ID"
}

# Function to generate vouchers
generate_vouchers() {
  CAMPAIGN_ID=$1
  echo "Generating Vouchers for Campaign ID: $CAMPAIGN_ID..."
  RESPONSE=$(curl -s -X POST "$SERVER_URL/campaigns/$CAMPAIGN_ID/vouchers" \
    -H 'Content-Type: application/json' \
    -d '{
      "count": 2
    }')
  echo "Response: $RESPONSE"
  # Assuming the service returns the list of vouchers
  VOUCHER_CODES=($(echo $RESPONSE | jq -r '.[].code'))
  echo "Voucher Codes: ${VOUCHER_CODES[@]}"
}

# Function to redeem a voucher
redeem_voucher() {
  VOUCHER_CODE=$1
  USER_ID=$2
  echo "Redeeming Voucher Code: $VOUCHER_CODE for User ID: $USER_ID..."
  RESPONSE=$(curl -s -X POST "$SERVER_URL/vouchers/redeem" \
    -H 'Content-Type: application/json' \
    -d "{
      \"code\": \"$VOUCHER_CODE\",
      \"user_id\": \"$USER_ID\"
    }")
  echo "Response: $RESPONSE"
}

# Function to process a purchase
process_purchase() {
  USER_ID=$1
  PLAN=$2
  VOUCHER_CODE=$3
  echo "Processing Purchase for User ID: $USER_ID with Plan: $PLAN and Voucher Code: $VOUCHER_CODE..."
  RESPONSE=$(curl -s -X POST "$SERVER_URL/purchases" \
    -H 'Content-Type: application/json' \
    -d "{
      \"user_id\": \"$USER_ID\",
      \"plan\": \"$PLAN\",
      \"voucher_code\": \"$VOUCHER_CODE\"
    }")
  echo "Response: $RESPONSE"
}

# Function to test exceeding voucher limit
test_exceed_voucher_limit() {
  CAMPAIGN_ID=$1
  echo "Testing Exceeding Voucher Limit..."
  RESPONSE=$(curl -s -X POST "$SERVER_URL/campaigns/$CAMPAIGN_ID/vouchers" \
    -H 'Content-Type: application/json' \
    -d '{
      "count": 1
    }')
  echo "Response: $RESPONSE"
}

# Function to redeem an already used voucher
test_redeem_already_used_voucher() {
  VOUCHER_CODE=$1
  USER_ID=$2
  echo "Redeeming Already Used Voucher Code: $VOUCHER_CODE for User ID: $USER_ID..."
  RESPONSE=$(curl -s -X POST "$SERVER_URL/vouchers/redeem" \
    -H 'Content-Type: application/json' \
    -d "{
      \"code\": \"$VOUCHER_CODE\",
      \"user_id\": \"$USER_ID\"
    }")
  echo "Response: $RESPONSE"
}

# Function to redeem an expired voucher
test_redeem_expired_voucher() {
  VOUCHER_CODE=$1
  USER_ID=$2
  echo "Redeeming Expired Voucher Code: $VOUCHER_CODE for User ID: $USER_ID..."
  RESPONSE=$(curl -s -X POST "$SERVER_URL/vouchers/redeem" \
    -H 'Content-Type: application/json' \
    -d "{
      \"code\": \"$VOUCHER_CODE\",
      \"user_id\": \"$USER_ID\"
    }")
  echo "Response: $RESPONSE"
}

# Function to process purchase with invalid voucher code
test_purchase_invalid_voucher() {
  USER_ID=$1
  PLAN=$2
  VOUCHER_CODE=$3
  echo "Processing Purchase with Invalid Voucher Code: $VOUCHER_CODE..."
  RESPONSE=$(curl -s -X POST "$SERVER_URL/purchases" \
    -H 'Content-Type: application/json' \
    -d "{
      \"user_id\": \"$USER_ID\",
      \"plan\": \"$PLAN\",
      \"voucher_code\": \"$VOUCHER_CODE\"
    }")
  echo "Response: $RESPONSE"
}

# Function to process purchase with invalid subscription plan
test_purchase_invalid_plan() {
  USER_ID=$1
  PLAN=$2
  VOUCHER_CODE=$3
  echo "Processing Purchase with Invalid Subscription Plan: $PLAN..."
  RESPONSE=$(curl -s -X POST "$SERVER_URL/purchases" \
    -H 'Content-Type: application/json' \
    -d "{
      \"user_id\": \"$USER_ID\",
      \"plan\": \"$PLAN\",
      \"voucher_code\": \"$VOUCHER_CODE\"
    }")
  echo "Response: $RESPONSE"
}

# Main Execution
echo "Starting Tests..."

# Step 1: Create Campaign
create_campaign_output=$(create_campaign)
CAMPAIGN_ID=$(echo "$create_campaign_output" | grep "Campaign ID" | awk '{print $3}')

# Step 2: Generate Vouchers
generate_vouchers_output=$(generate_vouchers "$CAMPAIGN_ID")
VOUCHER1=$(echo "$generate_vouchers_output" | grep -oP 'VOUCHER\d+')
VOUCHER2=$(echo "$generate_vouchers_output" | grep -oP 'VOUCHER\d+')

# Step 3: Redeem Vouchers
redeem_voucher "$VOUCHER1" "user1"
redeem_voucher "$VOUCHER2" "user2"

# Step 4: Process Purchase with Voucher
process_purchase "user1" "silver" "$VOUCHER1"

# Step 5: Test Exceeding Voucher Limit
test_exceed_voucher_limit "$CAMPAIGN_ID"

# Step 6: Redeem Already Used Voucher
test_redeem_already_used_voucher "$VOUCHER1" "user3"

# Step 7: Redeem Expired Voucher
# To simulate, manually set the voucher's expiry_date to a past date in the database
# For demonstration, assuming "VOUCHER3" is expired
test_redeem_expired_voucher "VOUCHER3" "user4"

# Step 8: Process Purchase with Invalid Voucher Code
test_purchase_invalid_voucher "user5" "silver" "INVALIDCODE"

# Step 9: Process Purchase with Invalid Subscription Plan
test_purchase_invalid_plan "user6" "platinum" "VOUCHER2"

echo "Tests Completed."
