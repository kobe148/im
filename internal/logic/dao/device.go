package dao

import (
	"im/internal/logic/model"
	"im/pkg/db"
	"im/pkg/gerrors"
	"time"

	"github.com/jinzhu/gorm"
)

type deviceDao struct{}

var DeviceDao = new(deviceDao)

// Insert 插入一条设备信息
func (*deviceDao) Add(device model.Device) (int64, error) {
	device.CreateTime = time.Now()
	device.UpdateTime = time.Now()
	err := db.DB.Create(&device).Error
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return device.Id, nil
}

// Get 获取设备
func (*deviceDao) Get(deviceId int64) (*model.Device, error) {
	var device = model.Device{Id: deviceId}
	err := db.DB.First(&device).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, gerrors.WrapError(err)
	}
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}
	return &device, err
}

// ListUserOnline 查询用户所有的在线设备
func (*deviceDao) ListOnlineByUserId(userId int64) ([]model.Device, error) {
	var devices []model.Device
	err := db.DB.Find(&devices, "user_id = ? and status = ?", userId, model.DeviceOnLine).Error
	if err != nil {
		return nil, gerrors.WrapError(err)
	}
	return devices, nil
}

// UpdateUserIdAndStatus 更新设备绑定用户和设备在线状态
func (*deviceDao) UpdateUserIdAndStatus(deviceId, userId int64, status int, connAddr string, connFd int64) error {
	err := db.DB.Exec("update device  set user_id = ?,status = ?,conn_addr = ?,conn_fd = ? where id = ? ",
		userId, status, connAddr, connFd, deviceId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// UpdateStatus 更新设备的在线状态
func (*deviceDao) UpdateStatus(deviceId int64, status int) error {
	err := db.DB.Exec("update device set status = ? where id = ?", status, deviceId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}

// Upgrade 升级设备
func (*deviceDao) Upgrade(deviceId int64, systemVersion, sdkVersion string) error {
	err := db.DB.Exec("update device set system_version = ?,sdk_version = ? where id = ? ",
		systemVersion, sdkVersion, deviceId).Error
	if err != nil {
		return gerrors.WrapError(err)
	}
	return nil
}
