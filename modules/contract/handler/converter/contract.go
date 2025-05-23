package converter

import (
	apiModel "diploma/modules/contract/handler/model"
	serviceModel "diploma/modules/contract/model"
)

func ToAPI(c *serviceModel.Contract) *apiModel.ContractResponse {
	return &apiModel.ContractResponse{
		ID:          c.ID,
		Content:     c.Content,
		Status:      c.Status,
		SupplierSig: c.SupplierSig,
		CustomerSig: c.CustomerSig,
	}
}
