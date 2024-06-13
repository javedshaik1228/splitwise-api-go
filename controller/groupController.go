package controller

import (
	"fmt"
	"log"
	"net/http"
	"splitwise/models"
	"splitwise/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GroupParams struct {
	GroupName  string   `json:"groupname" binding:"required"`
	MemberList []string `json:"members" binding:"required"`
	AdminID    uint
}

type GroupDetailsRs struct {
	GroupName string   `json:"GroupName"`
	Admin     string   `json:"Admin"`
	Members   []string `json:"Members"`
}

// Handlers

func CreateGroup(db *gorm.DB, c *gin.Context) {
	claims, err := utils.ValidateUserLogin(c)
	if err != nil {
		log.Println("Invalid login")
		SendError(ErrUnauthorized, err.Error(), c)
		return
	}

	var groupParams GroupParams
	if err := c.BindJSON(&groupParams); err != nil {
		SendError(ErrBadRequest, err.Error(), c)
		return
	}
	groupParams.AdminID = claims.UserID

	if err := db.Transaction(func(tx *gorm.DB) error {
		groupID, err := insertGroupRow(tx, groupParams.GroupName, groupParams.AdminID)
		if err != nil {
			return err
		}
		return insertGroupMembers(tx, groupParams.MemberList, groupID)
	}); err != nil {
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Status": "OK"})
}

func DeleteGroup(db *gorm.DB, c *gin.Context) {
	claims, err := utils.ValidateUserLogin(c)
	if err != nil {
		log.Println("Invalid login")
		SendError(ErrUnauthorized, err.Error(), c)
		return
	}

	groupName := c.Query("groupname")
	if groupName == "" {
		SendError(ErrBadRequest, "Invalid groupname", c)
		return
	}

	var group models.Group
	if err := db.Where("group_name = ? AND admin_id = ?", groupName, claims.UserID).First(&group).Error; err != nil {
		log.Println(err.Error())
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", group.GroupID).Delete(&models.Groupmember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&group).Error
	}); err != nil {
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": fmt.Sprintf("Deleted group: %v", groupName)})
}

func AddMembersToGroup(db *gorm.DB, c *gin.Context) {
	claims, err := utils.ValidateUserLogin(c)
	if err != nil {
		log.Println("Invalid login")
		SendError(ErrUnauthorized, err.Error(), c)
		return
	}

	var groupParams GroupParams
	if err := c.BindJSON(&groupParams); err != nil {
		SendError(ErrBadRequest, err.Error(), c)
		return
	}
	groupParams.AdminID = claims.UserID

	var group models.Group
	if err := db.Where("group_name = ? AND admin_id = ?", groupParams.GroupName, groupParams.AdminID).First(&group).Error; err != nil {
		log.Println(err.Error())
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		return insertGroupMembers(tx, groupParams.MemberList, group.GroupID)
	}); err != nil {
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Status": "OK"})
}

func GetGroupDetails(db *gorm.DB, c *gin.Context) {
	claims, err := utils.ValidateUserLogin(c)
	if err != nil {
		log.Println("Invalid login")
		SendError(ErrUnauthorized, err.Error(), c)
		return
	}

	groupName := c.Query("groupname")
	if groupName == "" {
		SendError(ErrBadRequest, "Invalid groupname", c)
		return
	}

	// Fetch group info:
	var group models.Group
	err = db.Table("groups g").
		Select("g.group_id, g.group_name, g.admin_id, g.created_at").
		Joins("LEFT JOIN groupmembers gm ON g.group_id = gm.group_id").
		Where("g.group_name = ? AND gm.user_id = ?", groupName, claims.UserID).
		First(&group).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			SendError(ErrNotFound, "No group found with the given data", c)
		} else {
			SendError(ErrInternalFailure, err.Error(), c)
		}
		return
	}

	log.Printf("retrieved group details: %v", group)

	// Fetch admin username
	adminUsername, err := getUsernameFromUserID(db, group.AdminID)
	if err != nil {
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	// Fetch members list
	var members []string
	err = db.Table("users u").
		Select("u.username").
		Joins("INNER JOIN groupmembers gm ON u.user_id = gm.user_id").
		Where("gm.group_id = ?", group.GroupID).
		Pluck("username", &members).Error

	if err != nil {
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	groupDetails := GroupDetailsRs{
		GroupName: groupName,
		Admin:     adminUsername,
		Members:   members,
	}

	c.JSON(http.StatusOK, groupDetails)
}

// DB helpers

func insertGroupMembers(tx *gorm.DB, membersList []string, groupID uint) error {
	for _, member := range membersList {
		userID, err := getUserIDfromUsername(tx, member)
		if err != nil {
			log.Printf("Error while getting user ID from username: %v\n", err)
			return err
		}
		groupMember := models.Groupmember{GroupID: groupID, UserID: userID}
		if err := tx.Create(&groupMember).Error; err != nil {
			log.Printf("Error while creating group member: %v\n", err)
			return err
		}
		log.Printf("Inserted group member row: %v\n", groupMember)
	}
	return nil
}

func insertGroupRow(tx *gorm.DB, groupName string, adminID uint) (uint, error) {
	group := models.Group{GroupName: groupName, AdminID: adminID}
	if err := tx.Create(&group).Error; err != nil {
		return 0, fmt.Errorf("failed to create group: %w", err)
	}
	log.Printf("Inserted group row: %v\n", group)

	if err := addAdminToGroup(tx, group.GroupID, adminID); err != nil {
		return 0, err
	}

	return group.GroupID, nil
}

func addAdminToGroup(tx *gorm.DB, groupID uint, adminID uint) error {
	groupMember := models.Groupmember{GroupID: groupID, UserID: adminID}
	if err := tx.Create(&groupMember).Error; err != nil {
		log.Printf("Error while creating group member: %v\n", err)
		return err
	}
	log.Println("Added admin to group")
	return nil
}
