package service

import (
	"context"
	"math"
	"sort"
	"strings"
)

type ModelSquarePricingInterval struct {
	MinTokens       int      `json:"min_tokens"`
	MaxTokens       *int     `json:"max_tokens"`
	TierLabel       string   `json:"tier_label,omitempty"`
	InputPrice      *float64 `json:"input_price"`
	OutputPrice     *float64 `json:"output_price"`
	CacheWritePrice *float64 `json:"cache_write_price"`
	CacheReadPrice  *float64 `json:"cache_read_price"`
	PerRequestPrice *float64 `json:"per_request_price"`
}

type ModelSquarePricing struct {
	BillingMode      string                       `json:"billing_mode"`
	InputPrice       *float64                     `json:"input_price"`
	OutputPrice      *float64                     `json:"output_price"`
	CacheWritePrice  *float64                     `json:"cache_write_price"`
	CacheReadPrice   *float64                     `json:"cache_read_price"`
	ImageOutputPrice *float64                     `json:"image_output_price"`
	PerRequestPrice  *float64                     `json:"per_request_price"`
	Intervals        []ModelSquarePricingInterval `json:"intervals"`
}

type ModelSquareCatalogRow struct {
	ChannelID        int64               `json:"channel_id"`
	ChannelName      string              `json:"channel_name"`
	Name             string              `json:"name"`
	Platform         string              `json:"platform"`
	GroupID          int64               `json:"group_id"`
	GroupName        string              `json:"group_name"`
	RateMultiplier   float64             `json:"rate_multiplier"`
	BasePricing      *ModelSquarePricing `json:"base_pricing"`
	EffectivePricing *ModelSquarePricing `json:"effective_pricing"`
	Pricing          *ModelSquarePricing `json:"pricing"`
	Enabled          bool                `json:"enabled"`
	SortOrder        int                 `json:"sort_order"`
	Configured       bool                `json:"configured"`
}

type modelSquareKey struct {
	channelID int64
	groupID   int64
	platform  string
	modelName string
}

// ListModelSquareCatalog returns model-square candidates. When includeDisabled
// is false and manual entries exist, only enabled configured rows are returned.
// When no manual entries exist, public callers get all derived candidates so
// the existing model square stays populated after the migration.
func (s *ChannelService) ListModelSquareCatalog(ctx context.Context, includeDisabled bool) ([]ModelSquareCatalogRow, error) {
	channels, err := s.ListAvailable(ctx)
	if err != nil {
		return nil, err
	}
	entries, err := s.ListModelSquareEntries(ctx)
	if err != nil {
		return nil, err
	}
	entryByKey := make(map[modelSquareKey]ModelSquareEntry, len(entries))
	for i := range entries {
		entryByKey[modelSquareEntryKey(entries[i])] = entries[i]
	}

	candidates := buildModelSquareCandidates(channels)
	hasManualConfig := len(entries) > 0
	for i := range candidates {
		key := modelSquareRowKey(candidates[i])
		if entry, ok := entryByKey[key]; ok {
			candidates[i].Enabled = entry.Enabled
			candidates[i].SortOrder = entry.SortOrder
			candidates[i].Configured = true
		} else if hasManualConfig {
			candidates[i].Enabled = false
		}
	}

	out := make([]ModelSquareCatalogRow, 0, len(candidates))
	for _, row := range candidates {
		if hasManualConfig && !row.Configured {
			if includeDisabled {
				out = append(out, row)
			}
			continue
		}
		if includeDisabled || row.Enabled {
			out = append(out, row)
		}
	}
	sortModelSquareRows(out)
	return out, nil
}

func buildModelSquareCandidates(channels []AvailableChannel) []ModelSquareCatalogRow {
	byKey := make(map[modelSquareKey]ModelSquareCatalogRow)
	for _, ch := range channels {
		if ch.Status != StatusActive {
			continue
		}
		for _, group := range ch.Groups {
			if !isModelSquarePublicGroup(group) {
				continue
			}
			platformSet := map[string]struct{}{group.Platform: {}}
			for _, model := range filterSupportedModels(ch.ExplicitPricedModels, platformSet) {
				if isModelSquarePerRequestAliasPricing(model.Pricing) {
					continue
				}
				addModelSquareCandidate(byKey, ch, group, model)
			}
			for _, model := range filterSupportedModels(ch.SupportedModels, platformSet) {
				if !isModelSquareGlobalFallbackModel(model) {
					continue
				}
				addModelSquareCandidate(byKey, ch, group, model)
			}
		}
	}
	out := make([]ModelSquareCatalogRow, 0, len(byKey))
	for _, row := range byKey {
		out = append(out, row)
	}
	return out
}

func addModelSquareCandidate(byKey map[modelSquareKey]ModelSquareCatalogRow, ch AvailableChannel, group AvailableGroupRef, model SupportedModel) {
	basePricing := modelSquarePricingFromChannel(model.Pricing)
	if basePricing == nil || !modelSquareHasAnyPricing(basePricing) {
		return
	}
	row := ModelSquareCatalogRow{
		ChannelID:        ch.ID,
		ChannelName:      ch.Name,
		Name:             model.Name,
		Platform:         model.Platform,
		GroupID:          group.ID,
		GroupName:        group.Name,
		RateMultiplier:   group.RateMultiplier,
		BasePricing:      basePricing,
		EffectivePricing: scaleModelSquarePricing(basePricing, group.RateMultiplier),
		Enabled:          true,
	}
	row.Pricing = row.EffectivePricing
	key := modelSquareRowKey(row)
	if current, ok := byKey[key]; !ok || modelSquareBetterPrice(row, current) {
		byKey[key] = row
	}
}

func filterSupportedModels(src []SupportedModel, allowedPlatforms map[string]struct{}) []SupportedModel {
	out := make([]SupportedModel, 0, len(src))
	for i := range src {
		if allowedPlatforms != nil {
			if _, ok := allowedPlatforms[src[i].Platform]; !ok {
				continue
			}
		}
		out = append(out, src[i])
	}
	return out
}

func modelSquarePricingFromChannel(p *ChannelModelPricing) *ModelSquarePricing {
	if p == nil {
		return nil
	}
	intervals := make([]ModelSquarePricingInterval, 0, len(p.Intervals))
	for _, iv := range p.Intervals {
		intervals = append(intervals, ModelSquarePricingInterval{
			MinTokens:       iv.MinTokens,
			MaxTokens:       iv.MaxTokens,
			TierLabel:       iv.TierLabel,
			InputPrice:      iv.InputPrice,
			OutputPrice:     iv.OutputPrice,
			CacheWritePrice: iv.CacheWritePrice,
			CacheReadPrice:  iv.CacheReadPrice,
			PerRequestPrice: iv.PerRequestPrice,
		})
	}
	billingMode := string(p.BillingMode)
	if billingMode == "" {
		billingMode = string(BillingModeToken)
	}
	return &ModelSquarePricing{
		BillingMode:      billingMode,
		InputPrice:       p.InputPrice,
		OutputPrice:      p.OutputPrice,
		CacheWritePrice:  p.CacheWritePrice,
		CacheReadPrice:   p.CacheReadPrice,
		ImageOutputPrice: p.ImageOutputPrice,
		PerRequestPrice:  p.PerRequestPrice,
		Intervals:        intervals,
	}
}

func scaleModelSquarePricing(p *ModelSquarePricing, multiplier float64) *ModelSquarePricing {
	if p == nil {
		return nil
	}
	intervals := make([]ModelSquarePricingInterval, 0, len(p.Intervals))
	for _, iv := range p.Intervals {
		intervals = append(intervals, ModelSquarePricingInterval{
			MinTokens:       iv.MinTokens,
			MaxTokens:       iv.MaxTokens,
			TierLabel:       iv.TierLabel,
			InputPrice:      scaleModelSquareFloat(iv.InputPrice, multiplier),
			OutputPrice:     scaleModelSquareFloat(iv.OutputPrice, multiplier),
			CacheWritePrice: scaleModelSquareFloat(iv.CacheWritePrice, multiplier),
			CacheReadPrice:  scaleModelSquareFloat(iv.CacheReadPrice, multiplier),
			PerRequestPrice: scaleModelSquareFloat(iv.PerRequestPrice, multiplier),
		})
	}
	return &ModelSquarePricing{
		BillingMode:      p.BillingMode,
		InputPrice:       scaleModelSquareFloat(p.InputPrice, multiplier),
		OutputPrice:      scaleModelSquareFloat(p.OutputPrice, multiplier),
		CacheWritePrice:  scaleModelSquareFloat(p.CacheWritePrice, multiplier),
		CacheReadPrice:   scaleModelSquareFloat(p.CacheReadPrice, multiplier),
		ImageOutputPrice: scaleModelSquareFloat(p.ImageOutputPrice, multiplier),
		PerRequestPrice:  scaleModelSquareFloat(p.PerRequestPrice, multiplier),
		Intervals:        intervals,
	}
}

func scaleModelSquareFloat(v *float64, multiplier float64) *float64 {
	if v == nil {
		return nil
	}
	scaled := *v * multiplier
	return &scaled
}

func modelSquareEntryKey(entry ModelSquareEntry) modelSquareKey {
	return modelSquareKey{
		channelID: entry.ChannelID,
		groupID:   entry.GroupID,
		platform:  entry.Platform,
		modelName: strings.ToLower(entry.ModelName),
	}
}

func modelSquareRowKey(row ModelSquareCatalogRow) modelSquareKey {
	return modelSquareKey{
		channelID: row.ChannelID,
		groupID:   row.GroupID,
		platform:  row.Platform,
		modelName: strings.ToLower(row.Name),
	}
}

func isModelSquarePublicGroup(g AvailableGroupRef) bool {
	return !g.IsExclusive &&
		g.Platform != "" &&
		(g.SubscriptionType == "" || g.SubscriptionType == SubscriptionTypeStandard)
}

func isModelSquareGlobalFallbackModel(model SupportedModel) bool {
	if model.Platform != "openai" || model.Pricing == nil {
		return false
	}
	return model.Pricing.BillingMode == BillingModeToken && modelSquareHasAnyPricing(modelSquarePricingFromChannel(model.Pricing))
}

func isModelSquarePerRequestAliasPricing(p *ChannelModelPricing) bool {
	if p == nil || p.BillingMode != BillingModePerRequest {
		return false
	}
	if p.InputPrice == nil || p.OutputPrice == nil {
		return false
	}
	return *p.InputPrice == 0 && *p.OutputPrice == 0
}

func modelSquareBetterPrice(candidate, current ModelSquareCatalogRow) bool {
	candidateScore := modelSquarePricingScore(candidate.EffectivePricing)
	currentScore := modelSquarePricingScore(current.EffectivePricing)
	if candidateScore != currentScore {
		return candidateScore < currentScore
	}
	if candidate.ChannelName != current.ChannelName {
		return candidate.ChannelName < current.ChannelName
	}
	return candidate.GroupName < current.GroupName
}

func modelSquarePricingScore(p *ModelSquarePricing) float64 {
	if p == nil {
		return math.Inf(1)
	}
	var score float64
	var count int
	add := func(v *float64) {
		if v == nil {
			return
		}
		score += *v
		count++
	}
	add(p.InputPrice)
	add(p.OutputPrice)
	add(p.CacheWritePrice)
	add(p.CacheReadPrice)
	add(p.ImageOutputPrice)
	add(p.PerRequestPrice)
	for _, iv := range p.Intervals {
		add(iv.InputPrice)
		add(iv.OutputPrice)
		add(iv.CacheWritePrice)
		add(iv.CacheReadPrice)
		add(iv.PerRequestPrice)
	}
	if count == 0 {
		return math.Inf(1)
	}
	return score
}

func modelSquareHasAnyPricing(p *ModelSquarePricing) bool {
	return !math.IsInf(modelSquarePricingScore(p), 1)
}

func sortModelSquareRows(rows []ModelSquareCatalogRow) {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].SortOrder != rows[j].SortOrder {
			return rows[i].SortOrder < rows[j].SortOrder
		}
		if rows[i].Platform != rows[j].Platform {
			return rows[i].Platform < rows[j].Platform
		}
		if strings.ToLower(rows[i].Name) != strings.ToLower(rows[j].Name) {
			return strings.ToLower(rows[i].Name) < strings.ToLower(rows[j].Name)
		}
		return rows[i].GroupName < rows[j].GroupName
	})
}
