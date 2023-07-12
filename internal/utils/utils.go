package utils

import (
	"math/rand"
	"time"

	"github.com/YuYAlexey/SecretSanta/internal/model"
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
	var indexGiver, indexRecepient int
	mPlaeyrs := map[uint64]*model.GiverRecipient{}
	rand.Seed(time.Now().UnixNano())
	orinLenGiver := len(giver)

	for len(giver) != 0 && len(recepient) != 0 {

		lenGiver := len(giver)
		lenRecepient := len(recepient)

		switch {
		case (orinLenGiver%2 != 0 && lenGiver == 1 && lenRecepient == 1):
			{
				indexGiver = 0
				indexRecepient = 0
			}
		case orinLenGiver%2 != 0:
			{
				indexGiver = 0
				indexRecepient = 1
			}
		case orinLenGiver%2 == 0:
			{
				indexGiver = rand.Intn(lenGiver)
				indexRecepient = rand.Intn(lenRecepient)
			}
		}

		selectGiver := giver[indexGiver]
		selectRecepient := recepient[indexRecepient]

		if selectGiver == selectRecepient {
			continue
		}

		mPlaeyrs[selectGiver] = &model.GiverRecipient{
			Giver:     selectGiver,
			Recipient: selectRecepient,
		}

		recepient = append(recepient[:indexRecepient], recepient[indexRecepient+1:]...)
		giver = append(giver[:indexGiver], giver[indexGiver+1:]...)
	}

	return &mPlaeyrs
}
