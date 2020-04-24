// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

import (
	"encoding/json"
	"io"
)

const (
	SidebarCategoryChannels       SidebarCategoryType = "C"
	SidebarCategoryDirectMessages SidebarCategoryType = "D"
	SidebarCategoryFavorites      SidebarCategoryType = "F"
	SidebarCategoryCustom         SidebarCategoryType = "C"
)

type SidebarCategoryType string

type SidebarCategory struct {
	Id          string              `json:"id"`
	UserId      string              `json:"user_id"`
	TeamId      string              `json:"team_id"`
	SortOrder   int64               `json:"-"`
	Type        SidebarCategoryType `json:"type"`
	DisplayName string              `json:"display_name"`
}

type SidebarCategoryWithChannels struct {
	SidebarCategory
	Channels []string `json:"channel_ids"`
}

type SidebarCategoryOrder []string

type OrderedSidebarCategories struct {
	Categories SidebarCategoriesWithChannels `json:"categories"`
	Order      SidebarCategoryOrder          `json:"order"`
}

type SidebarChannel struct {
	ChannelId  string `json:"channel_id"`
	UserId     string `json:"user_id"`
	CategoryId string `json:"category_id"`
	SortOrder  int64  `json:"-"`
}

type SidebarChannels []*SidebarChannel
type SidebarCategoriesWithChannels []*SidebarCategoryWithChannels

func SidebarCategoryFromJson(data io.Reader) (*SidebarCategoryWithChannels, error) {
	var o *SidebarCategoryWithChannels
	err := json.NewDecoder(data).Decode(&o)
	return o, err
}

func SidebarCategoriesFromJson(data io.Reader) ([]*SidebarCategoryWithChannels, error) {
	var o []*SidebarCategoryWithChannels
	err := json.NewDecoder(data).Decode(&o)
	return o, err
}

func OrderedSidebarCategoriesFromJson(data io.Reader) (*OrderedSidebarCategories, error) {
	var o *OrderedSidebarCategories
	err := json.NewDecoder(data).Decode(&o)
	return o, err
}

func (o SidebarCategoryWithChannels) ToJson() []byte {
	b, _ := json.Marshal(o)
	return b
}

func SidebarCategoryWithChannelsToJson(o []*SidebarCategoryWithChannels) []byte {
	if b, err := json.Marshal(o); err != nil {
		return []byte("[]")
	} else {
		return b
	}
}

func (o OrderedSidebarCategories) ToJson() []byte {
	if b, err := json.Marshal(o); err != nil {
		return []byte("[]")
	} else {
		return b
	}
}
