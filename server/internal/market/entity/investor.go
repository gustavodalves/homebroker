package entity

type InvestorAssetPosition struct {
	AssetID string
	Shares  int
}

func NewInvestorAssetPosition(assetID string, qntdShares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares:  qntdShares,
	}
}

type Investor struct {
	ID            string
	Name          string
	AssetPosition []*InvestorAssetPosition
}

func NewInvestor(id string) *Investor {
	return &Investor{
		ID:            id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}

func (i *Investor) AddAssetPosition(
	assetPosition *InvestorAssetPosition,
) {
	i.AssetPosition = append(i.AssetPosition, assetPosition)
}

func (i *Investor) GetAssetPosition(assetID string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetID == assetID {
			return assetPosition
		}
	}

	return nil
}

func (i *Investor) UpdateAssetPosition(assetID string, qtdShares int) {
	assetPosition := i.GetAssetPosition(assetID)

	if assetPosition == nil {
		i.AddAssetPosition(
			NewInvestorAssetPosition(
				assetID,
				qtdShares,
			),
		)
		return
	}

	assetPosition.Shares += qtdShares
}
