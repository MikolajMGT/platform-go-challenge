package assets_itc

import (
	assets_dm "assets/internal/core/domain/assets"
	"assets/internal/core/ports"
	"context"
	"sync"
)

func (i *Interactor) createDependencies(ctx context.Context, params ...ports.InsertAssetItcParams) (charts []assets_dm.ChartEntity, insights []assets_dm.InsightEntity, audiences []assets_dm.AudienceEntity, mapper map[ports.InsertAssetItcParams]string, err error) {
	mapper = make(map[ports.InsertAssetItcParams]string)

	for _, param := range params {
		switch param.Type {
		case assets_dm.TypeChart:
			if param.AssetData.Chart != nil {
				chart := assets_dm.NewChartEntity()
				chart.Chart = *param.AssetData.Chart
				charts = append(charts, chart)
				mapper[param] = chart.Id
			}
			break
		case assets_dm.TypeInsight:
			if param.AssetData.Insight != nil {
				insight := assets_dm.NewInsightEntity()
				insight.Insight = *param.AssetData.Insight
				insights = append(insights, insight)
				mapper[param] = insight.Id
			}
			break
		case assets_dm.TypeAudience:
			if param.AssetData.Audience != nil {
				audience := assets_dm.NewAudienceEntity()
				audience.Audience = *param.AssetData.Audience
				audiences = append(audiences, audience)
				mapper[param] = audience.Id
			}
			break
		}
	}

	errChan := make(chan error, 3)
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		_, err = i.chartsRepo.Insert(ctx, charts...)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		_, err = i.insightsRepo.Insert(ctx, insights...)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		_, err = i.audiencesRepo.Insert(ctx, audiences...)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()

	close(errChan)

	for err = range errChan {
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}

	return charts, insights, audiences, mapper, nil
}

func (i *Interactor) deleteDependencies(ctx context.Context, models ...assets_dm.AssetEntity) (charts []assets_dm.ChartEntity, insights []assets_dm.InsightEntity, audiences []assets_dm.AudienceEntity, err error) {

	for _, model := range models {
		switch model.Type {
		case assets_dm.TypeChart:
			if model.AssetData.Chart != nil {
				charts = append(charts, *model.AssetData.Chart)
			}
			break
		case assets_dm.TypeInsight:
			if model.AssetData.Insight != nil {
				insights = append(insights, *model.AssetData.Insight)
			}
			break
		case assets_dm.TypeAudience:
			if model.AssetData.Audience != nil {
				audiences = append(audiences, *model.AssetData.Audience)
			}
			break
		}
	}

	errChan := make(chan error, 3)
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		_, err = i.chartsRepo.Delete(ctx, charts...)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		_, err = i.insightsRepo.Delete(ctx, insights...)
		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		_, err = i.audiencesRepo.Delete(ctx, audiences...)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()

	close(errChan)

	for err = range errChan {
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return charts, insights, audiences, nil
}
