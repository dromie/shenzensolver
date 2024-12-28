package solver

import (
	"testing"
)

func Test_OBLOCK_cardConstruct(t *testing.T) {
	card := constructCard("O")
	if card.Suit != FLOWER {
		t.Errorf("Card suit is not FLOWER")
	}
	if card.Value != BLOCK {
		t.Errorf("Card value is not 0")
	}
}

func Test_solution(t *testing.T) {
	card := constructCard("r3")
	if card.is_solution(constructCard("g1")) {
		t.Errorf("Card is not a solution")
	}
	if !card.is_solution(constructCard("r2")) {
		t.Errorf("Card is a solution")
	}
}

func Test_can_be_put_over(t *testing.T) {
	card := constructCard("r3")
	if card.can_be_put_over(constructCard("r2")) {
		t.Errorf("Card can't be put over")
	}
	if !card.can_be_put_over(constructCard("g4")) {
		t.Errorf("Card can be put over")
	}
}

func Test_block(t *testing.T) {
	card := constructCard("r3")
	bcard := constructCard("G")
	if card.can_be_put_over(bcard) {
		t.Errorf("Card can't be put over")
	}
	if bcard.can_be_put_over(card) {
		t.Errorf("Card can't be put over")
	}
}
