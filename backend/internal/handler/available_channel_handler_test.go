package handler

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestBuildPublicModelPricing_IncludesOpenAIFallbackButSkipsClaudeAliases(t *testing.T) {
	zero := 0.0
	openAIInput := 0.1
	openAIOutput := 0.2
	claudeRequest := 0.03

	channels := []service.AvailableChannel{
		{
			Status: service.StatusActive,
			Groups: []service.AvailableGroupRef{
				{
					Name:             "public-openai",
					Platform:         "openai",
					SubscriptionType: service.SubscriptionTypeStandard,
					RateMultiplier:   2,
				},
				{
					Name:             "public-claude",
					Platform:         "anthropic",
					SubscriptionType: service.SubscriptionTypeStandard,
					RateMultiplier:   1,
				},
			},
			SupportedModels: []service.SupportedModel{
				{
					Name:     "gpt-5.4",
					Platform: "openai",
					Pricing: &service.ChannelModelPricing{
						BillingMode: service.BillingModeToken,
						InputPrice:  &openAIInput,
						OutputPrice: &openAIOutput,
					},
				},
				{
					Name:     "claude-3-5-sonnet-latest",
					Platform: "anthropic",
					Pricing: &service.ChannelModelPricing{
						BillingMode: service.BillingModePerRequest,
						InputPrice:  &zero,
						OutputPrice: &zero,
					},
				},
			},
			ExplicitPricedModels: []service.SupportedModel{
				{
					Name:     "claude-sonnet-4",
					Platform: "anthropic",
					Pricing: &service.ChannelModelPricing{
						BillingMode:     service.BillingModePerRequest,
						PerRequestPrice: &claudeRequest,
					},
				},
			},
		},
	}

	got := buildPublicModelPricing(channels)

	byName := make(map[string]publicModelPricing, len(got))
	for _, row := range got {
		byName[row.Name] = row
	}

	gpt, ok := byName["gpt-5.4"]
	require.True(t, ok)
	require.Equal(t, "openai", gpt.Platform)
	require.NotNil(t, gpt.Pricing)
	require.NotNil(t, gpt.Pricing.InputPrice)
	require.Equal(t, 0.2, *gpt.Pricing.InputPrice)

	claude, ok := byName["claude-sonnet-4"]
	require.True(t, ok)
	require.Equal(t, "anthropic", claude.Platform)

	_, aliasPresent := byName["claude-3-5-sonnet-latest"]
	require.False(t, aliasPresent)
}
