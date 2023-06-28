package assets_db

import (
	assets_dm "assets/internal/core/domain/assets"
	"assets/internal/core/ports"
	"context"
	"errors"
	"fmt"
	"sync"
)

func (cr *CassandraRepo) populate(ctx context.Context, assets ...assets_dm.AssetEntity) (results []assets_dm.AssetEntity, err error) {

	if len(assets) == 0 {
		return results, nil
	}

	var chartIds []string
	var insightIds []string
	var audienceIds []string

	for _, asset := range assets {
		switch asset.Type {
		case assets_dm.TypeChart:
			chartIds = append(chartIds, asset.ContentId)
			break
		case assets_dm.TypeInsight:
			insightIds = append(insightIds, asset.ContentId)
			break
		case assets_dm.TypeAudience:
			audienceIds = append(audienceIds, asset.ContentId)
			break
		}
	}

	var charts []assets_dm.ChartEntity
	var insights []assets_dm.InsightEntity
	var audiences []assets_dm.AudienceEntity

	errChan := make(chan error, 3)
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		charts, _, err = cr.chartsRepo.Select(ctx, ports.SelectChartsRepoParams{Ids: chartIds})
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		insights, _, err = cr.insightsRepo.Select(ctx, ports.SelectInsightsRepoParams{Ids: insightIds})
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		audiences, _, err = cr.audiencesRepo.Select(ctx, ports.SelectAudiencesRepoParams{Ids: audienceIds})
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()

	close(errChan)

	for err = range errChan {
		if err != nil {
			return nil, err
		}
	}

	chartsMap := make(map[string]assets_dm.ChartEntity)
	for _, chart := range charts {
		chartsMap[chart.Id] = chart
	}

	insightsMap := make(map[string]assets_dm.InsightEntity)
	for _, insight := range insights {
		insightsMap[insight.Id] = insight
	}

	audiencesMap := make(map[string]assets_dm.AudienceEntity)
	for _, audience := range audiences {
		audiencesMap[audience.Id] = audience
	}

	for idx, asset := range assets {
		switch asset.Type {
		case assets_dm.TypeChart:
			if val, ok := chartsMap[asset.ContentId]; ok {
				assets[idx].AssetData.Chart = &val
			} else {
				fmt.Println("herexd")
				return nil, errors.New("failed to populate data")
			}
			break
		case assets_dm.TypeInsight:
			if val, ok := insightsMap[asset.ContentId]; ok {
				assets[idx].AssetData.Insight = &val
			} else {
				return nil, errors.New("failed to populate data")
			}
			break
		case assets_dm.TypeAudience:
			if val, ok := audiencesMap[asset.ContentId]; ok {
				assets[idx].AssetData.Audience = &val
			} else {
				return nil, errors.New("failed to populate data")
			}
			break
		}
	}

	return assets, nil
}
