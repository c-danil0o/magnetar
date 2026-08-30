package services

import (
	"log"

	"github.com/c12s/magnetar/internal/domain"
)

type NodeUpdateService struct {
	nodeRepo domain.NodeRepo
}

func NewNodeUpdateService(nodeRepo domain.NodeRepo) (*NodeUpdateService, error) {
	return &NodeUpdateService{
		nodeRepo: nodeRepo,
	}, nil
}

func (s *NodeUpdateService) Update(req domain.NodeUpdateReq) (*domain.NodeUpdateResp, error) {
	node, err := s.nodeRepo.FindById(req.NodeId)
	if err != nil {
		log.Printf("node update skipped, node %s not found", req.NodeId.Value)
		return &domain.NodeUpdateResp{}, nil
	}
	// only overwrite fields the reporter actually sent: a nil field means
	// "unchanged", otherwise every status tick would clobber capabilities
	if req.Capabilities != nil {
		node.Capabilities = req.Capabilities
	}
	if req.OsStatus != nil {
		node.OsStatus = req.OsStatus
	}
	err = s.nodeRepo.Put(*node)
	if err != nil {
		return nil, err
	}
	return &domain.NodeUpdateResp{}, nil
}
