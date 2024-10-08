<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <cmdb-sticky-layout class="details-layout">
    <slot name="prepend"></slot>
    <div ref="detailsWrapper">
      <slot name="details-header"></slot>
      <template v-for="(group, groupIndex) in $sortedGroups">
        <div class="property-group"
          :key="groupIndex"
          v-if="$groupedProperties[groupIndex].length">
          <cmdb-collapse
            :label="group['bk_group_name']"
            :collapse.sync="groupState[group['bk_group_id']]">
            <ul class="property-list clearfix">
              <li :class="['property-item clearfix fl', { flex: flexProperties.includes(property['bk_property_id']) }]"
                v-for="property in $groupedProperties[groupIndex]"
                :key="`${property['bk_obj_id']}-${property['bk_property_id']}`">
                <span class="property-name fl"
                  v-if="!invisibleNameProperties.includes(property['bk_property_id'])"
                  :title="property['bk_property_name']">{{property['bk_property_name']}}
                </span>
                <slot :name="property['bk_property_id']">
                  <span class="property-value clearfix fl"
                    v-if="property.unit"
                    v-bk-overflow-tips>
                    <span class="property-value-text fl">{{getValue(property)}}</span>
                    <span class="property-value-unit fl" v-if="getValue(property) !== '--'">{{property.unit}}</span>
                  </span>
                  <cmdb-property-value
                    v-else
                    v-bk-overflow-tips
                    :class="'property-value fl'"
                    :value="inst[property.bk_property_id]"
                    :property="property">
                  </cmdb-property-value>
                </slot>
              </li>
            </ul>
          </cmdb-collapse>
        </div>
      </template>
    </div>
    <div class="details-options" slot="footer" slot-scope="{ sticky }"
      v-if="showOptions"
      :class="{ sticky: sticky }">
      <slot name="details-options">
        <cmdb-auth v-if="showEdit" class="inline-block-middle" :auth="editAuth">
          <bk-button slot-scope="{ disabled }"
            class="button-edit"
            theme="primary"
            :disabled="disabled"
            @click="handleEdit">
            {{editText}}
          </bk-button>
        </cmdb-auth>
        <cmdb-auth v-if="showDelete" class="inline-block-middle" :auth="deleteAuth">
          <bk-button slot-scope="{ disabled }"
            hover-theme="danger"
            class="button-delete"
            :disabled="disabled"
            @click="handleDelete">
            {{deleteText}}
          </bk-button>
        </cmdb-auth>
      </slot>
    </div>
  </cmdb-sticky-layout>
</template>

<script>
  import formMixins from '@/mixins/form'
  import formatter from '@/filters/formatter.js'
  export default {
    name: 'cmdb-details',
    mixins: [formMixins],
    props: {
      inst: {
        type: Object,
        required: true
      },
      showOptions: {
        type: Boolean,
        default: true
      },
      editButtonText: {
        type: String,
        default: ''
      },
      deleteButtonText: {
        type: String,
        default: ''
      },
      showEdit: {
        type: Boolean,
        default: true
      },
      showDelete: {
        type: Boolean,
        default: true
      },
      editAuth: {
        type: Object,
        default: null
      },
      deleteAuth: {
        type: [Object, Array],
        default: null
      },
      flexProperties: {
        type: Array,
        default: () => []
      },
      invisibleNameProperties: {
        type: Array,
        default: () => []
      }
    },
    data() {
      return {
        resizeEvent: null
      }
    },
    computed: {
      editText() {
        return this.editButtonText || this.$t('编辑')
      },
      deleteText() {
        return this.deleteButtonText || this.$t('删除')
      }
    },
    methods: {
      handleToggleGroup(group) {
        const groupId = group.bk_group_id
        const collapse = !!this.groupState[groupId]
        this.$set(this.groupState, groupId, !collapse)
      },
      getTitle(inst, property) {
        return `${property.bk_property_name}: ${inst[property.bk_property_id] || '--'} ${property.unit}`
      },
      getValue(property) {
        return formatter(this.inst[property.bk_property_id], property)
      },
      handleEdit() {
        this.$emit('on-edit', this.inst)
      },
      handleDelete() {
        this.$emit('on-delete', this.inst)
      }
    }
  }
</script>

<style lang="scss" scoped>
    .details-layout {
        height: 100%;
        padding: 0 0 0 32px;
        @include scrollbar-y;
    }
    .property-group {
        padding: 7px 0 10px 0;
        &:first-child{
            padding: 28px 0 10px 0;
        }
    }
    .group-name {
        font-size: 16px;
        line-height: 16px;
        color: #333948;
        overflow: visible;
        .group-toggle {
            cursor: pointer;
            &.collapse .bk-icon {
                transform: rotate(-90deg);
            }
            .bk-icon {
                vertical-align: baseline;
                font-size: 12px;
                font-weight: bold;
                transition: transform .2s ease-in-out;
            }
        }
    }
    .property-list {
        padding: 4px 0;
        .property-item {
            width: 50%;
            max-width: 400px;
            margin: 12px 0 0;
            font-size: 14px;
            line-height: 26px;
            .property-name {
                position: relative;
                width: 35%;
                padding: 0 16px 0 0;
                color: #63656e;
                text-align: right;
                @include ellipsis;
                &:after{
                    content: ":";
                    position: absolute;
                    right: 10px;
                }
            }
            .property-value {
                width: 65%;
                padding: 0 15px 0 0;
                color: #313238;
                @include ellipsis;
                &-text {
                    display: block;
                    max-width: calc(100% - 60px);
                    @include ellipsis;
                }
                &-unit {
                    display: block;
                    width: 60px;
                    padding: 0 0 0 5px;
                    @include ellipsis;
                }
            }

            &.flex {
                display: flex;
                width: 100%;
                max-width: unset;
                padding-right: 15px;
            }
        }
    }
    .details-options {
        padding: 10px 18px;
        &.sticky {
            width: calc(100% + 32px);
            margin: 0 0 0 -40px;
            padding: 10px 50px;
            background-color: #fff;
            border-top: 1px solid $cmdbBorderColor;
        }
        .button-edit {
            min-width: 76px;
            margin-right: 4px;
        }
        .button-delete {
            min-width: 76px;
        }
    }
</style>
