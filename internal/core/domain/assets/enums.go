package assets_dm

/*
 * Type
 */

type (
	Type = string
)

const (
	TypeChart    Type = "CHART"
	TypeInsight  Type = "INSIGHT"
	TypeAudience Type = "AUDIENCE"
)

func Types() []Type {
	return []Type{TypeChart, TypeInsight, TypeAudience}
}

/*
 * Gender
 */

type (
	Gender = string
)

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

func Genders() []Gender {
	return []Gender{GenderMale, GenderFemale}
}

/*
 * AgeGroup
 */

type (
	AgeGroup = string
)

const (
	AgeGroup18TO24    AgeGroup = "18-23"
	AgeGroup24TO35    AgeGroup = "24-35"
	AgeGroup35To45    AgeGroup = "36-45"
	AgeGroup46AndMore AgeGroup = "46+"
)

func AgeGroups() []AgeGroup {
	return []AgeGroup{AgeGroup18TO24, AgeGroup24TO35, AgeGroup35To45, AgeGroup46AndMore}
}
