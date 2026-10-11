package services

import (
	"errors"
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
	if errors.Is(err, domain.ErrNodeNotFound) {
		log.Printf("node update skipped, node %s not found", req.NodeId.Value)
		return &domain.NodeUpdateResp{NotFound: true}, nil
	}
	if err != nil {
		return nil, err
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
