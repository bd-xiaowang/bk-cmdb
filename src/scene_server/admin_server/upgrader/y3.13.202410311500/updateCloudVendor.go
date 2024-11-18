/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package y3_13_202410311500

import (
	"context"
	"errors"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

// attribute object attribute
type attribute struct {
	ID     int64        `json:"id" bson:"id"`
	Option []enumOption `json:"option" bson:"option"`
}

// enumOption enum option
type enumOption struct {
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	Type      string `json:"type" bson:"type"`
	IsDefault bool   `json:"is_default" bson:"is_default"`
}

func getOptionMap(vendors []string) map[string]string {
	optionMap := make(map[string]string, len(vendors))
	for index, cloudVendor := range vendors {
		optionMap[strconv.Itoa(index+1)] = cloudVendor
	}
	return optionMap
}

var (
	oldCloudVendors = []string{"AWS", "腾讯云", "GCP", "Azure", "企业私有云", "SalesForce", "Oracle Cloud", "IBM Cloud",
		"阿里云", "中国电信", "UCloud", "美团云", "金山云", "百度云", "华为云", "首都云"}

	newCloudVendors = []string{"亚马逊云", "腾讯云", "谷歌云", "微软云", "企业私有云", "SalesForce", "Oracle Cloud", "IBM Cloud",
		"阿里云", "中国电信", "UCloud", "美团云", "金山云", "百度云", "华为云", "首都云", "腾讯自研云", "Zenlayer"}

	oldOptionMap = getOptionMap(oldCloudVendors)

	newOptionMap = getOptionMap(newCloudVendors)
)

func updateCloudVendor(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	cond := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			common.BKDBIN: []string{common.BKInnerObjIDHost, common.BKInnerObjIDPlat},
		},
		common.BKPropertyIDField:   common.BKCloudVendor,
		common.BKPropertyTypeField: common.FieldTypeEnum,
	}

	objAttrs := make([]attribute, 0)

	err := db.Table(common.BKTableNameObjAttDes).Find(cond).All(ctx, &objAttrs)
	if err != nil {
		blog.Errorf("get cloud vendor field failed, err: %v", err)
		return err
	}
	if len(objAttrs) != 2 {
		blog.Errorf("get cloud vendor field failed, count cloud vendor field not equal 2, objAttrs: %v", objAttrs)
		return errors.New("count cloud vendor field not equal 2")
	}
	for _, attr := range objAttrs {
		for _, enum := range attr.Option {
			if oldOptionMap[enum.ID] != enum.Name && newOptionMap[enum.ID] != enum.Name {
				blog.Errorf("compare enum failed, currentEnum: %v", enum)
				return errors.New("compare enum failed")
			}
		}
	}

	newEnumOption := make([]enumOption, len(newCloudVendors))
	for index, cloudVendor := range newCloudVendors {
		newEnumOption[index] = enumOption{
			ID:        strconv.Itoa(index + 1),
			Name:      cloudVendor,
			Type:      "text",
			IsDefault: false,
		}
	}
	updateData := map[string]interface{}{common.BKOptionField: newEnumOption}

	if err := db.Table(common.BKTableNameObjAttDes).Update(ctx, cond, updateData); err != nil {
		blog.Errorf("update cloud vendor attribute failed, err: %v, cond: %v, updateData: %v", err, cond, updateData)
		return err
	}
	return nil
}
