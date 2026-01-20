package constant

var (
	AdminAllowedSorts = map[string]string{
		"date":   "u.created_at",
		"name":   "ud.full_name",
		"assign": "is_assigned",
	}
	AllowedOrderDir = map[string]struct{}{
		"asc":  {},
		"desc": {},
	}
)
