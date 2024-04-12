package blockchainhash

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ChengjinWu/aescrypto"
	"github.com/yusufwib/blockchain-medical-record/models/dblockchain"
	"github.com/yusufwib/blockchain-medical-record/models/dmedicalrecord"
)

const secretKey = "1234123412341234123412341234abcd"

func CalculateHash(block dblockchain.Block) string {
	record := fmt.Sprintf("%d%s%d%d%d%v%s", block.Index, block.Timestamp, block.PatientID,
		block.DoctorID, block.AppointmentID, block.Data, block.PrevHash)

	hash := sha256.New()
	hash.Write([]byte(record))
	return hex.EncodeToString(hash.Sum(nil))
}

func EncryptStruct(src dmedicalrecord.MedicalRecord) (res string, err error) {
	srcBytes, err := json.Marshal(src)
	if err != nil {
		return
	}

	crypted, err := aescrypto.AesCbcPkcs7Encrypt([]byte(srcBytes), []byte(secretKey), nil)
	if err != nil {
		return
	}

	return base64.URLEncoding.EncodeToString(crypted), nil
}

func DecryptStruct(src string) (res dmedicalrecord.MedicalRecord, err error) {
	crypted, _ := base64.URLEncoding.DecodeString(src)

	data, err := aescrypto.AesCbcPkcs7Decrypt(crypted, []byte(secretKey), nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &res); err != nil {
		return
	}

	return
}
