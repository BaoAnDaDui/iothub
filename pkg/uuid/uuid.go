package uuid

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"iothub"
)

var ErrGeneratingID = errors.New("failed to generate uuid")

var _ iothub.IdProvider = (*uuidProvider)(nil)

type uuidProvider struct{}

func New() iothub.IdProvider {
	return &uuidProvider{}
}

func (up *uuidProvider) ID() (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(ErrGeneratingID, err.Error())
	}

	return id.String(), nil
}
