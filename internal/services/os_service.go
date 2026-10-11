package services

import "github.com/c12s/magnetar/internal/domain"

// OsService holds the operator's declared OS image per node. Observed OS
// state (os_status) is reported by the node through NodeUpdateService.
type OsService struct {
	nodeRepo domain.NodeRepo
}

func NewOsService(nodeRepo domain.NodeRepo) (*OsService, error) {
	return &OsService{nodeRepo: nodeRepo}, nil
}

// SetDesiredOsImage records (or, with nil, clears) the declared image of an
// existing node.
func (s *OsService) SetDesiredOsImage(nodeId domain.NodeId, image *domain.DesiredOsImage) error {
	if _, err := s.nodeRepo.FindById(nodeId); err != nil {
		return err
	}
	return s.nodeRepo.PutDesiredOsImage(nodeId, image)
}

// GetDesiredOsImage returns the declared image, nil if none;
// domain.ErrNodeNotFound if the node is unknown.
func (s *OsService) GetDesiredOsImage(nodeId domain.NodeId) (*domain.DesiredOsImage, error) {
	node, err := s.nodeRepo.FindById(nodeId)
	if err != nil {
		return nil, err
	}
	if node.DesiredOsImage == nil {
		return nil, nil
	}
	return node.DesiredOsImage, nil
}
