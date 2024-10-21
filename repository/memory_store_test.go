package repository

import (
	"github.com/david-pawlowski/giveaway/models"
	"slices"
	"testing"
)

func TestRepository(t *testing.T) {
	t.Run("Two codes properly adds", func(t *testing.T) {
		store := InMemoryStore{}
		giveaway1 := models.Giveaway{Game: "TestGame", Code: "1234-1234", Claimed: false}
		giveaway2 := models.Giveaway{Game: "TestGame2", Code: "4321-4321", Claimed: false}
		store.Add(giveaway1)
		store.Add(giveaway2)
		if len(store) != 2 {
			t.Errorf("got %d len, expected %d len", len(store), 2)
		}
		if store[0].Code != "1234-1234" && store[1].Code != "4321-4321" {
			t.Errorf("Wrong codes in store")
		}
	})
	t.Run("Get all codes from store", func(t *testing.T) {
		store := InMemoryStore{}
		giveaway1 := models.Giveaway{Game: "TestGame", Code: "1234-1234", Claimed: false}
		giveaway2 := models.Giveaway{Game: "TestGame2", Code: "4321-4321", Claimed: false}
		giveaway3 := models.Giveaway{Game: "TestGame3", Code: "321-321", Claimed: false}
		giveaway4 := models.Giveaway{Game: "TestGame4", Code: "21-21", Claimed: false}
		store.Add(giveaway1)
		store.Add(giveaway2)
		store.Add(giveaway3)
		store.Add(giveaway4)
		keys := []string{"1234-1234", "4321-4321"}
		ga1, _ := store.GetRandomCode()
		ga2, _ := store.GetRandomCode()
		ga3, _ := store.GetRandomCode()
		ga4, _ := store.GetRandomCode()
		if !slices.Contains(keys, ga1.Code) || !slices.Contains(keys, ga2.Code) {
			t.Errorf("Didnt get expected codes")
		}
		if ga1.Claimed == false || ga2.Claimed == false || ga3.Claimed == false || ga4.Claimed == false {
			t.Errorf("Codes should be marked as claimed after being randomly returned")
		}

	})
}
