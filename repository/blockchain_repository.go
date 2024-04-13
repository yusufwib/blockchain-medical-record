package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/yusufwib/blockchain-medical-record/config"
	"github.com/yusufwib/blockchain-medical-record/models/dblockchain"
	"github.com/yusufwib/blockchain-medical-record/models/dmedicalrecord"
	blockchainhash "github.com/yusufwib/blockchain-medical-record/utils/blockchain_hash"
	"github.com/yusufwib/blockchain-medical-record/utils/randstr"
)

type BlockchainRepository struct {
	LevelDB *leveldb.DB
	Config  *config.ConfigGroup
}

func NewBlockchainRepository(levelDB *leveldb.DB, cfg *config.ConfigGroup) BlockchainRepository {
	return BlockchainRepository{levelDB, cfg}
}

// Blockchain represents the blockchain as a slice of blocks.
type Blockchain struct {
	Chain []dblockchain.Block
	sync.Mutex
}

var (
	bc            Blockchain
	blockchainKey = "chain"
)

func (r *BlockchainRepository) AddBlockMedicalRecord(req dmedicalrecord.MedicalRecordRequest) dblockchain.Block {
	bc.Lock()
	defer bc.Unlock()

	var prevBlock dblockchain.Block
	blockByPatient := r.getBlocksByPatientID(req.PatientID)
	if len(blockByPatient) == 0 {
		genesis := dblockchain.Block{
			Index:     0,
			Timestamp: time.Now().String(),
			PatientID: req.PatientID,
			Message:   fmt.Sprintf("Genesis Block for Patient ID: %d", req.PatientID),
		}
		genesis.Hash = blockchainhash.CalculateHash(genesis)

		bc.Chain = append(bc.Chain, genesis)
		prevBlock = genesis
		r.saveBlockchain()
	} else {
		prevBlock = blockByPatient[len(blockByPatient)-1]
	}

	// encrypt block
	encryptedData, _ := blockchainhash.EncryptStruct(dmedicalrecord.MedicalRecord{
		DoctorID:               req.DoctorID,
		PatientID:              req.PatientID,
		AppointmentID:          req.AppointmentID,
		Diagnose:               req.Diagnose,
		Notes:                  req.Notes,
		AdditionalDocumentPath: req.AdditionalDocumentPath,
		Prescription:           req.Prescription,
		CreatedAt:              time.Now(),
		MedicalRecordNumber:    randstr.GenerateRandomString("EMR"),
	})

	newBlock := dblockchain.Block{
		Index:         prevBlock.Index + 1,
		Timestamp:     time.Now().String(),
		PatientID:     req.PatientID,
		DoctorID:      req.DoctorID,
		AppointmentID: req.AppointmentID,
		EncryptedData: encryptedData,
		PrevHash:      prevBlock.Hash,
	}
	newBlock.Hash = blockchainhash.CalculateHash(newBlock)

	bc.Chain = append(bc.Chain, newBlock)
	r.saveBlockchain()
	return newBlock
}

func (r *BlockchainRepository) saveBlockchain() {
	data, err := json.Marshal(bc.Chain)
	if err != nil {
		log.Println(err)
		return
	}

	if err := r.LevelDB.Put([]byte(blockchainKey), data, nil); err != nil {
		log.Println(err)
	}
}

func (r *BlockchainRepository) getBlocksByPatientID(patientID uint64) (res []dblockchain.Block) {
	blocks := r.GetAllBlocks()

	for _, block := range blocks {
		if block.PatientID == patientID {
			if block.EncryptedData == "" {
				continue
			}
			decryptedData, _ := blockchainhash.DecryptStruct(block.EncryptedData)
			block.Data = decryptedData

			res = append(res, block)
		}
	}
	return res
}

func (r *BlockchainRepository) GetBlocksByAppointmentID(appointmentID uint64) (res dmedicalrecord.MedicalRecord) {
	data := r.GetAllBlocks()

	for _, block := range data {
		if block.AppointmentID == appointmentID {
			decryptedData, _ := blockchainhash.DecryptStruct(block.EncryptedData)
			return decryptedData
		}
	}

	return
}

func (r *BlockchainRepository) GetAllBlocks() (response []dblockchain.Block) {
	data, _ := r.LevelDB.Get([]byte(blockchainKey), nil)
	json.Unmarshal(data, &response)

	return
}
