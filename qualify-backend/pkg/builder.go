package pkg

import (
	"strings"

	"github.com/lib/pq"
	"github.com/n-r-w/squirrel"
)

func BuildFilterUser(builder squirrel.SelectBuilder, filters UserFilter) squirrel.SelectBuilder {
	toReturn := builder
	filters.Normalize()
	if filters.Name != "" {
		toReturn = toReturn.Where(squirrel.ILike{"u.name": PutPercent(filters.Name)})
	}

	if filters.Email != "" {
		toReturn = toReturn.Where(squirrel.ILike{"u.email": PutPercent(filters.Email)})
	}

	if filters.Country != "" {
		toReturn = toReturn.Where(squirrel.ILike{"u.country_name": PutPercent(filters.Country)})
	}

	if filters.CountryCode != "" {
		toReturn = toReturn.Where(squirrel.ILike{"u.country_code": PutPercent(filters.CountryCode)})
	}

	if filters.CountryState != "" {
		toReturn = toReturn.Where(squirrel.ILike{"u.country_state": PutPercent(filters.CountryState)})
	}

	if filters.City != "" {
		toReturn = toReturn.Where(squirrel.ILike{"u.city": PutPercent(filters.City)})
	}

	if filters.Timezone != "" {
		toReturn = toReturn.Where(squirrel.ILike{"u.timezone": PutPercent(filters.Timezone)})
	}

	return toReturn
}

func BuildFilterAnalyst(builder squirrel.SelectBuilder, filters AnalystFilter) squirrel.SelectBuilder {
	filters.Normalize()
	toReturn := BuildFilterUser(builder, filters.UserFilter)

	if filters.Skills != "" {
		var skills = strings.Split(filters.Skills, ",")
		var trimmed []string
		for i := range skills {
			trimmed = append(trimmed, PutPercent(strings.TrimSpace(skills[i])))
		}
		toReturn = toReturn.Where(
			`a.id IN (
            SELECT ask.analyst_id
            FROM analyst_skill ask
            JOIN skill s ON s.id = ask.skill_id
            WHERE s.name ILIKE ANY(?)
            GROUP BY ask.analyst_id
            HAVING COUNT(DISTINCT s.id) = ?)`,
			pq.Array(trimmed),
			len(trimmed),
		)
	}

	if filters.MinHourlyRate != nil {
		toReturn = toReturn.Where(squirrel.GtOrEq{"a.hourly_rate": *filters.MinHourlyRate})
	}

	if filters.MaxHourlyRate != nil {
		toReturn = toReturn.Where(squirrel.LtOrEq{"a.hourly_rate": *filters.MaxHourlyRate})
	}

	if filters.MinTotalReviews != nil {
		toReturn = toReturn.Where(squirrel.GtOrEq{"a.total_reviews": *filters.MinTotalReviews})
	}

	if filters.MaxTotalReviews != nil {
		toReturn = toReturn.Where(squirrel.LtOrEq{"a.total_reviews": *filters.MaxTotalReviews})
	}

	if filters.MinRating != nil {
		toReturn = toReturn.Where(squirrel.GtOrEq{"a.mean_rating": *filters.MinRating})
	}

	if filters.MaxRating != nil {
		toReturn = toReturn.Where(squirrel.LtOrEq{"a.mean_rating": *filters.MaxRating})
	}

	orderClause := filters.SortOptions.ValidateSort(AnalystSortFields)

	if orderClause != "" {
		toReturn = toReturn.OrderBy(orderClause)
	} else {
		toReturn = toReturn.OrderBy("u.time_created DESC")
	}

	toReturn = toReturn.Limit(uint64(filters.PageSize)).Offset(uint64(filters.Offset()))
	return toReturn
}

func BuildFilterClient(builder squirrel.SelectBuilder, filters ClientFilter) squirrel.SelectBuilder {
	filters.Normalize()
	toReturn := BuildFilterUser(builder, filters.UserFilter)

	if filters.MinProposedBudget != nil {
		toReturn = toReturn.Where(squirrel.GtOrEq{"c.proposed_budget": *filters.MinProposedBudget})
	}

	if filters.MaxProposedBudget != nil {
		toReturn = toReturn.Where(squirrel.LtOrEq{"c.proposed_budget": *filters.MaxProposedBudget})
	}

	orderClause := filters.SortOptions.ValidateSort(ClientSortFields)

	if orderClause != "" {
		toReturn = toReturn.OrderBy(orderClause)
	} else {
		toReturn = toReturn.OrderBy("u.time_created DESC")
	}

	toReturn = toReturn.Limit(uint64(filters.PageSize)).Offset(uint64(filters.Offset()))
	return toReturn
}

func BuildUpdateUser(req UserUpdateRequest) map[string]any {
	fields := map[string]any{}
	if req.Name != nil {
		fields["name"] = *req.Name
	}
	if req.Email != nil {
		fields["email"] = *req.Email
	}
	if req.Phone != nil {
		fields["phone"] = *req.Phone
	}
	if req.Country_code != nil {
		fields["country_code"] = *req.Country_code
	}
	if req.Country_name != nil {
		fields["country_name"] = *req.Country_name
	}
	if req.Country_state != nil {
		fields["country_state"] = *req.Country_state
	}
	if req.City != nil {
		fields["city"] = *req.City
	}
	if req.Timezone != nil {
		fields["timezone"] = *req.Timezone
	}
	return fields
}

func BuildUpdateAnalyst(req AnalystUpdateRequest) (userFields map[string]any, analystFields map[string]any) {
	userFields = BuildUpdateUser(req.UserUpdateRequest)
	analystFields = map[string]any{}
	if req.Hourly_rate != nil {
		analystFields["hourly_rate"] = *req.Hourly_rate
	}
	if req.Total_reviews != nil {
		analystFields["total_reviews"] = *req.Total_reviews
	}
	if req.Mean_rating != nil {
		analystFields["mean_rating"] = *req.Mean_rating
	}
	return userFields, analystFields
}

func BuildUpdateClient(req ClientUpdateRequest) (userFields map[string]any, clientFields map[string]any) {
	userFields = BuildUpdateUser(req.UserUpdateRequest)
	clientFields = map[string]any{}
	if req.Proposed_budget != nil {
		clientFields["proposed_budget"] = *req.Proposed_budget
	}
	return userFields, clientFields
}
