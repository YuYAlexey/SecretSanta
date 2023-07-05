package utils

import (
	"math/rand"
	"time"

	"github.com/adYushinW/SecretSanta/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SecretSanta(giver, recepient []uint64) (giverRecepient *map[uint64]*model.GiverRecipient) {

	rand.Seed(time.Now().UnixNano())
	// Установи https://marketplace.visualstudio.com/items?itemName=streetsidesoftware.code-spell-checker
	// и проблем с неймингом будешь меньше
	// mPlaeyrs => mPlayers
	mPlaeyrs := map[uint64]*model.GiverRecipient{}

	if len(giver)%2 != 0 {
		return nil
	}

	// recepient => recipient
	for len(giver) != 0 && len(recepient) != 0 {
		lenGiver := len(giver)
		// lenRecepient => lenRecipient
		lenRecepient := len(recepient)
		indexGiver := rand.Intn(lenGiver)
		// indexRecepient => indexRecipient
		indexRecepient := rand.Intn(lenRecepient)
		selectGiver := giver[indexGiver]
		// recepient => recipient
		selectRecepient := recepient[indexRecepient]

		if selectGiver == selectRecepient {
			continue
		}
		// mPlaeyrs => mPlayers
		mPlaeyrs[selectGiver] = &model.GiverRecipient{
			Giver:     selectGiver,
			Recipient: selectRecepient,
		}
		// recepient => recipient
		recepient = append(recepient[:indexRecepient], recepient[indexRecepient+1:]...)
		giver = append(giver[:indexGiver], giver[indexGiver+1:]...)
	}
	// mPlaeyrs => mPlayers
	return &mPlaeyrs
}
