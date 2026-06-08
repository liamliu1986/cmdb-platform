package core

type CreateCITypeRequest struct {
	Name         string `json:"name" binding:"required"`
	Alias        string `json:"alias"`
	UniqueAttrID uint   `json:"unique_attr_id"`
	Icon         string `json:"icon"`
}

type CreateAttributeRequest struct {
	Name        string `json:"name" binding:"required"`
	Alias       string `json:"alias"`
	ValueType   string `json:"value_type" binding:"required"`
	IsChoice    bool   `json:"is_choice"`
	IsList      bool   `json:"is_list"`
	IsUnique    bool   `json:"is_unique"`
	IsIndex     bool   `json:"is_index"`
	IsPassword  bool   `json:"is_password"`
	IsComputed  bool   `json:"is_computed"`
	ComputeExpr string `json:"compute_expr"`
}

type CreateCIRequest struct {
	CITypeID   uint                   `json:"ci_type_id" binding:"required"`
	AttrValues map[string]interface{} `json:"attr_values"`
}

type CISearchRequest struct {
	Q        string `form:"q"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=25"`
	Sort     string `form:"sort"`
}
