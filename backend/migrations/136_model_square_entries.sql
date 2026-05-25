-- Manual display control for the public/user model pricing catalog.
-- The table stores selection and ordering only; prices are derived from
-- channel_model_pricing, channel_pricing_intervals, and groups.rate_multiplier.

SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '10min';

CREATE TABLE IF NOT EXISTS model_square_entries (
    id          BIGSERIAL   PRIMARY KEY,
    channel_id  BIGINT      NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    group_id    BIGINT      NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    platform    VARCHAR(50) NOT NULL,
    model_name  TEXT        NOT NULL,
    enabled     BOOLEAN     NOT NULL DEFAULT true,
    sort_order  INT         NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_model_square_entries_unique
    ON model_square_entries (channel_id, group_id, platform, lower(model_name));

CREATE INDEX IF NOT EXISTS idx_model_square_entries_enabled_sort
    ON model_square_entries (enabled, sort_order, id);

COMMENT ON TABLE model_square_entries IS '模型广场展示控制：只保存展示选择和排序，价格实时从渠道定价与分组倍率计算';
COMMENT ON COLUMN model_square_entries.channel_id IS '所属渠道';
COMMENT ON COLUMN model_square_entries.group_id IS '所属分组';
COMMENT ON COLUMN model_square_entries.platform IS '所属平台';
COMMENT ON COLUMN model_square_entries.model_name IS '用户侧展示/请求模型名';
COMMENT ON COLUMN model_square_entries.enabled IS '是否展示在模型广场';
COMMENT ON COLUMN model_square_entries.sort_order IS '模型广场排序，升序';
