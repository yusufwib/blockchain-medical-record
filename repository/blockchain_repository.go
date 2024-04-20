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
	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"
	"github.com/yusufwib/blockchain-medical-record/utils/randstr"
)

type BlockchainRepository struct {
	LevelDB *leveldb.DB
	Config  *config.ConfigGroup
	Logger  mlog.Logger
}

func NewBlockchainRepository(levelDB *leveldb.DB, cfg *config.ConfigGroup, log mlog.Logger) BlockchainRepository {
	return BlockchainRepository{levelDB, cfg, log}
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

	log.Printf("adding block for patient ID: %d", req.PatientID)
	var prevBlock dblockchain.Block
	blockByPatient := r.GetBlocksByPatientID(req.PatientID)
	if len(blockByPatient) == 0 {
		log.Printf("creating genesis block for patient ID: %d", req.PatientID)
		genesis := dblockchain.Block{
			Index:     0,
			Timestamp: time.Now().String(),
			PatientID: req.PatientID,
			Message:   fmt.Sprintf("Genesis Block for Patient ID: %d", req.PatientID),
		}
		genesis.Hash = blockchainhash.CalculateHash(genesis)

		bc.Chain = append(bc.Chain, genesis)
		prevBlock = genesis
		r.saveBlockchain(genesis)
	} else {
		prevBlock = blockByPatient[len(blockByPatient)-1]
	}

	log.Printf("encyrpting data for patient ID: %d", req.PatientID)
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

	log.Printf("creating block for patient ID: %d", req.PatientID)

	log.Println("checking previous block hash...")
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
	r.saveBlockchain(newBlock)
	// log.Println("previous block hash is not valid")
	// log.Println("failed to create block")

	log.Println("previous block hash is valid")
	log.Println("block added successfully")
	return newBlock
}

func (r *BlockchainRepository) saveBlockchain(new dblockchain.Block) {
	// data, err := json.Marshal(bc.Chain)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	blocks := r.GetAllBlocks()
	blocks = append(blocks, new)

	data, err := json.Marshal(blocks)
	if err != nil {
		log.Println(err)
		return
	}

	if err := r.LevelDB.Put([]byte(blockchainKey), data, nil); err != nil {
		log.Println(err)
	}
}

func (r *BlockchainRepository) GetBlocksByPatientID(patientID uint64) (res []dblockchain.Block) {
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

func (r *BlockchainRepository) GetAllBlocksDecrypted() (res []dblockchain.Block) {
	blocks := r.GetAllBlocks()
	for _, block := range blocks {
		if block.EncryptedData == "" {
			continue
		}
		decryptedData, _ := blockchainhash.DecryptStruct(block.EncryptedData)
		block.Data = decryptedData

		res = append(res, block)
	}

	return res
}
