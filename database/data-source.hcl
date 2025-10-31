table "api_keys" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "api_key" {
    null = false
    type = character_varying(255)
  }
  column "hashed_api_key" {
    null = false
    type = character_varying(255)
  }
  column "partner_id" {
    null = false
    type = uuid
  }
  column "is_active" {
    null    = true
    type    = boolean
    default = true
  }
  column "created_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "api_keys_partner_id_fkey" {
    columns     = [column.partner_id]
    ref_columns = [table.partners.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_api_keys_created_at" {
    columns = [column.created_at]
  }
  unique "api_keys_api_key_key" {
    columns = [column.api_key]
  }
}
table "crypto_author_port_nav" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  
  column "datetime" {
    null = false
    type = timestamp
  }
  column "author_username" {
    null = false
    type = character_varying(100)
  }
  column "nav_24" {
    null = true
    type = double_precision
  }
  column "nav_48" {
    null = true
    type = double_precision
  }
  column "nav_72" {
    null = true
    type = double_precision
  }
  column "nav_96" {
    null = true
    type = double_precision
  }
  column "nav_120" {
    null = true
    type = double_precision
  }
  column "nav_144" {
    null = true
    type = double_precision
  }
  column "nav_168" {
    null = true
    type = double_precision
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  unique "uniq_crypto_author_port_nav" {
    columns = [column.author_username, column.datetime]
  }
}
table "crypto_author_nav" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "versioning" {
    null = false
    type = character_varying(10)
  }
  column "datetime" {
    null = false
    type = timestamp
  }
  column "author_username" {
    null = false
    type = character_varying(100)
  }
  column "nav" {
    null = false
    type = double_precision
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  unique "uniq_crypto_author_nav" {
    columns = [column.author_username, column.datetime, column.versioning]
  }
}
table "crypto_bot_wallet_history" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "value" {
    null = false
    type = double_precision
  }
  column "fee" {
    null = false
    type = double_precision
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "bot_name" {
    null = false
    type = character_varying(255)
  }
  column "transfer_type" {
    null = false
    type = character_varying(255)
  }
  primary_key {
    columns = [column.id]
  }
}
table "crypto_copytrade_authors" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "crypto_user_wallet_id" {
    null = false
    type = uuid
  }
  column "author_username" {
    null = false
    type = character_varying(255)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  unique "uniq_crypto_copytrade_authors" {
    columns = [column.crypto_user_wallet_id, column.author_username]
  }
}
table "crypto_copytrade_wallet" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "crypto_user_id" {
    null = false
    type = uuid
  }
  column "wallet_address" {
    null = false
    type = character_varying(255)
  }
  column "private_key" {
    null = false
    type = character_varying(255)
  }
  column "priority" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "execution_fee" {
    null    = false
    type    = boolean
    default = false
  }
  column "hyperliquid_basecode" {
    null    = false
    type    = boolean
    default = false
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "wallet_name" {
    null = true
    type = character_varying(255)
  }
  primary_key {
    columns = [column.id]
  }
  unique "uniq_crypto_copytrade_wallet" {
    columns = [column.wallet_address, column.private_key]
  }
}
table "crypto_copytrade_authors_privy" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "crypto_user_wallet_id_privy" {
    null = false
    type = uuid
  }
  column "author_username" {
    null = false
    type = character_varying(255)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  unique "uniq_crypto_copytrade_authors_privy" {
    columns = [column.crypto_user_wallet_id_privy, column.author_username]
  }
}

table "crypto_copytrade_wallet_privy" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "crypto_user_id" {
    null = false
    type = uuid
  }
  column "wallet_address" {
    null = false
    type = character_varying(255)
  }
  column "wallet_id" {
    null = false
    type = character_varying(255)
  }
  column "priority" {
    null = false
    type = integer
  }
  column "exchange" {
    null = false
    type = character_varying(255)
    default = "hyperliquid"
  }
  
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "execution_fee" {
    type    = numeric(10,6)
    default = 0.1
  }
  column "hyperliquid_basecode" {
    null    = false
    type    = boolean
    default = false
  }
  column "leverage" {
    type    = integer
    default = 1
    null = false
  }
  column "position_size_percentage" {
    type    = double_precision
    default = 0.1
    null = false
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "wallet_name" {
    null = true
    type = character_varying(255)
  }
  column "wallet_type" {
    null = true
    type = character_varying(255)
    default = "copytrade"
  }
  column "tp_percentage" {
    null = true
    type = double_precision
  }
  column "sl_percentage" {
    null = true
    type = double_precision
  }
  column "holding_hour_period" {
    null = true
    type = integer
    default = 48
  }
  column "nft_id" {
    null = true
    type = integer
  }
  primary_key {
    columns = [column.id]
  }
  unique "uniq_crypto_copytrade_wallet_privy" {
    columns = [column.wallet_address, column.wallet_id]
  }
}

table "crypto_copytrade_wallet_cex" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "crypto_user_id" {
    null = false
    type = uuid
  }
  column "wallet_address" {
    null = false
    type = character_varying(255)
  }
  column "api_key" {
    null = true
    type = character_varying(255)
  }
  column "api_secret" {
    null = true
    type = character_varying(255)
  }
  column "priority" {
    null = false
    type = integer
  }
  column "exchange" {
    null = false
    type = character_varying(255)
    default = "hyperliquid"
  }
  
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "execution_fee" {
    type    = numeric(10,6)
    default = 0
  }
  column "hyperliquid_basecode" {
    null    = false
    type    = boolean
    default = false
  }
  column "leverage" {
    type    = integer
    default = 1
    null = false
  }
  column "position_size_percentage" {
    type    = double_precision
    default = 0.1
    null = false
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "wallet_name" {
    null = true
    type = character_varying(255)
  }
  column "tp_percentage" {
    null = true
    type = double_precision
  }
  column "sl_percentage" {
    null = true
    type = double_precision
  }
  column "holding_hour_period" {
    null = true
    type = integer
    default = 48
  }
  
  primary_key {
    columns = [column.id]
  }
  unique "uniq_crypto_copytrade_wallet_cex" {
    columns = [column.wallet_address]
  }
}


table "crypto_copytrade_wallet_dex" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "crypto_user_id" {
    null = false
    type = uuid
  }
  column "wallet_address" {
    null = false
    type = character_varying(255)
  }
  column "api_key" {
    null = true
    type = character_varying(255)
  }
  column "private_key" {
    null = true
    type = character_varying(255)
  }
  column "trading_account" {
    null = true
    type = character_varying(255)
  }
  column "priority" {
    null = false
    type = integer
  }
  column "exchange" {
    null = false
    type = character_varying(255)
    default = "hyperliquid"
  }
  
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "execution_fee" {
    type    = numeric(10,6)
    default = 0
  }
  
  column "leverage" {
    type    = integer
    default = 1
    null = false
  }
  column "position_size_percentage" {
    type    = double_precision
    default = 0.1
    null = false
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "wallet_name" {
    null = true
    type = character_varying(255)
  }
  column "tp_percentage" {
    null = true
    type = double_precision
  }
  column "sl_percentage" {
    null = true
    type = double_precision
  }
  column "holding_hour_period" {
    null = true
    type = integer
    default = 48
  }
  
  primary_key {
    columns = [column.id]
  }
  unique "uniq_crypto_copytrade_wallet_dex" {
    columns = [column.wallet_address]
  }
}
table "crypto_copytrade_fillter" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "fillter_name" {
    null    = false
    type    = character_varying(255)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
}

table "crypto_copytrade_fillter_wallet" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "fillter_id" {
    null    = false
    type    = uuid
  }
  column "wallet_uuid" {
    null    = false
    type    = uuid
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fillter_id_fkey" {
    columns     = [column.fillter_id]
    ref_columns = [table.crypto_copytrade_fillter.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "crypto_crm_user" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "username" {
    null = false
    type = character_varying(255)
  }
  column "password" {
    null = false
    type = character_varying(255)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_crm_user_uuid" {
    columns = [column.id, column.username]
  }
}
table "crypto_features_toggle" {
  schema = schema.public
  column "feature_name" {
    null = false
    type = character_varying(50)
  }
  column "feature_enable" {
    null = false
    type = boolean
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.feature_name]
  }
}
table "crypto_notification_authors_list" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "group_id" {
    null = false
    type = uuid
  }
  column "authors_username" {
    null = false
    type = character_varying(255)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_notification_authors_list_id" {
    columns = [column.group_id]
  }
}
table "crypto_notification_group" {
  schema = schema.public
  column "group_id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "group_name" {
    null = false
    type = character_varying(255)
  }
  column "crypto_user_id" {
    null = false
    type = character_varying(255)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.group_id]
  }
  index "idx_notification_group_uuid" {
    columns = [column.group_id, column.crypto_user_id]
  }
}
table "crypto_price" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "ticker" {
    null = false
    type = character_varying(255)
  }
  column "datetime" {
    null = false
    type = timestamp
  }
  column "open" {
    null = false
    type = double_precision
  }
  column "high" {
    null = false
    type = double_precision
  }
  column "low" {
    null = false
    type = double_precision
  }
  column "close" {
    null = false
    type = double_precision
  }
  column "volume" {
    null = false
    type = double_precision
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  unique "uniq_crypto_price" {
    columns = [column.ticker, column.datetime]
  }
}
table "crypto_refferal_score" {
  schema = schema.public
  column "crypto_user_id" {
    null = false
    type = character_varying(255)
  }
  column "date" {
    null = false
    type = timestamp
  }
  column "direct_points" {
    null = true
    type = integer
  }
  column "indirect_points" {
    null = true
    type = integer
  }
  column "total_points" {
    null = true
    type = integer
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.crypto_user_id, column.date]
  }
}
table "crypto_trading_bot" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "free_margin" {
    null = false
    type = double_precision
  }
  column "used_margin" {
    null = false
    type = double_precision
  }
  column "opening_positions" {
    null = false
    type = integer
  }
  column "closed_positions" {
    null = false
    type = integer
  }
  column "win_positions" {
    null = false
    type = integer
  }
  column "loss_positions" {
    null = false
    type = integer
  }
  column "bot_name" {
    null = false
    type = character_varying(255)
  }
  column "realized_pnl" {
    null = false
    type = double_precision
  }
  column "unrealized_pnl" {
    null = false
    type = double_precision
  }
  column "percentage_return" {
    null = true
    type = double_precision
  }
  column "percentage_maximum_drawdown" {
    null = true
    type = double_precision
  }
  primary_key {
    columns = [column.id]
  }
  unique "uniq_crypto_trading_bot" {
    columns = [column.bot_name, column.created_at]
  }
}
table "crypto_user" {
  schema = schema.public
  column "uuid" {
    null = false
    type = character_varying(255)
  }
  column "email" {
    null = true
    type = character_varying(50)
  }
  column "last_update" {
    null = false
    type = character_varying(50)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "binding_wallet" {
    null = true
    type = character_varying(255)
  }
  column "twitter_uid" {
    null = true
    type = character_varying(255)
  }
  column "twitter_name" {
    null = true
    type = character_varying(255)
  }
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "address" {
    null = true
    type = character_varying(255)
  }
  column "hashed_rt" {
    null = true
    type = character_varying(255)
  }
  column "method_login" {
    null    = false
    type    = character_varying(20)
    default = "email"
  }
  column "stripe_customer_id" {
    null = true
    type = character_varying(64)
  }
  column "telegram_chat_id" {
    null = true
    type = character_varying(20)
  }
  column "telegram_user_id" {
    null = true
    type = character_varying(20)
  }
  column "privy_id" {
    null = true
    type = character_varying(255)
  }
  column "is_copytrade_approved" {
    null = false
    type = boolean
    default = false
  }
  column "waiting_list_timestamp" {
    null = true
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_user_uuid" {
    columns = [column.id]
  }
}
table "crypto_user_kol_code" {
  schema = schema.public
  column "crypto_user_id" {
    null = false
    type = character_varying(255)
  }
  column "display_code" {
    null = false
    type = character_varying(50)
  }
  column "refcode" {
    null = false
    type = character_varying(50)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.crypto_user_id, column.refcode]
  }
  unique "uniq_crypto_user_kol_displaycode" {
    columns = [column.display_code]
  }
  unique "uniq_crypto_user_kol_refcode" {
    columns = [column.refcode]
  }
}
table "crypto_user_refcode" {
  schema = schema.public
  column "crypto_user_id" {
    null = false
    type = character_varying(255)
  }
  column "refcode" {
    null = false
    type = character_varying(50)
  }
  column "crypto_ref_user_id" {
    null = true
    type = character_varying(255)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.crypto_user_id, column.refcode, column.created_at]
  }
  unique "uniq_crypto_user_refcode" {
    columns = [column.crypto_user_id, column.refcode, column.created_at]
  }
}
table "crypto_wallet_user" {
  schema = schema.public
  column "address" {
    null = false
    type = character_varying(255)
  }
  column "binding_email" {
    null = true
    type = character_varying(255)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "hashed_rt" {
    null = true
    type = character_varying(255)
  }
  primary_key {
    columns = [column.address]
  }
  unique "uniq_crypto_wallet_user_address" {
    columns = [column.address]
  }
}
table "logs" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "api_key_id" {
    null = false
    type = uuid
  }
  column "endpoint" {
    null = false
    type = character_varying(255)
  }
  column "method" {
    null = false
    type = character_varying(10)
  }
  column "ip_address" {
    null = false
    type = character_varying(45)
  }
  column "user_agent" {
    null = true
    type = text
  }
  column "request_payload" {
    null = true
    type = jsonb
  }
  column "response_status" {
    null = false
    type = integer
  }
  column "response_time_ms" {
    null = true
    type = integer
  }
  column "account_info" {
    null = true
    type = text
  }
  column "timestamp" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
    comment = "When the API request was made"
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "logs_api_key_id_fkey" {
    columns     = [column.api_key_id]
    ref_columns = [table.api_keys.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_logs_ip_address" {
    columns = [column.ip_address]
  }
  check "logs_method_check" {
    expr = "((method)::text = ANY ((ARRAY['GET'::character varying, 'POST'::character varying, 'PUT'::character varying, 'PATCH'::character varying, 'DELETE'::character varying, 'HEAD'::character varying, 'OPTIONS'::character varying])::text[]))"
  }
}
table "partners" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "company_name" {
    null = false
    type = character_varying(255)
  }
  column "contact_email" {
    null = false
    type = character_varying(255)
  }
  column "purpose" {
    null = true
    type = text
  }
  column "is_active" {
    null    = true
    type    = boolean
    default = true
  }
  column "is_approved" {
    null    = true
    type    = boolean
    default = false
  }
  column "created_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = true
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_partners_updated_at" {
    columns = [column.updated_at]
  }
  check "partners_contact_email_check" {
    expr = "((contact_email)::text ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$'::text)"
  }
  unique "partners_contact_email_key" {
    columns = [column.contact_email]
  }
}
table "trade_logs" {
  schema = schema.public

  column "id" {
    null    = false
    type    = bigserial
  }

  column "source" {
    null = false
    type = text
  }

  column "account_id" {
    null = false
    type = text
  }

  column "wallet_id" {
    null = true
    type = text
    comment = "Optional for Privy"
  }

  column "symbol" {
    null = false
    type = text
  }

  column "side" {
    null = false
    type = text
  }

  column "base_size" {
    null = true
    type = numeric
    comment = "Token size"
  }

  column "usdc_value" {
    null = true
    type = numeric
    comment = "Target value in USDC"
  }

  column "price" {
    null = true
    type = numeric
    comment = "Average or reference price"
  }

  column "leverage" {
    null = true
    type = integer
  }

  column "event" {
    null    = false
    type    = text
    default = "main"
  }

  column "status" {
    null = false
    type = text
  }

  column "executed_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }

  primary_key {
    columns = [column.id]
  }
}
schema "public" {
  comment = "standard public schema"
}
