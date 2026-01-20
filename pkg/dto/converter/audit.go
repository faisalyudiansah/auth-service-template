package converter

import (
	"github.com/faisalyudiansah/auth-service-template/pkg/dto"
	"github.com/faisalyudiansah/auth-service-template/pkg/entity"
)

func ConvertAudit(e entity.Audit) dto.Audit {
	return dto.Audit{
		CreatedAt: e.CreatedAt,
		CreatedBy: e.CreatedBy,
		UpdatedAt: e.UpdatedAt,
		UpdatedBy: e.UpdatedBy,
		DeletedAt: e.DeletedAt,
		DeletedBy: e.DeletedBy,
	}
}
