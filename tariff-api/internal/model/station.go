package model

type Station struct {
	ID        int    `json:"id"`
	Kod       string `json:"stan_kod"`
	Name      string `json:"stan_name"`
	Priznak   int    `json:"stan_priznak"`
	Paragraph string `json:"paragraph"`
}
