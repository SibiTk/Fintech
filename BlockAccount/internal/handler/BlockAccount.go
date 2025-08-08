package handler

import (
	"context"
	"fmt"
	"log"

	v1 "BlockAccount/api/helloworld/v1"
	accountpb "account/api/helloworld/v1"

	"gorm.io/gorm"
)



func (h *AccountBlockHandler) AccountBlock(ctx context.Context, in *v1.SaveAccBlockRequest) (*v1.SaveAccBlockReply, error) {
    
    // Validate input source
    validSources := []string{"risk block", "case block", "admin block"}
    isValidSource := false
    for _, validSource := range validSources {
        if in.Source == validSource {
            isValidSource = true
            break
        }
    }

    if !isValidSource {
        return nil, fmt.Errorf("invalid source: %s. Valid sources are: risk block, case block, admin block", in.Source)
    }
    

    accountResp, err := h.AccountClient.GetCustomerWithId(ctx, &accountpb.GetCustomerWithIdRequest{
        CustomerId: in.CustomerId,
    })
    if err != nil {
        log.Println("Failed to fetch customer from account service:", err)
        return nil, fmt.Errorf("failed to verify customer: %w", err)
    }
    if len(accountResp.Accounts) == 0 {
        return nil, fmt.Errorf("no account found for customer ID %d", in.CustomerId)
    }
    
    account := accountResp.Accounts[0]
    
   
    existingBlock, err := h.repo.GetByCustomerAndSource(ctx, in.CustomerId, in.Source)
    if err != nil && err != gorm.ErrRecordNotFound {
        return nil, fmt.Errorf("failed to check existing blocks: %w", err)
    }
  
    if existingBlock != nil {
        return nil, fmt.Errorf("already %s is done for customer ID %d", in.Source, in.CustomerId)
    }
    
   
    var status, description string
    if account.Status == "active" {
        status = "deactive"
        description = "Account is Deactive"
    } else {
        status = "active"
        description = "Account is Active"
    }
    

    block := &v1.SaveAccBlockRequest{
        CustomerId:  in.CustomerId,
        Description: description,
        Source:      in.Source,  
        Status:      status,
    }
    
    if err := h.repo.Save(ctx, block); err != nil {
        return nil, fmt.Errorf("failed to save block record: %w", err)
    }
    
    return &v1.SaveAccBlockReply{
        CustomerId:  block.CustomerId,
        Description: block.Description,
        Source:      block.Source,
        Status:      block.Status,
        Message:     fmt.Sprintf("%s created successfully", block.Source),
    }, nil
}


// UpdateAccBlock: Update existing block by customer ID
func (h *AccountBlockHandler) UpdateAccBlock(ctx context.Context, in *v1.UpdateAccBlockRequest) (*v1.UpdateAccBlockReply, error) {
	updated, err := h.repo.UpdateAccBlock(ctx, in)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateAccBlockReply{
		CustomerId:  updated.CustomerId,
		Description: updated.Description,
		Source:      updated.Source,
		Status: updated.Status,
		Message:     "Updated Block Info",
	}, nil
}

// Mapper: convert input to domain model (optional)
func Mapper(in *v1.SaveAccBlockRequest) *v1.SaveAccBlockRequest {
	return &v1.SaveAccBlockRequest{
		CustomerId:  in.CustomerId,
		Description: in.Description,
		Source:      in.Source,
	}
}
