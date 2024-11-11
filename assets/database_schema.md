## Database Schema

### Campaigns

| Field         | Type       | Description                                 |
| ------------- | ---------- | ------------------------------------------- |
| `_id`         | `string`   | Unique identifier for the campaign.         |
| `name`        | `string`   | Name of the campaign.                       |
| `discount`    | `float64`  | Discount percentage offered.                |
| `max_users`   | `int`      | Maximum number of users eligible.           |
| `used_users`  | `int`      | Number of users who have utilized vouchers. |
| `start_date`  | `datetime` | Campaign start date and time.               |
| `end_date`    | `datetime` | Campaign end date and time.                 |
| `description` | `string`   | Description of the campaign.                |

### Vouchers

| Field         | Type       | Description                                  |
| ------------- | ---------- | -------------------------------------------- |
| `_id`         | `string`   | Unique identifier for the voucher.           |
| `code`        | `string`   | Unique voucher code.                         |
| `campaign_id` | `string`   | Reference to the associated campaign.        |
| `user_id`     | `string`   | ID of the user who redeemed the voucher.     |
| `used`        | `bool`     | Indicates whether the voucher has been used. |
| `expiry_date` | `datetime` | Expiry date and time of the voucher.         |

### Purchases

| Field             | Type       | Description                               |
| ----------------- | ---------- | ----------------------------------------- |
| `_id`             | `string`   | Unique identifier for the purchase.       |
| `user_id`         | `string`   | ID of the user making the purchase.       |
| `subscription_id` | `string`   | Reference to the subscription plan.       |
| `amount`          | `float64`  | Original amount before discount.          |
| `discount`        | `float64`  | Discount applied to the purchase.         |
| `total`           | `float64`  | Total amount after applying the discount. |
| `voucher_code`    | `string`   | Voucher code applied (if any).            |
| `purchase_date`   | `datetime` | Date and time of the purchase.            |

### Subscriptions

| Field        | Type       | Description                               |
| ------------ | ---------- | ----------------------------------------- |
| `_id`        | `string`   | Unique identifier for the subscription.   |
| `plan`       | `string`   | Subscription plan name (e.g., silver).    |
| `user_id`    | `string`   | ID of the user who owns the subscription. |
| `start_date` | `datetime` | Subscription start date and time.         |
| `end_date`   | `datetime` | Subscription end date and time.           |
| `status`     | `string`   | Current status of the subscription.       |

### Users

| Field        | Type       | Description                     |
| ------------ | ---------- | ------------------------------- |
| `_id`        | `string`   | Unique identifier for the user. |
| `username`   | `string`   | Username of the user.           |
| `email`      | `string`   | Email address of the user.      |
| `password`   | `string`   | Hashed password of the user.    |
| `created_at` | `datetime` | Account creation date and time. |
