package controllers

type rolePaging struct {
	Items  []roleResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}
