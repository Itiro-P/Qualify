package pkg

const AnalystSelect = `u.id, u.name, u.email, u.phone, u.time_created,
	u.country_code, u.country_name, u.country_state, u.city, u.timezone,
	a.hourly_rate, a.total_reviews, a.mean_rating`

const AnalystJoin = `"user" u JOIN analyst a ON a.id = u.id`

const CertificationSelect = `id, name, year, description, institution`

const ClientSelect = `u.id, u.name, u.email, u.phone, u.time_created, 
	u.country_code, u.country_name, u.country_state, u.city, u.timezone, 
	c.proposed_budget`

const ClientJoin = `"user" u JOIN analyst a ON a.id = u.id`
