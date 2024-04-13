package service

import (
	"github.com/yusufwib/blockchain-medical-record/models/dblockchain"
	"github.com/yusufwib/blockchain-medical-record/repository"
)

type BlockchainService struct {
	BlockchainRepository repository.BlockchainRepository
}

func NewBlockchainService(r repository.BlockchainRepository) BlockchainService {
	return BlockchainService{
		BlockchainRepository: r,
	}
}

func (s BlockchainService) MineAll() (res []dblockchain.Block) {
	return s.BlockchainRepository.GetAllBlocks()
}
