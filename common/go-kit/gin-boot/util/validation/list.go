package validation

import (
	"fmt"
	"regexp"
	"strconv"
)

// Query Query
func Query(query string) error {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, query)
	if !match {
		return fmt.Errorf("Query '%v' should be made up of a-zA-z0-9_", query)
	}

	return nil
}

// Order Order
func Order(order string) error {
	match, _ := regexp.MatchString(`^[a-zA-Z_]+$`, order)
	if !match {
		return fmt.Errorf("Order '%v' should be made up of a-zA-z_", order)
	}

	return nil
}

// ListReq ListReq
type ListReq struct {
	Query    string
	Order    string
	Desc     bool
	Page     int
	PageSize int
}

// List List
func List(query string, order string, desc string, page string, pageSize string) (req ListReq, err error) {
	if query != "" {
		if err := Query(query); err != nil {
			return req, err
		}
	}
	req.Query = query

	if err := Order(order); err != nil {
		return req, err
	}
	req.Order = order

	req.Desc, err = strconv.ParseBool(desc)
	if err != nil {
		return req, err
	}

	req.Page, err = strconv.Atoi(page)
	if err != nil {
		return req, err
	}
	if err := Positive(req.Page); err != nil {
		return req, err
	}

	req.PageSize, err = strconv.Atoi(pageSize)
	if err != nil {
		return req, err
	}
	if err := Positive(req.PageSize); err != nil {
		return req, err
	}

	return req, nil
}
