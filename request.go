package main

import (
	"fmt"
)

type CardListRequest struct {
	Classes  []string `json:"classes"`
	ManaCost int      `json:"manaCost"`
	Rarity   string   `json:"rarity"`
}

func (c *CardListRequest) validate() error {
	if len(c.Classes) == 0 {
		return fmt.Errorf("fatal: Classes length is zero")
	}

	if !(c.ManaCost > 0) {
		return fmt.Errorf("fatal: Mana cost must be greater than 0")
	}

	if c.Rarity == "" {
		return fmt.Errorf("fatal: Rarity must be set")
	}

	return nil
}
