-- Seed the deploy-ready public model square defaults.
-- This intentionally creates only platform/group/channel/pricing metadata.
-- Upstream account credentials, account routing ids, runtime caches, and logs are
-- environment-specific and must be configured after deployment.

SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '10min';

WITH upsert_chat_gpt AS (
    INSERT INTO groups (
        name,
        description,
        platform,
        status,
        subscription_type,
        default_validity_days,
        rate_multiplier,
        is_exclusive,
        allow_image_generation,
        image_rate_independent,
        supported_model_scopes,
        allow_messages_dispatch,
        messages_dispatch_model_config,
        model_routing,
        model_routing_enabled,
        updated_at
    )
    VALUES (
        'Chat GPT',
        'OpenAI compatible models for ChatGPT, Codex, image generation, and Claude Messages dispatch.',
        'openai',
        'active',
        'standard',
        0,
        1,
        false,
        true,
        true,
        '["codex-auto-review","gpt-5.2","gpt-5.2-2025-12-11","gpt-5.2-chat-latest","gpt-5.2-pro","gpt-5.2-pro-2025-12-11","gpt-5.3-codex","gpt-5.4","gpt-5.4-2026-03-05","gpt-5.4-mini","gpt-5.5"]'::jsonb,
        true,
        '{"opus_mapped_model":"gpt-5.4","sonnet_mapped_model":"gpt-5.3-codex","haiku_mapped_model":"gpt-5.4-mini"}'::jsonb,
        '{}'::jsonb,
        false,
        NOW()
    )
    ON CONFLICT (name) WHERE deleted_at IS NULL DO UPDATE
    SET description = EXCLUDED.description,
        platform = EXCLUDED.platform,
        status = EXCLUDED.status,
        subscription_type = EXCLUDED.subscription_type,
        default_validity_days = EXCLUDED.default_validity_days,
        rate_multiplier = EXCLUDED.rate_multiplier,
        is_exclusive = EXCLUDED.is_exclusive,
        allow_image_generation = EXCLUDED.allow_image_generation,
        image_rate_independent = EXCLUDED.image_rate_independent,
        supported_model_scopes = EXCLUDED.supported_model_scopes,
        allow_messages_dispatch = EXCLUDED.allow_messages_dispatch,
        messages_dispatch_model_config = EXCLUDED.messages_dispatch_model_config,
        model_routing = EXCLUDED.model_routing,
        model_routing_enabled = EXCLUDED.model_routing_enabled,
        updated_at = NOW()
    RETURNING id
),
upsert_claude AS (
    INSERT INTO groups (
        name,
        description,
        platform,
        status,
        subscription_type,
        default_validity_days,
        rate_multiplier,
        is_exclusive,
        allow_image_generation,
        image_rate_independent,
        supported_model_scopes,
        allow_messages_dispatch,
        messages_dispatch_model_config,
        model_routing,
        model_routing_enabled,
        updated_at
    )
    VALUES (
        'Claude',
        'Anthropic Claude model family.',
        'anthropic',
        'active',
        'standard',
        0,
        1,
        false,
        false,
        false,
        '["claude","gemini_text","gemini_image"]'::jsonb,
        false,
        '{}'::jsonb,
        '{}'::jsonb,
        false,
        NOW()
    )
    ON CONFLICT (name) WHERE deleted_at IS NULL DO UPDATE
    SET description = EXCLUDED.description,
        platform = EXCLUDED.platform,
        status = EXCLUDED.status,
        subscription_type = EXCLUDED.subscription_type,
        default_validity_days = EXCLUDED.default_validity_days,
        rate_multiplier = EXCLUDED.rate_multiplier,
        is_exclusive = EXCLUDED.is_exclusive,
        allow_image_generation = EXCLUDED.allow_image_generation,
        image_rate_independent = EXCLUDED.image_rate_independent,
        supported_model_scopes = EXCLUDED.supported_model_scopes,
        allow_messages_dispatch = EXCLUDED.allow_messages_dispatch,
        messages_dispatch_model_config = EXCLUDED.messages_dispatch_model_config,
        model_routing = EXCLUDED.model_routing,
        model_routing_enabled = EXCLUDED.model_routing_enabled,
        updated_at = NOW()
    RETURNING id
),
upsert_channel AS (
    INSERT INTO channels (
        name,
        description,
        status,
        billing_model_source,
        restrict_models,
        features,
        features_config,
        model_mapping,
        apply_pricing_to_account_stats,
        updated_at
    )
    VALUES (
        U&'\8BCD\5143\6E20\9053',
        U&'\4E0A\6E38\8BCD\5143',
        'active',
        'requested',
        true,
        '',
        '{"web_search_emulation":{"anthropic":false},"codex_image_generation_bridge":{"openai":true}}'::jsonb,
        '{}'::jsonb,
        false,
        NOW()
    )
    ON CONFLICT (name) DO UPDATE
    SET description = EXCLUDED.description,
        status = EXCLUDED.status,
        billing_model_source = EXCLUDED.billing_model_source,
        restrict_models = EXCLUDED.restrict_models,
        features = EXCLUDED.features,
        features_config = EXCLUDED.features_config,
        model_mapping = EXCLUDED.model_mapping,
        apply_pricing_to_account_stats = EXCLUDED.apply_pricing_to_account_stats,
        updated_at = NOW()
    RETURNING id
),
target_groups AS (
    SELECT id FROM upsert_chat_gpt
    UNION ALL
    SELECT id FROM upsert_claude
)
INSERT INTO channel_groups (channel_id, group_id)
SELECT (SELECT id FROM upsert_channel), id
FROM target_groups
ON CONFLICT (group_id) DO UPDATE SET channel_id = EXCLUDED.channel_id;

WITH target_channel AS (
    SELECT id FROM channels WHERE name = U&'\8BCD\5143\6E20\9053'
),
deleted_intervals AS (
    DELETE FROM channel_pricing_intervals
    WHERE pricing_id IN (
        SELECT cmp.id
        FROM channel_model_pricing cmp
        JOIN target_channel tc ON tc.id = cmp.channel_id
    )
    RETURNING pricing_id
),
deleted_pricing AS (
    DELETE FROM channel_model_pricing
    WHERE channel_id = (SELECT id FROM target_channel)
    RETURNING id
)
INSERT INTO channel_model_pricing (
    channel_id,
    platform,
    models,
    billing_mode,
    input_price,
    output_price,
    cache_write_price,
    cache_read_price,
    image_output_price,
    per_request_price
)
VALUES
    (
        (SELECT id FROM target_channel),
        'anthropic',
        '["claude-opus-4-6","claude-opus-4-7"]'::jsonb,
        'per_request',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        0.013
    ),
    (
        (SELECT id FROM target_channel),
        'anthropic',
        '["claude-sonnet-4"]'::jsonb,
        'per_request',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        0.010
    ),
    (
        (SELECT id FROM target_channel),
        'anthropic',
        '["claude-3-opus-20240229","claude-opus-4-5","claude-opus-4-5-20251101","claude-opus-4-5-thinking","claude-opus-4-6-thinking"]'::jsonb,
        'per_request',
        0,
        0,
        0,
        0,
        0,
        0.013
    ),
    (
        (SELECT id FROM target_channel),
        'anthropic',
        '["claude-3-5-sonnet-20241022","claude-3-5-sonnet-latest","claude-3-7-sonnet-20250219","claude-3-7-sonnet-latest","claude-3-haiku-20240307","claude-3-sonnet-20240229","claude-haiku-4-5","claude-haiku-4-5-20251001","claude-sonnet-4-20250514","claude-sonnet-4-5","claude-sonnet-4-5-20250929","claude-sonnet-4-5-thinking","claude-sonnet-4-6","claude-sonnet-4-6-thinking"]'::jsonb,
        'per_request',
        0,
        0,
        0,
        0,
        0,
        0.010
    ),
    (
        (SELECT id FROM target_channel),
        'anthropic',
        '["claude-opus-*"]'::jsonb,
        'per_request',
        0,
        0,
        0,
        0,
        0,
        0.013
    ),
    (
        (SELECT id FROM target_channel),
        'anthropic',
        '["claude-sonnet-*","claude-haiku-*","claude-3-*"]'::jsonb,
        'per_request',
        0,
        0,
        0,
        0,
        0,
        0.010
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["codex-auto-review","gpt-5.2","gpt-5.2-2025-12-11","gpt-5.2-chat-latest","gpt-5.2-pro","gpt-5.2-pro-2025-12-11","gpt-5.3-codex","gpt-5.4","gpt-5.4-2026-03-05","gpt-5.4-mini","gpt-5.5"]'::jsonb,
        'token',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["grok-4.20-0309-reasoning-super","ZhipuAI/GLM-5","moonshotai/Kimi-K2.5","grok-4.20-0309-non-reasoning","deepseek-ai/DeepSeek-V4-Flash","inclusionAI/Ling-2.6-flash","MiniMax/MiniMax-M2.5","gemini-3.1-flash-lite-preview"]'::jsonb,
        'per_request',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        0.007
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["gemini-3.1-pro-preview"]'::jsonb,
        'per_request',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        0.020
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["ZhipuAI/GLM-5.1","qwen3.5-plus"]'::jsonb,
        'per_request',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        0.009
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["stepfun-ai/Step-3.5-Flash"]'::jsonb,
        'per_request',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        0.006
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["inclusionAI/Ling-2.6-1T"]'::jsonb,
        'per_request',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        0.008
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["gemini-3-flash-preview"]'::jsonb,
        'per_request',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        0.010
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["deepseek-ai/DeepSeek-V4-Pro"]'::jsonb,
        'per_request',
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        0.011
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["gpt-5.3-codex-spark"]'::jsonb,
        'token',
        0.000000100000,
        0.000000800000,
        NULL,
        0.000000100000,
        NULL,
        NULL
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["pro/gpt-5.5"]'::jsonb,
        'token',
        0.000017500000,
        0.000105000000,
        NULL,
        0.000001750000,
        NULL,
        NULL
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["pro/gpt-5.2"]'::jsonb,
        'token',
        0.000006125000,
        0.000049000000,
        NULL,
        0.000000612500,
        NULL,
        NULL
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["pro/gpt-5.4-mini"]'::jsonb,
        'token',
        0.000002625000,
        0.000015750000,
        NULL,
        0.000000262500,
        NULL,
        NULL
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["pro/gpt-5.3-codex-spark"]'::jsonb,
        'token',
        0.000000350000,
        0.000003500000,
        NULL,
        0.000000350000,
        NULL,
        NULL
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["pro/gpt-5.3-codex"]'::jsonb,
        'token',
        0.000006125000,
        0.000049000000,
        NULL,
        0.000000612500,
        NULL,
        NULL
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["pro/gpt-5.4"]'::jsonb,
        'token',
        0.000008750000,
        0.000052500000,
        NULL,
        0.000000875000,
        NULL,
        NULL
    ),
    (
        (SELECT id FROM target_channel),
        'openai',
        '["gpt-image-2"]'::jsonb,
        'image',
        NULL,
        NULL,
        NULL,
        NULL,
        0.05500000,
        0.0550000000
    );
