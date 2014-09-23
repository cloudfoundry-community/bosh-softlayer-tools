package stemcells

type LightStemcellCmd interface {
	GetLigthStemcellMF() LightStemcellMF

	GetStemcellsPath() string
	Create(vdImageId int) (string, error)
}
