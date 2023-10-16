package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Security struct {
	SecurityID   primitive.ObjectID `bson:"_id,omitempty"`
	Symbol       *string            `bson:"symbol,omitempty"`
	SecurityName *string            `bson:"security_name,omitempty"`
	SecurityType *uint32            `bson:"security_type,omitempty"`
	Region       *string            `bson:"region,omitempty"`
	Currency     *string            `bson:"currency,omitempty"`
}

func ToSecurityModelFromEntity(s *entity.Security) *Security {
	if s == nil {
		return nil
	}

	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(s.GetSecurityID()) {
		objID, _ = primitive.ObjectIDFromHex(s.GetSecurityID())
	}

	return &Security{
		SecurityID:   objID,
		Symbol:       s.Symbol,
		SecurityName: s.SecurityName,
		SecurityType: s.SecurityType,
		Region:       s.Region,
		Currency:     s.Currency,
	}
}

func ToSecurityEntity(s *Security) *entity.Security {
	if s == nil {
		return nil
	}

	return entity.NewSecurity(
		s.GetSymbol(),
		entity.WithSecurityID(goutil.String(s.GetSecurityID())),
		entity.WithSecurityName(s.SecurityName),
		entity.WithSecurityType(s.SecurityType),
		entity.WithSecurityRegion(s.Region),
		entity.WithSecurityCurrency(s.Currency),
	)
}

func (s *Security) GetSecurityID() string {
	if s != nil {
		return s.SecurityID.Hex()
	}
	return ""
}

func (s *Security) GetSymbol() string {
	if s != nil && s.Symbol != nil {
		return *s.Symbol
	}
	return ""
}

func (s *Security) GetSecurityName() string {
	if s != nil && s.SecurityName != nil {
		return *s.SecurityName
	}
	return ""
}

func (s *Security) GetSecurityType() uint32 {
	if s != nil && s.SecurityType != nil {
		return *s.SecurityType
	}
	return 0
}

func (s *Security) GetRegion() string {
	if s != nil && s.Region != nil {
		return *s.Region
	}
	return ""
}

func (s *Security) GetCurrency() string {
	if s != nil && s.Currency != nil {
		return *s.Currency
	}
	return ""
}
