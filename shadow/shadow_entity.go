package shadow

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"gorm.io/datatypes"
	"iothub/pkg/model"
)

type Entity struct {
	ThingId  string         `gorm:"primaryKey;size=64" json:"thingId"`
	Desired  datatypes.JSON `gorm:"desired" json:"desired"`
	Reported datatypes.JSON `gorm:"reported" json:"-"`
	Metadata datatypes.JSON `gorm:"metadata" json:"-"`
	Tags     datatypes.JSON `gorm:"tags" json:"-"`
	Version  int64          `json:"version"`

	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

func (t Entity) TableName() string {
	return "shadow"
}

func toEntity(s Shadow) (Entity, error) {
	r, err := json.Marshal(s.State.Reported)
	if err != nil {
		return Entity{}, model.ErrShadowFormat
	}
	d, err := json.Marshal(s.State.Desired)
	if err != nil {
		return Entity{}, model.ErrShadowFormat
	}
	m, err := json.Marshal(s.Metadata)
	if err != nil {
		return Entity{}, model.ErrShadowFormat
	}
	t, err := json.Marshal(s.Tags)
	if err != nil {
		return Entity{}, model.ErrShadowFormat
	}

	return Entity{
		Version:   s.Version,
		ThingId:   s.ThingId,
		Desired:   d,
		Reported:  r,
		Metadata:  m,
		Tags:      t,
		CreatedAt: s.CreatedAt,
	}, nil
}

func toShadow(en Entity) (Shadow, error) {
	var d StateValue
	var r StateValue
	var m Metadata
	var t TagsValue
	if en.Desired != nil {
		err := json.Unmarshal(en.Desired, &d)
		if err != nil {
			return Shadow{}, errors.Wrap(err, "unmarshal desired field")
		}
	}
	if en.Reported != nil {
		err := json.Unmarshal(en.Reported, &r)
		if err != nil {
			return Shadow{}, errors.Wrap(err, "unmarshal reported field")
		}
	}
	if en.Metadata != nil {
		err := json.Unmarshal(en.Metadata, &m)
		if err != nil {
			return Shadow{}, errors.Wrap(err, "unmarshal metadata field")
		}
	}
	if en.Tags != nil {
		err := json.Unmarshal(en.Tags, &t)
		if err != nil {
			return Shadow{}, errors.Wrap(err, "unmarshal tags field")
		}
	}

	return Shadow{
		Version:   en.Version,
		ThingId:   en.ThingId,
		State:     StateDR{Desired: d, Reported: r},
		Metadata:  m,
		Tags:      t,
		CreatedAt: en.CreatedAt,
		UpdatedAt: en.UpdatedAt,
	}, nil
}
