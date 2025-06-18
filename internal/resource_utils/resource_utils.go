// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package resource_utils

// ToPointer - Utility to obtain a pointer to any rvalue without having to declare a local variable.
func ToPointer[T interface{}](object interface{}) *T {
	if object == nil {
		return nil
	}
	objectAsT := object.(T)
	return &objectAsT
}

func ToBoolPointer(object interface{}) *bool {
	if object == nil {
		return nil
	}
	objectAsBool := object.(bool)
	return &objectAsBool
}

func ToStringPointer(object interface{}) *string {
	if object == nil {
		return nil
	}
	objectAsString := object.(string)
	return &objectAsString
}

func ToInt32Pointer(object interface{}) *int32 {
	if object == nil {
		return nil
	}
	objectAsInt32 := int32(object.(int))
	return &objectAsInt32
}

func ToIntPointer(object interface{}) *int {
	if object == nil {
		return nil
	}
	objectAsInt := object.(int)
	return &objectAsInt
}

func ToStringSlice(params []interface{}) []string {
	var paramSlice []string
	for _, p := range params {
		if param, ok := p.(string); ok {
			paramSlice = append(paramSlice, param)
		}
	}
	return paramSlice
}

// CreateIdToObjectMap Creates a Map with string ID index to Object.
func CreateIdToObjectMap(objectsList []interface{}) map[string]interface{} {
	// crete a map of new host id -> host
	result := make(map[string]interface{})
	for _, listEntryRaw := range objectsList {
		listEntry := listEntryRaw.(map[string]interface{})
		id := listEntry["id"].(string)
		result[id] = listEntry
	}
	return result
}

// CalculateAddedRemovedResources utility method that provides the newly created or removed
// resources as a separate list, provided the new and old values of the resource list.
func CalculateAddedRemovedResources(newResourcesList, oldResourcesList []interface{}) (
	addedResources []map[string]interface{}, removedResources []map[string]interface{}) {
	isAddingResources := len(newResourcesList) > len(oldResourcesList)
	if isAddingResources {
		oldResourcesMap := CreateIdToObjectMap(oldResourcesList)
		for _, newHostListEntryRaw := range newResourcesList {
			newResourceListEntry := newHostListEntryRaw.(map[string]interface{})
			newHostEntryId := newResourceListEntry["id"].(string)
			_, currentResourceAlreadyPresent := oldResourcesMap[newHostEntryId]
			if !currentResourceAlreadyPresent {
				addedResources = append(addedResources, newResourceListEntry)
			}
		}
	} else {
		newResourcesMap := CreateIdToObjectMap(newResourcesList)
		for _, oldHostListEntryRaw := range oldResourcesList {
			oldResourceListEntry := oldHostListEntryRaw.(map[string]interface{})
			oldHostEntryId := oldResourceListEntry["id"].(string)
			_, currentResourceAlreadyPresent := newResourcesMap[oldHostEntryId]
			if !currentResourceAlreadyPresent {
				removedResources = append(removedResources, oldResourceListEntry)
			}
		}
	}

	return addedResources, removedResources
}
