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

type GroupParms struct {
	GroupName  string   `json:"groupname" binding:"required"`
	MemberList []string `json:"members" binding:"required"`
	AdminId    uint
}

func CreateGroup(db *gorm.DB, c *gin.Context) {

	// check if valid login
	claims, err := utils.ValidateUserLogin(c)
	if err != nil {
		log.Printf("Invalid login")
		SendError(ErrUnauthorized, err.Error(), c)
		return
	}

	var groupParams GroupParms

	// check if valid request
	if err := c.BindJSON(&groupParams); err != nil {
		SendError(ErrBadRequest, err.Error(), c)
		return
	}

	groupParams.AdminId = claims.UserID

	// Begin Transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Insert group row
	var groupID uint

	if groupID, err = InsertGroupRow(tx, groupParams.GroupName, groupParams.AdminId); err != nil {
		tx.Rollback()
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	membersList := groupParams.MemberList

	// Insert group members rows
	for _, member := range membersList {
		if err := InsertGroupMember(tx, member, groupID); err != nil {
			tx.Rollback()
			SendError(ErrInternalFailure, err.Error(), c)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Status:": "OK"})
}

func DeleteGroup(db *gorm.DB, c *gin.Context) {
	// check if valid login
	claims, err := utils.ValidateUserLogin(c)
	if err != nil {
		log.Printf("Invalid login")
		SendError(ErrUnauthorized, err.Error(), c)
		return
	}

	groupName := c.Query("groupname")

	if groupName == "" {
		SendError(ErrBadRequest, "Invalid groupname", c)
		return
	}

	// check if the user is admin of that group and get group model
	var group models.Group
	if err := db.Find(&group).Where("group_name=? AND admin_id=?", groupName, claims.UserID).Error; err != nil {
		log.Print(err.Error())
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	log.Printf("Deleting group: %v\n", group)

	// Begin Transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			SendError(ErrInternalFailure, fmt.Sprintf("panic occured: %v", r), c)
		}
	}()

	// Delete groupmember rows
	if err := tx.Where("group_id = ?", group.GroupID).Delete(&models.Groupmember{}).Error; err != nil {
		log.Print(err.Error())
		SendError(ErrInternalFailure, err.Error(), c)
		tx.Rollback()
		return
	}

	// Delete group row
	if err := tx.Delete(&group).Error; err != nil {
		log.Print(err.Error())
		SendError(ErrInternalFailure, err.Error(), c)
		tx.Rollback()
		return
	}

	// Commit and check for errors
	if err := tx.Commit().Error; err != nil {
		log.Print(err.Error())
		SendError(ErrInternalFailure, err.Error(), c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": fmt.Sprintf("Deleted group: %v", groupName)})
}

func InsertGroupRow(db *gorm.DB, groupname string, adminId uint) (uint, error) {
	var group models.Group
	group.AdminId = adminId
	group.GroupName = groupname

	if db.Create(&group).Error != nil {
		return 0, fmt.Errorf("failed to create group")
	}

	log.Printf("Inserted group row: %v\n", group)

	if err := AddAdminToGroup(db, group.GroupID, adminId); err != nil {
		return 0, err
	}

	return group.GroupID, nil
}

func InsertGroupMember(db *gorm.DB, memberUsername string, groupID uint) error {
	userId, err := GetUserIDfromUsername(db, memberUsername)
	if err != nil {
		log.Printf("Errow while getting user id from username")
		return err
	}
	var groupMember models.Groupmember
	groupMember.GroupID = groupID
	groupMember.UserID = userId

	err = db.Model(&models.Groupmember{}).Create(&groupMember).Error
	if err != nil {
		log.Printf("Error while creating groupmember")
		return err
	}

	log.Printf("Inserted groupmember row: %v\n", groupMember)
	return nil
}

func AddAdminToGroup(db *gorm.DB, groupID uint, adminId uint) error {
	// Insert a row in groupmembers for the group admin
	var groupMember models.Groupmember
	groupMember.GroupID = groupID
	groupMember.UserID = adminId

	err := db.Model(&models.Groupmember{}).Create(&groupMember).Error
	if err != nil {
		log.Printf("Error while creating groupmember")
		return err
	}

	log.Printf("Added admin to group")
	return nil
}
